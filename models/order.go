package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID                uint        `json:"user_id"`
	TotalAmount           float64     `json:"total_amount"`
	Status                string      `json:"status"` // e.g., "Pending", "Shipped", "Delivered", "Cancelled"
	Items                 []OrderItem `gorm:"foreignKey:OrderID"`
	User                  User        `gorm:"foreignKey:UserID"`
	TrackingNumber        string      `json:"tracking_number"`
	Courier               string      `json:"courier"`
	EstimatedDeliveryDate string      `json:"estimated_delivery_date"`
}

type OrderItem struct {
	gorm.Model
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"` // Price at the time of order
	Product   Product `gorm:"foreignKey:ProductID"`
}
