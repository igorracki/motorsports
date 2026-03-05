package main

import (
	"fmt"
	"log"

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
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	f1DataClient := clients.NewF1DataClient(configuration.ExternalAPIURL)
	f1MemoryCache := cache.NewMemoryCache()
	f1DataService := services.NewF1Service(f1DataClient, f1MemoryCache)
	f1DataHandler := handlers.NewF1Handler(f1DataService)

	userRepository := repository.NewUserRepository(databaseManager.DB())
	scoreRepository := repository.NewScoreRepository(databaseManager.DB())
	userService := services.NewUserService(userRepository, scoreRepository)
	userHandler := handlers.NewUserHandler(userService)

	predictionRepository := repository.NewPredictionRepository(databaseManager.DB())
	predictionService := services.NewPredictionService(predictionRepository)
	predictionHandler := handlers.NewPredictionHandler(predictionService)

	apiGroup := server.Group("/api")
	apiGroup.GET("/schedule/:year", f1DataHandler.GetSchedule)
	apiGroup.GET("/schedule/:year/:round/:session/results", f1DataHandler.GetSessionResults)
	apiGroup.GET("/schedule/:year/:round/circuit", f1DataHandler.GetCircuit)
	apiGroup.GET("/schedule/:year/:round/drivers", f1DataHandler.GetDrivers)

	apiGroup.POST("/users", userHandler.RegisterUser)
	apiGroup.GET("/users/:id", userHandler.GetUserProfile)
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
