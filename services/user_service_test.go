package services

import (
	"testing"

	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
)

func TestUserService_CreateUser(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	userService := NewUserService()

	// Test creating a user
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Phone:    "+1555000001",
		Password: "hashedpassword",
		Role:     "customer",
	}

	err := userService.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if user.Username != "testuser" {
		t.Errorf("Expected username %s, got %s", "testuser", user.Username)
	}
	if user.Email != "test@example.com" {
		t.Errorf("Expected email %s, got %s", "test@example.com", user.Email)
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	userService := NewUserService()

	// Create test user
	user := models.User{
		Username: "testuser2",
		Email:    "test2@example.com",
		Phone:    "+1555000002",
		Password: "hashedpassword",
		Role:     "customer",
	}
	testDB.Create(&user)

	// Test getting user by ID
	retrievedUser, err := userService.GetUserByID(user.ID)
	if err != nil {
		t.Fatalf("Failed to get user by ID: %v", err)
	}

	if retrievedUser.ID != user.ID {
		t.Errorf("Expected ID %d, got %d", user.ID, retrievedUser.ID)
	}
	if retrievedUser.Username != user.Username {
		t.Errorf("Expected username %s, got %s", user.Username, retrievedUser.Username)
	}
}

func TestUserService_GetUserByID_NotFound(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	userService := NewUserService()

	// Test getting non-existent user
	_, err := userService.GetUserByID(999)
	if err == nil {
		t.Fatalf("Expected error for non-existent user")
	}
}

func TestUserService_GetUserByUsername_NotFound(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	userService := NewUserService()

	// Test getting non-existent user by username
	_, err := userService.GetUserByUsername("nonexistent")
	if err == nil {
		t.Fatalf("Expected error for non-existent username")
	}
}

func TestUserService_GetUserByEmail(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	userService := NewUserService()

	// Create test user
	user := models.User{
		Username: "testuser3",
		Email:    "test3@example.com",
		Phone:    "+1555000003",
		Password: "hashedpassword",
		Role:     "customer",
	}
	testDB.Create(&user)

	// Test getting user by email
	retrievedUser, err := userService.GetUserByEmail(user.Email)
	if err != nil {
		t.Fatalf("Failed to get user by email: %v", err)
	}

	if retrievedUser.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, retrievedUser.Email)
	}
	if retrievedUser.Username != user.Username {
		t.Errorf("Expected username %s, got %s", user.Username, retrievedUser.Username)
	}
}

func TestUserService_GetUserByEmail_NotFound(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	userService := NewUserService()

	// Test getting non-existent user by email
	_, err := userService.GetUserByEmail("nonexistent@example.com")
	if err == nil {
		t.Fatalf("Expected error for non-existent email")
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	userService := NewUserService()

	// Create test user
	user := &models.User{
		Username: "testuser4",
		Phone:    "+1555000004",
		Email:    "test4@example.com",
		Password: "hashedpassword",
		Role:     "customer",
	}
	testDB.Create(user)

	// Update user
	user.Username = "updateduser"
	user.Email = "updated@example.com"

	err := userService.UpdateUser(user)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	if user.Username != "updateduser" {
		t.Errorf("Expected username 'updateduser', got %s", user.Username)
	}
	if user.Email != "updated@example.com" {
		t.Errorf("Expected email 'updated@example.com', got %s", user.Email)
	}
}

func TestUserService_UpdateUser_NotFound(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	userService := NewUserService()

	// Try to update non-existent user
	user := &models.User{
		Username: "nonexistent",
		Email:    "nonexistent@example.com",
	}
	user.ID = 999

	err := userService.UpdateUser(user)
	if err == nil {
		t.Fatalf("Expected error for updating non-existent user")
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	userService := NewUserService()

	// Create test user
	user := models.User{
		Username: "testuser5",
		Email:    "test5@example.com",
		Phone:    "+1555000005",
		Password: "hashedpassword",
		Role:     "customer",
	}
	testDB.Create(&user)

	// Delete user
	err := userService.DeleteUser(user.ID)
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	// Verify user is deleted
	_, err = userService.GetUserByID(user.ID)
	if err == nil {
		t.Error("Expected error when getting deleted user, but got none")
	}
}

func TestUserService_DeleteUser_NotFound(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	userService := NewUserService()

	// Try to delete non-existent user
	err := userService.DeleteUser(999)
	if err == nil {
		t.Fatalf("Expected error for deleting non-existent user")
	}
}

func TestUserService_AuthenticateUser(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	userService := NewUserService()

	// Create test user with known password
	user := &models.User{
		Username: "testuser6",
		Email:    "test6@example.com",
		Phone:    "+1555000006",
		Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // "password" hashed
		Role:     "customer",
	}
	testDB.Create(user)

	// Test authentication
	authenticatedUser, err := userService.AuthenticateUser("testuser6", "password")
	if err != nil {
		t.Fatalf("Expected authentication to succeed: %v", err)
	}
	if authenticatedUser == nil {
		t.Error("Expected authenticated user, got nil")
	}

	// Test wrong password
	_, err = userService.AuthenticateUser("testuser6", "wrongpassword")
	if err == nil {
		t.Error("Expected authentication to fail with wrong password")
	}
}

// Benchmark tests for performance
func BenchmarkUserService_CreateUser(b *testing.B) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(b)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	userService := NewUserService()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := &models.User{
			Username: "benchuser",
			Email:    "bench@example.com",
			Password: "hashedpassword",
			Role:     "customer",
		}
		userService.CreateUser(user)
		// Clean up for next iteration
		testDB.Where("username = ?", "benchuser").Delete(&models.User{})
	}
}

func BenchmarkUserService_GetUserByID(b *testing.B) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(b)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	userService := NewUserService()

	// Create test user
	user := models.User{
		Username: "benchuser2",
		Email:    "bench2@example.com",
		Password: "hashedpassword",
		Role:     "customer",
	}
	testDB.Create(&user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userService.GetUserByID(user.ID)
	}
}

func BenchmarkUserService_GetUserByEmail(b *testing.B) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(b)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	userService := NewUserService()

	// Create test user
	user := models.User{
		Username: "benchuser3",
		Email:    "bench3@example.com",
		Password: "hashedpassword",
		Role:     "customer",
	}
	testDB.Create(&user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userService.GetUserByEmail(user.Email)
	}
}
