package main

import (
	"fmt"
	"log"

	"github.com/igorracki/f1/backend/internal/clients"
	"github.com/igorracki/f1/backend/internal/config"
	"github.com/igorracki/f1/backend/internal/handlers"
	"github.com/igorracki/f1/backend/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	configuration := config.Load()
	server := echo.New()

	server.Use(middleware.RequestLogger())
	server.Use(middleware.Recover())

	externalClient := clients.NewExternalAPIClient(configuration.ExternalAPIURL)
	eventsService := services.NewEventsService(externalClient)
	eventsHandler := handlers.NewEventsHandler(eventsService)

	api := server.Group("/api")
	api.GET("/events", eventsHandler.GetEvents)

	server.GET("/health", func(context echo.Context) error {
		return context.JSON(200, map[string]string{"status": "ok"})
	})

	address := fmt.Sprintf(":%d", configuration.ServerPort)
	log.Printf("Starting server on %s", address)
	log.Fatal(server.Start(address))
}
