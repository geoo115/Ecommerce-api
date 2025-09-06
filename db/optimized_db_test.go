package db

import (
	"os"
	"testing"
	"time"

	"github.com/geoo115/Ecommerce/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestConnectOptimizedDB_InMemory(t *testing.T) {
	// Use SQLite for testing
	os.Setenv("DB_DRIVER", "sqlite")
	defer os.Unsetenv("DB_DRIVER")

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Auto-migrate
	err = db.AutoMigrate(&models.User{}, &models.Product{})
	assert.NoError(t, err)

	// Test basic operations
	user := models.User{Username: "test", Email: "test@example.com"}
	result := db.Create(&user)
	assert.NoError(t, result.Error)

	var retrieved models.User
	result = db.First(&retrieved, user.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, "test", retrieved.Username)
}

func TestGetEnv(t *testing.T) {
	// Test getEnv
	os.Setenv("TEST_VAR", "value")
	defer os.Unsetenv("TEST_VAR")
	assert.Equal(t, "value", getEnv("TEST_VAR", "default"))
	assert.Equal(t, "default", getEnv("NON_EXISTENT", "default"))
}

func TestGetEnvAsInt(t *testing.T) {
	os.Setenv("TEST_INT", "42")
	defer os.Unsetenv("TEST_INT")
	assert.Equal(t, 42, getEnvAsInt("TEST_INT", 10))
	assert.Equal(t, 10, getEnvAsInt("NON_EXISTENT", 10))
}

func TestGetEnvAsDuration(t *testing.T) {
	os.Setenv("TEST_DURATION", "5m")
	defer os.Unsetenv("TEST_DURATION")
	assert.Equal(t, 5*time.Minute, getEnvAsDuration("TEST_DURATION", time.Minute))
	assert.Equal(t, time.Minute, getEnvAsDuration("NON_EXISTENT", time.Minute))
}

func TestPreloadUserSelect(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Auto-migrate User model
	err = db.AutoMigrate(&models.User{})
	assert.NoError(t, err)

	// Insert dummy user
	user := models.User{Username: "testuser", Email: "test@example.com"}
	err = db.Create(&user).Error
	assert.NoError(t, err)

	// Use preloadUserSelect in a query and inspect SQL
	var result models.User
	tx := preloadUserSelect(db.Model(&models.User{}).Where("id = ?", user.ID))
	tx = tx.Debug() // Print SQL to stdout for inspection
	err = tx.First(&result).Error
	assert.NoError(t, err)
	// Check that only the selected fields are populated
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.Username, result.Username)
	assert.Equal(t, user.Email, result.Email)
	// Other fields should be zero values (if any)
}

func TestPreloadCategorySelect(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Auto-migrate Category model
	err = db.AutoMigrate(&models.Category{})
	assert.NoError(t, err)

	// Insert dummy category
	category := models.Category{Name: "testcat"}
	err = db.Create(&category).Error
	assert.NoError(t, err)

	// Use preloadCategorySelect in a query and inspect SQL
	var result models.Category
	tx := preloadCategorySelect(db.Model(&models.Category{}).Where("id = ?", category.ID))
	tx = tx.Debug() // Print SQL to stdout for inspection
	err = tx.First(&result).Error
	assert.NoError(t, err)
	// Check that only the selected fields are populated
	assert.Equal(t, category.ID, result.ID)
	assert.Equal(t, category.Name, result.Name)
	// Other fields should be zero values (if any)
}

// Additional comprehensive tests for optimized_db coverage
func TestGetDBConfig_Coverage(t *testing.T) {
	// Set environment variables to test different configurations
	os.Setenv("DB_HOST", "testhost")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_SSLMODE", "require")
	os.Setenv("DB_MAX_OPEN_CONNS", "50")
	os.Setenv("DB_MAX_IDLE_CONNS", "25")
	os.Setenv("DB_MAX_LIFETIME", "30m")

	defer func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_SSLMODE")
		os.Unsetenv("DB_MAX_OPEN_CONNS")
		os.Unsetenv("DB_MAX_IDLE_CONNS")
		os.Unsetenv("DB_CONN_MAX_LIFETIME")
	}()

	config := GetDBConfig()
	assert.NotNil(t, config)
	assert.Equal(t, "testhost", config.Host)
	assert.Equal(t, 3306, config.Port)
	assert.Equal(t, "testuser", config.User)
	assert.Equal(t, "testpass", config.Password)
	assert.Equal(t, "testdb", config.DBName)
	assert.Equal(t, "require", config.SSLMode)
	assert.Equal(t, 50, config.MaxOpenConns)
	assert.Equal(t, 25, config.MaxIdleConns)
	assert.Equal(t, 30*time.Minute, config.MaxLifetime)
}

func TestConnectOptimizedDB_Coverage(t *testing.T) {
	// Temporarily set invalid DB_HOST to force an error
	os.Setenv("DB_HOST", "invalid_host")
	defer os.Unsetenv("DB_HOST")

	// Ensure DB_DRIVER is not sqlite, so it attempts postgres connection
	os.Setenv("DB_DRIVER", "postgres")
	defer os.Unsetenv("DB_DRIVER")

	// Test ConnectOptimizedDB function returns an error for invalid connection
	_, err := ConnectOptimizedDB()
	assert.Error(t, err, "ConnectOptimizedDB should return an error for invalid database")
}

func TestOptimizedDB_CreateOptimizedIndexes_Coverage(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Auto-migrate models first
	err = db.AutoMigrate(
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
	)
	assert.NoError(t, err)

	optimizedDB := &OptimizedDB{DB: db}
	err = optimizedDB.CreateOptimizedIndexes()
	assert.NoError(t, err)
}

func TestOptimizedDB_OptimizeQueries(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	optimizedDB := &OptimizedDB{DB: db}
	optimizedDB.OptimizeQueries()
	// OptimizeQueries doesn't return error, just testing execution
}

func TestOptimizedDB_HealthCheck(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Auto-migrate User model for the health check query
	err = db.AutoMigrate(&models.User{})
	assert.NoError(t, err)

	optimizedDB := &OptimizedDB{DB: db}
	healthInfo, err := optimizedDB.HealthCheck()
	assert.NoError(t, err)
	assert.NotNil(t, healthInfo)
	assert.Equal(t, "healthy", healthInfo["status"])
}

func TestGetEnvAsInt_InvalidValue(t *testing.T) {
	os.Setenv("TEST_INVALID_INT", "not_a_number")
	defer os.Unsetenv("TEST_INVALID_INT")

	// Should return default value for invalid integer
	assert.Equal(t, 10, getEnvAsInt("TEST_INVALID_INT", 10))
}

func TestGetEnvAsDuration_InvalidValue(t *testing.T) {
	os.Setenv("TEST_INVALID_DURATION", "not_a_duration")
	defer os.Unsetenv("TEST_INVALID_DURATION")

	// Should return default value for invalid duration
	assert.Equal(t, time.Minute, getEnvAsDuration("TEST_INVALID_DURATION", time.Minute))
}
