package db

import (
	"github/geoo115/Ecommerce/config"
	"github/geoo115/Ecommerce/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := config.GetDatabaseURL()
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	database.AutoMigrate(&models.User{}, &models.Product{}, &models.Cart{}, &models.Address{})
	DB = database
}
