package db

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/geoo115/Ecommerce/config"
	"github.com/geoo115/Ecommerce/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

var DB *gorm.DB

func createDatabaseIfNotExists(dsn string) error {
	// Skip database creation for managed services
	// Look for indicators of managed database services
	if strings.Contains(dsn, "render.com") ||
		strings.Contains(dsn, "amazonaws.com") ||
		strings.Contains(dsn, "railway.app") ||
		strings.Contains(dsn, "planetscale.com") ||
		strings.Contains(dsn, "sslmode=require") {
		log.Printf("Detected managed database service, skipping database creation check")
		return nil
	}

	// Parse the DSN to extract components
	u, err := url.Parse(dsn)
	if err != nil {
		return fmt.Errorf("failed to parse DSN: %w", err)
	}

	// Basic validation for parsed DSN
	if u.Host == "" || u.Path == "" {
		return fmt.Errorf("invalid DSN: missing host or path")
	}

	// Extract database name from DSN
	dbName := u.Path[1:] // Remove leading '/'
	if dbName == "" {
		return fmt.Errorf("database name is missing in DSN")
	}

	// Remove dbname from DSN for the connection
	q := u.Query()
	q.Del("dbname")
	u.RawQuery = q.Encode()
	u.Path = "/postgres" // Use the "postgres" system database for initial connection

	// Only set sslmode=disable if not already specified and not a managed service
	if _, ok := q["sslmode"]; !ok {
		q.Set("sslmode", "disable")
		u.RawQuery = q.Encode()
	}

	connStr := u.String()

	// Connect to PostgreSQL using the "postgres" database
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}
	defer conn.Close()

	// Check if the database exists using parameterized query
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)"
	err = conn.QueryRow(query, dbName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	// Create the database if it does not exist
	if !exists {
		log.Printf("Database %s does not exist. Creating it...", dbName)
		// Use parameterized query for database creation
		_, err = conn.Exec("CREATE DATABASE " + dbName)
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
	}

	return nil
}

// ConnectDatabase connects to the configured database. It returns an error
// instead of exiting the program so callers (and tests) can handle failures.
func ConnectDatabase() error {
	dsn, err := config.GetDatabaseURL()
	if err != nil {
		return fmt.Errorf("failed to get database URL: %w", err)
	}

	// Log connection attempt (but mask sensitive details)
	maskedDSN := dsn
	if strings.Contains(dsn, "://") {
		parts := strings.Split(dsn, "://")
		if len(parts) == 2 {
			// Extract just the host part for logging
			hostPart := strings.Split(parts[1], "@")
			if len(hostPart) > 1 {
				maskedDSN = parts[0] + "://*****@" + hostPart[len(hostPart)-1]
			}
		}
	}
	log.Printf("Attempting to connect to database: %s", maskedDSN)

	// If the DSN looks like Postgres, attempt to ensure database exists
	if strings.HasPrefix(dsn, "postgres://") || strings.Contains(dsn, "host=") {
		if err := createDatabaseIfNotExists(dsn); err != nil {
			return fmt.Errorf("failed to ensure database exists: %w", err)
		}
	}

	// Choose driver based on DSN heuristics (postgres vs sqlite)
	var dialector gorm.Dialector
	if strings.HasPrefix(dsn, "postgres://") || strings.Contains(dsn, "host=") {
		dialector = postgres.Open(dsn)
	} else {
		// Treat non-postgres DSNs as sqlite paths (file: or memory)
		dialector = sqlite.Open(dsn)
	}

	log.Printf("Opening database connection...")
	database, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Printf("Database connection successful, running auto-migration...")
	// Auto-migrate the models
	if err := database.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.Cart{},
		&models.Order{},
		&models.OrderItem{},
		&models.Payment{},
		&models.Address{},
		&models.Review{},
		&models.Wishlist{},
		&models.Inventory{},
	); err != nil {
		// AutoMigrate failing is not fatal for tests, but log it
		log.Printf("auto migrate failed: %v", err)
	}

	DB = database
	return nil
}
