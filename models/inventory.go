package models

import "gorm.io/gorm"

type Inventory struct {
	gorm.Model
	ProductID uint     `json:"product_id"`
	Stock     int      `json:"stock"`
	Product   *Product `json:"-" gorm:"foreignKey:ProductID"`
}
