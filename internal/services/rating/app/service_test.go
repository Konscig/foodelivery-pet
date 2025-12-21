package app_test

import (
	"errors"
	"testing"

	"github.com/Konscig/foodelivery-pet/internal/services/rating/app"
	models "github.com/Konscig/foodelivery-pet/internal/storage/models"
	"github.com/stretchr/testify/assert"
)

type mockReviewRepo struct {
	createdReviews []models.Review
	updateCalled   []string

	createErr bool
	updateErr bool
}

func (m *mockReviewRepo) CreateReview(review *models.Review) error {
	if m.createErr {
		return errors.New("failed to create review")
	}
	m.createdReviews = append(m.createdReviews, *review)
	return nil
}

func (m *mockReviewRepo) UpdateRestaurantStats(restaurantID string) error {
	if m.updateErr {
		return errors.New("failed to update stats")
	}
	m.updateCalled = append(m.updateCalled, restaurantID)
	return nil
}

func (m *mockReviewRepo) GetRestaurantStats(restaurantID string) (*models.RestaurantStats, error) {
	return &models.RestaurantStats{}, nil
}

func TestAddReview_Success(t *testing.T) {
	repo := &mockReviewRepo{}
	service := app.NewService(repo)

	orderID := "order123"
	restaurantID := "rest456"
	rating := uint32(5)
	comment := "Excellent!"

	err := service.AddReview(orderID, restaurantID, rating, comment)
	assert.NoError(t, err)

	assert.Len(t, repo.createdReviews, 1)
	assert.Equal(t, orderID, repo.createdReviews[0].OrderID)
	assert.Equal(t, restaurantID, repo.createdReviews[0].RestaurantID)
	assert.Equal(t, int32(rating), repo.createdReviews[0].Rating)
	assert.Equal(t, comment, repo.createdReviews[0].Comment)

	assert.Len(t, repo.updateCalled, 1)
	assert.Equal(t, restaurantID, repo.updateCalled[0])
}

func TestAddReview_CreateReviewError(t *testing.T) {
	repo := &mockReviewRepo{createErr: true}
	service := app.NewService(repo)

	err := service.AddReview("order123", "rest456", 4, "Good")
	assert.Error(t, err)
}

func TestAddReview_UpdateStatsError(t *testing.T) {
	repo := &mockReviewRepo{updateErr: true}
	service := app.NewService(repo)

	err := service.AddReview("order123", "rest456", 4, "Good")
	assert.Error(t, err)
}
