package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
	"github.com/geoo115/Ecommerce/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Enhanced user registration tests
func TestSignup_Success(t *testing.T) {
	SetupTestDB(t)

	userData := map[string]interface{}{
		"username": "testuser",
		"password": "TestPass123",
		"email":    "test@example.com",
		"phone":    "+1234567890",
		"role":     "user",
	}

	jsonData, _ := json.Marshal(userData)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Method: "POST",
		Body:   io.NopCloser(bytes.NewReader(jsonData)),
		Header: http.Header{"Content-Type": []string{"application/json"}},
	}

	Signup(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "User created successfully")
}

func TestSignup_InvalidEmail(t *testing.T) {
	SetupTestDB(t)

	userData := map[string]interface{}{
		"username": "testuser2",
		"password": "TestPass123",
		"email":    "invalid-email",
		"phone":    "+1234567891",
	}

	jsonData, _ := json.Marshal(userData)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Method: "POST",
		Body:   io.NopCloser(bytes.NewReader(jsonData)),
		Header: http.Header{"Content-Type": []string{"application/json"}},
	}

	Signup(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid email format")
}

func TestSignup_WeakPassword(t *testing.T) {
	SetupTestDB(t)

	userData := map[string]interface{}{
		"username": "testuser3",
		"password": "123",
		"email":    "test3@example.com",
		"phone":    "+1234567892",
	}

	jsonData, _ := json.Marshal(userData)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Method: "POST",
		Body:   io.NopCloser(bytes.NewReader(jsonData)),
		Header: http.Header{"Content-Type": []string{"application/json"}},
	}

	Signup(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Password must be at least")
}

func TestSignup_DuplicateUser(t *testing.T) {
	SetupTestDB(t)

	// Create first user
	userData := map[string]interface{}{
		"username": "duplicate_user",
		"password": "TestPass123",
		"email":    "duplicate@example.com",
		"phone":    "+1234567893",
	}

	jsonData, _ := json.Marshal(userData)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Method: "POST",
		Body:   io.NopCloser(bytes.NewReader(jsonData)),
		Header: http.Header{"Content-Type": []string{"application/json"}},
	}

	Signup(c)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Try to create duplicate
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request = &http.Request{
		Method: "POST",
		Body:   io.NopCloser(bytes.NewReader(jsonData)),
		Header: http.Header{"Content-Type": []string{"application/json"}},
	}

	Signup(c2)
	assert.Equal(t, http.StatusConflict, w2.Code)
	assert.Contains(t, w2.Body.String(), "User already exists")
}

func TestSignup_Integration(t *testing.T) {
	SetupTestDB(t)

	// Test data
	userData := map[string]interface{}{
		"username": "signup_testuser",
		"password": "TestPass123",
		"email":    "signup_test@example.com",
		"phone":    "+1234567890",
		"role":     "user",
	}

	jsonData, _ := json.Marshal(userData)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Method: "POST",
		Body:   io.NopCloser(bytes.NewReader(jsonData)),
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
	}

	Signup(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "User created successfully")
}

func TestLogin_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := testDB.AutoMigrate(&models.User{}, &models.Address{}, &models.Cart{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}
	db.DB = testDB

	// Create user
	user := models.User{
		Username: "login_testuser",
		Password: "TestPass123",
		Email:    "login_test@example.com",
		Phone:    "+19999999999",
	}
	hashedPass, _ := utils.HashPassword(user.Password)
	user.Password = hashedPass
	if err := db.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	loginData := map[string]interface{}{
		"username": "login_testuser",
		"password": "TestPass123",
	}

	jsonData, _ := json.Marshal(loginData)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Method: "POST",
		Body:   io.NopCloser(bytes.NewReader(jsonData)),
		Header: http.Header{"Content-Type": []string{"application/json"}},
	}

	Login(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Login successful")
	assert.Contains(t, w.Body.String(), "token")
}

func TestLogin_InvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	loginData := map[string]interface{}{
		"username": "nonexistent",
		"password": "wrongpass",
	}

	jsonData, _ := json.Marshal(loginData)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Method: "POST",
		Body:   io.NopCloser(bytes.NewReader(jsonData)),
		Header: http.Header{"Content-Type": []string{"application/json"}},
	}

	Login(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid credentials")
}

func TestLogin_EmptyFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	loginData := map[string]interface{}{
		"username": "",
		"password": "",
	}

	jsonData, _ := json.Marshal(loginData)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Method: "POST",
		Body:   io.NopCloser(bytes.NewReader(jsonData)),
		Header: http.Header{"Content-Type": []string{"application/json"}},
	}

	Login(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Username and password are required")
}

func TestLogin_WrongPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := testDB.AutoMigrate(&models.User{}, &models.Address{}, &models.Cart{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}
	db.DB = testDB

	// Create user with known password
	user := models.User{
		Username: "wrongpassuser",
		Password: "CorrectPass123",
		Email:    "wrongpass@example.com",
		Phone:    "+15555555555",
	}
	hashedPass, _ := utils.HashPassword(user.Password)
	user.Password = hashedPass
	db.DB.Create(&user)

	// Try login with wrong password
	loginData := map[string]interface{}{
		"username": "wrongpassuser",
		"password": "WrongPassword",
	}

	jsonData, _ := json.Marshal(loginData)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Method: "POST",
		Body:   io.NopCloser(bytes.NewReader(jsonData)),
		Header: http.Header{"Content-Type": []string{"application/json"}},
	}

	Login(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid credentials")
}

func TestLogin_DatabaseError(t *testing.T) {
	SetupTestDB(t)

	// Simulate database error by setting db.DB to nil
	originalDB := db.DB
	db.DB = nil
	defer func() { db.DB = originalDB }()

	loginData := map[string]interface{}{
		"username": "testuser",
		"password": "TestPass123",
	}

	jsonData, _ := json.Marshal(loginData)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Method: "POST",
		Body:   io.NopCloser(bytes.NewReader(jsonData)),
		Header: http.Header{"Content-Type": []string{"application/json"}},
	}

	Login(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Database error")
}

func TestLogin_Integration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// ensure test DB is isolated for this test (avoid other tests overwriting db.DB)
	testDB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := testDB.AutoMigrate(&models.User{}, &models.Address{}, &models.Cart{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}
	db.DB = testDB

	// First, create a user
	user := models.User{
		Username: "login_testuser",
		Password: "TestPass123",
		Email:    "login_test@example.com",
		Phone:    "+19999999999",
	}
	hashedPass, _ := utils.HashPassword(user.Password)
	user.Password = hashedPass
	if err := db.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Test login
	loginData := map[string]interface{}{
		"username": "login_testuser",
		"password": "TestPass123",
	}

	jsonData, _ := json.Marshal(loginData)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Method: "POST",
		Body:   io.NopCloser(bytes.NewReader(jsonData)),
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
	}

	Login(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Login successful")
	assert.Contains(t, w.Body.String(), "token")
}

func TestSignup_InvalidUsername(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	userData := map[string]interface{}{
		"username": "u", // Too short
		"password": "TestPass123",
		"email":    "test@example.com",
	}

	jsonData, _ := json.Marshal(userData)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Method: "POST",
		Body:   io.NopCloser(bytes.NewReader(jsonData)),
		Header: http.Header{"Content-Type": []string{"application/json"}},
	}

	Signup(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSignup_InvalidEmailFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	userData := map[string]interface{}{
		"username": "validuser",
		"password": "TestPass123",
		"email":    "invalid-email",
	}

	jsonData, _ := json.Marshal(userData)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Method: "POST",
		Body:   io.NopCloser(bytes.NewReader(jsonData)),
		Header: http.Header{"Content-Type": []string{"application/json"}},
	}

	Signup(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSignup_InvalidPhone(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	userData := map[string]interface{}{
		"username": "validuser",
		"password": "TestPass123",
		"email":    "test@example.com",
		"phone":    "invalid-phone",
	}

	jsonData, _ := json.Marshal(userData)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Method: "POST",
		Body:   io.NopCloser(bytes.NewReader(jsonData)),
		Header: http.Header{"Content-Type": []string{"application/json"}},
	}

	Signup(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSignup_InvalidRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	userData := map[string]interface{}{
		"username": "validuser",
		"password": "TestPass123",
		"email":    "test@example.com",
		"role":     "invalid_role",
	}

	jsonData, _ := json.Marshal(userData)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Method: "POST",
		Body:   io.NopCloser(bytes.NewReader(jsonData)),
		Header: http.Header{"Content-Type": []string{"application/json"}},
	}

	Signup(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_InvalidUsernameFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	loginData := map[string]interface{}{
		"username": "invalid@username",
		"password": "TestPass123",
	}

	jsonData, _ := json.Marshal(loginData)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Method: "POST",
		Body:   io.NopCloser(bytes.NewReader(jsonData)),
		Header: http.Header{"Content-Type": []string{"application/json"}},
	}

	Login(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_UserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	loginData := map[string]interface{}{
		"username": "nonexistent_user",
		"password": "TestPass123",
	}

	jsonData, _ := json.Marshal(loginData)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Method: "POST",
		Body:   io.NopCloser(bytes.NewReader(jsonData)),
		Header: http.Header{"Content-Type": []string{"application/json"}},
	}

	Login(c)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid credentials")
}

func TestLogout_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set user in context (simulating authenticated user)
	user := models.User{Username: "testuser", Email: "test@example.com"}
	c.Set("user", user)

	Logout(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "User logged out successfully")
}

func TestLogout_Unauthenticated(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Don't set user in context (simulating unauthenticated request)
	Logout(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "User not authenticated")
}
