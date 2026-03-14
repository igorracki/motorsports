package main

import (
	"fmt"
	"log"

	f1middleware "github.com/igorracki/f1/backend/internal/api/middleware"
	"github.com/igorracki/f1/backend/internal/cache"
	"github.com/igorracki/f1/backend/internal/clients"
	"github.com/igorracki/f1/backend/internal/config"
	"github.com/igorracki/f1/backend/internal/database"
	"github.com/igorracki/f1/backend/internal/handlers"
	"github.com/igorracki/f1/backend/internal/repository"
	"github.com/igorracki/f1/backend/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	configuration := config.Load()

	databaseManager, err := database.NewManager(configuration.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer databaseManager.Close()

	server := echo.New()

	server.Use(middleware.RequestLogger())
	server.Use(middleware.Recover())
	server.Use(middleware.Secure())
	server.Use(middleware.BodyLimit("1M"))
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     configuration.AllowedOrigins,
		AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	f1DataClient := clients.NewF1DataClient(configuration.ExternalAPIURL)
	f1MemoryCache := cache.NewMemoryCache()
	f1DataService := services.NewF1Service(f1DataClient, f1MemoryCache)

	userRepository := repository.NewUserRepository(databaseManager.DB())
	scoreRepository := repository.NewScoreRepository(databaseManager.DB())
	predictionRepository := repository.NewPredictionRepository(databaseManager.DB())

	scoringService := services.NewScoringService()
	predictionService := services.NewPredictionService(predictionRepository, f1DataService, scoringService)
	userService := services.NewUserService(userRepository, scoreRepository, predictionService)
	authService := services.NewAuthService(userRepository)

	f1DataHandler := handlers.NewF1Handler(f1DataService)
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService, userService)
	predictionHandler := handlers.NewPredictionHandler(predictionService, scoringService)

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

	apiGroup.GET("/users/:id", userHandler.GetUserProfile)
	apiGroup.GET("/users/:id/stats/seasons", userHandler.GetSeasonScores)
	apiGroup.POST("/users/:id/predictions", predictionHandler.SubmitPrediction)
	apiGroup.GET("/users/:id/predictions", predictionHandler.GetUserPredictions)
	apiGroup.GET("/users/:id/predictions/:year/:round", predictionHandler.GetRoundPredictions)

	server.GET("/health", func(context echo.Context) error {
		return context.JSON(200, map[string]string{"status": "ok"})
	})

	address := fmt.Sprintf(":%d", configuration.ServerPort)
	log.Printf("Starting server on %s", address)
	log.Fatal(server.Start(address))
}
