package handlers

import (
	"net/http"

	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
	"github.com/gin-gonic/gin"
)

// AddToWishlist adds a product to the user's wishlist
func AddToWishlist(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var wishlistRequest struct {
		ProductID uint `json:"product_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&wishlistRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the product exists
	var product models.Product
	if err := db.DB.Preload("Category").Preload("Inventory").First(&product, wishlistRequest.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Check if the product is already in the wishlist
	var existingWishlist models.Wishlist
	if err := db.DB.Where("user_id = ? AND product_id = ?", userID, wishlistRequest.ProductID).First(&existingWishlist).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Product already in wishlist"})
		return
	}

	// Add to wishlist
	wishlist := models.Wishlist{
		UserID:    userID.(uint),
		ProductID: wishlistRequest.ProductID,
	}

	if err := db.DB.Create(&wishlist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Fetch the complete wishlist record with preloaded data
	if err := db.DB.Preload("Product.Category").Preload("Product.Inventory").Preload("User").First(&wishlist, wishlist.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product added to wishlist", "wishlist": wishlist})
}

// ListWishlist retrieves the user's wishlist
func ListWishlist(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var wishlist []models.Wishlist
	if err := db.DB.Where("user_id = ?", userID).
		Preload("Product.Category").
		Preload("Product.Inventory").
		Preload("User").
		Find(&wishlist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(wishlist) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No items in wishlist"})
		return
	}

	c.JSON(http.StatusOK, wishlist)
}

// RemoveFromWishlist removes a product from the user's wishlist
func RemoveFromWishlist(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	wishlistID := c.Param("id")
	var wishlistItem models.Wishlist

	// Check if the wishlist item exists and belongs to the user
	if err := db.DB.Where("id = ? AND user_id = ?", wishlistID, userID).First(&wishlistItem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wishlist item not found"})
		return
	}

	// Remove the item from the wishlist
	if err := db.DB.Delete(&wishlistItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product removed from wishlist"})
}
