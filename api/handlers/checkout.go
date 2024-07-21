package handlers

import (
	"net/http"

	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
	"github.com/gin-gonic/gin"
)

func Checkout(c *gin.Context) {
	var cartItems []models.Cart
	userID, _ := c.Get("userID")

	if err := db.DB.Where("user_id = ?", userID).Preload("Product").Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart items"})
		return
	}

	// Here you would typically process the payment and create an order record
	// For simplicity, we will just clear the cart

	if err := db.DB.Where("user_id = ?", userID).Delete(&models.Cart{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Checkout successful"})
}
