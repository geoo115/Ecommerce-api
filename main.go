package main

import (
	"github/geoo115/Ecommerce/api"
	"github/geoo115/Ecommerce/config"
	"github/geoo115/Ecommerce/db"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize Gin router
	r := gin.Default()

	// Initialize database
	db.ConnectDatabase()

	// Set up routes
	api.SetupRoutes(r)

	// Run the server
	r.Run(":8080")
}
