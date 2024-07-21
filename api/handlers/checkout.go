package handlers

import (
    "ecommerce/db"
    "ecommerce/models"
    "github.com/gin-gonic/gin"
    "net/http"
)

func Checkout(c *gin.Context) {
    var cart []models.Cart
    userID := c.Query("user_id")

    if err := db.DB.Where("user_id = ?", userID).Find(&cart).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    if len(cart) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
        return
    }

    // Here you could implement payment processing logic
    // For simplicity, we'll just clear the cart

    if err := db.DB.Where("user_id = ?", userID).Delete(&models.Cart{}).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Checkout successful"})
}
