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

	"github.com/igorracki/motorsports/backend/internal/api"
	"github.com/igorracki/motorsports/backend/internal/auth"
	"github.com/igorracki/motorsports/backend/internal/clients"
	"github.com/igorracki/motorsports/backend/internal/config"
	"github.com/igorracki/motorsports/backend/internal/database"
	"github.com/igorracki/motorsports/backend/internal/handlers"
	"github.com/igorracki/motorsports/backend/internal/repository"
	"github.com/igorracki/motorsports/backend/internal/services"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	configuration := config.Load()
	if err := configuration.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	tokenManager, err := auth.NewTokenManager(configuration.JWTSecret)
	if err != nil {
		log.Fatalf("Failed to initialize token manager: %v", err)
	}

	databaseManager, err := database.NewManager(configuration.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer databaseManager.Close()

	// Clients
	f1DataClient := clients.NewF1DataClient(configuration.ExternalAPIURL, clients.WithTimeout(configuration.ExternalAPITimeout))

	// Repositories
	userRepository := repository.NewUserRepository(databaseManager)
	predictionRepository := repository.NewPredictionRepository(databaseManager)
	friendRepository := repository.NewFriendRepository(databaseManager)

	// Services
	predictionPolicy := services.NewPredictionPolicy()

	f1BaseService := services.NewF1Service(f1DataClient, predictionPolicy)
	f1DataService := services.NewF1CachingService(f1BaseService)
	defer f1DataService.Close()

	configService := services.NewConfigService()
	scoringService := services.NewScoringService()
	predictionService := services.NewPredictionService(predictionRepository, f1DataService, scoringService, predictionPolicy, configService)
	userService := services.NewUserService(userRepository, predictionService)
	authService := services.NewAuthService(userRepository, tokenManager)
	friendService := services.NewFriendService(friendRepository, userRepository)
	leaderboardService := services.NewLeaderboardService(friendRepository, userRepository, predictionRepository)

	// Handlers
	f1DataHandler := handlers.NewF1Handler(f1DataService)
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService, userService, handlers.WithCookieSecure(configuration.CookieSecure))
	predictionHandler := handlers.NewPredictionHandler(predictionService, scoringService)
	friendHandler := handlers.NewFriendHandler(friendService)
	leaderboardHandler := handlers.NewLeaderboardHandler(leaderboardService)
	configHandler := handlers.NewConfigHandler(configService)

	// API
	server := api.NewServer(configuration, tokenManager)
	server.RegisterRoutes(
		f1DataHandler,
		userHandler,
		authHandler,
		predictionHandler,
		friendHandler,
		leaderboardHandler,
		configHandler,
	)

	address := fmt.Sprintf(":%d", configuration.ServerPort)
	go func() {
		if err := server.Start(address); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Shutting down the server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Shutting down gracefully...")
	if err := server.Echo().Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server gracefully stopped")
}
