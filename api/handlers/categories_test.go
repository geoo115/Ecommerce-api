package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
	"github.com/geoo115/Ecommerce/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestListCategories_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test categories using unique data
	cat1 := generateUniqueCategory()
	db.DB.Create(&cat1)

	cat2 := generateUniqueCategory()
	db.DB.Create(&cat2)

	router := gin.New()
	router.GET("/categories", ListCategories)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/categories", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Categories retrieved successfully", response["message"])
}

func TestListCategories_NoCategoriesFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Ensure no categories exist in the database
	db.DB.Exec("DELETE FROM categories")

	router := gin.New()
	router.GET("/categories", ListCategories)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/categories", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Categories retrieved successfully", response["message"])
	data, ok := response["data"].([]interface{})
	assert.True(t, ok)
	assert.Empty(t, data)
}

func TestListCategories_DatabaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Setup a test DB, but then close its underlying connection to simulate an error
	testDB := SetupTestDB(t)
	sqlDB, err := testDB.DB()
	assert.NoError(t, err)
	sqlDB.Close()

	// Use a local router and DB instance to avoid affecting global state
	router := gin.New()
	router.GET("/categories", ListCategories)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/categories", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Failed to fetch categories", response["error"])
}

func TestAddCategory_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Set up JWT secret
	os.Setenv("JWT_SECRET", "test_secret_key")
	utils.AppLogger = utils.NewLogger(utils.INFO)

	router := gin.New()
	router.POST("/categories", func(c *gin.Context) {
		// Simulate admin middleware setting user role
		c.Set("userRole", "admin")
		AddCategory(c)
	})

	category := generateUniqueCategory()
	jsonData, _ := json.Marshal(category)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/categories", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Category added successfully", response["message"])
}

func TestAddCategory_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	router := gin.New()
	router.POST("/categories", func(c *gin.Context) {
		c.Set("userRole", "admin")
		AddCategory(c)
	})

	// Invalid JSON
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/categories", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "invalid character 'i' looking for beginning of value")
}

func TestAddCategory_DatabaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Setup a test DB, but then close its underlying connection to simulate an error
	testDB := SetupTestDB(t)
	sqlDB, err := testDB.DB()
	assert.NoError(t, err)
	sqlDB.Close()

	// Temporarily set db.DB to this errored instance
	originalDB := db.DB
	defer func() { db.DB = originalDB }()
	db.DB = testDB

	router := gin.New()
	router.POST("/categories", func(c *gin.Context) {
		c.Set("userRole", "admin")
		AddCategory(c)
	})

	category := generateUniqueCategory()
	jsonData, _ := json.Marshal(category)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/categories", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Failed to create category", response["error"])
}

func TestAddCategory_InvalidName(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	router := gin.New()
	router.POST("/categories", func(c *gin.Context) {
		c.Set("userRole", "admin")
		AddCategory(c)
	})

	categoryData := map[string]interface{}{
		"name": "", // Invalid empty name
	}
	jsonData, _ := json.Marshal(categoryData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/categories", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddCategory_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	router := gin.New()
	router.POST("/categories", AddCategory)

	categoryData := map[string]interface{}{
		"name": "New Category",
	}
	jsonData, _ := json.Marshal(categoryData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/categories", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDeleteCategory_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test category
	cat := models.Category{Name: "Test Category"}
	db.DB.Create(&cat)

	// Get the ID of the created category
	var createdCat models.Category
	db.DB.Last(&createdCat)
	utils.AppLogger.Info("Category ID after creation: %d", createdCat.ID)

	router := gin.New()
	router.DELETE("/categories/:id", func(c *gin.Context) {
		c.Set("userRole", "admin")
		DeleteCategory(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/categories/"+strconv.FormatUint(uint64(createdCat.ID), 10), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Category deleted successfully", response["message"])
}

func TestDeleteCategory_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	router := gin.New()
	router.DELETE("/categories/:id", func(c *gin.Context) {
		c.Set("userRole", "admin")
		DeleteCategory(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/categories/999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteCategory_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	router := gin.New()
	router.DELETE("/categories/:id", func(c *gin.Context) {
		c.Set("userRole", "admin")
		DeleteCategory(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/categories/invalid", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "Invalid id")
}

func TestDeleteCategory_DatabaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Setup a test DB, but then close its underlying connection to simulate an error
	testDB := SetupTestDB(t)
	sqlDB, err := testDB.DB()
	assert.NoError(t, err)
	sqlDB.Close()

	// Temporarily set db.DB to this errored instance
	originalDB := db.DB
	defer func() { db.DB = originalDB }()
	db.DB = testDB

	router := gin.New()
	router.DELETE("/categories/:id", func(c *gin.Context) {
		c.Set("userRole", "admin")
		DeleteCategory(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/categories/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Failed to delete category", response["error"])
}

func TestDeleteCategory_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	cat := models.Category{Name: "Test Category"}
	db.DB.Create(&cat)

	// Get the ID of the created category
	var createdCat models.Category
	db.DB.Last(&createdCat)

	router := gin.New()
	router.DELETE("/categories/:id", DeleteCategory)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/categories/"+strconv.FormatUint(uint64(createdCat.ID), 10), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func BenchmarkListCategories(b *testing.B) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(b)

	// Setup test data - create multiple categories
	for i := 0; i < 50; i++ {
		cat := models.Category{Name: "Bench Category " + fmt.Sprintf("%d", i)}
		db.DB.Create(&cat)
	}

	router := gin.New()
	router.GET("/categories", ListCategories)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/categories", nil)
		router.ServeHTTP(w, req)
	}
}

func TestAddCategory_DuplicateName(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create first category
	cat1 := generateUniqueCategory()
	db.DB.Create(&cat1)

	router := gin.New()
	router.POST("/categories", func(c *gin.Context) {
		c.Set("userRole", "admin")
		AddCategory(c)
	})

	// Try to create category with same name
	categoryData := map[string]interface{}{
		"name": cat1.Name,
	}
	jsonData, _ := json.Marshal(categoryData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/categories", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Should fail due to unique constraint
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAddCategory_SpecialCharacters(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	router := gin.New()
	router.POST("/categories", func(c *gin.Context) {
		c.Set("userRole", "admin")
		AddCategory(c)
	})

	categoryData := map[string]interface{}{
		"name": "Electronics & Gadgets!",
	}
	jsonData, _ := json.Marshal(categoryData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/categories", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}
