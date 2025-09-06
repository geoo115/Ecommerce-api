package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/geoo115/Ecommerce/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGetPaginationParams(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name     string
		query    string
		expPage  int
		expLimit int
	}{
		{
			name:     "default values",
			query:    "",
			expPage:  1,
			expLimit: 10,
		},
		{
			name:     "valid page and limit",
			query:    "page=2&limit=20",
			expPage:  2,
			expLimit: 20,
		},
		{
			name:     "invalid page - non-numeric",
			query:    "page=abc&limit=20",
			expPage:  1, // Should default to 1
			expLimit: 20,
		},
		{
			name:     "invalid limit - non-numeric",
			query:    "page=2&limit=xyz",
			expPage:  2,
			expLimit: 10, // Should default to 10
		},
		{
			name:     "invalid page - zero",
			query:    "page=0&limit=20",
			expPage:  1, // Should default to 1
			expLimit: 20,
		},
		{
			name:     "invalid limit - zero",
			query:    "page=2&limit=0",
			expPage:  2,
			expLimit: 10, // Should default to 10
		},
		{
			name:     "invalid limit - too high",
			query:    "page=2&limit=200",
			expPage:  2,
			expLimit: 10, // Should default to 10
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodGet, "/?"+tt.query, nil)

			handler := &HandlerBase{}
			params := handler.GetPaginationParams(c)

			assert.Equal(t, tt.expPage, params.Page)
			assert.Equal(t, tt.expLimit, params.Limit)
		})
	}
}

func TestApplyPagination(t *testing.T) {
	// Setup a test database with DryRun to inspect the statement without hitting a real DB
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{DryRun: true})
	assert.NoError(t, err)

	handler := &HandlerBase{}

	tests := []struct {
		name           string
		page           int
		limit          int
		expectedOffset int
		expectedLimit  int
	}{
		{
			name:           "page 1, limit 10",
			page:           1,
			limit:          10,
			expectedOffset: 0,
			expectedLimit:  10,
		},
		{
			name:           "page 2, limit 20",
			page:           2,
			limit:          20,
			expectedOffset: 20,
			expectedLimit:  20,
		},
		{
			name:           "page 5, limit 5",
			page:           5,
			limit:          5,
			expectedOffset: 20,
			expectedLimit:  5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := PaginationParams{Page: tt.page, Limit: tt.limit}

			// Apply pagination to a dummy query using the User model
			var users []models.User
			result := handler.ApplyPagination(db.Model(&models.User{}), params)

			// Convert to SQL to inspect the query
			sql := result.ToSQL(func(tx *gorm.DB) *gorm.DB {
				return tx.Find(&users)
			})

			// Check that LIMIT and OFFSET are in the SQL
			assert.Contains(t, sql, "LIMIT")
			if tt.expectedOffset > 0 {
				assert.Contains(t, sql, "OFFSET")
			}
		})
	}
}

func TestSendCreatedResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	response := &ResponseHelper{}
	response.SendCreatedResponse(c, "Resource created successfully", gin.H{"id": 1})
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Resource created successfully")
}

func TestSendUpdatedResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	response := &ResponseHelper{}
	response.SendUpdatedResponse(c, "Resource updated successfully", gin.H{"id": 1})
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Resource updated successfully")
}

func TestSendDeletedResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	response := &ResponseHelper{}
	response.SendDeletedResponse(c, "Resource deleted successfully")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Resource deleted successfully")
}

func TestSendListResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	response := &ResponseHelper{}
	response.SendListResponse(c, "Resources retrieved successfully", []gin.H{{"id": 1}})
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Resources retrieved successfully")
}

func TestHandleDBError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Test case 1: Record not found error
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	handler := &HandlerBase{}
	handler.HandleDBError(c, gorm.ErrRecordNotFound, "resource", "finding")
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Test case 2: Other database error
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	handler.HandleDBError(c, errors.New("some db error"), "resource", "deleting")
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Test case 3: No error - should not send any response
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	handler.HandleDBError(c, nil, "resource", "creating")
	// When no error, no response is sent, so status should remain 200 (default)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCheckOwnership(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup test database
	testDB := SetupTestDB(t)

	// Create test user and address
	user := CreateTestUser(t, testDB, "ownership")

	address := models.Address{UserID: user.ID, Address: "123 Test St", City: "Test City", ZipCode: "12345"}
	err := testDB.Create(&address).Error
	assert.NoError(t, err)

	// Test case 1: Ownership matches
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", user.ID)
	handler := &HandlerBase{}
	err = handler.CheckOwnership(c, &models.Address{}, address.ID)
	assert.NoError(t, err)

	// Test case 2: Ownership does not match
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Set("userID", uint(999)) // Different user
	err = handler.CheckOwnership(c, &models.Address{}, address.ID)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, w.Code) // Should be 404 for not found/access denied

	// Test case 3: UserID not in context
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	err = handler.CheckOwnership(c, &models.Address{}, address.ID)
	assert.Error(t, err)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestTransactionWrapper_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup test database
	testDB := SetupTestDB(t)

	handler := &HandlerBase{}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Generate unique test data
	username, email, phone := UniqueTestData("transaction_success")

	// Define a function to be wrapped that simulates a successful operation
	err := handler.TransactionWrapper(c, func(tx *gorm.DB) error {
		return tx.Create(&models.User{Username: username, Email: email, Phone: phone}).Error
	})

	assert.NoError(t, err)

	// Verify the user was created and committed
	var user models.User
	result := testDB.First(&user, "username = ?", username)
	assert.NoError(t, result.Error)
	assert.Equal(t, username, user.Username)
}

func TestTransactionWrapper_FnReturnsError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup test database
	testDB := SetupTestDB(t)

	handler := &HandlerBase{}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Define a function to be wrapped that simulates an error
	expectedErr := errors.New("simulated error from fn")
	err := handler.TransactionWrapper(c, func(tx *gorm.DB) error {
		tx.Create(&models.User{Username: "should-be-rolled-back", Email: "rollback@test.com", Phone: "456"})
		return expectedErr
	})

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

	// Verify that the transaction was rolled back
	var user models.User
	result := testDB.First(&user, "username = ?", "should-be-rolled-back")
	assert.Error(t, result.Error)
	assert.Equal(t, gorm.ErrRecordNotFound, result.Error)
}

func TestTransactionWrapper_FnPanics(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup test database
	testDB := SetupTestDB(t)

	handler := &HandlerBase{}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Define a function to be wrapped that panics
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code did not recover from panic: %v", r)
		}
	}()

	_ = handler.TransactionWrapper(c, func(tx *gorm.DB) error {
		tx.Create(&models.User{Username: "should-be-rolled-back-on-panic", Email: "panic@test.com", Phone: "789"})
		panic("simulated panic from fn")
	})

	// Verify that an internal server error was sent (due to recover() and SendInternalError)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Verify that the transaction was rolled back
	var user models.User
	result := testDB.First(&user, "username = ?", "should-be-rolled-back-on-panic")
	assert.Error(t, result.Error)
	assert.Equal(t, gorm.ErrRecordNotFound, result.Error)
}

func TestHandlerBase_SendCreatedResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	handler := &HandlerBase{}
	handler.SendCreatedResponse(c, "Resource created successfully", gin.H{"id": 1})
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Resource created successfully")
}

func TestHandlerBase_SendUpdatedResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	handler := &HandlerBase{}
	handler.SendUpdatedResponse(c, "Resource updated successfully", gin.H{"id": 1})
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Resource updated successfully")
}

func TestHandlerBase_SendDeletedResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	handler := &HandlerBase{}
	handler.SendDeletedResponse(c, "Resource deleted successfully")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Resource deleted successfully")
}

func TestHandlerBase_SendListResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	handler := &HandlerBase{}
	handler.SendListResponse(c, "Resources retrieved successfully", []gin.H{{"id": 1}})
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Resources retrieved successfully")
}
