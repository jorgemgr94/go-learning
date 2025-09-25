package config

import (
	"log"
	"os"
	"strconv"

	"go-learning/internal/db"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
)

type Config struct {
	Port int
	DB   db.Config
}

func LoadConfig() Config {
	// Try to load .env file (optional for local development)
	// Don't fail if .env file doesn't exist (for production deployment)
	_ = godotenv.Load()

	// Get port from environment variable (set by Fly.io)
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080"
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Printf("Error converting port to int: %v, using default port 8080", err)
		port = 8080
	}

	// Database configuration
	dbConfig := db.Config{
		Name:        getEnvWithDefault("APP_NAME", "go-learning-db"),
		Environment: getEnvWithDefault("ENVIRONMENT", "development"),
		Database:    os.Getenv("DB_NAME"),
		DBHost:      os.Getenv("DB_HOST"),
		DBPort:      os.Getenv("DB_PORT"),
		DBUser:      os.Getenv("DB_USER"),
		DBSecret:    os.Getenv("DB_PASS"),
		SSLMode:     getEnvWithDefault("DB_SSL_MODE", "disable"),
		Metrics:     prometheus.NewRegistry(),
	}

	return Config{
		Port: port,
		DB:   dbConfig,
	}
}

// getEnvWithDefault returns the environment variable value or a default if not set
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
