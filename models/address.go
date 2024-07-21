package models

import "gorm.io/gorm"

type Address struct {
	gorm.Model
	UserID  uint   `json:"user_id"`
	Address string `json:"address"`
	City    string `json:"city"`
	ZipCode string `json:"zip_code"`
	User    User   `gorm:"foreignKey:UserID"`
}
