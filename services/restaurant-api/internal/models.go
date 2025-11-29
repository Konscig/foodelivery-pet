package internal

import (
	"gorm.io/gorm"
)

type Restaurant struct {
	gorm.Model
	rest_id int    `gorm:"primaryKey"`
	name    string `gorm:"not null"`
	cords   [2]int `gorm:"not null"`
}
