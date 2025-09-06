package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupReportsTestDB(t *testing.T) *gorm.DB {
	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = testDB.AutoMigrate(&models.User{}, &models.Category{}, &models.Product{},
		&models.Inventory{}, &models.Order{}, &models.OrderItem{})
	assert.NoError(t, err)

	db.DB = testDB
	return testDB
}

func TestSalesReport_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupReportsTestDB(t)
	defer func() { db.DB = nil }()

	// Create test data
	user := models.User{Username: "testuser", Phone: "+15550000001"}
	testDB.Create(&user)

	category := models.Category{Name: "Electronics"}
	testDB.Create(&category)

	product1 := models.Product{
		Name:        "Laptop",
		Price:       999.99,
		Description: "Gaming laptop",
		CategoryID:  category.ID,
	}
	product2 := models.Product{
		Name:        "Mouse",
		Price:       29.99,
		Description: "Wireless mouse",
		CategoryID:  category.ID,
	}
	testDB.Create(&product1)
	testDB.Create(&product2)

	// Create orders
	order1 := models.Order{
		UserID:      user.ID,
		TotalAmount: 999.99,
		Status:      "completed",
	}
	order2 := models.Order{
		UserID:      user.ID,
		TotalAmount: 29.99,
		Status:      "completed",
	}
	testDB.Create(&order1)
	testDB.Create(&order2)

	// Create order items
	orderItem1 := models.OrderItem{
		OrderID:   order1.ID,
		ProductID: product1.ID,
		Quantity:  1,
		Price:     999.99,
	}
	orderItem2 := models.OrderItem{
		OrderID:   order2.ID,
		ProductID: product2.ID,
		Quantity:  2,
		Price:     29.99,
	}
	testDB.Create(&orderItem1)
	testDB.Create(&orderItem2)

	// Test request
	router := gin.New()
	router.GET("/reports/sales", SalesReport)

	req := httptest.NewRequest("GET", "/reports/sales", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "report")
}

func TestSalesReport_WithDateRange(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupReportsTestDB(t)
	defer func() { db.DB = nil }()

	// Create test data
	user := models.User{Username: "testuser", Phone: "+15550000001"}
	testDB.Create(&user)

	category := models.Category{Name: "Electronics"}
	testDB.Create(&category)

	product := models.Product{
		Name:       "Laptop",
		Price:      999.99,
		CategoryID: category.ID,
	}
	testDB.Create(&product)

	// Create order with specific date
	order := models.Order{
		UserID:      user.ID,
		TotalAmount: 999.99,
		Status:      "completed",
	}
	testDB.Create(&order)

	orderItem := models.OrderItem{
		OrderID:   order.ID,
		ProductID: product.ID,
		Quantity:  1,
		Price:     999.99,
	}
	testDB.Create(&orderItem)

	// Test request with date range
	router := gin.New()
	router.GET("/reports/sales", SalesReport)

	today := time.Now().Format("2006-01-02")
	req := httptest.NewRequest("GET", "/reports/sales?start_date="+today+"&end_date="+today, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "report")
}

func TestSalesReport_InvalidStartDate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_ = setupReportsTestDB(t)
	defer func() { db.DB = nil }()

	router := gin.New()
	router.GET("/reports/sales", SalesReport)

	req := httptest.NewRequest("GET", "/reports/sales?start_date=invalid-date", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
	assert.Contains(t, response["error"].(string), "Invalid start_date format")
}

func TestSalesReport_InvalidEndDate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_ = setupReportsTestDB(t)
	defer func() { db.DB = nil }()

	router := gin.New()
	router.GET("/reports/sales", SalesReport)

	req := httptest.NewRequest("GET", "/reports/sales?end_date=invalid-date", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
	assert.Contains(t, response["error"].(string), "Invalid end_date format")
}

func TestSalesReport_DatabaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db.DB = nil // Simulate database error

	router := gin.New()
	router.GET("/reports/sales", SalesReport)

	req := httptest.NewRequest("GET", "/reports/sales", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestInventoryReport_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupReportsTestDB(t)
	defer func() { db.DB = nil }()

	// Create test data
	category := models.Category{Name: "Electronics"}
	testDB.Create(&category)

	product1 := models.Product{
		Name:        "Laptop",
		Price:       999.99,
		Description: "Gaming laptop",
		CategoryID:  category.ID,
	}
	product2 := models.Product{
		Name:        "Mouse",
		Price:       29.99,
		Description: "Wireless mouse",
		CategoryID:  category.ID,
	}
	testDB.Create(&product1)
	testDB.Create(&product2)

	// Create inventory
	inventory1 := models.Inventory{
		ProductID: product1.ID,
		Stock:     10,
	}
	inventory2 := models.Inventory{
		ProductID: product2.ID,
		Stock:     5,
	}
	testDB.Create(&inventory1)
	testDB.Create(&inventory2)

	// Test request
	router := gin.New()
	router.GET("/reports/inventory", InventoryReport)

	req := httptest.NewRequest("GET", "/reports/inventory", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "inventory")
}

func TestInventoryReport_LowStockFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupReportsTestDB(t)
	defer func() { db.DB = nil }()

	// Create test data
	category := models.Category{Name: "Electronics"}
	testDB.Create(&category)

	product1 := models.Product{
		Name:       "Laptop",
		Price:      999.99,
		CategoryID: category.ID,
	}
	product2 := models.Product{
		Name:       "Mouse",
		Price:      29.99,
		CategoryID: category.ID,
	}
	testDB.Create(&product1)
	testDB.Create(&product2)

	// Create inventory - one with low stock
	inventory1 := models.Inventory{
		ProductID: product1.ID,
		Stock:     2, // Low stock
	}
	inventory2 := models.Inventory{
		ProductID: product2.ID,
		Stock:     50, // High stock
	}
	testDB.Create(&inventory1)
	testDB.Create(&inventory2)

	// Test request with low stock filter
	router := gin.New()
	router.GET("/reports/inventory", InventoryReport)

	req := httptest.NewRequest("GET", "/reports/inventory?low_stock=true", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "inventory")
}

func TestInventoryReport_CategoryFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testDB := setupReportsTestDB(t)
	defer func() { db.DB = nil }()

	// Create test data
	category1 := models.Category{Name: "Electronics"}
	category2 := models.Category{Name: "Books"}
	testDB.Create(&category1)
	testDB.Create(&category2)

	product1 := models.Product{
		Name:       "Laptop",
		Price:      999.99,
		CategoryID: category1.ID,
	}
	product2 := models.Product{
		Name:       "Novel",
		Price:      19.99,
		CategoryID: category2.ID,
	}
	testDB.Create(&product1)
	testDB.Create(&product2)

	// Create inventory
	inventory1 := models.Inventory{
		ProductID: product1.ID,
		Stock:     10,
	}
	inventory2 := models.Inventory{
		ProductID: product2.ID,
		Stock:     20,
	}
	testDB.Create(&inventory1)
	testDB.Create(&inventory2)

	// Test request with category filter
	router := gin.New()
	router.GET("/reports/inventory", InventoryReport)

	req := httptest.NewRequest("GET", "/reports/inventory?category_id=1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "inventory")
}

func TestInventoryReport_DatabaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db.DB = nil // Simulate database error

	router := gin.New()
	router.GET("/reports/inventory", InventoryReport)

	req := httptest.NewRequest("GET", "/reports/inventory", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestInventoryReport_EmptyInventory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_ = setupReportsTestDB(t)
	defer func() { db.DB = nil }()

	// No inventory data - test empty response
	router := gin.New()
	router.GET("/reports/inventory", InventoryReport)

	req := httptest.NewRequest("GET", "/reports/inventory", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "inventory")
}

func TestSalesReport_EmptyResults(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_ = setupReportsTestDB(t)
	defer func() { db.DB = nil }()

	// No sales data - test empty response
	router := gin.New()
	router.GET("/reports/sales", SalesReport)

	req := httptest.NewRequest("GET", "/reports/sales", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "report")
}
