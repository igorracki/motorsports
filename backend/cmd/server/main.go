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

	f1Client := clients.NewF1DataClient(configuration.ExternalAPIURL)
	f1Service := services.NewF1Service(f1Client)
	f1Handler := handlers.NewF1Handler(f1Service)

	api := server.Group("/api")
	api.GET("/events", f1Handler.GetEvents)

	server.GET("/health", func(context echo.Context) error {
		return context.JSON(200, map[string]string{"status": "ok"})
	})

	address := fmt.Sprintf(":%d", configuration.ServerPort)
	log.Printf("Starting server on %s", address)
	log.Fatal(server.Start(address))
}
