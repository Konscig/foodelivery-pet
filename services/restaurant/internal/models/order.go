package models

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	ID     string `gorm:"primaryKey"`
	UserID string
	RestID string
	Status string
}
