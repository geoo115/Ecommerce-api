package models

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name     string    `json:"name" gorm:"unique"`
	Products []Product `gorm:"foreignKey:CategoryID"`
}

type Product struct {
	gorm.Model
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	CategoryID  uint      `json:"category_id"`
	Description string    `json:"description"`
	Category    Category  `json:"category" gorm:"foreignKey:CategoryID"`
	Cart        []Cart    `json:"-" gorm:"foreignKey:ProductID"` // Hide in JSON
	Inventory   Inventory `json:"inventory" gorm:"foreignKey:ProductID"`
}
