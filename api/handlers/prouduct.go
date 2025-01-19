package handlers

import (
	"net/http"
	"strconv"

	"github.com/geoo115/Ecommerce/api/middlewares"
	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddProduct(c *gin.Context) {
	productInput, exists := c.Get("product_input")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Product data not found in context",
		})
		return
	}

	input := productInput.(middlewares.ProductInput)

	tx := db.DB.Begin()

	// Create the product
	product := models.Product{
		Name:        input.Name,
		Price:       input.Price,
		CategoryID:  input.CategoryID,
		Description: input.Description,
	}

	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create product",
			"details": err.Error(),
		})
		return
	}

	// Create inventory record
	inventory := models.Inventory{
		ProductID: product.ID,
		Stock:     input.Stock,
	}

	if err := tx.Create(&inventory).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create inventory",
			"details": err.Error(),
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to commit transaction",
			"details": err.Error(),
		})
		return
	}

	// Load the complete product with relationships
	var completeProduct models.Product
	if err := db.DB.Preload("Category").Preload("Inventory").First(&completeProduct, product.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to load product details",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"data":    completeProduct,
	})
}

func ListProducts(c *gin.Context) {
	var products []models.Product
	query := db.DB.Preload("Category").Preload("Inventory")
	query = paginate(c, query)

	if err := query.Find(&products).Error; err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to fetch products")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": products})
}

func GetProduct(c *gin.Context) {
	var product models.Product
	id := c.Param("id")

	if err := db.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func EditProduct(c *gin.Context) {
	var product models.Product
	id := c.Param("id")

	// Check if product exists
	if err := db.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Bind updated data
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the product
	if err := db.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := db.DB.First(&product, id).Error; err != nil {
		errorResponse(c, http.StatusNotFound, "Product not found")
		return
	}

	if err := db.DB.Delete(&product).Error; err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to delete product")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func SearchProducts(c *gin.Context) {
	query := c.Query("q")
	category := c.Query("category")
	var products []models.Product

	dbQuery := db.DB.Where("LOWER(name) LIKE ?", "%"+query+"%")
	if category != "" {
		dbQuery = dbQuery.Joins("Category").Where("LOWER(categories.name) = ?", category)
	}

	if err := dbQuery.Preload("Category").Preload("Inventory").Find(&products).Error; err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to fetch products")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": products})
}

func paginate(c *gin.Context, query *gorm.DB) *gorm.DB {
	limit := 10
	page := c.Query("page")
	if p, err := strconv.Atoi(page); err == nil {
		query = query.Offset((p - 1) * limit).Limit(limit)
	} else {
		query = query.Limit(limit)
	}
	return query
}
