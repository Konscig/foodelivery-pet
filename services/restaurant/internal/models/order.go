package models

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	ID     string `gorm:"primaryKey"`
	UserID string `gorm:"column:user_id"`
	RestID string `gorm:"column:restaurant_id"`
	Status string `gorm:"column:status"`
}
