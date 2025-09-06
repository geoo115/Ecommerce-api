package db

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/geoo115/Ecommerce/models"
	"github.com/geoo115/Ecommerce/utils"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// OptimizedDB represents an optimized database connection
type OptimizedDB struct {
	*gorm.DB
}

// DBConfig holds database configuration
type DBConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	DBName       string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

// GetDBConfig returns optimized database configuration
func GetDBConfig() *DBConfig {
	config := &DBConfig{
		Host:         getEnv("DB_HOST", "localhost"),
		Port:         getEnvAsInt("DB_PORT", 5432),
		User:         getEnv("DB_USER", "postgres"),
		Password:     getEnv("DB_PASSWORD", "password"),
		DBName:       getEnv("DB_NAME", "ecommerce"),
		SSLMode:      getEnv("DB_SSLMODE", "disable"),
		MaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 50),
		MaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 25),
		MaxLifetime:  getEnvAsDuration("DB_MAX_LIFETIME", 5*time.Minute),
	}
	return config
}

// ConnectOptimizedDB creates an optimized database connection
func ConnectOptimizedDB() (*OptimizedDB, error) {
	config := GetDBConfig()

	var dialector gorm.Dialector
	if os.Getenv("DB_DRIVER") == "sqlite" {
		dialector = sqlite.Open("ecommerce_optimized.db")
	} else {
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)
		dialector = postgres.Open(dsn)
	}

	// Configure GORM with performance optimizations
	gormConfig := &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Warn,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		),
		PrepareStmt:       true,  // Cache prepared statements
		QueryFields:       true,  // Select by fields
		AllowGlobalUpdate: false, // Prevent accidental global updates
	}

	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL DB for connection pool configuration
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying SQL DB: %w", err)
	}

	// Configure connection pool with optimized settings for high load
	sqlDB.SetMaxOpenConns(getEnvAsInt("DB_MAX_OPEN_CONNS", 100))
	sqlDB.SetMaxIdleConns(getEnvAsInt("DB_MAX_IDLE_CONNS", 50))
	sqlDB.SetConnMaxLifetime(getEnvAsDuration("DB_MAX_LIFETIME", 10*time.Minute))

	utils.Info("Database connected with optimized settings:")
	utils.Info("Max Open Connections: %d", config.MaxOpenConns)
	utils.Info("Max Idle Connections: %d", config.MaxIdleConns)
	utils.Info("Connection Max Lifetime: %v", config.MaxLifetime)

	return &OptimizedDB{DB: db}, nil
}

// CreateOptimizedIndexes creates performance indexes
func (odb *OptimizedDB) CreateOptimizedIndexes() error {
	var indexes []string

	// Check if we're using SQLite (for tests) or PostgreSQL
	isSQLite := os.Getenv("DB_DRIVER") == "sqlite"

	if isSQLite {
		// SQLite-compatible indexes
		indexes = []string{
			// User indexes
			"CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)",
			"CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at)",

			// Product indexes
			"CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id)",
			"CREATE INDEX IF NOT EXISTS idx_products_price ON products(price)",
			"CREATE INDEX IF NOT EXISTS idx_products_created_at ON products(created_at)",

			// Order indexes
			"CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id)",
			"CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status)",
			"CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at)",
			"CREATE INDEX IF NOT EXISTS idx_orders_user_status ON orders(user_id, status)",

			// Order items indexes
			"CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id)",
			"CREATE INDEX IF NOT EXISTS idx_order_items_product_id ON order_items(product_id)",

			// Cart indexes
			"CREATE INDEX IF NOT EXISTS idx_carts_user_id ON carts(user_id)",
			"CREATE INDEX IF NOT EXISTS idx_carts_product_id ON carts(product_id)",
			"CREATE INDEX IF NOT EXISTS idx_carts_user_product ON carts(user_id, product_id)",

			// Payment indexes
			"CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments(order_id)",
			"CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status)",
			"CREATE INDEX IF NOT EXISTS idx_payments_created_at ON payments(created_at)",

			// Review indexes
			"CREATE INDEX IF NOT EXISTS idx_reviews_product_id ON reviews(product_id)",
			"CREATE INDEX IF NOT EXISTS idx_reviews_user_id ON reviews(user_id)",
			"CREATE INDEX IF NOT EXISTS idx_reviews_rating ON reviews(rating)",
		}
	} else {
		// PostgreSQL indexes with advanced features
		indexes = []string{
			// User indexes
			"CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)",
			"CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at)",

			// Product indexes
			"CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id)",

			"CREATE INDEX IF NOT EXISTS idx_products_price ON products(price)",
			"CREATE INDEX IF NOT EXISTS idx_products_created_at ON products(created_at)",

			// Order indexes
			"CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id)",
			"CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status)",
			"CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at)",
			"CREATE INDEX IF NOT EXISTS idx_orders_user_status ON orders(user_id, status)",

			// Order items indexes
			"CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id)",
			"CREATE INDEX IF NOT EXISTS idx_order_items_product_id ON order_items(product_id)",

			// Cart indexes
			"CREATE INDEX IF NOT EXISTS idx_carts_user_id ON carts(user_id)",
			"CREATE INDEX IF NOT EXISTS idx_carts_product_id ON carts(product_id)",
			"CREATE INDEX IF NOT EXISTS idx_carts_user_product ON carts(user_id, product_id)",

			// Payment indexes
			"CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments(order_id)",
			"CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status)",
			"CREATE INDEX IF NOT EXISTS idx_payments_created_at ON payments(created_at)",

			// Review indexes
			"CREATE INDEX IF NOT EXISTS idx_reviews_product_id ON reviews(product_id)",
			"CREATE INDEX IF NOT EXISTS idx_reviews_user_id ON reviews(user_id)",
			"CREATE INDEX IF NOT EXISTS idx_reviews_rating ON reviews(rating)",
		}
	}

	for _, index := range indexes {
		if err := odb.Exec(index).Error; err != nil {
			utils.Warn("Failed to create index: %v", err)
			// Continue with other indexes even if one fails
		}
	}

	utils.Info("Database indexes optimization completed")
	return nil
}

// OptimizeQueries enables query optimizations
func (odb *OptimizedDB) OptimizeQueries() {
	// Preload commonly accessed relationships using reusable callback functions
	odb.DB = odb.DB.Preload("User", preloadUserSelect).Preload("Category", preloadCategorySelect)
}

// preloadUserSelect selects a minimal user projection for preloads
func preloadUserSelect(db *gorm.DB) *gorm.DB {
	return db.Select("id, username, email")
}

// preloadCategorySelect selects a minimal category projection for preloads
func preloadCategorySelect(db *gorm.DB) *gorm.DB {
	return db.Select("id, name")
}

// HealthCheck performs database health check with performance metrics
func (odb *OptimizedDB) HealthCheck() (map[string]interface{}, error) {
	sqlDB, err := odb.DB.DB()
	if err != nil {
		return nil, err
	}

	stats := sqlDB.Stats()

	health := map[string]interface{}{
		"status":               "healthy",
		"open_connections":     stats.OpenConnections,
		"in_use_connections":   stats.InUse,
		"idle_connections":     stats.Idle,
		"max_open_connections": stats.MaxOpenConnections,
		"max_idle_connections": stats.MaxIdleClosed,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration.String(),
		"max_idle_time_closed": stats.MaxIdleTimeClosed,
	}

	// Test query performance
	start := time.Now()
	var count int64
	if err := odb.DB.Model(&models.User{}).Count(&count).Error; err != nil {
		health["status"] = "unhealthy"
		health["error"] = err.Error()
	}
	health["query_time"] = time.Since(start).String()

	return health, nil
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets environment variable as integer with default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsDuration gets environment variable as duration with default value
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
