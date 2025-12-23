package app_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	orderpb "github.com/Konscig/foodelivery-pet/internal/pb/orderpb"
	"github.com/Konscig/foodelivery-pet/internal/services/order/app"
)

type RedisMock struct {
	mock.Mock
	t *testing.T
}

func NewRedisMock(t *testing.T) *RedisMock {
	return &RedisMock{t: t}
}

func (m *RedisMock) SetOrderStatus(key, status string) error {
	m.t.Logf("мок редис статуса SetOrderStatus key=%s status=%s", key, status)
	args := m.Called(key, status)
	return args.Error(0)
}

type PublisherMock struct {
	mock.Mock
	t *testing.T
}

func NewPublisherMock(t *testing.T) *PublisherMock {
	return &PublisherMock{t: t}
}

func (m *PublisherMock) PublishOrderCreated(
	orderID, userID, restID string,
	items []string,
) error {
	m.t.Logf(
		"мок кафка паблишера PublishOrderCreated orderID=%s userID=%s restID=%s items=%v",
		orderID, userID, restID, items,
	)
	args := m.Called(orderID, userID, restID, items)
	return args.Error(0)
}

func generateCreateOrderRequest(t *testing.T) *orderpb.CreateOrderRequest {
	rand.Seed(time.Now().UnixNano())

	n := rand.Intn(4) + 1
	items := make([]*orderpb.OrderItem, n)
	names := make([]string, 0, n)

	for i := 0; i < n; i++ {
		name := "item-" + uuid.NewString()[:6]
		items[i] = &orderpb.OrderItem{Name: name}
		names = append(names, name)
	}

	req := &orderpb.CreateOrderRequest{
		UserId: "user-" + uuid.NewString(),
		RestId: "rest-" + uuid.NewString(),
		Items:  items,
	}

	t.Logf(
		"мокнутые данные user=%s rest=%s items=%v",
		req.UserId,
		req.RestId,
		names,
	)

	return req
}

func TestCreateOrder_Success(t *testing.T) {
	redis := NewRedisMock(t)
	publisher := NewPublisherMock(t)

	service := app.NewService(redis, publisher)

	req := generateCreateOrderRequest(t)

	expectedItems := make([]string, 0, len(req.Items))
	for _, it := range req.Items {
		expectedItems = append(expectedItems, it.Name)
	}

	publisher.
		On(
			"PublishOrderCreated",
			mock.Anything,
			req.UserId,
			req.RestId,
			expectedItems,
		).
		Return(nil).
		Once()

	redis.
		On(
			"SetOrderStatus",
			mock.Anything,
			"CREATED",
		).
		Return(nil).
		Once()

	resp, err := service.CreateOrder(context.Background(), req)

	t.Logf("результат теста resp=%+v err=%v", resp, err)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.OrderId)

	publisher.AssertExpectations(t)
	redis.AssertExpectations(t)
}

func TestCreateOrder_PublisherError(t *testing.T) {
	redis := NewRedisMock(t)
	publisher := NewPublisherMock(t)

	service := app.NewService(redis, publisher)

	req := generateCreateOrderRequest(t)

	expectedItems := make([]string, 0, len(req.Items))
	for _, it := range req.Items {
		expectedItems = append(expectedItems, it.Name)
	}
	redis.
		On(
			"SetOrderStatus",
			mock.Anything,
			"CREATED",
		).
		Return(nil).
		Once()

	publisher.
		On(
			"PublishOrderCreated",
			mock.Anything,
			req.UserId,
			req.RestId,
			expectedItems,
		).
		Return(assert.AnError).
		Once()

	resp, err := service.CreateOrder(context.Background(), req)

	t.Logf("ошибка теста resp=%+v err=%v", resp, err)

	assert.Error(t, err)
	assert.Nil(t, resp)

	redis.AssertExpectations(t)
	publisher.AssertExpectations(t)
}
