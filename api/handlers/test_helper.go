package handlers

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// UniqueTestData generates unique test data to avoid constraint violations
func UniqueTestData(prefix string) (username, email, phone string) {
	timestamp := time.Now().UnixNano()
	random := rand.Intn(10000)

	username = fmt.Sprintf("%s_user_%d_%d", prefix, timestamp, random)
	if len(username) > 30 {
		username = username[:30]
	}
	email = fmt.Sprintf("%s_test_%d_%d@example.com", prefix, timestamp, random)
	phone = fmt.Sprintf("+1%d%04d", timestamp%1000000000, random)

	fmt.Printf("Generated Username: %s (length: %d)\n", username, len(username))
	fmt.Printf("Generated Email: %s (length: %d)\n", email, len(email))
	fmt.Printf("Generated Phone: %s (length: %d)\n", phone, len(phone))

	return username, email, phone
}

// CreateTestUser creates a unique test user to avoid constraint violations
func CreateTestUser(tb testing.TB, testDB *gorm.DB, prefix string) *models.User {
	tb.Helper()

	username, email, phone := UniqueTestData(prefix)

	user := &models.User{
		Username: username,
		Email:    email,
		Phone:    phone,
		Password: "hashedpassword123",
	}

	if err := testDB.Create(user).Error; err != nil {
		tb.Fatalf("failed to create test user: %v", err)
	}

	return user
}

// SetupTestDB initializes a new in-memory SQLite database for testing.
func SetupTestDB(tb testing.TB) *gorm.DB {
	tb.Helper()
	// Create a unique in-memory SQLite database for each test
	dsn := fmt.Sprintf("file:test_%d.db?mode=memory&cache=shared", time.Now().UnixNano())
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

	db.DB = testDB

	// Configure underlying sql.DB to be completely isolated
	if sqlDB, err := testDB.DB(); err == nil {
		// Use single connection for complete isolation
		sqlDB.SetMaxOpenConns(1)
		sqlDB.SetMaxIdleConns(1)
		sqlDB.SetConnMaxLifetime(time.Second * 30)
	} else {
		tb.Fatalf("failed to get sql DB: %v", err)
	}

	tb.Cleanup(func() {
		// Close the database connection
		if sqlDB, err := testDB.DB(); err == nil {
			sqlDB.Close()
		}
	})

	return testDB
}

func TestMain(m *testing.M) {
	// TestMain intentionally left empty. Individual tests call SetupTestDB(t)
	// to ensure each test gets an isolated in-memory database.
	m.Run()
}
