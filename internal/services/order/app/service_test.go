package app

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/Konscig/foodelivery-pet/internal/services/order/app/mocks"
)

type OrderServiceSuite struct {
	suite.Suite
	redis     *mocks.RedisMock
	publisher *mocks.PublisherMock
	svc       *Service
}

func (s *OrderServiceSuite) SetupTest() {
	s.redis = mocks.NewRedisMock(s.T())
	s.publisher = mocks.NewPublisherMock(s.T())
	s.svc = NewService(s.redis, s.publisher)
}

func (s *OrderServiceSuite) TestCreateOrder_Success() {
	req := mocks.GenerateCreateOrderRequest(s.T())

	expectedItems := make([]string, 0, len(req.Items))
	for _, it := range req.Items {
		expectedItems = append(expectedItems, it.Name)
	}

	s.redis.
		On(
			"SetOrderStatus",
			mock.Anything,
			"CREATED",
		).
		Return(nil).
		Once()

	s.publisher.
		On(
			"PublishOrderCreated",
			mock.Anything,
			req.UserId,
			req.RestId,
			expectedItems,
		).
		Return(nil).
		Once()

	resp, err := s.svc.CreateOrder(context.Background(), req)

	s.NoError(err)
	s.NotNil(resp)
	s.NotEmpty(resp.OrderId)

	s.redis.AssertExpectations(s.T())
	s.publisher.AssertExpectations(s.T())
}

func (s *OrderServiceSuite) TestCreateOrder_PublisherError() {
	req := mocks.GenerateCreateOrderRequest(s.T())

	expectedItems := make([]string, 0, len(req.Items))
	for _, it := range req.Items {
		expectedItems = append(expectedItems, it.Name)
	}

	s.redis.
		On(
			"SetOrderStatus",
			mock.Anything,
			"CREATED",
		).
		Return(nil).
		Once()

	s.publisher.
		On(
			"PublishOrderCreated",
			mock.Anything,
			req.UserId,
			req.RestId,
			expectedItems,
		).
		Return(assert.AnError).
		Once()

	resp, err := s.svc.CreateOrder(context.Background(), req)

	s.Error(err)
	s.Nil(resp)

	s.redis.AssertExpectations(s.T())
	s.publisher.AssertExpectations(s.T())
}

func TestOrderServiceSuite(t *testing.T) {
	suite.Run(t, new(OrderServiceSuite))
}
