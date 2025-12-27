package app

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/Konscig/foodelivery-pet/internal/services/rating/app/mocks"
	models "github.com/Konscig/foodelivery-pet/internal/storage/models"
)

type ReviewServiceSuite struct {
	suite.Suite
	repo *mocks.ReviewRepoMock
	svc  *Service
}

func (s *ReviewServiceSuite) SetupTest() {
	s.repo = mocks.NewReviewRepoMock(s.T())
	s.svc = NewService(s.repo)
}

func (s *ReviewServiceSuite) TestAddReview_Success() {
	data := mocks.GenerateReviewData(s.T())

	s.repo.
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

	s.repo.
		On(
			"UpdateRestaurantStats",
			data.RestaurantID,
		).
		Return(nil).
		Once()

	err := s.svc.AddReview(
		data.OrderID,
		data.RestaurantID,
		data.Rating,
		data.Comment,
	)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *ReviewServiceSuite) TestAddReview_CreateReviewError() {
	data := mocks.GenerateReviewData(s.T())

	s.repo.
		On(
			"CreateReview",
			mock.Anything,
		).
		Return(errors.New("failed to create review")).
		Once()

	err := s.svc.AddReview(
		data.OrderID,
		data.RestaurantID,
		data.Rating,
		data.Comment,
	)

	s.Error(err)
	s.repo.AssertExpectations(s.T())
}

func (s *ReviewServiceSuite) TestAddReview_UpdateStatsError() {
	data := mocks.GenerateReviewData(s.T())

	s.repo.
		On(
			"CreateReview",
			mock.Anything,
		).
		Return(nil).
		Once()

	s.repo.
		On(
			"UpdateRestaurantStats",
			data.RestaurantID,
		).
		Return(errors.New("failed to update stats")).
		Once()

	err := s.svc.AddReview(
		data.OrderID,
		data.RestaurantID,
		data.Rating,
		data.Comment,
	)

	s.Error(err)
	s.repo.AssertExpectations(s.T())
}

func TestReviewServiceSuite(t *testing.T) {
	suite.Run(t, new(ReviewServiceSuite))
}
