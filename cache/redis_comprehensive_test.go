package cache

import (
	"testing"
	"time"

	"github.com/geoo115/Ecommerce/models"
	"github.com/stretchr/testify/assert"
)

func TestNewInMemoryCache(t *testing.T) {
	cache := NewInMemoryCache()
	assert.NotNil(t, cache)
	assert.NotNil(t, cache.data)
	assert.True(t, cache.ttl > 0)
}

func TestInMemoryCache_SetAndGet(t *testing.T) {
	cache := NewInMemoryCache()

	// Test basic string set/get
	err := cache.Set("test_key", "test_value")
	assert.NoError(t, err)

	var result string
	err = cache.Get("test_key", &result)
	assert.NoError(t, err)
	assert.Equal(t, "test_value", result)
}

func TestInMemoryCache_GetNonExistentKey(t *testing.T) {
	cache := NewInMemoryCache()

	var result string
	err := cache.Get("non_existent", &result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "key not found")
}

func TestInMemoryCache_SetAndGetStruct(t *testing.T) {
	cache := NewInMemoryCache()

	// Test struct set/get
	user := models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Phone:    "+15550000001",
	}

	err := cache.Set("user_1", user)
	assert.NoError(t, err)

	var retrievedUser models.User
	err = cache.Get("user_1", &retrievedUser)
	assert.NoError(t, err)
	assert.Equal(t, user.Username, retrievedUser.Username)
	assert.Equal(t, user.Email, retrievedUser.Email)
}

func TestInMemoryCache_Delete(t *testing.T) {
	cache := NewInMemoryCache()

	// Set a value
	err := cache.Set("test_key", "test_value")
	assert.NoError(t, err)

	// Delete it
	err = cache.Delete("test_key")
	assert.NoError(t, err)

	// Try to get it - should fail
	var result string
	err = cache.Get("test_key", &result)
	assert.Error(t, err)
}

func TestInMemoryCache_Exists(t *testing.T) {
	cache := NewInMemoryCache()

	// Key doesn't exist initially
	exists := cache.Exists("test_key")
	assert.False(t, exists)

	// Set the key
	err := cache.Set("test_key", "test_value")
	assert.NoError(t, err)

	// Now it should exist
	exists = cache.Exists("test_key")
	assert.True(t, exists)
}

func TestInMemoryCache_Clear(t *testing.T) {
	cache := NewInMemoryCache()

	// Set multiple values
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	// Clear all
	err := cache.Clear()
	assert.NoError(t, err)

	// Check that all keys are gone
	assert.False(t, cache.Exists("key1"))
	assert.False(t, cache.Exists("key2"))
	assert.False(t, cache.Exists("key3"))
}

func TestInMemoryCache_GetStatsExtended(t *testing.T) {
	cache := NewInMemoryCache()

	// Set some values
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	stats := cache.GetStats()
	assert.NotNil(t, stats)
	assert.Contains(t, stats, "total_keys")
	assert.Equal(t, 2, stats["total_keys"])
}

func TestInMemoryCache_UserOperations(t *testing.T) {
	cache := NewInMemoryCache()

	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Phone:    "+15550000001",
		Role:     "customer",
	}
	user.ID = 1

	// Test SetUser
	err := cache.SetUser(user)
	assert.NoError(t, err)

	// Test GetUser
	retrievedUser, err := cache.GetUser(1)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedUser)
	assert.Equal(t, user.Username, retrievedUser.Username)
	assert.Equal(t, user.Email, retrievedUser.Email)
}

func TestInMemoryCache_GetUserNotFound(t *testing.T) {
	cache := NewInMemoryCache()

	// Try to get non-existent user
	user, err := cache.GetUser(999)
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestInMemoryCache_ProductOperations(t *testing.T) {
	cache := NewInMemoryCache()

	product := &models.Product{
		Name:        "Test Product",
		Price:       99.99,
		Description: "Test description",
		CategoryID:  1,
	}
	product.ID = 1

	// Test SetProduct
	err := cache.SetProduct(product)
	assert.NoError(t, err)

	// Test GetProduct
	retrievedProduct, err := cache.GetProduct(1)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedProduct)
	assert.Equal(t, product.Name, retrievedProduct.Name)
	assert.Equal(t, product.Price, retrievedProduct.Price)
}

func TestInMemoryCache_GetProductNotFound(t *testing.T) {
	cache := NewInMemoryCache()

	// Try to get non-existent product
	product, err := cache.GetProduct(999)
	assert.Error(t, err)
	assert.Nil(t, product)
}

func TestInMemoryCache_ProductsByCategoryOperations(t *testing.T) {
	cache := NewInMemoryCache()

	products := []models.Product{
		{Name: "Product 1", Price: 99.99, CategoryID: 1},
		{Name: "Product 2", Price: 199.99, CategoryID: 1},
	}
	products[0].ID = 1
	products[1].ID = 2

	// Test SetProductsByCategory
	err := cache.SetProductsByCategory(1, 1, 10, products)
	assert.NoError(t, err)

	// Test GetProductsByCategory
	retrievedProducts, err := cache.GetProductsByCategory(1, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, retrievedProducts, 2)
	assert.Equal(t, products[0].Name, retrievedProducts[0].Name)
	assert.Equal(t, products[1].Name, retrievedProducts[1].Name)
}

func TestInMemoryCache_GetProductsByCategoryNotFound(t *testing.T) {
	cache := NewInMemoryCache()

	// Try to get non-existent category products
	products, err := cache.GetProductsByCategory(999, 1, 10)
	assert.Error(t, err)
	assert.Nil(t, products)
}

func TestInMemoryCache_CartOperations(t *testing.T) {
	cache := NewInMemoryCache()

	cart := []models.Cart{
		{UserID: 1, ProductID: 1, Quantity: 2},
		{UserID: 1, ProductID: 2, Quantity: 1},
	}
	cart[0].ID = 1
	cart[1].ID = 2

	// Test SetCart
	err := cache.SetCart(1, cart)
	assert.NoError(t, err)

	// Test GetCart
	retrievedCart, err := cache.GetCart(1)
	assert.NoError(t, err)
	assert.Len(t, retrievedCart, 2)
	assert.Equal(t, cart[0].Quantity, retrievedCart[0].Quantity)
	assert.Equal(t, cart[1].Quantity, retrievedCart[1].Quantity)
}

func TestInMemoryCache_GetCartNotFound(t *testing.T) {
	cache := NewInMemoryCache()

	// Try to get non-existent cart
	cart, err := cache.GetCart(999)
	assert.Error(t, err)
	assert.Nil(t, cart)
}

func TestInMemoryCache_SetLargeData(t *testing.T) {
	cache := NewInMemoryCache()

	// Create large data
	largeData := make([]byte, 1024*1024) // 1MB
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	// Test setting large data
	err := cache.Set("large_key", largeData)
	assert.NoError(t, err)

	// Test getting large data
	var retrieved []byte
	err = cache.Get("large_key", &retrieved)
	assert.NoError(t, err)
	assert.Equal(t, largeData, retrieved)
}

func TestInMemoryCache_InvalidateUserCache(t *testing.T) {
	cache := NewInMemoryCache()

	// Set user data
	user := &models.User{Username: "testuser"}
	user.ID = 1
	cache.SetUser(user)

	// Set cart data for same user
	cart := []models.Cart{{UserID: 1, ProductID: 1, Quantity: 2}}
	cache.SetCart(1, cart)

	// Invalidate user cache
	err := cache.InvalidateUserCache(1)
	assert.NoError(t, err)

	// User and cart should be gone
	retrievedUser, err := cache.GetUser(1)
	assert.Error(t, err)
	assert.Nil(t, retrievedUser)

	retrievedCart, err := cache.GetCart(1)
	assert.Error(t, err)
	assert.Nil(t, retrievedCart)
}

func TestInMemoryCache_InvalidateProductCache(t *testing.T) {
	cache := NewInMemoryCache()

	// Set product data
	product := &models.Product{Name: "Test Product", CategoryID: 1}
	product.ID = 1
	cache.SetProduct(product)

	// Set category products
	products := []models.Product{*product}
	cache.SetProductsByCategory(1, 1, 10, products)

	// Invalidate product cache
	err := cache.InvalidateProductCache(1)
	assert.NoError(t, err)

	// Product should be gone
	retrievedProduct, err := cache.GetProduct(1)
	assert.Error(t, err)
	assert.Nil(t, retrievedProduct)

	// Category products should also be invalidated
	retrievedProducts, err := cache.GetProductsByCategory(1, 1, 10)
	assert.Error(t, err)
	assert.Nil(t, retrievedProducts)
}

func TestInMemoryCache_TTLExpiration(t *testing.T) {
	// Create cache with very short TTL for testing
	cache := &InMemoryCache{
		data: make(map[string]cacheItem),
		ttl:  10 * time.Millisecond,
	}

	// Set a value
	err := cache.Set("test_key", "test_value")
	assert.NoError(t, err)

	// Should exist immediately
	assert.True(t, cache.Exists("test_key"))

	// Wait for expiration
	time.Sleep(20 * time.Millisecond)

	// Should be expired and not exist
	var result string
	err = cache.Get("test_key", &result)
	assert.Error(t, err)
}

func TestInMemoryCache_ConcurrentAccess(t *testing.T) {
	cache := NewInMemoryCache()

	// Test concurrent writes and reads
	go func() {
		for i := 0; i < 100; i++ {
			cache.Set("key1", "value1")
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			cache.Set("key2", "value2")
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			var result string
			cache.Get("key1", &result)
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			var result string
			cache.Get("key2", &result)
		}
	}()

	// Give goroutines time to complete
	time.Sleep(100 * time.Millisecond)

	// Cache should still be functional
	cache.Set("final_test", "final_value")
	var result string
	err := cache.Get("final_test", &result)
	assert.NoError(t, err)
	assert.Equal(t, "final_value", result)
}

func TestRedisCache_Comprehensive_SetAndGet(t *testing.T) {
	cache := GetCache()
	if cache == nil {
		t.Skip("Cache not initialized")
	}

	// Test basic string set/get
	err := cache.Set("redis_test_key", "redis_test_value")
	assert.NoError(t, err)

	var result string
	err = cache.Get("redis_test_key", &result)
	assert.NoError(t, err)
	assert.Equal(t, "redis_test_value", result)
}

func TestRedisCache_Comprehensive_GetNonExistentKey(t *testing.T) {
	cache := GetCache()
	if cache == nil {
		t.Skip("Cache not initialized")
	}

	var result string
	err := cache.Get("non_existent_redis_key", &result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "key not found")
}

func TestRedisCache_Comprehensive_Delete(t *testing.T) {
	cache := GetCache()
	if cache == nil {
		t.Skip("Cache not initialized")
	}

	// Set a value
	err := cache.Set("redis_delete_key", "redis_delete_value")
	assert.NoError(t, err)

	// Delete it
	err = cache.Delete("redis_delete_key")
	assert.NoError(t, err)

	// Try to get it - should fail
	var result string
	err = cache.Get("redis_delete_key", &result)
	assert.Error(t, err)
}

func TestRedisCache_Comprehensive_Exists(t *testing.T) {
	cache := GetCache()
	if cache == nil {
		t.Skip("Cache not initialized")
	}

	// Key doesn't exist initially
	exists := cache.Exists("redis_exists_key")
	assert.False(t, exists)

	// Set the key
	err := cache.Set("redis_exists_key", "redis_exists_value")
	assert.NoError(t, err)

	// Now it should exist
	exists = cache.Exists("redis_exists_key")
	assert.True(t, exists)
}

func TestRedisCache_Comprehensive_Clear(t *testing.T) {
	cache := GetCache()
	if cache == nil {
		t.Skip("Cache not initialized")
	}

	// Set multiple values
	cache.Set("redis_clear_key1", "value1")
	cache.Set("redis_clear_key2", "value2")
	cache.Set("redis_clear_key3", "value3")

	// Clear all
	err := cache.Clear()
	assert.NoError(t, err)

	// Check that all keys are gone
	assert.False(t, cache.Exists("redis_clear_key1"))
	assert.False(t, cache.Exists("redis_clear_key2"))
	assert.False(t, cache.Exists("redis_clear_key3"))
}

func TestRedisCache_Comprehensive_UserOperations(t *testing.T) {
	cache := GetCache()
	if cache == nil {
		t.Skip("Cache not initialized")
	}

	user := &models.User{
		Username: "redis_testuser",
		Email:    "redis_test@example.com",
		Phone:    "+15550000002",
		Role:     "customer",
	}
	user.ID = 100

	// Test SetUser
	err := cache.SetUser(user)
	assert.NoError(t, err)

	// Test GetUser
	retrievedUser, err := cache.GetUser(100)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedUser)
	assert.Equal(t, user.Username, retrievedUser.Username)
	assert.Equal(t, user.Email, retrievedUser.Email)
}

func TestRedisCache_Comprehensive_ProductOperations(t *testing.T) {
	cache := GetCache()
	if cache == nil {
		t.Skip("Cache not initialized")
	}

	product := &models.Product{
		Name:        "Redis Test Product",
		Price:       199.99,
		Description: "Redis test description",
		CategoryID:  2,
	}
	product.ID = 200

	// Test SetProduct
	err := cache.SetProduct(product)
	assert.NoError(t, err)

	// Test GetProduct
	retrievedProduct, err := cache.GetProduct(200)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedProduct)
	assert.Equal(t, product.Name, retrievedProduct.Name)
	assert.Equal(t, product.Price, retrievedProduct.Price)
}

func TestRedisCache_Comprehensive_CartOperations(t *testing.T) {
	cache := GetCache()
	if cache == nil {
		t.Skip("Cache not initialized")
	}

	cart := []models.Cart{
		{UserID: 100, ProductID: 200, Quantity: 3},
		{UserID: 100, ProductID: 201, Quantity: 1},
	}
	cart[0].ID = 300
	cart[1].ID = 301

	// Test SetCart
	err := cache.SetCart(100, cart)
	assert.NoError(t, err)

	// Test GetCart
	retrievedCart, err := cache.GetCart(100)
	assert.NoError(t, err)
	assert.Len(t, retrievedCart, 2)
	assert.Equal(t, cart[0].Quantity, retrievedCart[0].Quantity)
	assert.Equal(t, cart[1].Quantity, retrievedCart[1].Quantity)
}

func TestRedisCache_Comprehensive_InvalidateUserCache(t *testing.T) {
	cache := GetCache()
	if cache == nil {
		t.Skip("Cache not initialized")
	}

	// Set user data
	user := &models.User{Username: "redis_invalidate_user"}
	user.ID = 400
	cache.SetUser(user)

	// Set cart data for same user
	cart := []models.Cart{{UserID: 400, ProductID: 200, Quantity: 2}}
	cache.SetCart(400, cart)

	// Invalidate user cache
	err := cache.InvalidateUserCache(400)
	assert.NoError(t, err)

	// User and cart should be gone
	retrievedUser, err := cache.GetUser(400)
	assert.Error(t, err)
	assert.Nil(t, retrievedUser)

	retrievedCart, err := cache.GetCart(400)
	assert.Error(t, err)
	assert.Nil(t, retrievedCart)
}
