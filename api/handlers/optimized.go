package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/geoo115/Ecommerce/cache"
	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
	"github.com/geoo115/Ecommerce/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OptimizedListProducts provides cached product listing
func OptimizedListProducts(c *gin.Context) {
	if db.DB == nil {
		Base.HandleDBError(c, fmt.Errorf("database connection is nil"), "Database error", "Database error")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	categoryID, _ := strconv.ParseUint(c.DefaultQuery("category_id", "0"), 10, 32)

	// Try cache first
	cache := cache.GetCache()
	if cache != nil {
		if categoryID > 0 {
			if products, err := cache.GetProductsByCategory(uint(categoryID), page, limit); err == nil {
				c.JSON(http.StatusOK, gin.H{
					"products": products,
					"cached":   true,
				})
				return
			}
		}
	}

	// Cache miss - query database
	var products []models.Product
	query := db.DB.Preload("Category").Offset((page - 1) * limit).Limit(limit)

	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	if err := query.Find(&products).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to fetch products")
		return
	}

	// Cache the result
	if cache != nil && categoryID > 0 {
		cache.SetProductsByCategory(uint(categoryID), page, limit, products)
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"cached":   false,
	})
}

// OptimizedGetProduct provides cached product retrieval
func OptimizedGetProduct(c *gin.Context) {
	if db.DB == nil {
		Base.HandleDBError(c, fmt.Errorf("database connection is nil"), "Database error", "Database error")
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.SendValidationError(c, "Invalid product ID")
		return
	}

	// Try cache first
	cache := cache.GetCache()
	if cache != nil {
		if product, err := cache.GetProduct(uint(id)); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"product": product,
				"cached":  true,
			})
			return
		}
	}

	// Cache miss - query database with optimized preloading
	var product models.Product
	if err := db.DB.Preload("Category").Preload("Inventory").First(&product, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.SendNotFound(c, "Product not found")
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "Failed to fetch product")
		return
	}

	// Cache the result
	if cache != nil {
		cache.SetProduct(&product)
	}

	c.JSON(http.StatusOK, gin.H{
		"product": product,
		"cached":  false,
	})
}

// OptimizedListCart provides cached cart retrieval
func OptimizedListCart(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if db.DB == nil {
		Base.HandleDBError(c, fmt.Errorf("database connection is nil"), "Database error", "Database error")
		return
	}

	// Try cache first
	cache := cache.GetCache()
	if cache != nil {
		if cart, err := cache.GetCart(userID.(uint)); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"cart":   cart,
				"cached": true,
			})
			return
		}
	}

	// Cache miss - query database with optimized joins
	var cartItems []models.Cart
	if err := db.DB.Where("user_id = ?", userID).Preload("Product.Category").Find(&cartItems).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to fetch cart")
		return
	}

	// Cache the result
	if cache != nil {
		cache.SetCart(userID.(uint), cartItems)
	}

	c.JSON(http.StatusOK, gin.H{
		"cart":   cartItems,
		"cached": false,
	})
}

// OptimizedGetUser provides cached user retrieval
func OptimizedGetUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if db.DB == nil {
		Base.HandleDBError(c, fmt.Errorf("database connection is nil"), "Database error", "Database error")
		return
	}

	// Try cache first
	cache := cache.GetCache()
	if cache != nil {
		if user, err := cache.GetUser(userID.(uint)); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"user":   user,
				"cached": true,
			})
			return
		}
	}

	// Cache miss - query database
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to fetch user")
		return
	}

	// Cache the result
	if cache != nil {
		cache.SetUser(&user)
	}

	c.JSON(http.StatusOK, gin.H{
		"user":   user,
		"cached": false,
	})
}

// OptimizedHealthCheck provides enhanced health check with performance metrics
func OptimizedHealthCheck(c *gin.Context) {
	start := time.Now()

	// Database health check
	dbStats := map[string]interface{}{
		"status": "unknown",
	}

	if db.DB == nil {
		dbStats = map[string]interface{}{
			"status": "unhealthy",
			"error":  "database connection is nil",
		}
	} else if sqlDB, err := db.DB.DB(); err == nil {
		stats := sqlDB.Stats()
		dbStats = map[string]interface{}{
			"status":               "healthy",
			"open_connections":     stats.OpenConnections,
			"in_use_connections":   stats.InUse,
			"idle_connections":     stats.Idle,
			"max_open_connections": stats.MaxOpenConnections,
		}
	} else {
		Base.HandleDBError(c, err, "Database error", "Database error")
		return
	}

	// Cache health check
	cacheStats := map[string]interface{}{
		"status": "disabled",
	}

	if cache := cache.GetCache(); cache != nil {
		cacheStats = cache.GetStats()
		cacheStats["status"] = "healthy"
	}

	// Response time
	responseTime := time.Since(start)

	c.JSON(http.StatusOK, gin.H{
		"status":        "healthy",
		"timestamp":     time.Now().Format(time.RFC3339),
		"response_time": responseTime.String(),
		"database":      dbStats,
		"cache":         cacheStats,
		"version":       "1.0.0",
	})
}

// OptimizedProductSearch provides full-text search with caching
func OptimizedProductSearch(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		utils.SendValidationError(c, "Search query is required")
		return
	}

	if db.DB == nil {
		Base.HandleDBError(c, fmt.Errorf("database connection is nil"), "Database error", "Database error")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// Try cache first
	cache := cache.GetCache()
	cacheKey := fmt.Sprintf("search:%s:page:%d:limit:%d", query, page, limit)
	if cache != nil {
		var products []models.Product
		if err := cache.Get(cacheKey, &products); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"products": products,
				"query":    query,
				"cached":   true,
			})
			return
		}
	}

	// Cache miss - perform search
	var products []models.Product
	searchQuery := "%" + strings.ToUpper(query) + "%"

	if err := db.DB.Where("UPPER(name) LIKE ? OR UPPER(description) LIKE ?", searchQuery, searchQuery).
		Preload("Category").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&products).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Search failed")
		return
	}

	// Cache the result
	if cache != nil {
		cache.Set(cacheKey, products)
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"query":    query,
		"cached":   false,
	})
}

// OptimizedOrderHistory provides paginated order history with caching
func OptimizedOrderHistory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if db.DB == nil {
		Base.HandleDBError(c, fmt.Errorf("database connection is nil"), "Database error", "Database error")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Try cache first
	cache := cache.GetCache()
	cacheKey := fmt.Sprintf("orders:user:%d:page:%d:limit:%d", userID.(uint), page, limit)
	if cache != nil {
		var orders []models.Order
		if err := cache.Get(cacheKey, &orders); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"orders": orders,
				"cached": true,
			})
			return
		}
	}

	// Cache miss - query database with optimized joins
	var orders []models.Order
	if err := db.DB.Where("user_id = ?", userID).
		Preload("Items.Product").
		Order("created_at DESC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&orders).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to fetch orders")
		return
	}

	// Cache the result
	if cache != nil {
		cache.Set(cacheKey, orders)
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"cached": false,
	})
}
