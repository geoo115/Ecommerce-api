package handlers

import (
	"net/http"

	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
	"github.com/geoo115/Ecommerce/utils"
	"github.com/gin-gonic/gin"
)

// PlaceOrder creates a new order for the authenticated user
func PlaceOrder(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.SendUnauthorized(c, "Unauthorized")
		return
	}

	var orderRequest struct {
		Items []struct {
			ProductID uint `json:"product_id" binding:"required"`
			Quantity  int  `json:"quantity" binding:"required,min=1"`
		} `json:"items" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		utils.SendValidationError(c, "Invalid order format")
		return
	}

	// Validate order items
	if len(orderRequest.Items) == 0 {
		utils.SendValidationError(c, "Order must contain at least one item")
		return
	}

	var totalAmount float64
	var orderItems []models.OrderItem

	for _, item := range orderRequest.Items {
		// Validate quantity
		if !utils.ValidateQuantity(item.Quantity) {
			utils.SendValidationError(c, "Quantity must be between 1 and 1000")
			return
		}

		var product models.Product
		if err := db.DB.First(&product, item.ProductID).Error; err != nil {
			utils.SendNotFound(c, "Product not found")
			return
		}

		// Check inventory
		var inventory models.Inventory
		if err := db.DB.Where("product_id = ?", item.ProductID).First(&inventory).Error; err != nil || inventory.Stock < item.Quantity {
			utils.SendValidationError(c, "Insufficient stock for product")
			return
		}

		// Deduct stock
		inventory.Stock -= item.Quantity
		db.DB.Save(&inventory)

		totalAmount += product.Price * float64(item.Quantity)
		orderItems = append(orderItems, models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})
	}

	// Create the order
	order := models.Order{
		UserID:      userID.(uint),
		TotalAmount: totalAmount,
		Status:      "Pending",
		Items:       orderItems,
	}

	if err := db.DB.Create(&order).Error; err != nil {
		utils.SendInternalError(c, "Failed to create order")
		return
	}

	// Load complete order with relationships
	if err := db.DB.Preload("Items.Product").Preload("Items.Product.Category").First(&order, order.ID).Error; err != nil {
		utils.SendInternalError(c, "Failed to load order details")
		return
	}

	utils.SendSuccess(c, http.StatusCreated, "Order placed successfully", order)
}

// ListOrders retrieves all orders for the authenticated user
func ListOrders(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.SendUnauthorized(c, "Unauthorized")
		return
	}

	var orders []models.Order
	if err := db.DB.Where("user_id = ?", userID).Preload("Items.Product").Find(&orders).Error; err != nil {
		utils.SendInternalError(c, "Failed to fetch orders")
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Orders retrieved successfully", orders)
}

// GetOrder retrieves details of a specific order for the authenticated user
func GetOrder(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.SendUnauthorized(c, "Unauthorized")
		return
	}

	orderID := c.Param("id")
	if orderID == "" {
		utils.SendValidationError(c, "Order ID is required")
		return
	}

	var order models.Order

	if err := db.DB.Where("id = ? AND user_id = ?", orderID, userID).
		Preload("Items.Product").
		Preload("Items.Product.Category").
		Preload("Items.Product.Inventory").
		Preload("User").First(&order).Error; err != nil {
		utils.SendNotFound(c, "Order not found")
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Order retrieved successfully", order)
}

// CancelOrder cancels an existing order if it is still pending
func CancelOrder(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.SendUnauthorized(c, "Unauthorized")
		return
	}

	orderID := c.Param("id")
	if orderID == "" {
		utils.SendValidationError(c, "Order ID is required")
		return
	}

	var order models.Order

	if err := db.DB.Where("id = ? AND user_id = ?", orderID, userID).Preload("Items").First(&order).Error; err != nil {
		utils.SendNotFound(c, "Order not found")
		return
	}

	if order.Status != "Pending" {
		utils.SendValidationError(c, "Order cannot be canceled")
		return
	}

	// Update status to "Cancelled"
	order.Status = "Cancelled"
	if err := db.DB.Save(&order).Error; err != nil {
		utils.SendInternalError(c, "Failed to cancel order")
		return
	}

	// Restock inventory
	for _, item := range order.Items {
		var inventory models.Inventory
		if err := db.DB.Where("product_id = ?", item.ProductID).First(&inventory).Error; err == nil {
			inventory.Stock += item.Quantity
			db.DB.Save(&inventory)
		}
	}

	utils.SendSuccess(c, http.StatusOK, "Order canceled successfully", nil)
}
