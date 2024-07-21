package handlers

import (
    "ecommerce/db"
    "ecommerce/models"
    "github.com/gin-gonic/gin"
    "net/http"
)

func AddToCart(c *gin.Context) {
    var cart models.Cart
    if err := c.ShouldBindJSON(&cart); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := db.DB.Create(&cart).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, cart)
}

func ListCart(c *gin.Context) {
    var cart []models.Cart
    userID := c.Query("user_id")

    if err := db.DB.Where("user_id = ?", userID).Find(&cart).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, cart)
}

func RemoveFromCart(c *gin.Context) {
    var cart models.Cart
    id := c.Param("id")

    if err := db.DB.Delete(&cart, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Item not found in cart"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart"})
}
