package models

import "gorm.io/gorm"

type Delivery struct {
	gorm.Model
	ID        string `gorm:"primaryKey"`
	OrderID   string
	CourierID string
	Status    string
}
