package handlers

import (
	"errors"
	"net/http"

	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
	"github.com/geoo115/Ecommerce/utils"
	"github.com/gin-gonic/gin"
)

type AddToCartInput struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

func AddToCart(c *gin.Context) {
	var input AddToCartInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendValidationError(c, "Invalid input")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.SendUnauthorized(c, "Unauthorized")
		return
	}

	// Validate quantity
	if !utils.ValidateQuantity(input.Quantity) {
		utils.SendValidationError(c, "Quantity must be between 1 and 1000")
		return
	}

	// Check stock availability
	if err := CheckStock(input.ProductID, input.Quantity); err != nil {
		utils.SendValidationError(c, err.Error())
		return
	}

	// Create the cart item
	cartItem := models.Cart{
		UserID:    userID.(uint),
		ProductID: input.ProductID,
		Quantity:  input.Quantity,
	}

	if err := db.DB.Create(&cartItem).Error; err != nil {
		utils.SendInternalError(c, "Failed to add to cart")
		return
	}

	// Reload the cart item with associated relationships
	if err := db.DB.
		Preload("Product.Category").
		Preload("Product.Inventory").
		Preload("User").
		First(&cartItem, cartItem.ID).Error; err != nil {
		utils.SendInternalError(c, "Failed to fetch cart item details")
		return
	}

	utils.SendSuccess(c, http.StatusCreated, "Item added to cart", cartItem)
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
		utils.SendUnauthorized(c, "Unauthorized")
		return
	}

	if err := db.DB.Where("user_id = ?", userID).
		Preload("Product.Category").
		Preload("Product.Inventory").
		Preload("User").
		Find(&cartItems).Error; err != nil {
		utils.SendInternalError(c, "Failed to list cart items")
		return
	}

	// Calculate total amount
	var totalAmount float64
	for _, item := range cartItems {
		totalAmount += float64(item.Quantity) * item.Product.Price
	}

	utils.SendSuccess(c, http.StatusOK, "Cart items retrieved successfully", gin.H{
		"cart_items":   cartItems,
		"total_amount": totalAmount,
	})
}

func RemoveFromCart(c *gin.Context) {
	var cartItem models.Cart
	id := c.Param("id")

	// Validate ID parameter
	if id == "" {
		utils.SendValidationError(c, "Cart item ID is required")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.SendUnauthorized(c, "Unauthorized")
		return
	}

	// Check if cart item exists and belongs to user
	if err := db.DB.Where("id = ? AND user_id = ?", id, userID).First(&cartItem).Error; err != nil {
		utils.SendNotFound(c, "Cart item not found")
		return
	}

	if err := db.DB.Delete(&cartItem).Error; err != nil {
		utils.SendInternalError(c, "Failed to remove from cart")
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Cart item removed successfully", nil)
}
