package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/geoo115/Ecommerce/api/middlewares"
	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Test ListProducts with different scenarios
func TestListProducts_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test products
	category := models.Category{Name: "Test Category"}
	db.DB.Create(&category)

	products := []models.Product{
		{Name: "Product 1", Price: 10.99, CategoryID: category.ID, Description: "Desc 1"},
		{Name: "Product 2", Price: 20.99, CategoryID: category.ID, Description: "Desc 2"},
		{Name: "Product 3", Price: 30.99, CategoryID: category.ID, Description: "Desc 3"},
	}

	for _, p := range products {
		db.DB.Create(&p)
		// Create inventory for each product
		inventory := models.Inventory{ProductID: p.ID, Stock: 10}
		db.DB.Create(&inventory)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create a proper request with URL and query parameters
	req, _ := http.NewRequest("GET", "/products?page=1&limit=10", nil)
	c.Request = req

	ListProducts(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Product 1")
	assert.Contains(t, w.Body.String(), "Product 2")
	assert.Contains(t, w.Body.String(), "Product 3")
}

func TestListProducts_EmptyDatabase(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create a proper request with URL and query parameters
	req, _ := http.NewRequest("GET", "/products?page=1&limit=10", nil)
	c.Request = req

	ListProducts(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Products retrieved successfully")
}

func TestGetProduct_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test product
	category := models.Category{Name: "Test Category"}
	db.DB.Create(&category)

	product := models.Product{
		Name:        "Test Product",
		Price:       15.99,
		CategoryID:  category.ID,
		Description: "Test Description",
	}
	db.DB.Create(&product)

	// Create inventory for the product
	inventory := models.Inventory{ProductID: product.ID, Stock: 10}
	db.DB.Create(&inventory)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{Method: "GET"}
	c.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(product.ID))}}

	GetProduct(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test Product")
	assert.Contains(t, w.Body.String(), "15.99")
}

func TestGetProduct_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{Method: "GET"}
	c.Params = []gin.Param{{Key: "id", Value: "999"}}

	GetProduct(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Product not found")
}

func TestGetProduct_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{Method: "GET"}
	c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

	GetProduct(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid product ID")
}

func TestAddProduct_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test category
	category := models.Category{Name: "Test Category"}
	db.DB.Create(&category)

	productInput := middlewares.ProductInput{
		Name:        "New Product",
		Price:       25.99,
		CategoryID:  category.ID,
		Description: "New product description",
		Stock:       10,
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{Method: "POST"}
	c.Set("product_input", productInput)

	AddProduct(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Product created successfully")
	assert.Contains(t, w.Body.String(), "New Product")
}

func TestAddProduct_InvalidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{Method: "POST"}
	// Don't set product_input

	AddProduct(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid input")
}

func TestAddProduct_InvalidCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	productInput := middlewares.ProductInput{
		Name:        "Test Product",
		Price:       29.99,
		CategoryID:  999, // Non-existent category
		Description: "Test description",
		Stock:       10,
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{Method: "POST"}
	c.Set("product_input", productInput)

	AddProduct(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAddProduct_InvalidName(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	category := models.Category{Name: "Test Category"}
	db.DB.Create(&category)

	productInput := middlewares.ProductInput{
		Name:        "", // Invalid name
		Price:       29.99,
		CategoryID:  category.ID,
		Description: "Test description",
		Stock:       10,
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{Method: "POST"}
	c.Set("product_input", productInput)

	AddProduct(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEditProduct_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test data
	category := models.Category{Name: "Test Category"}
	db.DB.Create(&category)

	product := models.Product{Name: "Test Product", Price: 29.99, CategoryID: category.ID, Description: "Test desc"}
	db.DB.Create(&product)

	inventory := models.Inventory{ProductID: product.ID, Stock: 10}
	db.DB.Create(&inventory)

	router := gin.New()
	router.PUT("/products/:id", EditProductHandlerWrapper(db.DB))

	// Create JSON payload
	updateData := map[string]interface{}{
		"name":        "Updated Product",
		"price":       39.99,
		"description": "Updated description",
		"stock":       15,
	}
	jsonData, _ := json.Marshal(updateData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/products/"+strconv.Itoa(int(product.ID)), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Product updated successfully", response["message"])
}

func TestEditProduct_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	router := gin.New()
	router.PUT("/products/:id", EditProductHandlerWrapper(db.DB))

	updateData := map[string]interface{}{
		"name":  "Updated Product",
		"price": 19.99,
	}
	jsonData, _ := json.Marshal(updateData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/products/999", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteProduct_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test product
	category := models.Category{Name: "Test Category"}
	db.DB.Create(&category)

	product := models.Product{
		Name:        "Product to Delete",
		Price:       15.99,
		CategoryID:  category.ID,
		Description: "Test description",
	}
	db.DB.Create(&product)

	// Create inventory for the product
	inventory := models.Inventory{ProductID: product.ID, Stock: 5}
	db.DB.Create(&inventory)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{Method: "DELETE"}
	c.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(product.ID))}}

	DeleteProduct(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Product deleted successfully")

	// Verify product is deleted
	var deletedProduct models.Product
	err := db.DB.First(&deletedProduct, product.ID).Error
	assert.Error(t, err) // Should not find the product
}

func TestDeleteProduct_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{Method: "DELETE"}
	c.Params = []gin.Param{{Key: "id", Value: "999"}}

	DeleteProduct(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Product not found")
}

func TestSearchProducts_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test data
	category := models.Category{Name: "Test Category"}
	db.DB.Create(&category)

	product := models.Product{Name: "Test Product", Price: 29.99, CategoryID: category.ID, Description: "Test desc"}
	db.DB.Create(&product)

	inventory := models.Inventory{ProductID: product.ID, Stock: 10}
	db.DB.Create(&inventory)

	router := gin.New()
	router.GET("/products/search", SearchProducts)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/products/search?q=Test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Search completed successfully", response["message"])
}

func TestSearchProducts_NoQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	router := gin.New()
	router.GET("/products/search", SearchProducts)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/products/search", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Benchmark tests
func BenchmarkAddProduct(b *testing.B) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(b)

	// Setup test data
	cat := models.Category{Name: "Bench Category"}
	db.DB.Create(&cat)

	router := gin.New()
	router.POST("/products", func(c *gin.Context) {
		input := middlewares.ProductInput{
			Name:        "Bench Product",
			Price:       29.99,
			CategoryID:  cat.ID,
			Description: "Bench description",
			Stock:       10,
		}
		c.Set("product_input", input)
		AddProduct(c)
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/products", nil)
		router.ServeHTTP(w, req)
	}
}

func BenchmarkListProducts(b *testing.B) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(b)

	// Setup test data
	cat := models.Category{Name: "Bench Category"}
	db.DB.Create(&cat)

	// Create multiple products
	for i := 0; i < 50; i++ {
		prod := models.Product{
			Name:        "Bench Product " + strconv.Itoa(i),
			Price:       29.99,
			CategoryID:  cat.ID,
			Description: "Bench description",
		}
		db.DB.Create(&prod)

		inv := models.Inventory{ProductID: prod.ID, Stock: 10}
		db.DB.Create(&inv)
	}

	router := gin.New()
	router.GET("/products", ListProducts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/products", nil)
		router.ServeHTTP(w, req)
	}
}

func BenchmarkSearchProducts(b *testing.B) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(b)

	// Setup test data
	cat := models.Category{Name: "Bench Category"}
	db.DB.Create(&cat)

	// Create multiple products
	for i := 0; i < 50; i++ {
		prod := models.Product{
			Name:        "Bench Product " + strconv.Itoa(i),
			Price:       29.99,
			CategoryID:  cat.ID,
			Description: "Bench description",
		}
		db.DB.Create(&prod)

		inv := models.Inventory{ProductID: prod.ID, Stock: 10}
		db.DB.Create(&inv)
	}

	router := gin.New()
	router.GET("/products/search", SearchProducts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/products/search?q=Bench", nil)
		router.ServeHTTP(w, req)
	}
}

// Additional edge case tests
func TestAddProduct_ZeroPrice(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	category := models.Category{Name: "Test Category"}
	db.DB.Create(&category)

	productInput := middlewares.ProductInput{
		Name:        "Free Product",
		Price:       0, // Zero price
		CategoryID:  category.ID,
		Description: "Free product",
		Stock:       10,
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{Method: "POST"}
	c.Set("product_input", productInput)

	AddProduct(c)

	// Should accept zero price for free products
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestAddProduct_NegativeStock(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	category := models.Category{Name: "Test Category"}
	db.DB.Create(&category)

	productInput := middlewares.ProductInput{
		Name:        "Invalid Stock Product",
		Price:       29.99,
		CategoryID:  category.ID,
		Description: "Invalid stock",
		Stock:       -5, // Negative stock
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{Method: "POST"}
	c.Set("product_input", productInput)

	AddProduct(c)

	// Should reject negative stock
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Performance edge cases
func TestListProducts_LargeDataset(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create large dataset
	category := models.Category{Name: "Large Category"}
	db.DB.Create(&category)

	// Create 100 products to test pagination/performance
	for i := 0; i < 100; i++ {
		product := models.Product{
			Name:        "Product " + strconv.Itoa(i),
			Price:       float64(i) * 1.5,
			CategoryID:  category.ID,
			Description: "Description " + strconv.Itoa(i),
		}
		db.DB.Create(&product)

		inventory := models.Inventory{ProductID: product.ID, Stock: i + 1}
		db.DB.Create(&inventory)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// Provide URL so c.Query works and doesn't panic
	req, _ := http.NewRequest("GET", "/products?page=1&limit=10", nil)
	c.Request = req

	ListProducts(c)

	assert.Equal(t, http.StatusOK, w.Code)
	// Response should contain paginated results
	assert.Contains(t, w.Body.String(), "Products retrieved successfully")
}

func TestListProductsHandlerWrapper(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test products
	category := models.Category{Name: "Test Category"}
	db.DB.Create(&category)

	products := []models.Product{
		{Name: "Product 1", Price: 10.99, CategoryID: category.ID, Description: "Desc 1"},
		{Name: "Product 2", Price: 20.99, CategoryID: category.ID, Description: "Desc 2"},
	}

	for _, p := range products {
		db.DB.Create(&p)
		// Create inventory for each product
		inventory := models.Inventory{ProductID: p.ID, Stock: 10}
		db.DB.Create(&inventory)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create a proper request with URL and query parameters
	req, _ := http.NewRequest("GET", "/products?page=1&limit=10", nil)
	c.Request = req

	// Test the wrapper function
	handler := ListProductsHandlerWrapper(db.DB)
	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Product 1")
	assert.Contains(t, w.Body.String(), "Product 2")
}

func TestGetProductHandlerWrapper(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test product
	category := models.Category{Name: "Test Category"}
	db.DB.Create(&category)

	product := models.Product{
		Name:        "Test Product",
		Price:       29.99,
		CategoryID:  category.ID,
		Description: "Test Description",
	}
	db.DB.Create(&product)

	// Create inventory for the product
	inventory := models.Inventory{ProductID: product.ID, Stock: 10}
	db.DB.Create(&inventory)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create a proper request with product ID parameter
	req, _ := http.NewRequest("GET", "/products/"+strconv.Itoa(int(product.ID)), nil)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(product.ID))}}

	// Test the wrapper function
	handler := GetProductHandlerWrapper(db.DB)
	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test Product")
}

func TestDeleteProductHandlerWrapper(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test product
	category := models.Category{Name: "Test Category"}
	db.DB.Create(&category)

	product := models.Product{
		Name:        "Test Product",
		Price:       29.99,
		CategoryID:  category.ID,
		Description: "Test Description",
	}
	db.DB.Create(&product)

	// Create inventory for the product
	inventory := models.Inventory{ProductID: product.ID, Stock: 10}
	db.DB.Create(&inventory)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create a proper request with product ID parameter
	req, _ := http.NewRequest("DELETE", "/products/"+strconv.Itoa(int(product.ID)), nil)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(product.ID))}}

	// Test the wrapper function
	handler := DeleteProductHandlerWrapper(db.DB)
	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Product deleted successfully")
}

func TestSearchProductsHandlerWrapper(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test products
	category := models.Category{Name: "Test Category"}
	db.DB.Create(&category)

	products := []models.Product{
		{Name: "Apple iPhone", Price: 999.99, CategoryID: category.ID, Description: "Latest iPhone"},
		{Name: "Samsung Galaxy", Price: 899.99, CategoryID: category.ID, Description: "Android phone"},
	}
	for _, p := range products {
		db.DB.Create(&p)
		// Create inventory for each product
		inventory := models.Inventory{ProductID: p.ID, Stock: 10}
		db.DB.Create(&inventory)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create a proper request with search query
	req, _ := http.NewRequest("GET", "/products/search?q=iPhone", nil)
	c.Request = req

	// Test the wrapper function
	handler := SearchProductsHandlerWrapper(db.DB)
	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Apple iPhone")
}

func TestEditProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create test product
	category := models.Category{Name: "Test Category"}
	db.DB.Create(&category)

	product := models.Product{
		Name:        "Original Product",
		Price:       19.99,
		CategoryID:  category.ID,
		Description: "Original Description",
	}
	db.DB.Create(&product)

	// Create inventory for the product
	inventory := models.Inventory{ProductID: product.ID, Stock: 5}
	db.DB.Create(&inventory)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create request with updated data
	updateData := map[string]interface{}{
		"name":        "Updated Product",
		"price":       29.99,
		"description": "Updated Description",
		"stock":       10,
	}
	jsonData, _ := json.Marshal(updateData)

	req, _ := http.NewRequest("PUT", "/products/"+strconv.Itoa(int(product.ID)), bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(product.ID))}}

	// Test the EditProduct function directly
	EditProduct(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Product updated successfully")
	assert.Contains(t, w.Body.String(), "Updated Product")
}

func TestAddProduct_InvalidPrice(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create a category
	category := models.Category{Name: "Test Category"}
	db.DB.Create(&category)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set invalid price (negative)
	input := middlewares.ProductInput{
		Name:        "Test Product",
		Price:       -10.0,
		Description: "Test Description",
		Stock:       10,
		CategoryID:  category.ID,
	}
	c.Set("product_input", input)

	AddProduct(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Price must be greater than 0 and less than 999999.99")
}

func TestAddProduct_InvalidDescription(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create a category
	category := models.Category{Name: "Test Category"}
	db.DB.Create(&category)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set invalid description (too long)
	longDesc := string(make([]byte, 1001)) // 1001 characters
	input := middlewares.ProductInput{
		Name:        "Test Product",
		Price:       10.0,
		Description: longDesc,
		Stock:       10,
		CategoryID:  category.ID,
	}
	c.Set("product_input", input)

	AddProduct(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Description must be less than 1000 characters")
}

func TestAddProduct_InvalidStock(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create a category
	category := models.Category{Name: "Test Category"}
	db.DB.Create(&category)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set invalid stock (too high)
	input := middlewares.ProductInput{
		Name:        "Test Product",
		Price:       10.0,
		Description: "Test Description",
		Stock:       100001,
		CategoryID:  category.ID,
	}
	c.Set("product_input", input)

	AddProduct(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Stock must be between 0 and 100000")
}

func TestAddProduct_DatabaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)

	// Create a category
	category := models.Category{Name: "Test Category"}
	db.DB.Create(&category)

	// Simulate database error by setting db.DB to nil
	originalDB := db.DB
	db.DB = nil
	defer func() { db.DB = originalDB }()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	input := middlewares.ProductInput{
		Name:        "Test Product",
		Price:       10.0,
		Description: "Test Description",
		Stock:       10,
		CategoryID:  category.ID,
	}
	c.Set("product_input", input)

	AddProduct(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Database error")
}
