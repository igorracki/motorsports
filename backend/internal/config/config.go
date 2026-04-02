package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Configuration struct {
	ServerPort         int
	ExternalAPIURL     string
	ExternalAPITimeout int
	DatabasePath       string
	AllowedOrigins     []string
	CookieSecure       bool
	JWTSecret          string
	TraceLogging       bool
}

func Load() *Configuration {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	serverPort := getEnvAsInt("SERVER_PORT", 8080)
	externalAPIURL := getEnv("EXTERNAL_API_URL", "http://localhost:8081/wrapper")
	externalAPITimeout := getEnvAsInt("EXTERNAL_API_TIMEOUT", 30)
	databasePath := getEnv("DATABASE_PATH", "motorsports_data.db")
	allowedOrigins := getEnvAsSlice("ALLOWED_ORIGINS", nil)
	cookieSecure := getEnvAsBool("COOKIE_SECURE", false)
	jwtSecret := getEnv("JWT_SECRET", "")
	traceLogging := getEnvAsBool("TRACE_LOGGING", false)

	return &Configuration{
		ServerPort:         serverPort,
		ExternalAPIURL:     externalAPIURL,
		ExternalAPITimeout: externalAPITimeout,
		DatabasePath:       databasePath,
		AllowedOrigins:     allowedOrigins,
		CookieSecure:       cookieSecure,
		JWTSecret:          jwtSecret,
		TraceLogging:       traceLogging,
	}
}

func (c *Configuration) Validate() error {
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET environment variable is required")
	}

	if c.ServerPort < 1 || c.ServerPort > 65535 {
		return fmt.Errorf("invalid SERVER_PORT: %d", c.ServerPort)
	}

	if _, err := url.ParseRequestURI(c.ExternalAPIURL); err != nil {
		return fmt.Errorf("invalid EXTERNAL_API_URL: %w", err)
	}

	if c.DatabasePath == "" {
		return fmt.Errorf("DATABASE_PATH cannot be empty")
	}

	return nil
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
