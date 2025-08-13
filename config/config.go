package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL       string
	JWTSecret         string
	Port              string
	TelegramBotToken  string
	Environment       string
}

func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config := &Config{
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/tma_db?sslmode=disable"),
		JWTSecret:        getEnv("JWT_SECRET", "default-jwt-secret-change-in-production"),
		Port:             getEnv("PORT", "8080"),
		TelegramBotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
		Environment:      getEnv("ENV", "development"),
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
