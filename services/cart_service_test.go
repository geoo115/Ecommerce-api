package services

import (
	"testing"

	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
	"github.com/stretchr/testify/assert"
)

func TestNewCartService(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	service := NewCartService()
	if service == nil {
		t.Fatal("Service is nil")
	}
}

func TestCartService_AddToCart_NewItem(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	// Create test data
	user := models.User{Username: "testuser", Email: "test@example.com"}
	testDB.Create(&user)

	category := models.Category{Name: "Test Category"}
	testDB.Create(&category)

	product := models.Product{Name: "Test Product", Price: 10.99, CategoryID: category.ID}
	testDB.Create(&product)

	inventory := models.Inventory{ProductID: product.ID, Stock: 10}
	testDB.Create(&inventory)

	// Create service with test DB
	service := NewCartService()

	// Test adding new item to cart
	err := service.AddToCart(user.ID, product.ID, 2)
	assert.NoError(t, err)

	// Verify cart item was created
	var cart models.Cart
	err = testDB.Where("user_id = ? AND product_id = ?", user.ID, product.ID).First(&cart).Error
	assert.NoError(t, err)
	assert.Equal(t, 2, cart.Quantity)
}

func TestCartService_AddToCart_ExistingItem(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	// Create test data
	user := models.User{Username: "testuser", Email: "test@example.com"}
	testDB.Create(&user)

	category := models.Category{Name: "Test Category"}
	testDB.Create(&category)

	product := models.Product{Name: "Test Product", Price: 10.99, CategoryID: category.ID}
	testDB.Create(&product)

	inventory := models.Inventory{ProductID: product.ID, Stock: 10}
	testDB.Create(&inventory)

	// Create service with test DB
	service := NewCartService()

	// Add initial item
	err := service.AddToCart(user.ID, product.ID, 2)
	assert.NoError(t, err)

	// Add more of the same item
	err = service.AddToCart(user.ID, product.ID, 3)
	assert.NoError(t, err)

	// Verify quantity was updated
	var cart models.Cart
	err = testDB.Where("user_id = ? AND product_id = ?", user.ID, product.ID).First(&cart).Error
	assert.NoError(t, err)
	assert.Equal(t, 5, cart.Quantity)
}

func TestCartService_AddToCart_InsufficientStock(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	// Create test data
	user := models.User{Username: "testuser", Email: "test@example.com"}
	testDB.Create(&user)

	category := models.Category{Name: "Test Category"}
	testDB.Create(&category)

	product := models.Product{Name: "Test Product", Price: 10.99, CategoryID: category.ID}
	testDB.Create(&product)

	inventory := models.Inventory{ProductID: product.ID, Stock: 5}
	testDB.Create(&inventory)

	// Create service with test DB
	service := NewCartService()

	// Try to add more than available stock
	err := service.AddToCart(user.ID, product.ID, 10)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient stock")
}

func TestCartService_GetUserCart(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	// Create test data
	user := models.User{Username: "testuser", Email: "test@example.com"}
	testDB.Create(&user)

	category := models.Category{Name: "Test Category"}
	testDB.Create(&category)

	product1 := models.Product{Name: "Product 1", Price: 10.99, CategoryID: category.ID}
	product2 := models.Product{Name: "Product 2", Price: 20.99, CategoryID: category.ID}
	testDB.Create(&product1)
	testDB.Create(&product2)

	// Create cart items
	cart1 := models.Cart{UserID: user.ID, ProductID: product1.ID, Quantity: 2}
	cart2 := models.Cart{UserID: user.ID, ProductID: product2.ID, Quantity: 1}
	testDB.Create(&cart1)
	testDB.Create(&cart2)

	// Create service with test DB
	service := NewCartService()

	// Get user cart
	cartItems, err := service.GetUserCart(user.ID)
	assert.NoError(t, err)
	assert.Len(t, cartItems, 2)

	// Verify cart items
	assert.Equal(t, user.ID, cartItems[0].UserID)
	assert.Equal(t, product1.ID, cartItems[0].ProductID)
	assert.Equal(t, 2, cartItems[0].Quantity)
}

func TestCartService_GetUserCart_Empty(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	// Create test user
	user := models.User{Username: "testuser", Email: "test@example.com"}
	testDB.Create(&user)

	// Create service with test DB
	service := NewCartService()

	// Test getting empty cart
	cartItems, err := service.GetUserCart(user.ID)
	assert.NoError(t, err)
	assert.Empty(t, cartItems)
}

func TestCartService_UpdateCartItem(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	// Create test data
	user := models.User{Username: "testuser", Email: "test@example.com"}
	testDB.Create(&user)

	category := models.Category{Name: "Test Category"}
	testDB.Create(&category)

	product := models.Product{Name: "Test Product", Price: 10.99, CategoryID: category.ID}
	testDB.Create(&product)

	inventory := models.Inventory{ProductID: product.ID, Stock: 10}
	testDB.Create(&inventory)

	// Create initial cart item
	cart := models.Cart{UserID: user.ID, ProductID: product.ID, Quantity: 2}
	testDB.Create(&cart)

	// Create service with test DB
	service := NewCartService()

	// Update cart item quantity
	err := service.UpdateCartItem(user.ID, product.ID, 5)
	assert.NoError(t, err)

	// Verify quantity was updated
	var updatedCart models.Cart
	err = testDB.Where("user_id = ? AND product_id = ?", user.ID, product.ID).First(&updatedCart).Error
	assert.NoError(t, err)
	assert.Equal(t, 5, updatedCart.Quantity)
}

func TestCartService_RemoveFromCart(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	// Create test data
	user := models.User{Username: "testuser", Email: "test@example.com"}
	testDB.Create(&user)

	category := models.Category{Name: "Test Category"}
	testDB.Create(&category)

	product := models.Product{Name: "Test Product", Price: 10.99, CategoryID: category.ID}
	testDB.Create(&product)

	// Create cart item
	cart := models.Cart{UserID: user.ID, ProductID: product.ID, Quantity: 2}
	testDB.Create(&cart)

	// Create service with test DB
	service := NewCartService()

	// Remove item from cart
	err := service.RemoveFromCart(user.ID, product.ID)
	assert.NoError(t, err)

	// Verify item was removed
	var count int64
	testDB.Model(&models.Cart{}).Where("user_id = ? AND product_id = ?", user.ID, product.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestCartService_CheckStock(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	// Create test data
	category := models.Category{Name: "Test Category"}
	testDB.Create(&category)

	product := models.Product{Name: "Test Product", Price: 10.99, CategoryID: category.ID}
	testDB.Create(&product)

	inventory := models.Inventory{ProductID: product.ID, Stock: 10}
	testDB.Create(&inventory)

	// Create service with test DB
	service := NewCartService()

	// Test sufficient stock
	available, err := service.CheckStock(product.ID, 5)
	assert.NoError(t, err)
	assert.True(t, available)

	// Test insufficient stock
	available, err = service.CheckStock(product.ID, 15)
	assert.NoError(t, err)
	assert.False(t, available)
}

func TestCartService_CalculateCartTotal(t *testing.T) {
	// Temporarily set DB for testing
	originalDB := db.DB
	testDB := db.SetupTestDB(t)
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	// Create test data
	category := models.Category{Name: "Test Category"}
	testDB.Create(&category)

	product1 := models.Product{Name: "Product 1", Price: 10.99, CategoryID: category.ID}
	product2 := models.Product{Name: "Product 2", Price: 20.99, CategoryID: category.ID}
	testDB.Create(&product1)
	testDB.Create(&product2)

	// Create cart items
	cartItems := []models.Cart{
		{Product: product1, Quantity: 2},
		{Product: product2, Quantity: 1},
	}

	// Create service with test DB
	service := NewCartService()

	// Calculate total
	total := service.CalculateCartTotal(cartItems)
	expectedTotal := (10.99 * 2) + (20.99 * 1)
	assert.Equal(t, expectedTotal, total)
}
