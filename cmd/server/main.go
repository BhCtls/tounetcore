package main

import (
	"log"
	"os"

	"tounetcore/internal/api"
	"tounetcore/internal/config"
	"tounetcore/internal/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	router := gin.Default()

	// Setup routes
	api.SetupRoutes(router, db, cfg)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "44544"
	}

	// Listen on all interfaces (0.0.0.0)
	address := "0.0.0.0:" + port
	log.Printf("Server starting on %s", address)
	if err := router.Run(address); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
