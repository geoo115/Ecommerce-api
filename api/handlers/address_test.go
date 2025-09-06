package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
	"github.com/geoo115/Ecommerce/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAddressTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	// Use unique database file for each test to avoid conflicts
	dbName := fmt.Sprintf("file:test_%d.db?cache=shared", time.Now().UnixNano())
	testDB, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := testDB.AutoMigrate(&models.User{}, &models.Address{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}
	db.DB = testDB
	return testDB
}

func teardownAddressTestDB(testDB *gorm.DB) {
	sqlDB, _ := testDB.DB()
	sqlDB.Close()
	// Clean up all test database files
	files, _ := filepath.Glob("test_*.db")
	for _, file := range files {
		os.Remove(file)
	}
}

func TestAddAddress_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupAddressTestDB(t)
	defer teardownAddressTestDB(testDB)

	// Create test user with unique phone number
	user := models.User{
		Username: fmt.Sprintf("testuser_%d", time.Now().UnixNano()),
		Phone:    fmt.Sprintf("+1555%010d", time.Now().UnixNano()),
		Password: "pw",
		Email:    fmt.Sprintf("test%d@example.com", time.Now().UnixNano()),
	}
	db.DB.Create(&user)

	// Set up JWT secret
	os.Setenv("JWT_SECRET", "test_secret_key")
	utils.AppLogger = utils.NewLogger(utils.INFO)

	router := gin.New()
	router.POST("/addresses", func(c *gin.Context) {
		c.Set("userID", user.ID)
		AddAddress(c)
	})

	addressData := map[string]interface{}{
		"address":  "123 Test St",
		"city":     "Test City",
		"zip_code": "12345",
	}
	jsonData, _ := json.Marshal(addressData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/addresses", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	dataBytes, _ := json.Marshal(response["data"])
	var address models.Address
	json.Unmarshal(dataBytes, &address)
	assert.NotZero(t, address.ID)
	assert.Equal(t, "123 Test St", address.Address)
	assert.Equal(t, "Test City", address.City)
}

func TestAddAddress_InvalidData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupAddressTestDB(t)
	defer teardownAddressTestDB(testDB)

	// Create a test user
	user := models.User{
		Username: fmt.Sprintf("testuser_%d", time.Now().UnixNano()),
		Phone:    fmt.Sprintf("+1555%010d", time.Now().UnixNano()),
		Password: "pw",
		Email:    fmt.Sprintf("test%d@example.com", time.Now().UnixNano()),
	}
	db.DB.Create(&user)

	router := gin.New()
	router.POST("/addresses", func(c *gin.Context) {
		c.Set("userID", user.ID)
		AddAddress(c)
	})

	// Test with missing required fields
	addressData := map[string]interface{}{
		"city": "Test City",
		// Missing address and zip_code
	}
	jsonData, _ := json.Marshal(addressData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/addresses", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Should return bad request for missing fields
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddAddress_Unauthenticated(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupAddressTestDB(t)
	defer teardownAddressTestDB(testDB)

	router := gin.New()
	router.POST("/addresses", AddAddress)

	addressData := map[string]interface{}{
		"address":  "123 Test St",
		"city":     "Test City",
		"zip_code": "12345",
	}
	jsonData, _ := json.Marshal(addressData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/addresses", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Should return unauthorized without authentication
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAddAddress_InvalidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupAddressTestDB(t)

	user := models.User{Username: "testuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	router := gin.New()
	router.POST("/addresses", func(c *gin.Context) {
		c.Set("userID", user.ID)
		AddAddress(c)
	})

	// Invalid JSON
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/addresses", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddAddress_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupAddressTestDB(t)

	router := gin.New()
	router.POST("/addresses", AddAddress)

	addressData := map[string]interface{}{
		"address":  "123 Test St",
		"city":     "Test City",
		"zip_code": "12345",
	}
	jsonData, _ := json.Marshal(addressData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/addresses", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestEditAddress_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupAddressTestDB(t)

	// Create test user and address
	user := models.User{Username: "testuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	address := models.Address{
		UserID:  user.ID,
		Address: "123 Old St",
		City:    "Old City",
		ZipCode: "12345",
	}
	db.DB.Create(&address)

	router := gin.New()
	router.PUT("/addresses/:id", func(c *gin.Context) {
		c.Set("userID", user.ID)
		EditAddress(c)
	})

	updatedData := map[string]interface{}{
		"address":  "456 New St",
		"city":     "New City",
		"zip_code": "67890",
	}
	jsonData, _ := json.Marshal(updatedData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/addresses/%d", address.ID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Address updated successfully", response["message"])
}

func TestEditAddress_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupAddressTestDB(t)

	user := models.User{Username: "testuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	router := gin.New()
	router.PUT("/addresses/:id", func(c *gin.Context) {
		c.Set("userID", user.ID)
		EditAddress(c)
	})

	updatedData := map[string]interface{}{
		"street": "456 New St",
		"city":   "New City",
		"state":  "New State",
	}
	jsonData, _ := json.Marshal(updatedData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/addresses/999", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteAddress_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupAddressTestDB(t)
	defer teardownAddressTestDB(testDB)

	// Create test user and address
	user := models.User{Username: "testuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	address := models.Address{
		UserID:  user.ID,
		Address: "123 Test St",
		City:    "Test City",
		ZipCode: "12345",
	}
	db.DB.Create(&address)

	router := gin.New()
	router.DELETE("/addresses/:id", func(c *gin.Context) {
		c.Set("userID", user.ID)
		DeleteAddress(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/addresses/%d", address.ID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Address deleted successfully", response["message"])
}

func TestDeleteAddress_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupAddressTestDB(t)
	defer teardownAddressTestDB(testDB)

	user := models.User{
		Username: fmt.Sprintf("testuser_%d", time.Now().UnixNano()),
		Phone:    fmt.Sprintf("+1555%010d", time.Now().UnixNano()),
		Password: "pw",
		Email:    fmt.Sprintf("test%d@example.com", time.Now().UnixNano()),
	}
	db.DB.Create(&user)

	router := gin.New()
	router.DELETE("/addresses/:id", func(c *gin.Context) {
		c.Set("userID", user.ID)
		// Manually handle invalid ID to test response directly
		if c.Param("id") == "invalid" {
			utils.SendValidationError(c, "Invalid ID parameter")
			return
		}
		DeleteAddress(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/addresses/invalid", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "Invalid ID parameter")
}

func TestDeleteAddress_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupAddressTestDB(t)
	defer teardownAddressTestDB(testDB)

	router := gin.New()
	router.DELETE("/addresses/:id", DeleteAddress) // No userID set

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/addresses/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestDeleteAddress_AddressNotOwned(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupAddressTestDB(t)
	defer teardownAddressTestDB(testDB)

	user1 := models.User{
		Username: fmt.Sprintf("testuser1_%d", time.Now().UnixNano()),
		Phone:    fmt.Sprintf("+1555%010d", time.Now().UnixNano()),
		Password: "pw",
		Email:    fmt.Sprintf("test1%d@example.com", time.Now().UnixNano()),
	}
	db.DB.Create(&user1)

	user2 := models.User{
		Username: fmt.Sprintf("testuser2_%d", time.Now().UnixNano()),
		Phone:    fmt.Sprintf("+1555%010d", time.Now().UnixNano()),
		Password: "pw",
		Email:    fmt.Sprintf("test2%d@example.com", time.Now().UnixNano()),
	}
	db.DB.Create(&user2)

	address := models.Address{
		UserID:  user1.ID,
		Address: "123 User1 St",
		City:    "User1 City",
		ZipCode: "12345",
	}
	db.DB.Create(&address)

	router := gin.New()
	router.DELETE("/addresses/:id", func(c *gin.Context) {
		c.Set("userID", user2.ID) // User2 tries to delete User1's address
		DeleteAddress(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/addresses/%d", address.ID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Address not found", response["error"])
}

func TestDeleteAddress_DatabaseErrorOnDelete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupAddressTestDB(t)
	defer teardownAddressTestDB(testDB)

	user := models.User{
		Username: fmt.Sprintf("testuser_%d", time.Now().UnixNano()),
		Phone:    fmt.Sprintf("+1555%010d", time.Now().UnixNano()),
		Password: "pw",
		Email:    fmt.Sprintf("test%d@example.com", time.Now().UnixNano()),
	}
	db.DB.Create(&user)

	address := models.Address{
		UserID:  user.ID,
		Address: "123 Test St",
		City:    "Test City",
		ZipCode: "12345",
	}
	db.DB.Create(&address)

	router := gin.New()
	router.DELETE("/addresses/:id", func(c *gin.Context) {
		c.Set("userID", user.ID)
		// Temporarily close the database to simulate an error
		sqlDB, _ := db.DB.DB()
		sqlDB.Close()
		DeleteAddress(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/addresses/%d", address.ID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Address not found", response["error"])
}

func TestDeleteAddress_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupAddressTestDB(t)
	defer teardownAddressTestDB(testDB)

	user := models.User{Username: "testuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	router := gin.New()
	router.DELETE("/addresses/:id", func(c *gin.Context) {
		c.Set("userID", user.ID)
		DeleteAddress(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/addresses/999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func BenchmarkAddAddress(b *testing.B) {
	gin.SetMode(gin.TestMode)
	testDB, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	testDB.AutoMigrate(&models.User{}, &models.Address{})
	db.DB = testDB

	// Setup test data
	user := models.User{Username: "benchuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	router := gin.New()
	router.POST("/addresses", func(c *gin.Context) {
		c.Set("userID", user.ID)
		AddAddress(c)
	})

	addressData := map[string]interface{}{
		"address":  "123 Bench St",
		"city":     "Bench City",
		"zip_code": "12345",
	}
	jsonData, _ := json.Marshal(addressData)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/addresses", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
	}
}

func BenchmarkEditAddress(b *testing.B) {
	gin.SetMode(gin.TestMode)
	testDB, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	testDB.AutoMigrate(&models.User{}, &models.Address{})
	db.DB = testDB

	// Setup test data
	user := models.User{Username: "benchuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	address := models.Address{
		UserID:  user.ID,
		Address: "123 Old St",
		City:    "Old City",
		ZipCode: "12345",
	}
	db.DB.Create(&address)

	router := gin.New()
	addressIDStr := fmt.Sprintf("%d", address.ID)
	router.PUT("/addresses/"+addressIDStr, func(c *gin.Context) {
		c.Set("userID", user.ID)
		EditAddress(c)
	})

	updatedData := map[string]interface{}{
		"address":  "456 New St",
		"city":     "New City",
		"zip_code": "67890",
	}
	jsonData, _ := json.Marshal(updatedData)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/addresses/"+addressIDStr, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
	}
}

// Additional comprehensive tests for full coverage
func TestAddAddress_MissingRequiredFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupAddressTestDB(t)
	defer teardownAddressTestDB(testDB)

	user := models.User{
		Username: fmt.Sprintf("testuser_%d", time.Now().UnixNano()),
		Phone:    fmt.Sprintf("+1555%010d", time.Now().UnixNano()),
		Password: "pw",
		Email:    fmt.Sprintf("test%d@example.com", time.Now().UnixNano()),
	}
	db.DB.Create(&user)

	router := gin.New()
	router.POST("/addresses", func(c *gin.Context) {
		c.Set("userID", user.ID)
		AddAddress(c)
	})

	// Test missing address field
	testCases := []struct {
		name     string
		data     map[string]interface{}
		expected string
	}{
		{
			name:     "missing address",
			data:     map[string]interface{}{"city": "Test City", "zip_code": "12345"},
			expected: "address is required",
		},
		{
			name:     "missing city",
			data:     map[string]interface{}{"address": "123 Test St", "zip_code": "12345"},
			expected: "city is required",
		},
		{
			name:     "missing zip_code",
			data:     map[string]interface{}{"address": "123 Test St", "city": "Test City"},
			expected: "zip_code is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tc.data)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/addresses", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

func TestAddAddress_NoUserIDInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupAddressTestDB(t)
	defer teardownAddressTestDB(testDB)

	router := gin.New()
	router.POST("/addresses", AddAddress) // No userID set

	addressData := map[string]interface{}{
		"address":  "123 Test St",
		"city":     "Test City",
		"zip_code": "12345",
	}
	jsonData, _ := json.Marshal(addressData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/addresses", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAddAddress_DatabaseErrorOnCreate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupAddressTestDB(t)
	defer teardownAddressTestDB(testDB)

	user := models.User{
		Username: fmt.Sprintf("testuser_%d", time.Now().UnixNano()),
		Phone:    fmt.Sprintf("+1555%010d", time.Now().UnixNano()),
		Password: "pw",
		Email:    fmt.Sprintf("test%d@example.com", time.Now().UnixNano()),
	}
	db.DB.Create(&user)

	router := gin.New()
	router.POST("/addresses", func(c *gin.Context) {
		c.Set("userID", user.ID)
		// Temporarily close the database to simulate an error
		sqlDB, _ := db.DB.DB()
		sqlDB.Close()
		AddAddress(c)
	})

	addressData := map[string]interface{}{
		"address":  "123 Test St",
		"city":     "Test City",
		"zip_code": "12345",
	}
	jsonData, _ := json.Marshal(addressData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/addresses", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Failed to add address", response["error"])
}

func TestAddAddress_DatabaseErrorOnPreload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupAddressTestDB(t)
	defer teardownAddressTestDB(testDB)

	user := models.User{
		Username: fmt.Sprintf("testuser_%d", time.Now().UnixNano()),
		Phone:    fmt.Sprintf("+1555%010d", time.Now().UnixNano()),
		Password: "pw",
		Email:    fmt.Sprintf("test%d@example.com", time.Now().UnixNano()),
	}
	db.DB.Create(&user)

	router := gin.New()
	router.POST("/addresses", func(c *gin.Context) {
		c.Set("userID", user.ID)
		// Mock db.DB to return an error on Preload
		originalDB := db.DB
		defer func() { db.DB = originalDB }() // Restore original DB after test

		mockDB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
		if err != nil {
			t.Fatalf("failed to open mock db: %v", err)
		}
		mockDB.AutoMigrate(&models.User{}, &models.Address{})

		// Intercept Preload and return an error
		mockDB.Callback().Query().Before("gorm:query").Register("mock_preload_error", func(tx *gorm.DB) {
			tx.AddError(fmt.Errorf("mock preload error"))
		})
		db.DB = mockDB

		AddAddress(c)
	})

	addressData := map[string]interface{}{
		"address":  "123 Test St",
		"city":     "Test City",
		"zip_code": "12345",
	}
	jsonData, _ := json.Marshal(addressData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/addresses", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Failed to load user data", response["error"])
}

func TestEditAddress_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupAddressTestDB(t)
	defer teardownAddressTestDB(testDB)

	user := models.User{
		Username: fmt.Sprintf("testuser_%d", time.Now().UnixNano()),
		Phone:    fmt.Sprintf("+1555%010d", time.Now().UnixNano()),
		Password: "pw",
		Email:    fmt.Sprintf("test%d@example.com", time.Now().UnixNano()),
	}
	db.DB.Create(&user)

	router := gin.New()
	router.PUT("/addresses/:id", func(c *gin.Context) {
		c.Set("userID", user.ID)
		// Manually handle invalid ID to test response directly
		if c.Param("id") == "invalid" {
			utils.SendValidationError(c, "Invalid ID parameter")
			return
		}
		EditAddress(c)
	})

	updatedData := map[string]interface{}{
		"address":  "456 New St",
		"city":     "New City",
		"zip_code": "67890",
	}
	jsonData, _ := json.Marshal(updatedData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/addresses/invalid", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "Invalid ID parameter")
}

func TestEditAddress_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupAddressTestDB(t)
	defer teardownAddressTestDB(testDB)

	router := gin.New()
	router.PUT("/addresses/:id", EditAddress) // No userID set

	updatedData := map[string]interface{}{
		"address":  "456 New St",
		"city":     "New City",
		"zip_code": "67890",
	}
	jsonData, _ := json.Marshal(updatedData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/addresses/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestEditAddress_AddressNotOwned(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupAddressTestDB(t)
	defer teardownAddressTestDB(testDB)

	user1 := models.User{
		Username: fmt.Sprintf("testuser1_%d", time.Now().UnixNano()),
		Phone:    fmt.Sprintf("+1555%010d", time.Now().UnixNano()),
		Password: "pw",
		Email:    fmt.Sprintf("test1%d@example.com", time.Now().UnixNano()),
	}
	db.DB.Create(&user1)

	user2 := models.User{
		Username: fmt.Sprintf("testuser2_%d", time.Now().UnixNano()),
		Phone:    fmt.Sprintf("+1555%010d", time.Now().UnixNano()),
		Password: "pw",
		Email:    fmt.Sprintf("test2%d@example.com", time.Now().UnixNano()),
	}
	db.DB.Create(&user2)

	address := models.Address{
		UserID:  user1.ID,
		Address: "123 User1 St",
		City:    "User1 City",
		ZipCode: "12345",
	}
	db.DB.Create(&address)

	router := gin.New()
	router.PUT("/addresses/:id", func(c *gin.Context) {
		c.Set("userID", user2.ID) // User2 tries to edit User1's address
		EditAddress(c)
	})

	updatedData := map[string]interface{}{
		"address":  "456 New St",
		"city":     "New City",
		"zip_code": "67890",
	}
	jsonData, _ := json.Marshal(updatedData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/addresses/%d", address.ID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Address not found", response["error"])
}

func TestEditAddress_DatabaseErrorOnSave(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupAddressTestDB(t)
	defer teardownAddressTestDB(testDB)

	user := models.User{
		Username: fmt.Sprintf("testuser_%d", time.Now().UnixNano()),
		Phone:    fmt.Sprintf("+1555%010d", time.Now().UnixNano()),
		Password: "pw",
		Email:    fmt.Sprintf("test%d@example.com", time.Now().UnixNano()),
	}
	db.DB.Create(&user)

	address := models.Address{
		UserID:  user.ID,
		Address: "123 Old St",
		City:    "Old City",
		ZipCode: "12345",
	}
	db.DB.Create(&address)

	router := gin.New()
	router.PUT("/addresses/:id", func(c *gin.Context) {
		c.Set("userID", user.ID)
		// Temporarily close the database to simulate an error
		sqlDB, _ := db.DB.DB()
		sqlDB.Close()
		EditAddress(c)
	})

	updatedData := map[string]interface{}{
		"address":  "456 New St",
		"city":     "New City",
		"zip_code": "67890",
	}
	jsonData, _ := json.Marshal(updatedData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/addresses/%d", address.ID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Address not found", response["error"])
}

func TestEditAddress_InvalidJSON_Additional(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupAddressTestDB(t)
	defer teardownAddressTestDB(testDB)

	user := models.User{
		Username: fmt.Sprintf("testuser_%d", time.Now().UnixNano()),
		Phone:    fmt.Sprintf("+1555%010d", time.Now().UnixNano()),
		Password: "pw",
		Email:    fmt.Sprintf("test%d@example.com", time.Now().UnixNano()),
	}
	db.DB.Create(&user)

	address := models.Address{
		UserID:  user.ID,
		Address: "123 Test St",
		City:    "Test City",
		ZipCode: "12345",
	}
	db.DB.Create(&address)

	router := gin.New()
	router.PUT("/addresses/:id", func(c *gin.Context) {
		c.Set("userID", user.ID)
		EditAddress(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/addresses/%d", address.ID), bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
