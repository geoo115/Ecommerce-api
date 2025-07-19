package main

import (
	"log"

	"github.com/geoo115/Ecommerce/api"
	"github.com/geoo115/Ecommerce/api/middlewares"
	"github.com/geoo115/Ecommerce/config"
	"github.com/geoo115/Ecommerce/db"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// Add CORS and security middleware
	r.Use(middlewares.CORSMiddleware())
	r.Use(middlewares.SecureHeadersMiddleware())

	// Initialize database
	db.ConnectDatabase()

	// Set up routes
	api.SetupRoutes(r)

	// Run the server
	log.Println("Server starting on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
