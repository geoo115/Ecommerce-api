package models

import "gorm.io/gorm"

type Review struct {
	gorm.Model
	ProductID uint    `json:"product_id"`
	UserID    uint    `json:"user_id"`
	Rating    int     `json:"rating"` // 1-5 stars
	Comment   string  `json:"comment"`
	Product   Product `gorm:"foreignKey:ProductID"`
	User      User    `gorm:"foreignKey:UserID"`
}
