package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Configuration struct {
	ServerPort     int
	ExternalAPIURL string
	DatabasePath   string
	AllowedOrigins []string
	CookieSecure   bool
}

func Load() *Configuration {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	serverPort := getEnvAsInt("SERVER_PORT", 8080)
	externalAPIURL := getEnv("EXTERNAL_API_URL", "http://localhost:8081/wrapper")
	databasePath := getEnv("DATABASE_PATH", "motorsports_data.db")
	allowedOrigins := getEnvAsSlice("ALLOWED_ORIGINS", nil)
	cookieSecure := getEnvAsBool("COOKIE_SECURE", false)

	if len(allowedOrigins) == 0 {
		log.Println("Warning: ALLOWED_ORIGINS is empty. CORS will block all requests.")
	}

	return &Configuration{
		ServerPort:     serverPort,
		ExternalAPIURL: externalAPIURL,
		DatabasePath:   databasePath,
		AllowedOrigins: allowedOrigins,
		CookieSecure:   cookieSecure,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	// Split by comma
	origins := strings.Split(valueStr, ",")
	for i, origin := range origins {
		origins[i] = strings.TrimSpace(origin)
	}
	return origins
}

func getEnvAsInt(key string, defaultValue int) int {
	if valueStr := os.Getenv(key); valueStr != "" {
		value, err := strconv.Atoi(valueStr)
		if err == nil {
			return value
		}
		log.Printf("Warning: Failed to parse environment variable %s=%s as int: %v. Using default: %d", key, valueStr, err, defaultValue)
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if valueStr := os.Getenv(key); valueStr != "" {
		value, err := strconv.ParseBool(valueStr)
		if err == nil {
			return value
		}
		log.Printf("Warning: Failed to parse environment variable %s=%s as bool: %v. Using default: %t", key, valueStr, err, defaultValue)
	}
	return defaultValue
}
