package cache

import (
	"os"
	"testing"
	"time"

	"github.com/geoo115/Ecommerce/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestInMemoryCache_SetGetExistsDeleteClear(t *testing.T) {
	// ensure a reasonably long TTL for basic operations
	os.Setenv("CACHE_TTL", "10m")
	c := NewInMemoryCache()

	if err := c.Set("k1", "v1"); err != nil {
		t.Fatalf("Set error: %v", err)
	}

	if !c.Exists("k1") {
		t.Fatalf("expected key to exist")
	}

	var got string
	if err := c.Get("k1", &got); err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if got != "v1" {
		t.Fatalf("expected v1, got %v", got)
	}

	if err := c.Delete("k1"); err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if c.Exists("k1") {
		t.Fatalf("expected key to be deleted")
	}

	// test Clear
	if err := c.Set("a", "b"); err != nil {
		t.Fatalf("Set error: %v", err)
	}
	if err := c.Clear(); err != nil {
		t.Fatalf("Clear error: %v", err)
	}
	if c.Exists("a") {
		t.Fatalf("expected cache to be cleared")
	}
}

func TestInitCache(t *testing.T) {
	// Test InitCache (may skip if Redis not available)
	err := InitCache()
	if err != nil {
		t.Skip("Redis not available:", err)
	}
	// If no error, cache is initialized
}

func TestRedisSetGet(t *testing.T) {
	cache := GetCache()
	if cache == nil {
		t.Skip("Cache not initialized")
	}

	err := cache.Set("redis_key", "redis_value")
	if err != nil {
		t.Skip("Redis set failed:", err)
	}

	var value string
	err = cache.Get("redis_key", &value)
	assert.NoError(t, err)
	assert.Equal(t, "redis_value", value)

	err = cache.Delete("redis_key")
	assert.NoError(t, err)

	exists := cache.Exists("redis_key")
	assert.False(t, exists)
}

func TestRedisClear(t *testing.T) {
	cache := GetCache()
	if cache == nil {
		t.Skip("Cache not initialized")
	}

	err := cache.Set("clear_key", "value")
	if err != nil {
		t.Skip("Redis set failed:", err)
	}

	err = cache.Clear()
	if err != nil {
		t.Skip("Redis clear failed:", err)
	}

	exists := cache.Exists("clear_key")
	assert.False(t, exists)
}

func TestInMemoryCache_Expiration(t *testing.T) {
	// short TTL to exercise expiration
	os.Setenv("CACHE_TTL", "200ms")
	c := NewInMemoryCache()

	if err := c.Set("exp", "x"); err != nil {
		t.Fatalf("Set error: %v", err)
	}

	// wait for expiration
	time.Sleep(350 * time.Millisecond)

	var dest string
	if err := c.Get("exp", &dest); err == nil {
		t.Fatalf("expected error for expired key, got value: %v", dest)
	}
	if c.Exists("exp") {
		t.Fatalf("expected expired key to not exist")
	}
}

func TestInMemoryCache_InvalidateUserProductCache(t *testing.T) {
	os.Setenv("CACHE_TTL", "10m")
	c := NewInMemoryCache()

	// set keys that match patterns used by the invalidation helpers
	if err := c.Set("user:1", "u"); err != nil {
		t.Fatalf("Set error: %v", err)
	}
	if err := c.Set("cart:1", "c"); err != nil {
		t.Fatalf("Set error: %v", err)
	}
	if err := c.Set("product:2", "p"); err != nil {
		t.Fatalf("Set error: %v", err)
	}
	if err := c.Set("products:category:3:page:1:limit:10", "list"); err != nil {
		t.Fatalf("Set error: %v", err)
	}

	// invalidate user cache
	if err := c.InvalidateUserCache(1); err != nil {
		t.Fatalf("InvalidateUserCache error: %v", err)
	}
	if c.Exists("user:1") || c.Exists("cart:1") {
		t.Fatalf("expected user-related keys to be removed")
	}

	// invalidate product cache
	if err := c.InvalidateProductCache(2); err != nil {
		t.Fatalf("InvalidateProductCache error: %v", err)
	}
	if c.Exists("product:2") || c.Exists("products:category:3:page:1:limit:10") {
		t.Fatalf("expected product-related keys to be removed")
	}
}

func TestInMemoryCache_GetStats(t *testing.T) {
	os.Setenv("CACHE_TTL", "10m")
	c := NewInMemoryCache()

	if err := c.Set("s1", "v1"); err != nil {
		t.Fatalf("Set error: %v", err)
	}

	stats := c.GetStats()
	if stats == nil {
		t.Fatalf("expected stats map, got nil")
	}
	if _, ok := stats["total_items"]; !ok {
		t.Fatalf("expected total_items in stats")
	}
}

func TestInMemoryCache_ModelHelpers(t *testing.T) {
	os.Setenv("CACHE_TTL", "10m")
	c := NewInMemoryCache()

	// Test User
	user := &models.User{Model: gorm.Model{ID: 1}, Username: "testuser"}
	if err := c.SetUser(user); err != nil {
		t.Fatalf("SetUser error: %v", err)
	}
	retrievedUser, err := c.GetUser(1)
	if err != nil {
		t.Fatalf("GetUser error: %v", err)
	}
	if retrievedUser.Username != "testuser" {
		t.Fatalf("expected username testuser, got %v", retrievedUser.Username)
	}

	// Test Product
	product := &models.Product{Model: gorm.Model{ID: 2}, Name: "testproduct", Price: 10.0}
	if err := c.SetProduct(product); err != nil {
		t.Fatalf("SetProduct error: %v", err)
	}
	retrievedProduct, err := c.GetProduct(2)
	if err != nil {
		t.Fatalf("GetProduct error: %v", err)
	}
	if retrievedProduct.Name != "testproduct" {
		t.Fatalf("expected name testproduct, got %v", retrievedProduct.Name)
	}

	// Test Cart
	cart := []models.Cart{{Model: gorm.Model{ID: 3}, UserID: 1, ProductID: 2, Quantity: 5}}
	if err := c.SetCart(1, cart); err != nil {
		t.Fatalf("SetCart error: %v", err)
	}
	retrievedCart, err := c.GetCart(1)
	if err != nil {
		t.Fatalf("GetCart error: %v", err)
	}
	if len(retrievedCart) != 1 || retrievedCart[0].Quantity != 5 {
		t.Fatalf("expected cart with quantity 5, got %v", retrievedCart)
	}

	// Test ProductsByCategory
	products := []models.Product{*product}
	if err := c.SetProductsByCategory(1, 1, 10, products); err != nil {
		t.Fatalf("SetProductsByCategory error: %v", err)
	}
	retrievedProducts, err := c.GetProductsByCategory(1, 1, 10)
	if err != nil {
		t.Fatalf("GetProductsByCategory error: %v", err)
	}
	if len(retrievedProducts) != 1 || retrievedProducts[0].Name != "testproduct" {
		t.Fatalf("expected products list, got %v", retrievedProducts)
	}
}

func TestNewRedisCache_Success(t *testing.T) {
	// This test requires a running Redis instance.
	// If Redis is not running, this test will fail.
	// You can run Redis using Docker: docker run --name some-redis -p 6379:6379 -d redis

	os.Setenv("REDIS_ADDR", "localhost:6379")
	os.Setenv("REDIS_PASSWORD", "")
	os.Setenv("REDIS_DB", "0")
	os.Setenv("CACHE_TTL", "1h")
	defer func() {
		os.Unsetenv("REDIS_ADDR")
		os.Unsetenv("REDIS_PASSWORD")
		os.Unsetenv("REDIS_DB")
		os.Unsetenv("CACHE_TTL")
	}()

	// Temporarily store the original newRedisCacheCreator
	oldNewRedisCacheCreator := newRedisCacheCreator
	defer func() {
		newRedisCacheCreator = oldNewRedisCacheCreator
	}()

	cache := newRedisCacheCreator()
	if cache == nil {
		t.Skip("Redis is not running or connection failed, skipping RedisCache tests.")
	}
	assert.NotNil(t, cache)
	assert.Equal(t, 1*time.Hour, cache.ttl)

	// Clean up Redis after test
	cache.Clear()
}

func TestNewRedisCache_Fallback(t *testing.T) {
	// This test simulates a Redis connection failure, ensuring fallback to in-memory.
	// Set an invalid Redis address to force connection error.
	os.Setenv("REDIS_ADDR", "localhost:9999") // Invalid port
	defer os.Unsetenv("REDIS_ADDR")

	// Temporarily store the original newRedisCacheCreator
	oldNewRedisCacheCreator := newRedisCacheCreator
	defer func() {
		newRedisCacheCreator = oldNewRedisCacheCreator
	}()

	cache := newRedisCacheCreator()
	assert.Nil(t, cache, "NewRedisCache should return nil on connection failure")
}

func TestRedisCache_SetAndGet(t *testing.T) {
	cache := newRedisCacheCreator()
	if cache == nil {
		t.Skip("Redis is not running or connection failed, skipping RedisCache tests.")
	}
	defer cache.Clear()

	expectedValue := "test_value"
	err := cache.Set("key1", expectedValue)
	assert.NoError(t, err)

	var actualValue string
	err = cache.Get("key1", &actualValue)
	assert.NoError(t, err)
	assert.Equal(t, expectedValue, actualValue)

	// Test non-existent key
	err = cache.Get("nonexistent", &actualValue)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "redis: nil")

	// Test expired key (requires setting a short TTL for the cache instance)
	cache.ttl = 1 * time.Millisecond
	cache.Set("key_expired", "expired_value")
	time.Sleep(2 * time.Millisecond)
	err = cache.Get("key_expired", &actualValue)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "redis: nil")
}

func TestRedisCache_Delete(t *testing.T) {
	cache := newRedisCacheCreator()
	if cache == nil {
		t.Skip("Redis is not running or connection failed, skipping RedisCache tests.")
	}
	defer cache.Clear()

	cache.Set("key_to_delete", "value")
	assert.True(t, cache.Exists("key_to_delete"))

	err := cache.Delete("key_to_delete")
	assert.NoError(t, err)
	assert.False(t, cache.Exists("key_to_delete"))

	// Deleting non-existent key should not return error
	err = cache.Delete("nonexistent")
	assert.NoError(t, err)
}

func TestRedisCache_Exists(t *testing.T) {
	cache := newRedisCacheCreator()
	if cache == nil {
		t.Skip("Redis is not running or connection failed, skipping RedisCache tests.")
	}
	defer cache.Clear()

	cache.Set("exists_key", "value")
	assert.True(t, cache.Exists("exists_key"))
	assert.False(t, cache.Exists("nonexistent_key"))

	// Test expired key
	cache.ttl = 1 * time.Millisecond
	cache.Set("exists_expired", "value")
	time.Sleep(2 * time.Millisecond)
	assert.False(t, cache.Exists("exists_expired"))
}

func TestRedisCache_Clear(t *testing.T) {
	cache := newRedisCacheCreator()
	if cache == nil {
		t.Skip("Redis is not running or connection failed, skipping RedisCache tests.")
	}
	defer cache.Clear()

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	assert.True(t, cache.Exists("key1"))
	assert.True(t, cache.Exists("key2"))

	err := cache.Clear()
	assert.NoError(t, err)

	assert.False(t, cache.Exists("key1"))
	assert.False(t, cache.Exists("key2"))
}

func TestRedisCache_GetStats(t *testing.T) {
	cache := newRedisCacheCreator()
	if cache == nil {
		t.Skip("Redis is not running or connection failed, skipping RedisCache tests.")
	}
	defer cache.Clear()

	stats := cache.GetStats()
	assert.NotNil(t, stats)
	assert.Contains(t, stats, "redis_info")
	assert.Contains(t, stats, "ttl")
}

func TestRedisCache_UserMethods(t *testing.T) {
	cache := newRedisCacheCreator()
	if cache == nil {
		t.Skip("Redis is not running or connection failed, skipping RedisCache tests.")
	}
	defer cache.Clear()

	user := &models.User{Model: gorm.Model{ID: 1}, Username: "testuser"}
	err := cache.SetUser(user)
	assert.NoError(t, err)

	retrievedUser, err := cache.GetUser(1)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, retrievedUser.ID)
	assert.Equal(t, user.Username, retrievedUser.Username)

	// Test non-existent user
	_, err = cache.GetUser(999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "redis: nil")
}

func TestRedisCache_ProductMethods(t *testing.T) {
	cache := newRedisCacheCreator()
	if cache == nil {
		t.Skip("Redis is not running or connection failed, skipping RedisCache tests.")
	}
	defer cache.Clear()

	product := &models.Product{Model: gorm.Model{ID: 1}, Name: "testproduct"}
	err := cache.SetProduct(product)
	assert.NoError(t, err)

	retrievedProduct, err := cache.GetProduct(1)
	assert.NoError(t, err)
	assert.Equal(t, product.ID, retrievedProduct.ID)
	assert.Equal(t, product.Name, retrievedProduct.Name)

	// Test non-existent product
	_, err = cache.GetProduct(999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "redis: nil")
}

func TestRedisCache_ProductsByCategoryMethods(t *testing.T) {
	cache := newRedisCacheCreator()
	if cache == nil {
		t.Skip("Redis is not running or connection failed, skipping RedisCache tests.")
	}
	defer cache.Clear()

	products := []models.Product{{Model: gorm.Model{ID: 1}, Name: "p1"}, {Model: gorm.Model{ID: 2}, Name: "p2"}}
	err := cache.SetProductsByCategory(1, 1, 10, products)
	assert.NoError(t, err)

	retrievedProducts, err := cache.GetProductsByCategory(1, 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, len(products), len(retrievedProducts))
	assert.Equal(t, products[0].ID, retrievedProducts[0].ID)

	// Test non-existent category products
	_, err = cache.GetProductsByCategory(999, 1, 10)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "redis: nil")
}

func TestRedisCache_CartMethods(t *testing.T) {
	cache := newRedisCacheCreator()
	if cache == nil {
		t.Skip("Redis is not running or connection failed, skipping RedisCache tests.")
	}
	defer cache.Clear()

	cart := []models.Cart{{Model: gorm.Model{ID: 1}, ProductID: 1, Quantity: 1}}
	err := cache.SetCart(1, cart)
	assert.NoError(t, err)

	retrievedCart, err := cache.GetCart(1)
	assert.NoError(t, err)
	assert.Equal(t, len(cart), len(retrievedCart))
	assert.Equal(t, cart[0].ID, retrievedCart[0].ID)

	// Test non-existent cart
	_, err = cache.GetCart(999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "redis: nil")
}

func TestRedisCache_InvalidateUserCache(t *testing.T) {
	cache := newRedisCacheCreator()
	if cache == nil {
		t.Skip("Redis is not running or connection failed, skipping RedisCache tests.")
	}
	defer cache.Clear()

	cache.SetUser(&models.User{Model: gorm.Model{ID: 1}, Username: "user1"})
	cache.SetCart(1, []models.Cart{{Model: gorm.Model{ID: 1}}})
	cache.Set("other:key", "value")

	assert.True(t, cache.Exists("user:1"))
	assert.True(t, cache.Exists("cart:1"))
	assert.True(t, cache.Exists("other:key"))

	err := cache.InvalidateUserCache(1)
	assert.NoError(t, err)

	assert.False(t, cache.Exists("user:1"))
	assert.False(t, cache.Exists("cart:1"))
	assert.True(t, cache.Exists("other:key")) // Should not invalidate other keys
}

func TestRedisCache_InvalidateProductCache(t *testing.T) {
	cache := newRedisCacheCreator()
	if cache == nil {
		t.Skip("Redis is not running or connection failed, skipping RedisCache tests.")
	}
	defer cache.Clear()

	cache.SetProduct(&models.Product{Model: gorm.Model{ID: 1}, Name: "product1"})
	cache.SetProductsByCategory(1, 1, 10, []models.Product{{Model: gorm.Model{ID: 1}}})
	cache.Set("other:key", "value")

	assert.True(t, cache.Exists("product:1"))
	assert.True(t, cache.Exists("products:category:1:page:1:limit:10"))
	assert.True(t, cache.Exists("other:key"))

	err := cache.InvalidateProductCache(1)
	assert.NoError(t, err)

	assert.False(t, cache.Exists("product:1"))
	assert.False(t, cache.Exists("products:category:1:page:1:limit:10"))
	assert.True(t, cache.Exists("other:key")) // Should not invalidate other keys
}
