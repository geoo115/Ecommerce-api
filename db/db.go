package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/geoo115/Ecommerce/config"
	"github.com/geoo115/Ecommerce/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

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

	// Configure GORM with connection timeout and retry settings
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn), // Reduce log noise
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	// Add connection timeout for PostgreSQL with enhanced Render optimizations
	if strings.HasPrefix(dsn, "postgres://") || strings.Contains(dsn, "host=") {
		// Parse the DSN and add timeout parameters if not already present
		if !strings.Contains(dsn, "connect_timeout=") {
			separator := "?"
			if strings.Contains(dsn, "?") {
				separator = "&"
			}
			dsn += separator + "connect_timeout=60" // Increased for Render cold starts
		}
		if !strings.Contains(dsn, "statement_timeout=") {
			dsn += "&statement_timeout=120000" // 2 minutes for Render stability
		}
		// Add additional Render-specific parameters for connection reliability
		if !strings.Contains(dsn, "pool_max_conns=") {
			dsn += "&pool_max_conns=5" // Reduced for Render free tier stability
		}
		if !strings.Contains(dsn, "pool_timeout=") {
			dsn += "&pool_timeout=60" // Increased pool timeout
		}
		// Add keepalive settings for Render
		if !strings.Contains(dsn, "tcp_keepalives_idle=") {
			dsn += "&tcp_keepalives_idle=600" // 10 minutes
		}
		if !strings.Contains(dsn, "tcp_keepalives_interval=") {
			dsn += "&tcp_keepalives_interval=30" // 30 seconds
		}
		if !strings.Contains(dsn, "tcp_keepalives_count=") {
			dsn += "&tcp_keepalives_count=3" // 3 probes
		}
		dialector = postgres.Open(dsn)
	}

	// Retry connection with exponential backoff - optimized for Render
	var database *gorm.DB
	maxRetries := 12             // Increased for Render cold starts
	baseDelay := 5 * time.Second // Longer base delay for stability

	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Printf("Database connection attempt %d/%d...", attempt, maxRetries)

		database, err = gorm.Open(dialector, gormConfig)
		if err == nil {
			// Test the connection with extended timeout for Render
			sqlDB, sqlErr := database.DB()
			if sqlErr == nil {
				// Set a longer timeout for Render's managed database
				ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
				defer cancel()

				if pingErr := sqlDB.PingContext(ctx); pingErr == nil {
					log.Printf("Database connection and ping successful after %d attempt(s)", attempt)
					break
				} else {
					err = fmt.Errorf("ping failed: %w", pingErr)
				}
			} else {
				err = fmt.Errorf("failed to get sql.DB: %w", sqlErr)
			}
		}

		if attempt < maxRetries {
			delay := time.Duration(attempt) * baseDelay

			// Special handling for "unexpected EOF" which is common with Render cold starts
			errorMsg := err.Error()
			if strings.Contains(errorMsg, "unexpected EOF") {
				log.Printf("Render database cold start detected (attempt %d/%d). Waiting %v for warm-up...", attempt, maxRetries, delay)
			} else {
				log.Printf("Connection failed (attempt %d/%d): %v. Retrying in %v...", attempt, maxRetries, err, delay)
			}
			time.Sleep(delay)
		} else {
			return fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
		}
	}

	// Configure connection pool for Render's managed PostgreSQL
	sqlDB, err := database.DB()
	if err == nil {
		// Conservative pool settings for Render free tier
		sqlDB.SetMaxIdleConns(2)                   // Reduced from 10
		sqlDB.SetMaxOpenConns(5)                   // Reduced from 100
		sqlDB.SetConnMaxLifetime(30 * time.Minute) // Shorter lifetime
		sqlDB.SetConnMaxIdleTime(5 * time.Minute)  // Shorter idle time
		log.Printf("Connection pool configured successfully")
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
