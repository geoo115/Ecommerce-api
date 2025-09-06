package handlers

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestOptimizedListProducts_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test data
	cat := models.Category{Name: "Test Category"}
	db.DB.Create(&cat)

	prod1 := models.Product{Name: "Product 1", Price: 10.0, CategoryID: cat.ID}
	db.DB.Create(&prod1)

	prod2 := models.Product{Name: "Product 2", Price: 20.0, CategoryID: cat.ID}
	db.DB.Create(&prod2)

	router := gin.New()
	router.GET("/products", OptimizedListProducts)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/products", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOptimizedListProducts_WithCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test data
	cat1 := models.Category{Name: "Category 1"}
	db.DB.Create(&cat1)

	cat2 := models.Category{Name: "Category 2"}
	db.DB.Create(&cat2)

	prod1 := models.Product{Name: "Product 1", Price: 10.0, CategoryID: cat1.ID}
	db.DB.Create(&prod1)

	prod2 := models.Product{Name: "Product 2", Price: 20.0, CategoryID: cat2.ID}
	db.DB.Create(&prod2)

	router := gin.New()
	router.GET("/products", OptimizedListProducts)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/products?category_id="+strconv.Itoa(int(cat1.ID)), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOptimizedGetProduct_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test data
	cat := models.Category{Name: "Test Category"}
	db.DB.Create(&cat)

	prod := models.Product{Name: "Test Product", Price: 10.0, CategoryID: cat.ID}
	db.DB.Create(&prod)

	router := gin.New()
	router.GET("/products/:id", OptimizedGetProduct)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/products/"+strconv.Itoa(int(prod.ID)), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOptimizedGetProduct_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	router := gin.New()
	router.GET("/products/:id", OptimizedGetProduct)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/products/999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestOptimizedListCart_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test data
	user := models.User{Username: "testuser", Email: "test@example.com"}
	db.DB.Create(&user)

	cat := models.Category{Name: "Test Category"}
	db.DB.Create(&cat)

	prod := models.Product{Name: "Test Product", Price: 10.0, CategoryID: cat.ID}
	db.DB.Create(&prod)

	cart := models.Cart{UserID: user.ID, ProductID: prod.ID, Quantity: 2}
	db.DB.Create(&cart)

	router := gin.New()
	router.GET("/cart", func(c *gin.Context) {
		c.Set("userID", user.ID)
		OptimizedListCart(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/cart", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOptimizedGetUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test data
	user := models.User{Username: "testuser", Email: "test@example.com"}
	db.DB.Create(&user)

	router := gin.New()
	router.GET("/users/:id", func(c *gin.Context) {
		c.Set("userID", user.ID) // Set userID in context
		OptimizedGetUser(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/"+strconv.Itoa(int(user.ID)), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOptimizedHealthCheck_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	router := gin.New()
	router.GET("/health", OptimizedHealthCheck)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOptimizedProductSearch_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test data
	cat := models.Category{Name: "Test Category"}
	db.DB.Create(&cat)

	prod := models.Product{Name: "Test Product", Price: 10.0, CategoryID: cat.ID}
	db.DB.Create(&prod)

	router := gin.New()
	router.GET("/products/search", OptimizedProductSearch)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/products/search?q=Test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOptimizedOrderHistory_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test data
	user := models.User{Username: "testuser", Email: "test@example.com"}
	db.DB.Create(&user)

	order := models.Order{UserID: user.ID, TotalAmount: 100.0, Status: "Completed"}
	db.DB.Create(&order)

	router := gin.New()
	router.GET("/orders/history", func(c *gin.Context) {
		c.Set("userID", user.ID)
		OptimizedOrderHistory(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/orders/history", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
