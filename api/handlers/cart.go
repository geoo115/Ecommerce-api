package handlers

import (
	"errors"
	"net/http"

	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
	"github.com/gin-gonic/gin"
)

type AddToCartInput struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

func AddToCart(c *gin.Context) {
	var input AddToCartInput
	if err := c.ShouldBindJSON(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid input")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		errorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Check stock availability
	if err := CheckStock(input.ProductID, input.Quantity); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Create the cart item
	cartItem := models.Cart{
		UserID:    userID.(uint),
		ProductID: input.ProductID,
		Quantity:  input.Quantity,
	}

	if err := db.DB.Create(&cartItem).Error; err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to add to cart")
		return
	}

	// Reload the cart item with associated relationships
	if err := db.DB.
		Preload("Product.Category").
		Preload("Product.Inventory").
		Preload("User").
		First(&cartItem, cartItem.ID).Error; err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to fetch cart item details")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Item added to cart", "data": cartItem})
}

func errorResponse(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

func CheckStock(productID uint, quantity int) error {
	var inventory models.Inventory
	if err := db.DB.Where("product_id = ?", productID).First(&inventory).Error; err != nil {
		return errors.New("product not found")
	}
	if inventory.Stock < quantity {
		return errors.New("insufficient stock")
	}
	return nil
}

func ListCart(c *gin.Context) {
	var cartItems []models.Cart
	userID, exists := c.Get("userID")
	if !exists {
		errorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := db.DB.Where("user_id = ?", userID).
		Preload("Product.Category").
		Preload("Product.Inventory").
		Preload("User").
		Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list cart items"})
		return
	}

	// Calculate total amount
	var totalAmount float64
	for _, item := range cartItems {
		totalAmount += float64(item.Quantity) * item.Product.Price
	}

	c.JSON(http.StatusOK, gin.H{
		"cart_items":   cartItems,
		"total_amount": totalAmount,
	})
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
