package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Role      string    `json:"role" gorm:"default:customer"` // Role field with default value "customer"
	Addresses []Address `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	Cart      []Cart    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}
