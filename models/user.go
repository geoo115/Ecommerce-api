package models

import (
    "gorm.io/gorm"
)

type User struct {
    gorm.Model
    Username  string    `json:"username" gorm:"unique"`
    Password  string    `json:"password"`
    Addresses []Address `gorm:"foreignKey:UserID"`
    Cart      []Cart    `gorm:"foreignKey:UserID"`
}
