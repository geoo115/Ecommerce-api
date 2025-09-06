package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/geoo115/Ecommerce/api/middlewares"
	"github.com/geoo115/Ecommerce/cache"
	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
	"github.com/geoo115/Ecommerce/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddProduct(c *gin.Context) {
	if db.DB == nil {
		Base.HandleDBError(c, fmt.Errorf("database connection is nil"), "Database error", "Database error")
		return
	}
	productInput, exists := c.Get("product_input")
	if !exists {
		// Keep message aligned with tests expecting a generic invalid input
		utils.SendValidationError(c, "Invalid input")
		return
	}

	input := productInput.(middlewares.ProductInput)

	// Validate product data
	if !utils.ValidateProductName(input.Name) {
		utils.SendValidationError(c, "Product name must be 2-200 characters long")
		return
	}

	// Accept zero-price (free) products per tests; otherwise enforce standard price validation
	if !(input.Price == 0 || utils.ValidatePrice(input.Price)) {
		utils.SendValidationError(c, "Price must be greater than 0 and less than 999999.99")
		return
	}

	if !utils.ValidateDescription(input.Description) {
		utils.SendValidationError(c, "Description must be less than 1000 characters")
		return
	}

	if !utils.ValidateStock(input.Stock) {
		utils.SendValidationError(c, "Stock must be between 0 and 100000")
		return
	}

	// Check if category exists
	var category models.Category
	if err := db.DB.First(&category, input.CategoryID).Error; err != nil {
		utils.SendNotFound(c, "Category not found")
		return
	}

	tx := db.DB.Begin()

	// Create the product
	product := models.Product{
		Name:        utils.SanitizeString(input.Name),
		Price:       input.Price,
		CategoryID:  input.CategoryID,
		Description: utils.SanitizeString(input.Description),
	}

	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		utils.SendInternalError(c, "Failed to create product")
		return
	}

	// Refresh product to ensure ID is populated in transaction
	if err := tx.First(&product, product.ID).Error; err != nil {
		tx.Rollback()
		utils.SendInternalError(c, "Failed to load product after creation")
		return
	}

	// Create inventory record
	inventory := models.Inventory{
		ProductID: product.ID,
		Stock:     input.Stock,
	}

	if err := tx.Create(&inventory).Error; err != nil {
		tx.Rollback()
		utils.SendInternalError(c, "Failed to create inventory: "+err.Error())
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		utils.SendInternalError(c, "Failed to commit transaction")
		return
	}

	// Load the complete product with relationships
	var completeProduct models.Product
	if err := db.DB.Preload("Category").Preload("Inventory").First(&completeProduct, product.ID).Error; err != nil {
		utils.SendInternalError(c, "Failed to load product details")
		return
	}

	// Invalidate product cache if cache is initialized
	if cch := cache.GetCache(); cch != nil {
		if err := cch.InvalidateProductCache(completeProduct.ID); err != nil {
			utils.Warn("Failed to invalidate product cache: %v", err)
		}
	}

	utils.SendSuccess(c, http.StatusCreated, "Product created successfully", completeProduct)
}

func ListProducts(c *gin.Context) {
	ListProductsWithDB(c, db.DB)
}

// ListProductsWithDB handles product listing with database injection
func ListProductsWithDB(c *gin.Context, dbInstance *gorm.DB) {
	if dbInstance == nil {
		Base.HandleDBError(c, fmt.Errorf("database connection is nil"), "Database error", "Database error")
		return
	}
	// Check cache first (avoid shadowing package name and handle nil cache)
	cch := cache.GetCache()
	page := 1
	limit := 10

	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		page = p
	}

	if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 {
		limit = l
	}

	cacheKey := fmt.Sprintf("products:list:page:%d:limit:%d", page, limit)

	// Try to get from cache
	var products []models.Product
	if cch != nil {
		if err := cch.Get(cacheKey, &products); err == nil {
			utils.Info("Products retrieved from cache")
			utils.SendSuccess(c, http.StatusOK, "Products retrieved successfully", products)
			return
		}
	}

	// Cache miss - fetch from database
	query := dbInstance.Preload("Category").Preload("Inventory")
	query = paginate(c, query)

	if err := query.Find(&products).Error; err != nil {
		utils.SendInternalError(c, "Failed to fetch products")
		return
	}

	// Cache the result if cache is available
	if cch != nil {
		if err := cch.Set(cacheKey, products); err != nil {
			utils.Warn("Failed to cache products: %v", err)
		}
	}

	utils.SendSuccess(c, http.StatusOK, "Products retrieved successfully", products)
}

// ListProductsHandlerWrapper wraps the ListProducts handler to inject the database instance
func ListProductsHandlerWrapper(dbInstance *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ListProductsWithDB(c, dbInstance)
	}
}

func GetProduct(c *gin.Context) {
	GetProductWithDB(c, db.DB)
}

// GetProductWithDB retrieves a product by ID with a specific database instance
func GetProductWithDB(c *gin.Context, dbInstance *gorm.DB) {
	if dbInstance == nil {
		Base.HandleDBError(c, fmt.Errorf("database connection is nil"), "Database error", "Database error")
		return
	}
	var product models.Product
	id := c.Param("id")

	// Validate ID parameter
	if id == "" {
		utils.SendValidationError(c, "Product ID is required")
		return
	}

	// Validate that ID is a valid number
	if _, err := strconv.Atoi(id); err != nil {
		utils.SendValidationError(c, "Invalid product ID format")
		return
	}

	if err := dbInstance.Preload("Category").Preload("Inventory").First(&product, id).Error; err != nil {
		utils.SendNotFound(c, "Product not found")
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Product retrieved successfully", product)
}

// GetProductHandlerWrapper wraps the GetProduct handler to inject the database instance
func GetProductHandlerWrapper(dbInstance *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		GetProductWithDB(c, dbInstance)
	}
}

func EditProduct(c *gin.Context) {
	EditProductWithDB(c, db.DB)
}

// EditProductWithDB provides database dependency injection for testing
func EditProductWithDB(c *gin.Context, dbInstance *gorm.DB) {
	if dbInstance == nil {
		Base.HandleDBError(c, fmt.Errorf("database connection is nil"), "Database error", "Database error")
		return
	}

	id := c.Param("id")
	var product models.Product

	// Validate ID parameter
	if id == "" {
		utils.SendValidationError(c, "Product ID is required")
		return
	}

	// Validate that ID is a valid number
	if _, err := strconv.Atoi(id); err != nil {
		utils.SendValidationError(c, "Invalid product ID format")
		return
	}

	// Find the product
	if err := dbInstance.Preload("Category").Preload("Inventory").First(&product, id).Error; err != nil {
		utils.SendNotFound(c, "Product not found")
		return
	}

	// Parse update data from request body
	var updateData struct {
		Name        string  `json:"name"`
		Price       float64 `json:"price"`
		Description string  `json:"description"`
		Stock       int     `json:"stock"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.SendValidationError(c, err.Error())
		return
	}

	// Validate updated data
	if updateData.Name != "" && !utils.ValidateProductName(updateData.Name) {
		utils.SendValidationError(c, "Product name must be 2-200 characters long")
		return
	}

	if updateData.Price != 0 && !utils.ValidatePrice(updateData.Price) {
		utils.SendValidationError(c, "Price must be greater than 0 and less than 999999.99")
		return
	}

	if updateData.Description != "" && !utils.ValidateDescription(updateData.Description) {
		utils.SendValidationError(c, "Description must be less than 1000 characters")
		return
	}

	if updateData.Stock < 0 || (updateData.Stock > 0 && !utils.ValidateStock(updateData.Stock)) {
		utils.SendValidationError(c, "Stock must be between 0 and 100000")
		return
	}

	// Update fields if provided
	if updateData.Name != "" {
		product.Name = utils.SanitizeString(updateData.Name)
	}
	if updateData.Price > 0 {
		product.Price = updateData.Price
	}
	if updateData.Description != "" {
		product.Description = utils.SanitizeString(updateData.Description)
	}

	// Update the product
	if err := dbInstance.Save(&product).Error; err != nil {
		utils.SendInternalError(c, "Failed to update product")
		return
	}

	// Update inventory if stock is provided
	if updateData.Stock >= 0 {
		var inventory models.Inventory
		if err := dbInstance.Where("product_id = ?", product.ID).First(&inventory).Error; err == nil {
			inventory.Stock = updateData.Stock
			dbInstance.Save(&inventory)
		}
	}

	// Load updated product with relationships
	if err := dbInstance.Preload("Category").Preload("Inventory").First(&product, product.ID).Error; err != nil {
		utils.SendInternalError(c, "Failed to load updated product")
		return
	}

	// Invalidate product cache if cache is initialized
	if cch := cache.GetCache(); cch != nil {
		if err := cch.InvalidateProductCache(product.ID); err != nil {
			utils.Warn("Failed to invalidate product cache: %v", err)
		}
	}

	utils.SendSuccess(c, http.StatusOK, "Product updated successfully", product)
}

// EditProductHandlerWrapper wraps the EditProduct handler to inject the database instance
func EditProductHandlerWrapper(dbInstance *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		EditProductWithDB(c, dbInstance)
	}
}

func DeleteProduct(c *gin.Context) {
	DeleteProductWithDB(c, db.DB)
}

// DeleteProductWithDB handles product deletion with database injection
func DeleteProductWithDB(c *gin.Context, dbInstance *gorm.DB) {
	if dbInstance == nil {
		Base.HandleDBError(c, fmt.Errorf("database connection is nil"), "Database error", "Database error")
		return
	}
	id := c.Param("id")
	var product models.Product

	// Validate ID parameter
	if id == "" {
		utils.SendValidationError(c, "Product ID is required")
		return
	}

	// Validate that ID is a valid number
	if _, err := strconv.Atoi(id); err != nil {
		utils.SendValidationError(c, "Invalid product ID format")
		return
	}

	if err := dbInstance.First(&product, id).Error; err != nil {
		utils.SendNotFound(c, "Product not found")
		return
	}

	// Store product ID for cache invalidation
	productID := product.ID

	if err := dbInstance.Delete(&product).Error; err != nil {
		utils.SendInternalError(c, "Failed to delete product")
		return
	}

	// Invalidate product cache if cache is initialized
	if cch := cache.GetCache(); cch != nil {
		if err := cch.InvalidateProductCache(productID); err != nil {
			utils.Warn("Failed to invalidate product cache: %v", err)
		}
	}

	utils.SendSuccess(c, http.StatusOK, "Product deleted successfully", nil)
}

// DeleteProductHandlerWrapper wraps the DeleteProduct handler to inject the database instance
func DeleteProductHandlerWrapper(dbInstance *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		DeleteProductWithDB(c, dbInstance)
	}
}

func SearchProducts(c *gin.Context) {
	SearchProductsWithDB(c, db.DB)
}

// SearchProductsWithDB handles product search with database injection
func SearchProductsWithDB(c *gin.Context, dbInstance *gorm.DB) {
	if dbInstance == nil {
		Base.HandleDBError(c, fmt.Errorf("database connection is nil"), "Database error", "Database error")
		return
	}
	query := c.Query("q")
	category := c.Query("category")

	// Sanitize search query (empty query is allowed to return all products)
	query = utils.SanitizeString(query)
	category = utils.SanitizeString(category)

	// Require a non-empty query per API tests
	if query == "" {
		utils.SendValidationError(c, "Search query is required")
		return
	}

	// Create cache key
	cacheKey := fmt.Sprintf("products:search:q:%s:category:%s", query, category)

	// Check cache first
	cch := cache.GetCache()
	var products []models.Product
	if cch != nil {
		if err := cch.Get(cacheKey, &products); err == nil {
			utils.Info("Search results retrieved from cache")
			utils.SendSuccess(c, http.StatusOK, "Search completed successfully", products)
			return
		}
	}

	// Cache miss - perform search
	dbQuery := dbInstance.Model(&models.Product{})

	// Apply search filter (case-insensitive)
	if query != "" {
		dbQuery = dbQuery.Where("LOWER(products.name) LIKE ?", "%"+query+"%")
	}

	// Apply category filter if provided
	if category != "" {
		dbQuery = dbQuery.Joins("JOIN categories ON products.category_id = categories.id").Where("LOWER(categories.name) = ?", category)
	}

	if err := dbQuery.Preload("Category").Preload("Inventory").Find(&products).Error; err != nil {
		utils.SendInternalError(c, "Failed to fetch products")
		return
	}

	// Cache the search results
	if cch != nil {
		if err := cch.Set(cacheKey, products); err != nil {
			utils.Warn("Failed to cache search results: %v", err)
		}
	}

	utils.SendSuccess(c, http.StatusOK, "Search completed successfully", products)
}

// SearchProductsHandlerWrapper wraps the SearchProducts handler to inject the database instance
func SearchProductsHandlerWrapper(dbInstance *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		SearchProductsWithDB(c, dbInstance)
	}
}

func paginate(c *gin.Context, query *gorm.DB) *gorm.DB {
	limit := 10
	page := c.Query("page")
	if p, err := strconv.Atoi(page); err == nil && p > 0 {
		query = query.Offset((p - 1) * limit).Limit(limit)
	} else {
		query = query.Limit(limit)
	}
	return query
}
