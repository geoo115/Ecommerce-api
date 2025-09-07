// Standalone database connection test
// Run with: go run test-db-connection.go
package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/geoo115/Ecommerce/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	log.Println("=== Database Connection Test ===")

	// Load config
	if err := config.LoadConfig(); err != nil {
		log.Printf("Config load error (non-fatal): %v", err)
	}

	// Get database URL
	dsn, err := config.GetDatabaseURL()
	if err != nil {
		log.Fatalf("Failed to get database URL: %v", err)
	}

	// Mask password for logging
	maskedDSN := dsn
	if len(dsn) > 50 {
		maskedDSN = dsn[:20] + "***MASKED***" + dsn[len(dsn)-20:]
	}
	log.Printf("Attempting connection with DSN: %s", maskedDSN)

	// Print individual environment variables (masked)
	log.Printf("Environment variables:")
	log.Printf("  DATABASE_HOST: %s", os.Getenv("DATABASE_HOST"))
	log.Printf("  DATABASE_PORT: %s", os.Getenv("DATABASE_PORT"))
	log.Printf("  DATABASE_USER: %s", os.Getenv("DATABASE_USER"))
	log.Printf("  DATABASE_PASSWORD: %s", maskString(os.Getenv("DATABASE_PASSWORD")))
	log.Printf("  DATABASE_NAME: %s", os.Getenv("DATABASE_NAME"))
	log.Printf("  DATABASE_SSLMODE: %s", os.Getenv("DATABASE_SSLMODE"))
	log.Printf("  DATABASE_URL: %s", maskString(os.Getenv("DATABASE_URL")))

	// Test connection with timeout
	log.Println("Testing database connection...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying sql.DB: %v", err)
	}
	defer sqlDB.Close()

	// Test ping
	if err := sqlDB.PingContext(ctx); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	log.Println("✅ Database connection successful!")

	// Test a simple query
	var version string
	if err := db.Raw("SELECT version()").Scan(&version).Error; err != nil {
		log.Printf("❌ Failed to query database version: %v", err)
	} else {
		log.Printf("✅ Database version: %.100s...", version)
	}

	// Test database stats
	stats := sqlDB.Stats()
	log.Printf("Connection stats: OpenConnections=%d, InUse=%d, Idle=%d",
		stats.OpenConnections, stats.InUse, stats.Idle)
}

func maskString(s string) string {
	if s == "" {
		return "<not set>"
	}
	if len(s) <= 8 {
		return "***"
	}
	return s[:3] + "***" + s[len(s)-3:]
}
