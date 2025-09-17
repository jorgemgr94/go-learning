package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port int
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

	return Config{
		Port: port,
	}
}
