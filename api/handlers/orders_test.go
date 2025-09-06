package handlers

import (
	"bytes"
	"encoding/json"
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
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestPlaceOrder_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test data
	user := models.User{Username: "testuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	cat := models.Category{Name: "Test Category"}
	db.DB.Create(&cat)

	prod := models.Product{Name: "Test Product", Price: 29.99, CategoryID: cat.ID, Description: "Test desc"}
	db.DB.Create(&prod)

	inv := models.Inventory{ProductID: prod.ID, Stock: 10}
	db.DB.Create(&inv)

	// Set up JWT secret
	os.Setenv("JWT_SECRET", "test_secret_key")
	utils.AppLogger = utils.NewLogger(utils.INFO)

	// Create order request payload
	orderRequest := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"product_id": prod.ID,
				"quantity":   2,
			},
		},
	}
	jsonData, _ := json.Marshal(orderRequest)

	router := gin.New()
	router.POST("/orders", func(c *gin.Context) {
		c.Set("userID", user.ID)
		PlaceOrder(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Order placed successfully", response["message"])

	// Verify inventory was reduced
	var updatedInventory models.Inventory
	db.DB.Where("product_id = ?", prod.ID).First(&updatedInventory)
	assert.Equal(t, 8, updatedInventory.Stock) // 10 - 2 = 8
}

func TestPlaceOrder_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	user := models.User{Username: "testuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	router := gin.New()
	router.POST("/orders", func(c *gin.Context) {
		c.Set("userID", user.ID)
		PlaceOrder(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/orders", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPlaceOrder_EmptyItems(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	user := models.User{Username: "testuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	// Create order request with empty items
	orderRequest := map[string]interface{}{
		"items": []map[string]interface{}{},
	}
	jsonData, _ := json.Marshal(orderRequest)

	router := gin.New()
	router.POST("/orders", func(c *gin.Context) {
		c.Set("userID", user.ID)
		PlaceOrder(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPlaceOrder_ProductNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	user := models.User{Username: "testuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	// Create order request with non-existent product
	orderRequest := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"product_id": 999,
				"quantity":   1,
			},
		},
	}
	jsonData, _ := json.Marshal(orderRequest)

	router := gin.New()
	router.POST("/orders", func(c *gin.Context) {
		c.Set("userID", user.ID)
		PlaceOrder(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPlaceOrder_InvalidQuantity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test data
	user := models.User{Username: "testuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	cat := models.Category{Name: "Test Category"}
	db.DB.Create(&cat)

	prod := models.Product{Name: "Test Product", Price: 29.99, CategoryID: cat.ID, Description: "Test desc"}
	db.DB.Create(&prod)

	inv := models.Inventory{ProductID: prod.ID, Stock: 10}
	db.DB.Create(&inv)

	// Create order request with invalid quantity (0)
	orderRequest := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"product_id": prod.ID,
				"quantity":   0,
			},
		},
	}
	jsonData, _ := json.Marshal(orderRequest)

	router := gin.New()
	router.POST("/orders", func(c *gin.Context) {
		c.Set("userID", user.ID)
		PlaceOrder(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPlaceOrder_EmptyCart(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	user := models.User{Username: "testuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	router := gin.New()
	router.POST("/orders", func(c *gin.Context) {
		c.Set("userID", user.ID)
		PlaceOrder(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/orders", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestListOrders_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test data
	user := models.User{Username: "testuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	cat := models.Category{Name: "Test Category"}
	db.DB.Create(&cat)

	prod := models.Product{Name: "Test Product", Price: 29.99, CategoryID: cat.ID, Description: "Test desc"}
	db.DB.Create(&prod)

	inv := models.Inventory{ProductID: prod.ID, Stock: 10}
	db.DB.Create(&inv)

	// Create order
	order := models.Order{
		UserID:      user.ID,
		TotalAmount: 59.98,
		Status:      "Pending",
	}
	db.DB.Create(&order)

	orderItem := models.OrderItem{
		OrderID:   order.ID,
		ProductID: prod.ID,
		Quantity:  2,
		Price:     29.99,
	}
	db.DB.Create(&orderItem)

	router := gin.New()
	router.GET("/orders", func(c *gin.Context) {
		c.Set("userID", user.ID)
		ListOrders(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/orders", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Orders retrieved successfully", response["message"])
}

func TestGetOrder_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test data
	user := models.User{Username: "testuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	cat := models.Category{Name: "Test Category"}
	db.DB.Create(&cat)

	prod := models.Product{Name: "Test Product", Price: 29.99, CategoryID: cat.ID, Description: "Test desc"}
	db.DB.Create(&prod)

	inv := models.Inventory{ProductID: prod.ID, Stock: 10}
	db.DB.Create(&inv)

	// Create order
	order := models.Order{
		UserID:      user.ID,
		TotalAmount: 59.98,
		Status:      "Pending",
	}
	db.DB.Create(&order)

	orderItem := models.OrderItem{
		OrderID:   order.ID,
		ProductID: prod.ID,
		Quantity:  2,
		Price:     29.99,
	}
	db.DB.Create(&orderItem)

	router := gin.New()
	router.GET("/orders/:id", func(c *gin.Context) {
		c.Set("userID", user.ID)
		GetOrder(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/orders/"+strconv.Itoa(int(order.ID)), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Order retrieved successfully", response["message"])
}

func TestGetOrder_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	user := models.User{Username: "testuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	router := gin.New()
	router.GET("/orders/:id", func(c *gin.Context) {
		c.Set("userID", user.ID)
		GetOrder(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/orders/999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCancelOrder_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test data
	user := models.User{Username: "testuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	cat := models.Category{Name: "Test Category"}
	db.DB.Create(&cat)

	prod := models.Product{Name: "Test Product", Price: 29.99, CategoryID: cat.ID, Description: "Test desc"}
	db.DB.Create(&prod)

	inv := models.Inventory{ProductID: prod.ID, Stock: 10}
	db.DB.Create(&inv)

	// Create order
	order := models.Order{
		UserID:      user.ID,
		TotalAmount: 59.98,
		Status:      "Pending",
	}
	db.DB.Create(&order)

	orderItem := models.OrderItem{
		OrderID:   order.ID,
		ProductID: prod.ID,
		Quantity:  2,
		Price:     29.99,
	}
	db.DB.Create(&orderItem)

	router := gin.New()
	router.PUT("/orders/:id/cancel", func(c *gin.Context) {
		c.Set("userID", user.ID)
		CancelOrder(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/orders/"+strconv.Itoa(int(order.ID))+"/cancel", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Order cancelled successfully", response["message"])
}

func TestCancelOrder_AlreadyShipped(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test data
	user := models.User{Username: "testuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	// Create shipped order
	order := models.Order{
		UserID:      user.ID,
		TotalAmount: 59.98,
		Status:      "Shipped", // Already shipped
	}
	db.DB.Create(&order)

	router := gin.New()
	router.PUT("/orders/:id/cancel", func(c *gin.Context) {
		c.Set("userID", user.ID)
		CancelOrder(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/orders/"+strconv.Itoa(int(order.ID))+"/cancel", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func BenchmarkPlaceOrder(b *testing.B) {
	gin.SetMode(gin.TestMode)
	testDB, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	testDB.AutoMigrate(&models.User{}, &models.Category{}, &models.Product{}, &models.Inventory{}, &models.Order{}, &models.OrderItem{}, &models.Cart{}, &models.Address{})
	db.DB = testDB

	// Setup test data
	user := models.User{Username: "benchuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	cat := models.Category{Name: "Bench Category"}
	db.DB.Create(&cat)

	prod := models.Product{Name: "Bench Product", Price: 29.99, CategoryID: cat.ID, Description: "Bench desc"}
	db.DB.Create(&prod)

	inv := models.Inventory{ProductID: prod.ID, Stock: 100}
	db.DB.Create(&inv)

	// Add item to cart
	cartItem := models.Cart{
		UserID:    user.ID,
		ProductID: prod.ID,
		Quantity:  1,
	}
	db.DB.Create(&cartItem)

	router := gin.New()
	router.POST("/orders", func(c *gin.Context) {
		c.Set("userID", user.ID)
		PlaceOrder(c)
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/orders", nil)
		router.ServeHTTP(w, req)
	}
}

func BenchmarkListOrders(b *testing.B) {
	gin.SetMode(gin.TestMode)
	testDB, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	testDB.AutoMigrate(&models.User{}, &models.Category{}, &models.Product{}, &models.Inventory{}, &models.Order{}, &models.OrderItem{}, &models.Cart{}, &models.Address{})
	db.DB = testDB

	// Setup test data
	user := models.User{Username: "benchuser", Phone: "+15550000001", Password: "pw"}
	db.DB.Create(&user)

	cat := models.Category{Name: "Bench Category"}
	db.DB.Create(&cat)

	prod := models.Product{Name: "Bench Product", Price: 29.99, CategoryID: cat.ID, Description: "Bench desc"}
	db.DB.Create(&prod)

	inv := models.Inventory{ProductID: prod.ID, Stock: 100}
	db.DB.Create(&inv)

	// Create multiple orders
	for i := 0; i < 20; i++ {
		order := models.Order{
			UserID:      user.ID,
			TotalAmount: 29.99,
			Status:      "Pending",
			Items: []models.OrderItem{
				{ProductID: prod.ID, Quantity: 1, Price: 29.99},
			},
		}
		db.DB.Create(&order)
	}

	router := gin.New()
	router.GET("/orders", func(c *gin.Context) {
		c.Set("userID", user.ID)
		ListOrders(c)
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/orders", nil)
		router.ServeHTTP(w, req)
	}
}
