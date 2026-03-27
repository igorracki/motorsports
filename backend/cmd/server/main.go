package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	f1middleware "github.com/igorracki/motorsports/backend/internal/api/middleware"
	"github.com/igorracki/motorsports/backend/internal/auth"
	"github.com/igorracki/motorsports/backend/internal/cache"

	"github.com/igorracki/motorsports/backend/internal/clients"
	"github.com/igorracki/motorsports/backend/internal/config"
	"github.com/igorracki/motorsports/backend/internal/database"
	"github.com/igorracki/motorsports/backend/internal/handlers"
	"github.com/igorracki/motorsports/backend/internal/repository"
	"github.com/igorracki/motorsports/backend/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	configuration := config.Load()

	if err := auth.InitJWTSecret(); err != nil {
		log.Fatalf("Failed to initialize JWT secret: %v", err)
	}

	databaseManager, err := database.NewManager(configuration.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer databaseManager.Close()

	server := echo.New()
	server.Validator = &CustomValidator{validator: validator.New()}

	server.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogMethod:   true,
		LogRemoteIP: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			slog.InfoContext(c.Request().Context(), "request",
				slog.Int("status", v.Status),
				slog.String("method", v.Method),
				slog.String("uri", v.URI),
				slog.String("remote_ip", v.RemoteIP),
			)
			return nil
		},
	}))
	server.Use(middleware.Recover())
	server.Use(middleware.Secure())
	server.Use(middleware.BodyLimit("1M"))
	server.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(30))) // 30 requests per second
	server.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "header:X-CSRF-Token",
		CookieName:     "csrf_token",
		CookieHTTPOnly: false, // Allow frontend to read the token
		CookieSameSite: http.SameSiteLaxMode,
		CookieSecure:   configuration.CookieSecure,
		CookiePath:     "/",
	}))
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     configuration.AllowedOrigins,
		AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	f1DataClient := clients.NewF1DataClient(configuration.ExternalAPIURL)
	f1MemoryCache := cache.NewMemoryCache()
	defer f1MemoryCache.Close()
	f1DataService := services.NewF1Service(f1DataClient, f1MemoryCache)

	userRepository := repository.NewUserRepository(databaseManager.DB())
	scoreRepository := repository.NewScoreRepository(databaseManager.DB())
	predictionRepository := repository.NewPredictionRepository(databaseManager.DB())
	friendRepository := repository.NewFriendRepository(databaseManager.DB())

	scoringService := services.NewScoringService()
	predictionService := services.NewPredictionService(predictionRepository, f1DataService, scoringService)
	userService := services.NewUserService(userRepository, scoreRepository, predictionService)
	authService := services.NewAuthService(userRepository)
	friendService := services.NewFriendService(friendRepository, userRepository)
	leaderboardService := services.NewLeaderboardService(friendRepository, userRepository, scoreRepository)

	f1DataHandler := handlers.NewF1Handler(f1DataService)
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService, userService, configuration.CookieSecure)
	predictionHandler := handlers.NewPredictionHandler(predictionService, scoringService)
	friendHandler := handlers.NewFriendHandler(friendService)
	leaderboardHandler := handlers.NewLeaderboardHandler(leaderboardService)

	apiGroup := server.Group("/api")
	apiGroup.GET("/schedule/:year", f1DataHandler.GetSchedule)
	apiGroup.GET("/schedule/:year/:round/:session/results", f1DataHandler.GetSessionResults)
	apiGroup.GET("/schedule/:year/:round/circuit", f1DataHandler.GetCircuit)
	apiGroup.GET("/schedule/:year/:round/drivers", f1DataHandler.GetDrivers)
	apiGroup.GET("/predictions/scoring-rules", predictionHandler.GetScoringRules)

	apiGroup.POST("/auth/register", authHandler.Register)
	apiGroup.POST("/auth/login", authHandler.Login)
	apiGroup.POST("/auth/logout", authHandler.Logout)
	apiGroup.GET("/auth/me", authHandler.Me, f1middleware.AuthMiddleware)

	apiGroup.GET("/users/:id", userHandler.GetUserProfile, f1middleware.AuthMiddleware)
	apiGroup.GET("/users/:id/stats/seasons", userHandler.GetSeasonScores, f1middleware.AuthMiddleware)
	apiGroup.POST("/users/:id/predictions", predictionHandler.SubmitPrediction, f1middleware.AuthMiddleware)
	apiGroup.GET("/users/:id/predictions", predictionHandler.GetUserPredictions, f1middleware.AuthMiddleware)
	apiGroup.GET("/users/:id/predictions/:year/:round", predictionHandler.GetRoundPredictions, f1middleware.AuthMiddleware)

	apiGroup.POST("/users/friends/request", friendHandler.SendFriendRequest, f1middleware.AuthMiddleware)
	apiGroup.GET("/users/friends/requests", friendHandler.GetPendingRequests, f1middleware.AuthMiddleware)
	apiGroup.PUT("/users/friends/requests/:id", friendHandler.HandleFriendRequest, f1middleware.AuthMiddleware)
	apiGroup.GET("/users/friends/leaderboard/:season", leaderboardHandler.GetLeaderboard, f1middleware.AuthMiddleware)

	server.HEAD("/health", func(context echo.Context) error {
		return context.JSON(200, map[string]string{"status": "ok"})
	})

	address := fmt.Sprintf(":%d", configuration.ServerPort)
	log.Printf("Starting server on %s", address)

	go func() {
		if err := server.Start(address); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Shutting down the server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Shutting down gracefully...")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server gracefully stopped")
}
