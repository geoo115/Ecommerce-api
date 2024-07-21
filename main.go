package main

import (
	"github.com/geoo115/Ecommerce/api"
	"github.com/geoo115/Ecommerce/config"
	"github.com/geoo115/Ecommerce/db"
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
