package mocks

import (
	"math/rand"
	"testing"
	"time"

	models "github.com/Konscig/foodelivery-pet/internal/storage/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type ReviewRepoMock struct {
	mock.Mock
	t *testing.T
}

func NewReviewRepoMock(t *testing.T) *ReviewRepoMock {
	return &ReviewRepoMock{t: t}
}

func (m *ReviewRepoMock) CreateReview(review *models.Review) error {
	m.t.Logf(
		"мок отзыва CreateReview orderID=%s restaurantID=%s rating=%d comment=%q",
		review.OrderID,
		review.RestaurantID,
		review.Rating,
		review.Comment,
	)
	args := m.Called(review)
	return args.Error(0)
}

func (m *ReviewRepoMock) UpdateRestaurantStats(restaurantID string) error {
	m.t.Logf(
		"мок отзыва UpdateRestaurantStats restaurantID=%s",
		restaurantID,
	)
	args := m.Called(restaurantID)
	return args.Error(0)
}

func (m *ReviewRepoMock) GetRestaurantStats(restaurantID string) (*models.RestaurantStats, error) {
	m.t.Logf(
		"мок отзывов по ресторану GetRestaurantStats restaurantID=%s",
		restaurantID,
	)
	args := m.Called(restaurantID)
	return args.Get(0).(*models.RestaurantStats), args.Error(1)
}

type reviewTestData struct {
	OrderID      string
	RestaurantID string
	Rating       uint32
	Comment      string
}

func GenerateReviewData(t *testing.T) reviewTestData {
	rand.Seed(time.Now().UnixNano())

	rating := uint32(rand.Intn(5) + 1)

	data := reviewTestData{
		OrderID:      "order-" + uuid.NewString(),
		RestaurantID: "rest-" + uuid.NewString(),
		Rating:       rating,
		Comment:      "comment-" + uuid.NewString()[:8],
	}

	t.Logf(
		"данные для теста отзыва orderID=%s restaurantID=%s rating=%d comment=%q",
		data.OrderID,
		data.RestaurantID,
		data.Rating,
		data.Comment,
	)

	return data
}
