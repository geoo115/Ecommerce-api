package db

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"

	"github.com/geoo115/Ecommerce/config"
	"github.com/geoo115/Ecommerce/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

var DB *gorm.DB

func createDatabaseIfNotExists(dsn string) error {
	// Parse the DSN to extract components
	u, err := url.Parse(dsn)
	if err != nil {
		return fmt.Errorf("failed to parse DSN: %w", err)
	}

	// Remove the dbname from DSN
	q := u.Query()
	q.Del("dbname")
	u.RawQuery = q.Encode()
	u.Path = ""

	// Add sslmode=disable if not present
	if _, ok := q["sslmode"]; !ok {
		q.Set("sslmode", "disable")
		u.RawQuery = q.Encode()
	}

	connStr := u.String()
	dbName := "ecommerce"

	// Connect to PostgreSQL without specifying a database
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}
	defer conn.Close()

	// Check if the database exists
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = '%s')", dbName)
	err = conn.QueryRow(query).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	// Create the database if it does not exist
	if !exists {
		_, err = conn.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
	}

	return nil
}

func ConnectDatabase() {
	dsn := config.GetDatabaseURL()

	// Attempt to create the database if it doesn't exist
	err := createDatabaseIfNotExists(dsn)
	if err != nil {
		log.Fatal("Failed to ensure database exists:", err)
	}

	// Connect to the specified database
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the models
	database.AutoMigrate(&models.User{}, &models.Product{}, &models.Cart{}, &models.Address{})
	DB = database
}
