package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
	"github.com/geoo115/Ecommerce/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AddToCartInput struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

func AddToCart(c *gin.Context) {
	var input AddToCartInput
	if err := Base.BindJSON(c, &input); err != nil {
		return
	}

	userID, err := Base.GetUserID(c)
	if err != nil {
		return
	}

	// Validate quantity
	if !utils.ValidateQuantity(input.Quantity) {
		utils.SendValidationError(c, "Invalid quantity")
		return
	}

	// Check stock availability
	hasStock, err := CheckStock(input.ProductID, input.Quantity)
	if err != nil || !hasStock {
		// Tests expect specific message
		if err != nil && err.Error() == "Insufficient stock for product" {
			utils.SendValidationError(c, "Insufficient stock for product")
		} else if err != nil && err.Error() == "product not found" {
			utils.SendNotFound(c, "product not found")
		} else {
			utils.SendValidationError(c, "Insufficient stock for product")
		}
		return
	}

	// If an item exists, increment quantity instead of creating duplicate
	var existing models.Cart
	if err := db.DB.Where("user_id = ? AND product_id = ?", userID, input.ProductID).First(&existing).Error; err == nil {
		newQty := existing.Quantity + input.Quantity
		// Stock check already done for input qty; ensure combined qty still within stock
		if ok, _ := CheckStock(input.ProductID, newQty); !ok {
			utils.SendValidationError(c, "Insufficient stock for product")
			return
		}
		existing.Quantity = newQty
		if err := db.DB.Save(&existing).Error; err != nil {
			utils.SendInternalError(c, "Failed to add to cart")
			return
		}
		utils.SendSuccess(c, http.StatusOK, "Cart item updated successfully", existing)
		return
	}

	// Create the cart item
	cartItem := models.Cart{
		UserID:    userID,
		ProductID: input.ProductID,
		Quantity:  input.Quantity,
	}

	if err := db.DB.Create(&cartItem).Error; err != nil {
		utils.SendInternalError(c, "Failed to add to cart")
		return
	}

	utils.SendSuccess(c, http.StatusCreated, "Item added to cart", cartItem)
}

func CheckStock(productID uint, quantity int) (bool, error) {
	var inventory models.Inventory
	if err := db.DB.Where("product_id = ?", productID).First(&inventory).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, errors.New("product not found")
		}
		return false, err
	}
	if inventory.Stock < quantity {
		return false, errors.New("Insufficient stock for product")
	}
	return true, nil
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
	// Validate ID parameter
	idUint, err := Base.ValidateIDParam(c, "id")
	if err != nil {
		return
	}

	userID, err := Base.GetUserID(c)
	if err != nil {
		return
	}

	// Check if cart item exists and belongs to user
	if err := db.DB.Where("id = ? AND user_id = ?", idUint, userID).First(&cartItem).Error; err != nil {
		utils.SendNotFound(c, "Cart item not found")
		return
	}

	if err := db.DB.Delete(&cartItem).Error; err != nil {
		utils.SendInternalError(c, "Failed to remove from cart")
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Cart item removed successfully", nil)
}

func UpdateCartItem(c *gin.Context) {
	cartItemID := c.Param("id")
	id, err := strconv.Atoi(cartItemID)
	if err != nil {
		utils.SendValidationError(c, "Invalid id")
		return
	}

	userID, err := Base.GetUserID(c)
	if err != nil {
		return
	}

	var input struct {
		Quantity int `json:"quantity" binding:"required,min=1"`
	}

	if err := Base.BindJSON(c, &input); err != nil {
		return
	}

	// Find the cart item
	var cartItem models.Cart
	if err := db.DB.Where("id = ? AND user_id = ?", id, userID).First(&cartItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendNotFound(c, "Resource not found or access denied")
		} else {
			utils.SendInternalError(c, "Database error")
		}
		return
	}

	// Check stock availability
	hasStock, err := CheckStock(cartItem.ProductID, input.Quantity)
	if err != nil || !hasStock {
		utils.SendValidationError(c, "Insufficient stock for product")
		return
	}

	// Update the cart item
	cartItem.Quantity = input.Quantity
	if err := db.DB.Save(&cartItem).Error; err != nil {
		utils.SendInternalError(c, "Failed to update cart item")
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Cart item updated successfully", cartItem)
}
