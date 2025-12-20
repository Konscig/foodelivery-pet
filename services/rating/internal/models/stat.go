package models

import "gorm.io/gorm"

type RestaurantStats struct {
	gorm.Model
	RestaurantID  string `gorm:"uniqueIndex"`
	AvgRating     float64
	ReviewsCount  uint32
	WordCloudJSON string `gorm:"type:jsonb"`
}
