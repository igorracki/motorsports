package api

import (
	"log"
	"net/http"

	f1middleware "github.com/igorracki/motorsports/backend/internal/api/middleware"
	"github.com/igorracki/motorsports/backend/internal/auth"
	"github.com/igorracki/motorsports/backend/internal/config"
	"github.com/igorracki/motorsports/backend/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	echo         *echo.Echo
	config       *config.Configuration
	tokenManager auth.TokenManager
}

type ServerOption func(*Server)

func NewServer(config *config.Configuration, tokenManager auth.TokenManager, options ...ServerOption) *Server {
	e := echo.New()
	e.Validator = NewCustomValidator()
	e.HTTPErrorHandler = HTTPErrorHandler

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("1M"))
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(30)))
	e.Use(f1middleware.TraceLogger(f1middleware.TraceConfig{Enabled: config.TraceLogging}))

	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "header:X-CSRF-Token",
		CookieName:     "csrf_token",
		CookieHTTPOnly: false,
		CookieSameSite: http.SameSiteLaxMode,
		CookieSecure:   config.CookieSecure,
		CookiePath:     "/",
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     config.AllowedOrigins,
		AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	server := &Server{
		echo:         e,
		config:       config,
		tokenManager: tokenManager,
	}

	for _, option := range options {
		option(server)
	}

	return server
}

func (server *Server) RegisterRoutes(
	f1Handler *handlers.F1Handler,
	userHandler *handlers.UserHandler,
	authHandler *handlers.AuthHandler,
	predictionHandler *handlers.PredictionHandler,
	friendHandler *handlers.FriendHandler,
	leaderboardHandler *handlers.LeaderboardHandler,
	configHandler *handlers.ConfigHandler,
) {
	apiGroup := server.echo.Group("/api")
	authMiddleware := f1middleware.AuthMiddleware(server.tokenManager)
	ownerMiddleware := f1middleware.RequireResourceOwnerMiddleware()

	apiGroup.GET("/config", configHandler.GetConfig)

	apiGroup.GET("/schedule/:year", f1Handler.GetSchedule)
	apiGroup.GET("/schedule/:year/:round/:session/results", f1Handler.GetSessionResults)
	apiGroup.GET("/schedule/:year/:round/circuit", f1Handler.GetCircuit)
	apiGroup.GET("/schedule/:year/:round/drivers", f1Handler.GetDrivers)

	apiGroup.POST("/auth/register", authHandler.Register)
	apiGroup.POST("/auth/login", authHandler.Login)
	apiGroup.POST("/auth/logout", authHandler.Logout)
	apiGroup.GET("/auth/me", authHandler.Me, authMiddleware)

	apiGroup.GET("/users/:id", userHandler.GetUserProfile, authMiddleware, ownerMiddleware)
	apiGroup.POST("/users/:id/predictions", predictionHandler.SubmitPrediction, authMiddleware, ownerMiddleware)
	apiGroup.GET("/users/:id/predictions/:year/:round", predictionHandler.GetRoundPredictions, authMiddleware, ownerMiddleware)
	apiGroup.GET("/predictions/scoring-rules", predictionHandler.GetScoringRules)
	apiGroup.GET("/predictions/policy", predictionHandler.GetPredictionPolicy)

	apiGroup.POST("/users/friends/request", friendHandler.SendFriendRequest, authMiddleware)
	apiGroup.GET("/users/friends/requests", friendHandler.GetPendingRequests, authMiddleware)
	apiGroup.PUT("/users/friends/requests/:id", friendHandler.HandleFriendRequest, authMiddleware)
	apiGroup.GET("/users/friends/leaderboard/:season", leaderboardHandler.GetLeaderboard, authMiddleware)

	server.echo.HEAD("/health", func(context echo.Context) error {
		return context.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})
}

func (server *Server) Start(address string) error {
	log.Printf("Starting server on %s", address)
	return server.echo.Start(address)
}

func (server *Server) Echo() *echo.Echo {
	return server.echo
}
