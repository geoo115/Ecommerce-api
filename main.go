package main

import (
	"log"
	"os"

	"github.com/geoo115/Ecommerce/api"
	"github.com/geoo115/Ecommerce/api/middlewares"
	"github.com/geoo115/Ecommerce/config"
	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize logger with environment configuration
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	utils.AppLogger.SetLogLevelFromString(logLevel)

	utils.Info("Starting Ecommerce API server...")

	// Initialize Gin router
	r := gin.Default()

	// Add middleware in order
	r.Use(middlewares.CORSMiddleware())
	r.Use(middlewares.SecureHeadersMiddleware())
	r.Use(middlewares.LoggingMiddleware())
	r.Use(middlewares.RequestTimingMiddleware())
	r.Use(middlewares.ErrorLoggingMiddleware())
	r.Use(middlewares.GeneralRateLimit())

	// Initialize database
	utils.Info("Connecting to database...")
	db.ConnectDatabase()
	utils.Info("Database connected successfully")

	// Set up routes
	utils.Info("Setting up routes...")
	api.SetupRoutes(r)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Run the server
	utils.Info("Server starting on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		utils.Fatal("Failed to start server: %v", err)
	}
}
