package handlers

import (
	"net/http"

	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
	"github.com/gin-gonic/gin"
)

// AddReview adds a new review for a product
func AddReview(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var reviewRequest struct {
		ProductID uint   `json:"product_id" binding:"required"`
		Rating    int    `json:"rating" binding:"required,min=1,max=5"`
		Comment   string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&reviewRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the product exists
	var product models.Product
	if err := db.DB.First(&product, reviewRequest.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Create the review
	review := models.Review{
		ProductID: reviewRequest.ProductID,
		UserID:    userID.(uint),
		Rating:    reviewRequest.Rating,
		Comment:   reviewRequest.Comment,
	}

	if err := db.DB.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := db.DB.Preload("Product.Category").
		Preload("Product.Inventory").
		Preload("User").
		First(&review, review.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load user data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review added successfully", "review": review})
}

// ListReviews lists all reviews for a specific product
func ListReviews(c *gin.Context) {
	productID := c.Param("product_id")

	var reviews []models.Review
	if err := db.DB.Where("product_id = ?", productID).Preload("User").Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(reviews) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No reviews found for this product"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}
