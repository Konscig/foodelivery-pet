package app

import (
	"github.com/Konscig/foodelivery-pet/internal/storage"
	models "github.com/Konscig/foodelivery-pet/internal/storage/models"
)

// Refactored Service to use a repository for database operations
type Service struct {
	repo storage.ReviewRepository
}

func NewService(repo storage.ReviewRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) AddReview(orderID, restaurantID string, rating uint32, comment string) error {
	review := models.Review{
		ID:           orderID,
		OrderID:      orderID,
		RestaurantID: restaurantID,
		Rating:       int32(rating), // Adjusted type
		Comment:      comment,
	}

	if err := s.repo.CreateReview(&review); err != nil {
		return err
	}

	return s.repo.UpdateRestaurantStats(restaurantID)
}
