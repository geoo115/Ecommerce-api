package handlers

import (
	"net/http"

	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
	"github.com/gin-gonic/gin"
)

// PlaceOrder creates a new order for the authenticated user
func PlaceOrder(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var orderRequest struct {
		Items []struct {
			ProductID uint `json:"product_id"`
			Quantity  int  `json:"quantity"`
		} `json:"items"`
	}

	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var totalAmount float64
	var orderItems []models.OrderItem

	for _, item := range orderRequest.Items {
		var product models.Product
		if err := db.DB.First(&product, item.ProductID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		// Check inventory
		var inventory models.Inventory
		if err := db.DB.Where("product_id = ?", item.ProductID).First(&inventory).Error; err != nil || inventory.Stock < item.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock for product ID", "product_id": item.ProductID})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// ListOrders retrieves all orders for the authenticated user
func ListOrders(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var orders []models.Order
	if err := db.DB.Where("user_id = ?", userID).Preload("Items.Product").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetOrder retrieves details of a specific order for the authenticated user
func GetOrder(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	orderID := c.Param("id")
	var order models.Order

	if err := db.DB.Where("id = ? AND user_id = ?", orderID, userID).
		Preload("Items.Product").
		Preload("Items.Product.Category").
		Preload("Items.Product.Inventory").
		Preload("User").First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// CancelOrder cancels an existing order if it is still pending
func CancelOrder(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	orderID := c.Param("id")
	var order models.Order

	if err := db.DB.Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if order.Status != "Pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order cannot be canceled"})
		return
	}

	// Update status to "Cancelled"
	order.Status = "Cancelled"
	if err := db.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	c.JSON(http.StatusOK, gin.H{"message": "Order canceled successfully"})
}
