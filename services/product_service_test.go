package services

import (
	"testing"

	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
)

func TestProductService_CreateProduct(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	productService := NewProductService()

	// Create test category
	category := &models.Category{
		Name: "Test Category",
	}
	testDB.Create(category)

	// Test creating a product
	product := &models.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		CategoryID:  category.ID,
	}

	err := productService.CreateProduct(product)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	if product.Name != "Test Product" {
		t.Errorf("Expected name %s, got %s", "Test Product", product.Name)
	}
	if product.Price != 99.99 {
		t.Errorf("Expected price %f, got %f", 99.99, product.Price)
	}
}

func TestProductService_GetProductByID(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	productService := NewProductService()

	// Create test product
	product := &models.Product{
		Name:        "Test Product 2",
		Description: "Test Description 2",
		Price:       149.99,
	}
	testDB.Create(product)

	// Test getting product by ID
	retrievedProduct, err := productService.GetProductByID(product.ID)
	if err != nil {
		t.Fatalf("Failed to get product by ID: %v", err)
	}

	if retrievedProduct.ID != product.ID {
		t.Errorf("Expected ID %d, got %d", product.ID, retrievedProduct.ID)
	}
	if retrievedProduct.Name != product.Name {
		t.Errorf("Expected name %s, got %s", product.Name, retrievedProduct.Name)
	}
}

func TestProductService_GetProductByID_NotFound(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	productService := NewProductService()

	// Test getting non-existent product
	_, err := productService.GetProductByID(999)
	if err == nil {
		t.Fatalf("Expected error for non-existent product")
	}
}

func TestProductService_GetAllProducts(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	productService := NewProductService()

	// Create test products
	for i := 1; i <= 3; i++ {
		product := &models.Product{
			Name:        "Test Product",
			Description: "Test Description",
			Price:       float64(i * 10),
		}
		testDB.Create(product)
	}

	// Test getting all products
	products, err := productService.GetAllProducts(1, 10)
	if err != nil {
		t.Fatalf("Failed to get all products: %v", err)
	}

	if len(products) != 3 {
		t.Errorf("Expected 3 products, got %d", len(products))
	}
}

func TestProductService_GetAllProducts_Empty(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	productService := NewProductService()

	// Test getting all products when none exist
	products, err := productService.GetAllProducts(1, 10)
	if err != nil {
		t.Fatalf("Failed to get all products: %v", err)
	}

	if len(products) != 0 {
		t.Errorf("Expected 0 products, got %d", len(products))
	}
}

func TestProductService_UpdateProduct(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	productService := NewProductService()

	// Create test product
	product := &models.Product{
		Name:        "Test Product 3",
		Description: "Test Description 3",
		Price:       199.99,
	}
	testDB.Create(product)

	// Update product
	product.Name = "Updated Product"
	product.Price = 299.99

	err := productService.UpdateProduct(product)
	if err != nil {
		t.Fatalf("Failed to update product: %v", err)
	}

	// Verify update
	updatedProduct, _ := productService.GetProductByID(product.ID)
	if updatedProduct.Name != "Updated Product" {
		t.Errorf("Expected name 'Updated Product', got %s", updatedProduct.Name)
	}
	if updatedProduct.Price != 299.99 {
		t.Errorf("Expected price 299.99, got %f", updatedProduct.Price)
	}
}

func TestProductService_UpdateProduct_NotFound(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	productService := NewProductService()

	// Try to update non-existent product
	product := &models.Product{
		Name:        "Non-existent Product",
		Description: "Test Description",
		Price:       99.99,
	}
	product.ID = 999

	err := productService.UpdateProduct(product)
	if err == nil {
		t.Fatalf("Expected error for updating non-existent product")
	}
}

func TestProductService_DeleteProduct(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	productService := NewProductService()

	// Create test product
	product := &models.Product{
		Name:        "Test Product 4",
		Description: "Test Description 4",
		Price:       399.99,
	}
	testDB.Create(product)

	// Delete product
	err := productService.DeleteProduct(product.ID)
	if err != nil {
		t.Fatalf("Failed to delete product: %v", err)
	}

	// Verify product is deleted
	_, err = productService.GetProductByID(product.ID)
	if err == nil {
		t.Error("Expected error when getting deleted product, but got none")
	}
}

func TestProductService_DeleteProduct_NotFound(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	productService := NewProductService()

	// Try to delete non-existent product
	err := productService.DeleteProduct(999)
	if err == nil {
		t.Fatalf("Expected error for deleting non-existent product")
	}
}

func TestProductService_SearchProducts(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	productService := NewProductService()

	// Create test products
	products := []*models.Product{
		{Name: "Laptop Computer", Description: "High performance laptop", Price: 999.99},
		{Name: "Gaming Mouse", Description: "RGB gaming mouse", Price: 49.99},
		{Name: "Keyboard", Description: "Mechanical keyboard", Price: 129.99},
	}

	for _, product := range products {
		testDB.Create(product)
	}

	// Test searching products
	results, err := productService.SearchProducts("laptop", 1, 10)
	if err != nil {
		t.Fatalf("Failed to search products: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 search result, got %d", len(results))
	}

	if results[0].Name != "Laptop Computer" {
		t.Errorf("Expected product name 'Laptop Computer', got %s", results[0].Name)
	}
}

func TestProductService_SearchProducts_NoResults(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	productService := NewProductService()

	// Test searching with no matching products
	results, err := productService.SearchProducts("nonexistent", 1, 10)
	if err != nil {
		t.Fatalf("Failed to search products: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 search results, got %d", len(results))
	}
}

// Benchmark tests for performance
func BenchmarkProductService_CreateProduct(b *testing.B) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(b)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	productService := NewProductService()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		product := &models.Product{
			Name:        "Bench Product",
			Description: "Bench Description",
			Price:       99.99,
		}
		productService.CreateProduct(product)
		// Clean up for next iteration
		testDB.Where("name = ?", "Bench Product").Delete(&models.Product{})
	}
}

func BenchmarkProductService_GetProductByID(b *testing.B) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(b)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	productService := NewProductService()

	// Create test product
	product := &models.Product{
		Name:        "Bench Product 2",
		Description: "Bench Description 2",
		Price:       149.99,
	}
	testDB.Create(product)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		productService.GetProductByID(product.ID)
	}
}

func BenchmarkProductService_SearchProducts(b *testing.B) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(b)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	productService := NewProductService()

	// Create test products
	for i := 1; i <= 10; i++ {
		product := &models.Product{
			Name:        "Search Product",
			Description: "Searchable product description",
			Price:       float64(i * 10),
		}
		testDB.Create(product)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		productService.SearchProducts("search", 1, 10)
	}
}

func TestProductService_GetProductsByCategory(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	productService := NewProductService()

	// Create test category
	category := &models.Category{
		Name: "Test Category",
	}
	testDB.Create(category)

	// Create test products in the category
	products := []models.Product{
		{
			Name:        "Product 1",
			Description: "Description 1",
			Price:       10.99,
			CategoryID:  category.ID,
		},
		{
			Name:        "Product 2",
			Description: "Description 2",
			Price:       20.99,
			CategoryID:  category.ID,
		},
		{
			Name:        "Product 3",
			Description: "Description 3",
			Price:       30.99,
			CategoryID:  category.ID,
		},
	}

	for _, p := range products {
		testDB.Create(&p)
	}

	// Test getting products by category
	result, err := productService.GetProductsByCategory(category.ID, 10, 0)
	if err != nil {
		t.Fatalf("Failed to get products by category: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 products, got %d", len(result))
	}

	// Verify products are in the correct category
	for _, product := range result {
		if product.CategoryID != category.ID {
			t.Errorf("Product %s has wrong category ID: expected %d, got %d", product.Name, category.ID, product.CategoryID)
		}
	}
}

func TestProductService_GetProductsByCategory_Empty(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	productService := NewProductService()

	// Test getting products by non-existent category
	result, err := productService.GetProductsByCategory(999, 10, 0)
	if err != nil {
		t.Fatalf("Failed to get products by category: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected 0 products for non-existent category, got %d", len(result))
	}
}
