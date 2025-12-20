package app

import (
	"encoding/json"
	"strings"

	"github.com/Konscig/foodelivery-pet/services/rating/internal/models"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) AddReview(
	orderID, restaurantID string,
	rating uint32,
	comment string,
) error {

	review := models.Review{
		ID:           orderID,
		OrderID:      orderID,
		RestaurantID: restaurantID,
		Rating:       rating,
		Comment:      comment,
	}

	if err := s.db.Create(&review).Error; err != nil {
		return err
	}

	return s.updateStats(restaurantID)
}

func (s *Service) updateStats(restaurantID string) error {
	var reviews []models.Review
	if err := s.db.
		Where("restaurant_id = ?", restaurantID).
		Find(&reviews).Error; err != nil {
		return err
	}

	var sum uint32
	var comments []string
	for _, r := range reviews {
		sum += r.Rating
		comments = append(comments, r.Comment)
	}

	avg := float64(sum) / float64(len(reviews))
	wordCloud := BuildWordCloud(strings.Join(comments, " "))

	wcJSON, _ := json.Marshal(wordCloud)

	stats := models.RestaurantStats{
		RestaurantID:  restaurantID,
		AvgRating:     avg,
		ReviewsCount:  uint32(len(reviews)),
		WordCloudJSON: string(wcJSON),
	}

	return s.db.
		Where("restaurant_id = ?", restaurantID).
		Assign(stats).
		FirstOrCreate(&stats).Error
}
