package handlers

import (
	"net/http"

	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
	"github.com/gin-gonic/gin"
)

// ProcessPayment handles payment processing
func ProcessPayment(c *gin.Context) {
	var paymentRequest struct {
		OrderID       uint    `json:"order_id" binding:"required"`
		PaymentMethod string  `json:"payment_method" binding:"required"`
		Amount        float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&paymentRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the order exists
	var order models.Order
	if err := db.DB.First(&order, paymentRequest.OrderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Validate the payment amount
	if paymentRequest.Amount != order.TotalAmount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment amount"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record payment"})
		return
	}

	// Update the order status
	order.Status = "Paid"
	if err := db.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	// Fetch the payment with related order and user details
	if err := db.DB.Preload("Order.User").First(&payment, payment.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment processed successfully", "payment": payment})
}

// GetPaymentStatus retrieves the payment status of a specific order
func GetPaymentStatus(c *gin.Context) {
	orderID := c.Param("order_id")

	var payment models.Payment
	if err := db.DB.Where("order_id = ?", orderID).First(&payment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found for the given order ID"})
		return
	}

	// Fetch the payment with related order and user details
	if err := db.DB.Preload("Order.User").First(&payment, payment.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// Checkout processes the checkout by clearing the cart and creating an order
func Checkout(c *gin.Context) {
	var cartItems []models.Cart
	userID, _ := c.Get("userID")

	// Fetch all cart items for the user
	if err := db.DB.Where("user_id = ?", userID).Preload("Product").Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart items"})
		return
	}

	// Calculate total amount
	var totalAmount float64
	for _, item := range cartItems {
		totalAmount += float64(item.Quantity) * item.Product.Price
	}

	// Create a new order
	order := models.Order{
		UserID:      userID.(uint),
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Clear the user's cart
	if err := db.DB.Where("user_id = ?", userID).Delete(&models.Cart{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Checkout successful", "order": order})
}
