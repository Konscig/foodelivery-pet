package models

import "gorm.io/gorm"

type Review struct {
	gorm.Model
	ID           string `gorm:"primaryKey"`
	OrderID      string `gorm:"uniqueIndex"`
	RestaurantID string `gorm:"index"`
	Rating       uint32
	Comment      string
}
