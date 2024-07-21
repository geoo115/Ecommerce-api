package handlers

import (
	"net/http"

	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
	"github.com/gin-gonic/gin"
)

func AddToCart(c *gin.Context) {
	var cartItem models.Cart
	if err := c.ShouldBindJSON(&cartItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	cartItem.UserID = userID.(uint)

	if err := db.DB.Create(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to cart"})
		return
	}

	c.JSON(http.StatusOK, cartItem)
}

func ListCart(c *gin.Context) {
	var cartItems []models.Cart
	userID, _ := c.Get("userID")

	if err := db.DB.Where("user_id = ?", userID).Preload("Product").Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list cart items"})
		return
	}

	c.JSON(http.StatusOK, cartItems)
}

func RemoveFromCart(c *gin.Context) {
	var cartItem models.Cart
	id := c.Param("id")

	if err := db.DB.Where("id = ?", id).First(&cartItem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}

	if err := db.DB.Delete(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove from cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart item removed"})
}
