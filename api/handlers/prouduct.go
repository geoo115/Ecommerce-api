package handlers

import (
	"net/http"
	"strconv"

	"github.com/geoo115/Ecommerce/api/middlewares"
	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
	"github.com/geoo115/Ecommerce/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddProduct(c *gin.Context) {
	productInput, exists := c.Get("product_input")
	if !exists {
		utils.SendValidationError(c, "Product data not found in context")
		return
	}

	input := productInput.(middlewares.ProductInput)

	// Validate product data
	if !utils.ValidateProductName(input.Name) {
		utils.SendValidationError(c, "Product name must be 2-200 characters long")
		return
	}

	if !utils.ValidatePrice(input.Price) {
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

	// Create inventory record
	inventory := models.Inventory{
		ProductID: product.ID,
		Stock:     input.Stock,
	}

	if err := tx.Create(&inventory).Error; err != nil {
		tx.Rollback()
		utils.SendInternalError(c, "Failed to create inventory")
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

	utils.SendSuccess(c, http.StatusCreated, "Product created successfully", completeProduct)
}

func ListProducts(c *gin.Context) {
	var products []models.Product
	query := db.DB.Preload("Category").Preload("Inventory")
	query = paginate(c, query)

	if err := query.Find(&products).Error; err != nil {
		utils.SendInternalError(c, "Failed to fetch products")
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Products retrieved successfully", products)
}

func GetProduct(c *gin.Context) {
	var product models.Product
	id := c.Param("id")

	// Validate ID parameter
	if id == "" {
		utils.SendValidationError(c, "Product ID is required")
		return
	}

	if err := db.DB.Preload("Category").Preload("Inventory").First(&product, id).Error; err != nil {
		utils.SendNotFound(c, "Product not found")
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Product retrieved successfully", product)
}

func EditProduct(c *gin.Context) {
	var product models.Product
	id := c.Param("id")

	// Validate ID parameter
	if id == "" {
		utils.SendValidationError(c, "Product ID is required")
		return
	}

	// Check if product exists
	if err := db.DB.First(&product, id).Error; err != nil {
		utils.SendNotFound(c, "Product not found")
		return
	}

	// Bind updated data
	var updateData struct {
		Name        string  `json:"name"`
		Price       float64 `json:"price"`
		Description string  `json:"description"`
		Stock       int     `json:"stock"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.SendValidationError(c, "Invalid input format")
		return
	}

	// Validate updated data
	if updateData.Name != "" && !utils.ValidateProductName(updateData.Name) {
		utils.SendValidationError(c, "Product name must be 2-200 characters long")
		return
	}

	if updateData.Price > 0 && !utils.ValidatePrice(updateData.Price) {
		utils.SendValidationError(c, "Price must be greater than 0 and less than 999999.99")
		return
	}

	if updateData.Description != "" && !utils.ValidateDescription(updateData.Description) {
		utils.SendValidationError(c, "Description must be less than 1000 characters")
		return
	}

	if updateData.Stock >= 0 && !utils.ValidateStock(updateData.Stock) {
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
	if err := db.DB.Save(&product).Error; err != nil {
		utils.SendInternalError(c, "Failed to update product")
		return
	}

	// Update inventory if stock is provided
	if updateData.Stock >= 0 {
		var inventory models.Inventory
		if err := db.DB.Where("product_id = ?", product.ID).First(&inventory).Error; err == nil {
			inventory.Stock = updateData.Stock
			db.DB.Save(&inventory)
		}
	}

	// Load updated product with relationships
	if err := db.DB.Preload("Category").Preload("Inventory").First(&product, product.ID).Error; err != nil {
		utils.SendInternalError(c, "Failed to load updated product")
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Product updated successfully", product)
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	// Validate ID parameter
	if id == "" {
		utils.SendValidationError(c, "Product ID is required")
		return
	}

	if err := db.DB.First(&product, id).Error; err != nil {
		utils.SendNotFound(c, "Product not found")
		return
	}

	if err := db.DB.Delete(&product).Error; err != nil {
		utils.SendInternalError(c, "Failed to delete product")
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Product deleted successfully", nil)
}

func SearchProducts(c *gin.Context) {
	query := c.Query("q")
	category := c.Query("category")
	var products []models.Product

	// Validate search query
	if query == "" {
		utils.SendValidationError(c, "Search query is required")
		return
	}

	// Sanitize search query
	query = utils.SanitizeString(query)
	category = utils.SanitizeString(category)

	dbQuery := db.DB.Where("LOWER(name) LIKE ?", "%"+query+"%")
	if category != "" {
		dbQuery = dbQuery.Joins("Category").Where("LOWER(categories.name) = ?", category)
	}

	if err := dbQuery.Preload("Category").Preload("Inventory").Find(&products).Error; err != nil {
		utils.SendInternalError(c, "Failed to fetch products")
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Search completed successfully", products)
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
