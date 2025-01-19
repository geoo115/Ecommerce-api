package models

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	OrderID     uint    `json:"order_id"`
	PaymentMode string  `json:"payment_mode"` // e.g., "Credit Card", "PayPal", "Cash on Delivery"
	Amount      float64 `json:"amount"`
	Status      string  `json:"status"` // e.g., "Success", "Failed", "Pending"
	Order       Order   `gorm:"foreignKey:OrderID"`
}
