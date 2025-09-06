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

// ProcessPayment handles payment processing
func ProcessPayment(c *gin.Context) {
	var paymentRequest struct {
		OrderID       uint    `json:"order_id" binding:"required"`
		PaymentMethod string  `json:"payment_method" binding:"required"`
		Amount        float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&paymentRequest); err != nil {
		utils.SendValidationError(c, err.Error())
		return
	}

	// Check if the order exists
	var order models.Order
	if err := db.DB.First(&order, paymentRequest.OrderID).Error; err != nil {
		// DB closed or other DB error => 500, record not found => 404
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendInternalError(c, "Internal server error")
			return
		}
		utils.SendNotFound(c, "Order not found")
		return
	}

	// Validate the payment amount
	if paymentRequest.Amount != order.TotalAmount {
		utils.SendValidationError(c, "Invalid payment amount")
		return
	}

	// Record the payment
	payment := models.Payment{
		OrderID:     paymentRequest.OrderID,
		PaymentMode: paymentRequest.PaymentMethod,
		Amount:      paymentRequest.Amount,
		Status:      "Success",
	}

	if err := db.DB.Create(&payment).Error; err != nil {
		utils.SendInternalError(c, "Failed to record payment")
		return
	}

	// Update the order status
	order.Status = "Paid"
	if err := db.DB.Save(&order).Error; err != nil {
		utils.SendInternalError(c, "Failed to update order status")
		return
	}

	// Fetch the payment with related order and user details
	if err := db.DB.Preload("Order.User").First(&payment, payment.ID).Error; err != nil {
		utils.SendInternalError(c, "Failed to load payment details")
		return
	}
	utils.SendSuccess(c, http.StatusOK, "Payment processed successfully", gin.H{"payment": payment})
}

// GetPaymentStatus retrieves the payment status of a specific order
func GetPaymentStatus(c *gin.Context) {
	orderIDStr := c.Param("order_id")
	orderIDUint64, convErr := strconv.ParseUint(orderIDStr, 10, 64)
	if convErr != nil {
		utils.SendNotFound(c, "Payment not found for the given order ID")
		return
	}

	var payment models.Payment
	if err := db.DB.Preload("Order.User").Where("order_id = ?", uint(orderIDUint64)).First(&payment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendNotFound(c, "Payment not found for the given order ID")
			return
		}
		utils.SendInternalError(c, "Internal server error")
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Payment status retrieved", gin.H{"payment": payment})
}

// Checkout processes the checkout by clearing the cart and creating an order
func Checkout(c *gin.Context) {
	var cartItems []models.Cart
	uid, err := Base.GetUserID(c)
	if err != nil {
		return
	}

	// Fetch all cart items for the user
	if err := db.DB.Where("user_id = ?", uid).Preload("Product").Find(&cartItems).Error; err != nil {
		utils.SendInternalError(c, "Failed to fetch cart items")
		return
	}

	// Calculate total amount
	var totalAmount float64
	for _, item := range cartItems {
		totalAmount += float64(item.Quantity) * item.Product.Price
	}

	// Create a new order
	order := models.Order{
		UserID:      uid,
		TotalAmount: totalAmount,
		Status:      "Pending",
		Items:       []models.OrderItem{},
	}

	for _, item := range cartItems {
		order.Items = append(order.Items, models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Product.Price,
		})
	}

	if err := db.DB.Create(&order).Error; err != nil {
		utils.SendInternalError(c, "Failed to create order")
		return
	}

	// Decrement inventory for purchased items (best-effort)
	for _, item := range cartItems {
		var inv models.Inventory
		if err := db.DB.Where("product_id = ?", item.ProductID).First(&inv).Error; err == nil {
			inv.Stock = inv.Stock - item.Quantity
			db.DB.Save(&inv)
		}
	}

	// Clear the user's cart
	if err := db.DB.Where("user_id = ?", uid).Delete(&models.Cart{}).Error; err != nil {
		utils.SendInternalError(c, "Failed to clear cart")
		return
	}
	utils.SendSuccess(c, http.StatusOK, "Checkout successful", gin.H{"order": order})
}
