package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/geoo115/Ecommerce/models"
	"github.com/geoo115/Ecommerce/utils"
	"github.com/go-redis/redis/v8"
)

// InMemoryCache provides simple in-memory caching
type InMemoryCache struct {
	data  map[string]cacheItem
	mutex sync.RWMutex
	ttl   time.Duration
}

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

// NewInMemoryCache creates a new in-memory cache
func NewInMemoryCache() *InMemoryCache {
	ttl := 30 * time.Minute // Default TTL
	if ttlStr := os.Getenv("CACHE_TTL"); ttlStr != "" {
		if parsed, err := time.ParseDuration(ttlStr); err == nil {
			ttl = parsed
		}
	}

	cache := &InMemoryCache{
		data: make(map[string]cacheItem),
		ttl:  ttl,
	}

	// Start cleanup goroutine
	go cache.cleanup()

	utils.Info("In-memory cache initialized with TTL: %v", ttl)
	return cache
}

// Set stores a value in cache
func (c *InMemoryCache) Set(key string, value interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(c.ttl),
	}

	return nil
}

// Get retrieves a value from cache
func (c *InMemoryCache) Get(key string, dest interface{}) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return fmt.Errorf("key not found")
	}

	if time.Now().After(item.expiration) {
		// Item expired, remove it
		delete(c.data, key)
		return fmt.Errorf("key expired")
	}

	// Convert to JSON and back to handle type conversion
	data, err := json.Marshal(item.value)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}

// Delete removes a key from cache
func (c *InMemoryCache) Delete(key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.data, key)
	return nil
}

// Exists checks if a key exists and is not expired
func (c *InMemoryCache) Exists(key string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return false
	}

	if time.Now().After(item.expiration) {
		// Clean up expired item
		delete(c.data, key)
		return false
	}

	return true
}

// Clear removes all items from cache
func (c *InMemoryCache) Clear() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data = make(map[string]cacheItem)
	return nil
}

// GetStats returns cache statistics
func (c *InMemoryCache) GetStats() map[string]interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	validItems := 0
	expiredItems := 0

	for _, item := range c.data {
		if time.Now().After(item.expiration) {
			expiredItems++
		} else {
			validItems++
		}
	}

	return map[string]interface{}{
		"total_items":   len(c.data),
		"total_keys":    len(c.data), // For compatibility with tests
		"valid_items":   validItems,
		"expired_items": expiredItems,
		"ttl":           c.ttl.String(),
	}
}

// cleanup removes expired items periodically
func (c *InMemoryCache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		for key, item := range c.data {
			if time.Now().After(item.expiration) {
				delete(c.data, key)
			}
		}
		c.mutex.Unlock()
	}
}

// GetUser retrieves user from cache
func (c *InMemoryCache) GetUser(userID uint) (*models.User, error) {
	key := fmt.Sprintf("user:%d", userID)
	var user models.User

	if err := c.Get(key, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// SetUser stores user in cache
func (c *InMemoryCache) SetUser(user *models.User) error {
	key := fmt.Sprintf("user:%d", user.ID)
	return c.Set(key, user)
}

// GetProduct retrieves product from cache
func (c *InMemoryCache) GetProduct(productID uint) (*models.Product, error) {
	key := fmt.Sprintf("product:%d", productID)
	var product models.Product

	if err := c.Get(key, &product); err != nil {
		return nil, err
	}

	return &product, nil
}

// SetProduct stores product in cache
func (c *InMemoryCache) SetProduct(product *models.Product) error {
	key := fmt.Sprintf("product:%d", product.ID)
	return c.Set(key, product)
}

// GetProductsByCategory retrieves products by category from cache
func (c *InMemoryCache) GetProductsByCategory(categoryID uint, page, limit int) ([]models.Product, error) {
	key := fmt.Sprintf("products:category:%d:page:%d:limit:%d", categoryID, page, limit)
	var products []models.Product

	if err := c.Get(key, &products); err != nil {
		return nil, err
	}

	return products, nil
}

// SetProductsByCategory stores products by category in cache
func (c *InMemoryCache) SetProductsByCategory(categoryID uint, page, limit int, products []models.Product) error {
	key := fmt.Sprintf("products:category:%d:page:%d:limit:%d", categoryID, page, limit)
	return c.Set(key, products)
}

// GetCart retrieves user cart from cache
func (c *InMemoryCache) GetCart(userID uint) ([]models.Cart, error) {
	key := fmt.Sprintf("cart:%d", userID)
	var cart []models.Cart

	if err := c.Get(key, &cart); err != nil {
		return nil, err
	}

	return cart, nil
}

// SetCart stores user cart in cache
func (c *InMemoryCache) SetCart(userID uint, cart []models.Cart) error {
	key := fmt.Sprintf("cart:%d", userID)
	return c.Set(key, cart)
}

// InvalidateUserCache invalidates all user-related cache
func (c *InMemoryCache) InvalidateUserCache(userID uint) error {
	patterns := []string{
		fmt.Sprintf("user:%d", userID),
		fmt.Sprintf("cart:%d", userID),
	}

	for _, pattern := range patterns {
		// Simple pattern matching for in-memory cache
		c.mutex.Lock()
		for key := range c.data {
			if strings.Contains(key, pattern) {
				delete(c.data, key)
			}
		}
		c.mutex.Unlock()
	}

	return nil
}

// InvalidateProductCache invalidates product-related cache
func (c *InMemoryCache) InvalidateProductCache(productID uint) error {
	patterns := []string{
		fmt.Sprintf("product:%d", productID),
		"products:",
	}

	for _, pattern := range patterns {
		c.mutex.Lock()
		for key := range c.data {
			if strings.Contains(key, pattern) {
				delete(c.data, key)
			}
		}
		c.mutex.Unlock()
	}

	return nil
}

// GetStats returns Redis cache statistics
func (c *RedisCache) GetStats() map[string]interface{} {
	info, err := c.client.Info(c.ctx, "memory", "stats").Result()
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	return map[string]interface{}{
		"redis_info": info,
		"ttl":        c.ttl.String(),
	}
}

// InvalidateProductCache invalidates product-related Redis cache
func (c *RedisCache) InvalidateProductCache(productID uint) error {
	patterns := []string{
		fmt.Sprintf("product:%d", productID),
		"products:",
	}

	for _, pattern := range patterns {
		keys, err := c.client.Keys(c.ctx, "*"+pattern+"*").Result()
		if err != nil {
			continue
		}
		if len(keys) > 0 {
			c.client.Del(c.ctx, keys...)
		}
	}

	return nil
}

// RedisCache provides Redis-based caching with fallback to in-memory
type RedisCache struct {
	client *redis.Client
	ctx    context.Context
	ttl    time.Duration
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache() *RedisCache {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := 0
	if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
		if parsed, err := strconv.Atoi(dbStr); err == nil {
			redisDB = parsed
		}
	}

	ttl := 30 * time.Minute // Default TTL
	if ttlStr := os.Getenv("CACHE_TTL"); ttlStr != "" {
		if parsed, err := time.ParseDuration(ttlStr); err == nil {
			ttl = parsed
		}
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:         redisAddr,
		Password:     redisPassword,
		DB:           redisDB,
		PoolSize:     15, // Connection pool size for high concurrency
		MinIdleConns: 5,  // Minimum idle connections
		MaxRetries:   3,  // Maximum retry attempts
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	cache := &RedisCache{
		client: rdb,
		ctx:    context.Background(),
		ttl:    ttl,
	}

	// Test connection
	if err := cache.client.Ping(cache.ctx).Err(); err != nil {
		utils.Warn("Redis connection failed, falling back to in-memory cache: %v", err)
		return nil // Will fallback to in-memory
	}

	utils.Info("Redis cache initialized with TTL: %v", ttl)
	return cache
}

// Set stores a value in Redis cache
func (c *RedisCache) Set(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(c.ctx, key, data, c.ttl).Err()
}

// Get retrieves a value from Redis cache
func (c *RedisCache) Get(key string, dest interface{}) error {
	data, err := c.client.Get(c.ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

// Delete removes a key from Redis cache
func (c *RedisCache) Delete(key string) error {
	return c.client.Del(c.ctx, key).Err()
}

// Exists checks if a key exists in Redis cache
func (c *RedisCache) Exists(key string) bool {
	count, err := c.client.Exists(c.ctx, key).Result()
	return err == nil && count > 0
}

// Clear removes all keys from Redis cache
func (c *RedisCache) Clear() error {
	return c.client.FlushDB(c.ctx).Err()
}

// GetUser retrieves user from Redis cache
func (c *RedisCache) GetUser(userID uint) (*models.User, error) {
	key := fmt.Sprintf("user:%d", userID)
	var user models.User

	if err := c.Get(key, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// SetUser stores user in Redis cache
func (c *RedisCache) SetUser(user *models.User) error {
	key := fmt.Sprintf("user:%d", user.ID)
	return c.Set(key, user)
}

// GetProduct retrieves product from Redis cache
func (c *RedisCache) GetProduct(productID uint) (*models.Product, error) {
	key := fmt.Sprintf("product:%d", productID)
	var product models.Product

	if err := c.Get(key, &product); err != nil {
		return nil, err
	}

	return &product, nil
}

// SetProduct stores product in Redis cache
func (c *RedisCache) SetProduct(product *models.Product) error {
	key := fmt.Sprintf("product:%d", product.ID)
	return c.Set(key, product)
}

// GetProductsByCategory retrieves products by category from Redis cache
func (c *RedisCache) GetProductsByCategory(categoryID uint, page, limit int) ([]models.Product, error) {
	key := fmt.Sprintf("products:category:%d:page:%d:limit:%d", categoryID, page, limit)
	var products []models.Product

	if err := c.Get(key, &products); err != nil {
		return nil, err
	}

	return products, nil
}

// SetProductsByCategory stores products by category in Redis cache
func (c *RedisCache) SetProductsByCategory(categoryID uint, page, limit int, products []models.Product) error {
	key := fmt.Sprintf("products:category:%d:page:%d:limit:%d", categoryID, page, limit)
	return c.Set(key, products)
}

// GetCart retrieves user cart from Redis cache
func (c *RedisCache) GetCart(userID uint) ([]models.Cart, error) {
	key := fmt.Sprintf("cart:%d", userID)
	var cart []models.Cart

	if err := c.Get(key, &cart); err != nil {
		return nil, err
	}

	return cart, nil
}

// SetCart stores user cart in Redis cache
func (c *RedisCache) SetCart(userID uint, cart []models.Cart) error {
	key := fmt.Sprintf("cart:%d", userID)
	return c.Set(key, cart)
}

// InvalidateUserCache invalidates all user-related Redis cache
func (c *RedisCache) InvalidateUserCache(userID uint) error {
	patterns := []string{
		fmt.Sprintf("user:%d", userID),
		fmt.Sprintf("cart:%d", userID),
	}

	for _, pattern := range patterns {
		keys, err := c.client.Keys(c.ctx, "*"+pattern+"*").Result()
		if err != nil {
			continue
		}
		if len(keys) > 0 {
			c.client.Del(c.ctx, keys...)
		}
	}

	return nil
}

// Cache interface for unified caching operations
type Cache interface {
	Set(key string, value interface{}) error
	Get(key string, dest interface{}) error
	Delete(key string) error
	Exists(key string) bool
	Clear() error
	GetStats() map[string]interface{}

	// Domain-specific methods
	GetUser(userID uint) (*models.User, error)
	SetUser(user *models.User) error
	GetProduct(productID uint) (*models.Product, error)
	SetProduct(product *models.Product) error
	GetProductsByCategory(categoryID uint, page, limit int) ([]models.Product, error)
	SetProductsByCategory(categoryID uint, page, limit int, products []models.Product) error
	GetCart(userID uint) ([]models.Cart, error)
	SetCart(userID uint, cart []models.Cart) error
	InvalidateUserCache(userID uint) error
	InvalidateProductCache(productID uint) error
}

// Global cache instance
var GlobalCache Cache

var newRedisCacheCreator = NewRedisCache

// InitCache initializes the global cache instance with Redis or fallback to in-memory
func InitCache() error {
	// Try Redis first
	if redisCache := newRedisCacheCreator(); redisCache != nil {
		GlobalCache = redisCache
		utils.Info("Using Redis cache for improved performance")
		return nil
	}

	// Fallback to in-memory cache
	GlobalCache = NewInMemoryCache()
	utils.Info("Using in-memory cache (Redis not available)")
	return nil
}

// GetCache returns the global cache instance
func GetCache() Cache {
	return GlobalCache
}
