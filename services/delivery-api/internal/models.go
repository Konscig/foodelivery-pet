package internal

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	user_id int `gorm:"primaryKey"`
}
