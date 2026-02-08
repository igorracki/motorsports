package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Configuration struct {
	ServerPort     int
	ExternalAPIURL string
}

func Load() *Configuration {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	serverPort := getEnvAsInt("SERVER_PORT", 8081)
	externalAPIURL := getEnv("EXTERNAL_API_URL", "http://localhost:8080/wrapper")

	return &Configuration{
		ServerPort:     serverPort,
		ExternalAPIURL: externalAPIURL,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if valueStr := os.Getenv(key); valueStr != "" {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}
