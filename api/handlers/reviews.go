package handlers

import (
	"net/http"
	"strconv"

	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
	"github.com/geoo115/Ecommerce/utils"
	"github.com/gin-gonic/gin"
)

// AddReview adds a new review for a product
func AddReview(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.SendUnauthorized(c, "Unauthorized")
		return
	}

	var reviewRequest struct {
		Rating  int    `json:"rating" binding:"required,min=1,max=5"`
		Comment string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&reviewRequest); err != nil {
		utils.SendValidationError(c, err.Error())
		return
	}

	// Get product ID from path param :id
	pidStr := c.Param("id")
	pid, err := strconv.ParseUint(pidStr, 10, 64)
	if err != nil {
		utils.SendValidationError(c, "Invalid product ID")
		return
	}

	// Check if the product exists
	var product models.Product
	if err := db.DB.First(&product, uint(pid)).Error; err != nil {
		utils.SendNotFound(c, "Product not found")
		return
	}

	// Create the review
	review := models.Review{
		ProductID: uint(pid),
		UserID:    userID.(uint),
		Rating:    reviewRequest.Rating,
		Comment:   reviewRequest.Comment,
	}

	if err := db.DB.Create(&review).Error; err != nil {
		utils.SendInternalError(c, "Failed to create review")
		return
	}
	if err := db.DB.Preload("Product.Category").
		Preload("User").
		First(&review, review.ID).Error; err != nil {
		utils.SendInternalError(c, "Failed to load user data")
		return
	}

	utils.SendSuccess(c, http.StatusCreated, "Review added successfully", gin.H{"review": review})
}

// ListReviews lists all reviews for a specific product
func ListReviews(c *gin.Context) {
	pidStr := c.Param("id")
	pid, err := strconv.ParseUint(pidStr, 10, 64)
	if err != nil {
		utils.SendValidationError(c, "Invalid product ID")
		return
	}

	// Check if the product exists
	var product models.Product
	if err := db.DB.First(&product, uint(pid)).Error; err != nil {
		utils.SendNotFound(c, "Product not found")
		return
	}

	var reviews []models.Review
	if err := db.DB.Where("product_id = ?", uint(pid)).Preload("User").Find(&reviews).Error; err != nil {
		utils.SendInternalError(c, "Failed to fetch reviews")
		return
	}

	if len(reviews) == 0 {
		utils.SendNotFound(c, "No reviews found for this product")
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Reviews retrieved successfully", reviews)
}
