package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
	"github.com/geoo115/Ecommerce/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestAddToCart_ValidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test data
	user := models.User{Username: "testuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	cat := models.Category{Name: "Test Category"}
	db.DB.Create(&cat)

	prod := models.Product{Name: "Test Product", Price: 29.99, CategoryID: cat.ID}
	db.DB.Create(&prod)

	inv := models.Inventory{ProductID: prod.ID, Stock: 10}
	db.DB.Create(&inv)

	// Set up JWT secret for token generation
	os.Setenv("JWT_SECRET", "test_secret_key")
	utils.AppLogger = utils.NewLogger(utils.INFO)

	router := gin.New()
	router.POST("/cart", func(c *gin.Context) {
		c.Set("userID", user.ID)
		AddToCart(c)
	})

	input := AddToCartInput{
		ProductID: prod.ID,
		Quantity:  2,
	}
	jsonData, _ := json.Marshal(input)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/cart", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Item added to cart", response["message"]) // Correct message for AddToCart operation
}

func TestAddToCart_InvalidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	router := gin.New()
	router.POST("/cart", func(c *gin.Context) {
		c.Set("userID", uint(1))
		AddToCart(c)
	})

	// Invalid JSON
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/cart", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Invalid request payload", response["error"])
}

func TestAddToCart_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	router := gin.New()
	router.POST("/cart", AddToCart)

	input := AddToCartInput{
		ProductID: 1,
		Quantity:  1,
	}
	jsonData, _ := json.Marshal(input)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/cart", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAddToCart_InsufficientStock(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test data with low stock
	user := models.User{Username: "testuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	cat := models.Category{Name: "Test Category"}
	db.DB.Create(&cat)

	prod := models.Product{Name: "Test Product", Price: 29.99, CategoryID: cat.ID}
	db.DB.Create(&prod)

	inv := models.Inventory{ProductID: prod.ID, Stock: 1} // Only 1 in stock
	db.DB.Create(&inv)

	router := gin.New()
	router.POST("/cart", func(c *gin.Context) {
		c.Set("userID", user.ID)
		AddToCart(c)
	})

	input := AddToCartInput{
		ProductID: prod.ID,
		Quantity:  5, // Request more than available
	}
	jsonData, _ := json.Marshal(input)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/cart", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Insufficient stock for product", response["error"])
}

func TestAddToCart_UpdateQuantity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test data
	user := models.User{Username: "testuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	cat := models.Category{Name: "Test Category"}
	db.DB.Create(&cat)

	prod := models.Product{Name: "Test Product", Price: 29.99, CategoryID: cat.ID}
	db.DB.Create(&prod)

	inv := models.Inventory{ProductID: prod.ID, Stock: 10}
	db.DB.Create(&inv)

	// Set up JWT secret for token generation
	os.Setenv("JWT_SECRET", "test_secret_key")
	utils.AppLogger = utils.NewLogger(utils.INFO)

	router := gin.New()
	router.POST("/cart", func(c *gin.Context) {
		c.Set("userID", user.ID)
		AddToCart(c)
	})

	// Add product to cart for the first time
	input1 := AddToCartInput{
		ProductID: prod.ID,
		Quantity:  2,
	}
	jsonData1, _ := json.Marshal(input1)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/cart", bytes.NewBuffer(jsonData1))
	req1.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	// Add the same product to the cart again
	input2 := AddToCartInput{
		ProductID: prod.ID,
		Quantity:  3,
	}
	jsonData2, _ := json.Marshal(input2)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/cart", bytes.NewBuffer(jsonData2))
	req2.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	// Check that the quantity has been updated
	var cartItem models.Cart
	db.DB.Where("user_id = ? AND product_id = ?", user.ID, prod.ID).First(&cartItem)
	assert.Equal(t, 5, cartItem.Quantity)

	// Check that a new cart item was not created
	var count int64
	db.DB.Model(&models.Cart{}).Where("user_id = ? AND product_id = ?", user.ID, prod.ID).Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestAddToCart_ProductNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test user
	user := models.User{Username: "testuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	// Set up JWT secret for token generation
	os.Setenv("JWT_SECRET", "test_secret_key")
	utils.AppLogger = utils.NewLogger(utils.INFO)

	router := gin.New()
	router.POST("/cart", func(c *gin.Context) {
		c.Set("userID", user.ID)
		AddToCart(c)
	})

	input := AddToCartInput{
		ProductID: 999, // Non-existent product
		Quantity:  2,
	}
	jsonData, _ := json.Marshal(input)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/cart", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "product not found")
}

func TestListCart_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Clean up any existing cart items
	db.DB.Where("1 = 1").Delete(&models.Cart{})

	// Create test data
	user := models.User{Username: "testuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	cat := models.Category{Name: "Test Category"}
	db.DB.Create(&cat)

	prod := models.Product{Name: "Test Product", Price: 29.99, CategoryID: cat.ID}
	db.DB.Create(&prod)

	inv := models.Inventory{ProductID: prod.ID, Stock: 10}
	db.DB.Create(&inv)

	// Add item to cart
	cartItem := models.Cart{
		UserID:    user.ID,
		ProductID: prod.ID,
		Quantity:  2,
	}
	db.DB.Create(&cartItem)

	router := gin.New()
	router.GET("/cart", func(c *gin.Context) {
		c.Set("userID", user.ID)
		ListCart(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/cart", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Cart items retrieved successfully", response["message"])

	// Check that data contains the expected fields
	data := response["data"].(map[string]interface{})
	assert.Equal(t, 59.98, data["total_amount"]) // 2 * 29.99
}

func TestListCart_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	router := gin.New()
	router.GET("/cart", ListCart)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/cart", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestListCart_Empty(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test user
	user := models.User{Username: "testuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	// Set up JWT secret for token generation
	os.Setenv("JWT_SECRET", "test_secret_key")
	utils.AppLogger = utils.NewLogger(utils.INFO)

	router := gin.New()
	router.GET("/cart", func(c *gin.Context) {
		c.Set("userID", user.ID)
		ListCart(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/cart", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Cart items retrieved successfully", response["message"])
	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok)
	cartItems, ok := data["cart_items"].([]interface{})
	assert.True(t, ok)
	assert.Empty(t, cartItems)
}

func TestRemoveFromCart_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test data
	user := models.User{Username: "testuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	cat := models.Category{Name: "Test Category"}
	db.DB.Create(&cat)

	prod := models.Product{Name: "Test Product", Price: 29.99, CategoryID: cat.ID}
	db.DB.Create(&prod)

	inv := models.Inventory{ProductID: prod.ID, Stock: 10}
	db.DB.Create(&inv)

	// Add item to cart
	cartItem := models.Cart{
		UserID:    user.ID,
		ProductID: prod.ID,
		Quantity:  2,
	}
	db.DB.Create(&cartItem)

	router := gin.New()
	router.DELETE("/cart/:id", func(c *gin.Context) {
		c.Set("userID", user.ID)
		RemoveFromCart(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/cart/%d", cartItem.ID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Cart item removed successfully", response["message"])
}

func TestRemoveFromCart_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	user := models.User{Username: "testuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	router := gin.New()
	router.DELETE("/cart/:id", func(c *gin.Context) {
		c.Set("userID", user.ID)
		RemoveFromCart(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/cart/999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateCartItem_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test data
	user := models.User{Username: "testuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	cat := models.Category{Name: "Test Category"}
	db.DB.Create(&cat)

	prod := models.Product{Name: "Test Product", Price: 29.99, CategoryID: cat.ID}
	db.DB.Create(&prod)

	inv := models.Inventory{ProductID: prod.ID, Stock: 10}
	db.DB.Create(&inv)

	// Add item to cart
	cartItem := models.Cart{
		UserID:    user.ID,
		ProductID: prod.ID,
		Quantity:  2,
	}
	db.DB.Create(&cartItem)

	router := gin.New()
	router.PUT("/cart/:id", func(c *gin.Context) {
		c.Set("userID", user.ID)
		UpdateCartItem(c)
	})

	updateInput := map[string]interface{}{
		"quantity": 5,
	}
	jsonData, _ := json.Marshal(updateInput)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/cart/%d", cartItem.ID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Cart item updated successfully", response["message"])

	// Verify the quantity in the database
	var updatedCartItem models.Cart
	db.DB.First(&updatedCartItem, cartItem.ID)
	assert.Equal(t, 5, updatedCartItem.Quantity)
}

func TestUpdateCartItem_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	user := models.User{Username: "testuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	router := gin.New()
	router.PUT("/cart/:id", func(c *gin.Context) {
		c.Set("userID", user.ID)
		UpdateCartItem(c)
	})

	updateInput := map[string]interface{}{
		"quantity": 5,
	}
	jsonData, _ := json.Marshal(updateInput)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/cart/invalid", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Invalid id", response["error"])
}

func TestUpdateCartItem_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	user := models.User{Username: "testuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	router := gin.New()
	router.PUT("/cart/:id", func(c *gin.Context) {
		c.Set("userID", user.ID)
		UpdateCartItem(c)
	})

	updateInput := map[string]interface{}{
		"quantity": 5,
	}
	jsonData, _ := json.Marshal(updateInput)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/cart/999", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Resource not found or access denied", response["error"])
}

func TestUpdateCartItem_InsufficientStock(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test data with low stock
	user := models.User{Username: "testuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	cat := models.Category{Name: "Test Category"}
	db.DB.Create(&cat)

	prod := models.Product{Name: "Test Product", Price: 29.99, CategoryID: cat.ID}
	db.DB.Create(&prod)

	inv := models.Inventory{ProductID: prod.ID, Stock: 3} // Only 3 in stock
	db.DB.Create(&inv)

	// Add item to cart with initial quantity
	cartItem := models.Cart{
		UserID:    user.ID,
		ProductID: prod.ID,
		Quantity:  1,
	}
	db.DB.Create(&cartItem)

	router := gin.New()
	router.PUT("/cart/:id", func(c *gin.Context) {
		c.Set("userID", user.ID)
		UpdateCartItem(c)
	})

	updateInput := map[string]interface{}{
		"quantity": 5, // Request more than available stock
	}
	jsonData, _ := json.Marshal(updateInput)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/cart/%d", cartItem.ID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Insufficient stock for product", response["error"])
}

func TestCheckStock_Sufficient(t *testing.T) {
	SetupTestDB(t)

	// Create test data
	cat := models.Category{Name: "Test Category"}
	db.DB.Create(&cat)

	prod := models.Product{Name: "Test Product", Price: 29.99, CategoryID: cat.ID}
	db.DB.Create(&prod)

	inv := models.Inventory{ProductID: prod.ID, Stock: 10}
	db.DB.Create(&inv)

	hasStock, err := CheckStock(prod.ID, 5)
	assert.NoError(t, err)
	assert.True(t, hasStock)
}

func TestCheckStock_Insufficient(t *testing.T) {
	SetupTestDB(t)

	// Create test data
	cat := models.Category{Name: "Test Category"}
	db.DB.Create(&cat)

	prod := models.Product{Name: "Test Product", Price: 29.99, CategoryID: cat.ID}
	db.DB.Create(&prod)

	inv := models.Inventory{ProductID: prod.ID, Stock: 3}
	db.DB.Create(&inv)

	hasStock, err := CheckStock(prod.ID, 5)
	assert.Error(t, err)
	assert.False(t, hasStock)
	assert.Contains(t, err.Error(), "Insufficient stock for product")
}

func TestCheckStock_ProductNotFound(t *testing.T) {
	SetupTestDB(t)

	hasStock, err := CheckStock(999, 1)
	assert.Error(t, err)
	assert.False(t, hasStock)
	assert.Contains(t, err.Error(), "product not found")
}

func BenchmarkAddToCart(b *testing.B) {
	gin.SetMode(gin.TestMode)
	testDB, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	testDB.AutoMigrate(&models.User{}, &models.Category{}, &models.Product{}, &models.Inventory{}, &models.Cart{})
	db.DB = testDB

	// Setup test data
	user := models.User{Username: "benchuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	cat := models.Category{Name: "Bench Category"}
	db.DB.Create(&cat)

	prod := models.Product{Name: "Bench Product", Price: 29.99, CategoryID: cat.ID}
	db.DB.Create(&prod)

	inv := models.Inventory{ProductID: prod.ID, Stock: 100}
	db.DB.Create(&inv)

	router := gin.New()
	router.POST("/cart", func(c *gin.Context) {
		c.Set("userID", user.ID)
		AddToCart(c)
	})

	input := AddToCartInput{
		ProductID: prod.ID,
		Quantity:  1,
	}
	jsonData, _ := json.Marshal(input)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/cart", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
	}
}

func BenchmarkListCart(b *testing.B) {
	gin.SetMode(gin.TestMode)
	testDB, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	testDB.AutoMigrate(&models.User{}, &models.Category{}, &models.Product{}, &models.Inventory{}, &models.Cart{})
	db.DB = testDB

	// Setup test data
	user := models.User{Username: "benchuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	cat := models.Category{Name: "Bench Category"}
	db.DB.Create(&cat)

	prod := models.Product{Name: "Bench Product", Price: 29.99, CategoryID: cat.ID}
	db.DB.Create(&prod)

	inv := models.Inventory{ProductID: prod.ID, Stock: 100}
	db.DB.Create(&inv)

	// Add items to cart
	for i := 0; i < 10; i++ {
		cartItem := models.Cart{
			UserID:    user.ID,
			ProductID: prod.ID,
			Quantity:  1,
		}
		db.DB.Create(&cartItem)
	}

	router := gin.New()
	router.GET("/cart", func(c *gin.Context) {
		c.Set("userID", user.ID)
		ListCart(c)
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/cart", nil)
		router.ServeHTTP(w, req)
	}
}
