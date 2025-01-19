package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/geoo115/Ecommerce/db"
	"github.com/geoo115/Ecommerce/models"
	"github.com/gin-gonic/gin"
)

// SalesReport generates a sales report with optional date filtering
func SalesReport(c *gin.Context) {
	// Parse date parameters
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	var start, end time.Time
	var err error

	// Parse dates if provided
	if startDate != "" {
		start, err = time.Parse("2006-01-02", startDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use YYYY-MM-DD"})
			return
		}
	}

	if endDate != "" {
		end, err = time.Parse("2006-01-02", endDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Use YYYY-MM-DD"})
			return
		}
		// Set end date to end of day
		end = end.Add(24*time.Hour - time.Second)
	}

	var report []struct {
		ProductName string  `json:"product_name"`
		TotalSold   int     `json:"total_sold"`
		TotalSales  float64 `json:"total_sales"`
		Period      string  `json:"period"`
	}

	// Create the period string
	periodStr := "All Time"
	if startDate != "" && endDate != "" {
		periodStr = fmt.Sprintf("%s to %s", startDate, endDate)
	}

	// Build the query with fixed period
	query := db.DB.Model(&models.OrderItem{}).
		Select(`
			products.name AS product_name, 
			SUM(order_items.quantity) AS total_sold, 
			SUM(order_items.quantity * order_items.price) AS total_sales,
			? AS period
		`, periodStr).
		Joins("INNER JOIN products ON products.id = order_items.product_id").
		Joins("INNER JOIN orders ON orders.id = order_items.order_id")

	// Add date filters if provided
	if !start.IsZero() {
		query = query.Where("orders.created_at >= ?", start)
	}
	if !end.IsZero() {
		query = query.Where("orders.created_at <= ?", end)
	}

	// Execute the query
	if err := query.Group("products.name").Scan(&report).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate sales report: " + err.Error()})
		return
	}

	// Calculate total summary
	var totalSummary struct {
		TotalQuantity int     `json:"total_quantity"`
		TotalRevenue  float64 `json:"total_revenue"`
	}
	for _, item := range report {
		totalSummary.TotalQuantity += item.TotalSold
		totalSummary.TotalRevenue += item.TotalSales
	}

	// Add period information to response
	response := gin.H{
		"sales_report": report,
		"summary":      totalSummary,
		"filters": gin.H{
			"start_date": startDate,
			"end_date":   endDate,
		},
	}

	c.JSON(http.StatusOK, response)
}

// InventoryReport generates an inventory report with optional date filtering
func InventoryReport(c *gin.Context) {
	// Parse date parameters
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	var start, end time.Time
	var err error

	// Parse dates if provided
	if startDate != "" {
		start, err = time.Parse("2006-01-02", startDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use YYYY-MM-DD"})
			return
		}
	}

	if endDate != "" {
		end, err = time.Parse("2006-01-02", endDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Use YYYY-MM-DD"})
			return
		}
		// Set end date to end of day
		end = end.Add(24*time.Hour - time.Second)
	}

	var report []struct {
		ProductName  string    `json:"product_name"`
		CurrentStock int       `json:"current_stock"`
		StockValue   float64   `json:"stock_value"`
		LastUpdated  time.Time `json:"last_updated"`
		Category     string    `json:"category"`
	}

	// Updated query to use inventories table
	query := db.DB.Model(&models.Product{}).
		Select(`
			products.name AS product_name,
			inventories.stock AS current_stock,
			inventories.stock * products.price AS stock_value,
			inventories.updated_at AS last_updated,
			categories.name AS category
		`).
		Joins("LEFT JOIN inventories ON inventories.product_id = products.id").
		Joins("LEFT JOIN categories ON categories.id = products.category_id")

	// Add date filters if provided
	if !start.IsZero() {
		query = query.Where("inventories.updated_at >= ?", start)
	}
	if !end.IsZero() {
		query = query.Where("inventories.updated_at <= ?", end)
	}

	// Execute the query
	if err := query.Scan(&report).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate inventory report: " + err.Error()})
		return
	}

	// Calculate summary statistics
	var totalItems int
	var totalValue float64
	for _, item := range report {
		totalItems += item.CurrentStock
		totalValue += item.StockValue
	}

	// Add period information and summary to response
	response := gin.H{
		"inventory_report": report,
		"summary": gin.H{
			"total_items": totalItems,
			"total_value": totalValue,
		},
		"filters": gin.H{
			"start_date": startDate,
			"end_date":   endDate,
		},
	}

	c.JSON(http.StatusOK, response)
}
