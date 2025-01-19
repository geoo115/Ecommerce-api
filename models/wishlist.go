package models

import "gorm.io/gorm"

type Wishlist struct {
	gorm.Model
	UserID    uint    `json:"user_id"`
	ProductID uint    `json:"product_id"`
	User      User    `gorm:"foreignKey:UserID"`
	Product   Product `gorm:"foreignKey:ProductID"`
}

// type User struct {
// 	gorm.Model
// 	Username  string    `json:"username" gorm:"unique"`
// 	Password  string    `json:"password"`
// 	Email     string    `json:"email" gorm:"unique"`
// 	Phone     string    `json:"phone" gorm:"unique"`
// 	Role      string    `json:"role" gorm:"default:customer"` // Role field with default value "customer"
// 	Addresses []Address `gorm:"foreignKey:UserID"`
// 	Cart      []Cart    `gorm:"foreignKey:UserID"`
// }
// type Review struct {
// 	gorm.Model
// 	ProductID uint    `json:"product_id"`
// 	UserID    uint    `json:"user_id"`
// 	Rating    int     `json:"rating"` // 1-5 stars
// 	Comment   string  `json:"comment"`
// 	Product   Product `gorm:"foreignKey:ProductID"`
// 	User      User    `gorm:"foreignKey:UserID"`
// }
// type Category struct {
// 	gorm.Model
// 	Name     string    `json:"name" gorm:"unique"`
// 	Products []Product `gorm:"foreignKey:CategoryID"`
// }

// type Product struct {
// 	gorm.Model
// 	Name        string    `json:"name"`
// 	Price       float64   `json:"price"`
// 	CategoryID  uint      `json:"category_id"`
// 	Description string    `json:"description"`
// 	Category    Category  `json:"category" gorm:"foreignKey:CategoryID"`
// 	Cart        []Cart    `json:"-" gorm:"foreignKey:ProductID"` // Hide in JSON
// 	Inventory   Inventory `json:"inventory" gorm:"foreignKey:ProductID"`
// }
// type Payment struct {
// 	gorm.Model
// 	OrderID     uint    `json:"order_id"`
// 	PaymentMode string  `json:"payment_mode"` // e.g., "Credit Card", "PayPal", "Cash on Delivery"
// 	Amount      float64 `json:"amount"`
// 	Status      string  `json:"status"` // e.g., "Success", "Failed", "Pending"
// 	Order       Order   `gorm:"foreignKey:OrderID"`
// }

// type Order struct {
// 	gorm.Model
// 	UserID                uint        `json:"user_id"`
// 	TotalAmount           float64     `json:"total_amount"`
// 	Status                string      `json:"status"` // e.g., "Pending", "Shipped", "Delivered", "Cancelled"
// 	Items                 []OrderItem `gorm:"foreignKey:OrderID"`
// 	User                  User        `gorm:"foreignKey:UserID"`
// 	TrackingNumber        string      `json:"tracking_number"`
// 	Courier               string      `json:"courier"`
// 	EstimatedDeliveryDate string      `json:"estimated_delivery_date"`
// }

// type OrderItem struct {
// 	gorm.Model
// 	OrderID   uint    `json:"order_id"`
// 	ProductID uint    `json:"product_id"`
// 	Quantity  int     `json:"quantity"`
// 	Price     float64 `json:"price"` // Price at the time of order
// 	Product   Product `gorm:"foreignKey:ProductID"`
// }
// type Inventory struct {
// 	gorm.Model
// 	ProductID uint     `json:"product_id"`
// 	Stock     int      `json:"stock"`
// 	Product   *Product `json:"-" gorm:"foreignKey:ProductID"`
// }
// type Cart struct {
//     gorm.Model
//     UserID    uint    `json:"user_id"`
//     ProductID uint    `json:"product_id"`
//     Quantity  int     `json:"quantity"`
//     User      User    `gorm:"foreignKey:UserID"`
//     Product   Product `gorm:"foreignKey:ProductID"`
// }
// type Address struct {
// 	gorm.Model
// 	UserID  uint   `json:"user_id"`
// 	Address string `json:"address"`
// 	City    string `json:"city"`
// 	ZipCode string `json:"zip_code"`
// 	User    User   `gorm:"foreignKey:UserID"`
// }
