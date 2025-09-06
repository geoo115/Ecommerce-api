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
	"time"

	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
	"github.com/geoo115/Ecommerce/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupReviewsTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	testDB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := testDB.AutoMigrate(&models.User{}, &models.Category{}, &models.Product{}, &models.Review{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}
	db.DB = testDB

	// Clear existing data
	tables := []string{"users", "categories", "products", "reviews"}
	for _, table := range tables {
		testDB.Exec(fmt.Sprintf("DELETE FROM %s", table))
	}

	// Create unique test data
	user := generateUniqueUser()
	db.DB.Create(&user)

	category := generateUniqueCategory()
	db.DB.Create(&category)

	product := generateUniqueProduct(category.ID)
	db.DB.Create(&product)

	t.Cleanup(func() {
		sqlDB, _ := testDB.DB()
		sqlDB.Close()
	})

	return testDB
}

func generateUniqueUser() models.User {
	return models.User{
		Username: fmt.Sprintf("testuser_%d", time.Now().UnixNano()),
		Email:    fmt.Sprintf("testuser_%d@example.com", time.Now().UnixNano()),
		Phone:    fmt.Sprintf("+1555%d", time.Now().UnixNano()%1000000),
		Password: "pw",
	}
}

func generateUniqueCategory() models.Category {
	return models.Category{
		Name: fmt.Sprintf("TestCategory_%d", time.Now().UnixNano()),
	}
}

func generateUniqueProduct(categoryID uint) models.Product {
	return models.Product{
		Name:        fmt.Sprintf("TestProduct_%d", time.Now().UnixNano()),
		Price:       29.99,
		CategoryID:  categoryID,
		Description: "Test desc",
	}
}

func TestAddReview_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupReviewsTestDB(t)

	// Create test data
	user := generateUniqueUser()
	db.DB.Create(&user)

	cat := generateUniqueCategory()
	db.DB.Create(&cat)

	prod := generateUniqueProduct(cat.ID)
	db.DB.Create(&prod)

	// Set up JWT secret
	os.Setenv("JWT_SECRET", "test_secret_key")
	utils.AppLogger = utils.NewLogger(utils.INFO)

	router := gin.New()
	router.POST("/products/:id/reviews", func(c *gin.Context) {
		c.Set("userID", user.ID)
		AddReview(c)
	})

	reviewData := map[string]interface{}{
		"rating":  5,
		"comment": "Great product!",
	}
	jsonData, _ := json.Marshal(reviewData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/products/"+strconv.Itoa(int(prod.ID))+"/reviews", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Review added successfully", response["message"])
}

func TestAddReview_InvalidRating(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupReviewsTestDB(t)

	user := models.User{Username: "testuser2", Phone: "+15550000002", Password: "pw"}
	db.DB.Create(&user)

	cat := models.Category{Name: "Test Category"}
	db.DB.Create(&cat)

	prod := models.Product{Name: "Test Product", Price: 29.99, CategoryID: cat.ID, Description: "Test desc"}
	db.DB.Create(&prod)

	router := gin.New()
	router.POST("/products/:id/reviews", func(c *gin.Context) {
		c.Set("userID", user.ID)
		AddReview(c)
	})

	reviewData := map[string]interface{}{
		"rating":  10, // Invalid rating (should be 1-5)
		"comment": "Great product!",
	}
	jsonData, _ := json.Marshal(reviewData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/products/"+strconv.Itoa(int(prod.ID))+"/reviews", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddReview_ProductNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupReviewsTestDB(t)

	user := models.User{Username: "testuser3", Phone: "+15550000003", Password: "pw"}
	db.DB.Create(&user)

	router := gin.New()
	router.POST("/products/:id/reviews", func(c *gin.Context) {
		c.Set("userID", user.ID)
		AddReview(c)
	})

	reviewData := map[string]interface{}{
		"rating":  5,
		"comment": "Great product!",
	}
	jsonData, _ := json.Marshal(reviewData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/products/999/reviews", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAddReview_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupReviewsTestDB(t)

	cat := models.Category{Name: "Test Category"}
	db.DB.Create(&cat)

	prod := models.Product{Name: "Test Product", Price: 29.99, CategoryID: cat.ID, Description: "Test desc"}
	db.DB.Create(&prod)

	router := gin.New()
	router.POST("/products/:id/reviews", AddReview)

	reviewData := map[string]interface{}{
		"rating":  5,
		"comment": "Great product!",
	}
	jsonData, _ := json.Marshal(reviewData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/products/"+strconv.Itoa(int(prod.ID))+"/reviews", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestListReviews_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupReviewsTestDB(t)

	// Create test data
	user := models.User{Username: "testuser4", Phone: "+15550000004", Password: "pw"}
	db.DB.Create(&user)

	cat := models.Category{Name: "Test Category"}
	db.DB.Create(&cat)

	prod := models.Product{Name: "Test Product", Price: 29.99, CategoryID: cat.ID, Description: "Test desc"}
	db.DB.Create(&prod)

	// Create review
	review := models.Review{
		UserID:    user.ID,
		ProductID: prod.ID,
		Rating:    5,
		Comment:   "Great product!",
	}
	db.DB.Create(&review)

	router := gin.New()
	router.GET("/products/:id/reviews", ListReviews)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/products/"+strconv.Itoa(int(prod.ID))+"/reviews", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Reviews retrieved successfully", response["message"])
}

func TestAddReview_DuplicateReview(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupReviewsTestDB(t)

	// Create test data
	user := generateUniqueUser()
	db.DB.Create(&user)

	cat := generateUniqueCategory()
	db.DB.Create(&cat)

	prod := generateUniqueProduct(cat.ID)
	db.DB.Create(&prod)

	// Create first review
	review := models.Review{
		UserID:    user.ID,
		ProductID: prod.ID,
		Rating:    4,
		Comment:   "Good product",
	}
	db.DB.Create(&review)

	router := gin.New()
	router.POST("/products/:id/reviews", func(c *gin.Context) {
		c.Set("userID", user.ID)
		AddReview(c)
	})

	// Try to add another review for the same product by the same user
	reviewData := map[string]interface{}{
		"rating":  5,
		"comment": "Updated review",
	}
	jsonData, _ := json.Marshal(reviewData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/products/"+strconv.Itoa(int(prod.ID))+"/reviews", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Should allow updating existing review or handle gracefully
	assert.True(t, w.Code == http.StatusCreated || w.Code == http.StatusOK)
}

func TestAddReview_MissingRating(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupReviewsTestDB(t)

	user := generateUniqueUser()
	db.DB.Create(&user)

	cat := generateUniqueCategory()
	db.DB.Create(&cat)

	prod := generateUniqueProduct(cat.ID)
	db.DB.Create(&prod)

	router := gin.New()
	router.POST("/products/:id/reviews", func(c *gin.Context) {
		c.Set("userID", user.ID)
		AddReview(c)
	})

	reviewData := map[string]interface{}{
		"comment": "Missing rating",
	}
	jsonData, _ := json.Marshal(reviewData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/products/"+strconv.Itoa(int(prod.ID))+"/reviews", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddReview_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupReviewsTestDB(t)

	user := generateUniqueUser()
	db.DB.Create(&user)

	cat := generateUniqueCategory()
	db.DB.Create(&cat)

	prod := generateUniqueProduct(cat.ID)
	db.DB.Create(&prod)

	router := gin.New()
	router.POST("/products/:id/reviews", func(c *gin.Context) {
		c.Set("userID", user.ID)
		AddReview(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/products/"+strconv.Itoa(int(prod.ID))+"/reviews", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestListReviews_NoReviews(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupReviewsTestDB(t)

	cat := generateUniqueCategory()
	db.DB.Create(&cat)

	prod := generateUniqueProduct(cat.ID)
	db.DB.Create(&prod)

	router := gin.New()
	router.GET("/products/:id/reviews", ListReviews)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/products/"+strconv.Itoa(int(prod.ID))+"/reviews", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "No reviews found for this product", response["error"])
}

func BenchmarkAddReview(b *testing.B) {
	gin.SetMode(gin.TestMode)
	testDB, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	testDB.AutoMigrate(&models.User{}, &models.Category{}, &models.Product{}, &models.Review{})
	db.DB = testDB

	// Setup test data
	user := models.User{Username: "benchuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	cat := models.Category{Name: "Bench Category"}
	db.DB.Create(&cat)

	prod := models.Product{Name: "Bench Product", Price: 29.99, CategoryID: cat.ID, Description: "Bench desc"}
	db.DB.Create(&prod)

	router := gin.New()
	router.POST("/products/"+strconv.Itoa(int(prod.ID))+"/reviews", func(c *gin.Context) {
		c.Set("userID", user.ID)
		AddReview(c)
	})

	reviewData := map[string]interface{}{
		"rating":  5,
		"comment": "Great product!",
	}
	jsonData, _ := json.Marshal(reviewData)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/products/"+strconv.Itoa(int(prod.ID))+"/reviews", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
	}
}

func BenchmarkListReviews(b *testing.B) {
	gin.SetMode(gin.TestMode)
	testDB, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	testDB.AutoMigrate(&models.User{}, &models.Category{}, &models.Product{}, &models.Review{})
	db.DB = testDB

	// Setup test data
	user := models.User{Username: "benchuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	cat := models.Category{Name: "Bench Category"}
	db.DB.Create(&cat)

	prod := models.Product{Name: "Bench Product", Price: 29.99, CategoryID: cat.ID, Description: "Bench desc"}
	db.DB.Create(&prod)

	// Create multiple reviews
	for i := 0; i < 20; i++ {
		review := models.Review{
			UserID:    user.ID,
			ProductID: prod.ID,
			Rating:    5,
			Comment:   "Great product!",
		}
		db.DB.Create(&review)
	}

	router := gin.New()
	router.GET("/products/"+strconv.Itoa(int(prod.ID))+"/reviews", ListReviews)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/products/"+strconv.Itoa(int(prod.ID))+"/reviews", nil)
		router.ServeHTTP(w, req)
	}
}
