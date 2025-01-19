package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string    `json:"username" gorm:"unique"`
	Password  string    `json:"password"`
	Email     string    `json:"email" gorm:"unique"`
	Phone     string    `json:"phone" gorm:"unique"`
	Role      string    `json:"role" gorm:"default:customer"` // Role field with default value "customer"
	Addresses []Address `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	Cart      []Cart    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}
