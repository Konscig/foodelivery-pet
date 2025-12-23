package app_test

import (
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Konscig/foodelivery-pet/internal/services/rating/app"
	models "github.com/Konscig/foodelivery-pet/internal/storage/models"
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

func generateReviewData(t *testing.T) reviewTestData {
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

func TestAddReview_Success(t *testing.T) {
	repo := NewReviewRepoMock(t)
	service := app.NewService(repo)

	data := generateReviewData(t)

	repo.
		On(
			"CreateReview",
			mock.MatchedBy(func(r *models.Review) bool {
				return r.OrderID == data.OrderID &&
					r.RestaurantID == data.RestaurantID &&
					r.Rating == int32(data.Rating) &&
					r.Comment == data.Comment
			}),
		).
		Return(nil).
		Once()

	repo.
		On(
			"UpdateRestaurantStats",
			data.RestaurantID,
		).
		Return(nil).
		Once()

	err := service.AddReview(
		data.OrderID,
		data.RestaurantID,
		data.Rating,
		data.Comment,
	)

	t.Logf("результат теста err=%v", err)

	assert.NoError(t, err)

	repo.AssertExpectations(t)
}

func TestAddReview_CreateReviewError(t *testing.T) {
	repo := NewReviewRepoMock(t)
	service := app.NewService(repo)

	data := generateReviewData(t)

	repo.
		On(
			"CreateReview",
			mock.Anything,
		).
		Return(errors.New("failed to create review")).
		Once()

	err := service.AddReview(
		data.OrderID,
		data.RestaurantID,
		data.Rating,
		data.Comment,
	)

	t.Logf("ошибка теста err=%v", err)

	assert.Error(t, err)

	repo.AssertExpectations(t)
}

func TestAddReview_UpdateStatsError(t *testing.T) {
	repo := NewReviewRepoMock(t)
	service := app.NewService(repo)

	data := generateReviewData(t)

	repo.
		On(
			"CreateReview",
			mock.Anything,
		).
		Return(nil).
		Once()

	repo.
		On(
			"UpdateRestaurantStats",
			data.RestaurantID,
		).
		Return(errors.New("failed to update stats")).
		Once()

	err := service.AddReview(
		data.OrderID,
		data.RestaurantID,
		data.Rating,
		data.Comment,
	)

	t.Logf("ошибка теста err=%v", err)

	assert.Error(t, err)

	repo.AssertExpectations(t)
}
