package models

import (
	"gorm.io/gorm"
)

type Rating struct {
	gorm.Model
	ID      string `gorm:"primaryKey"`
	OrderID string
	UserID  string
	Rating  uint8
	Comment string
}
