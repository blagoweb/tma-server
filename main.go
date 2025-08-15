package main

import (
	"log"
	"net/http"

	"tma/auth"
	"tma/config"
	"tma/database"
	"tma/handlers"
	"tma/routes"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(cfg.JWTSecret)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db.DB, jwtManager, cfg.TelegramBotToken)
	pagesHandler := handlers.NewpagesHandler(db.DB)

	// Setup routes
	router := routes.SetupRoutes(authHandler, pagesHandler, jwtManager)

	// Start server
	log.Printf("Starting server on port %s", cfg.Port)
	log.Printf("Environment: %s", cfg.Environment)
	
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
