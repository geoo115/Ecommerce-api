package db

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/geoo115/Ecommerce/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupTestDB initializes a new in-memory SQLite database for testing.
func SetupTestDB(tb testing.TB) *gorm.DB {
	tb.Helper()
	// Create a completely isolated in-memory SQLite DSN for each test
	dsn := fmt.Sprintf("file:test_%d_%d?mode=memory&cache=private", time.Now().UnixNano(), rand.Intn(10000))
	testDB, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		tb.Fatalf("failed to open test db: %v", err)
	}

	// Automigrate all models
	err = testDB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.Inventory{},
		&models.Cart{},
		&models.Order{},
		&models.OrderItem{},
		&models.Address{},
		&models.Review{},
		&models.Wishlist{},
		&models.Payment{},
	)
	if err != nil {
		tb.Fatalf("auto migrate failed: %v", err)
	}

	return testDB
}
