package app_test

import (
	"context"
	"testing"

	"github.com/Konscig/foodelivery-pet/internal/pb/orderpb"
	"github.com/Konscig/foodelivery-pet/internal/services/order/app"
	"github.com/stretchr/testify/assert"
)

type mockRedis struct {
	status map[string]string
}

func (m *mockRedis) SetOrderStatus(key, status string) error {
	if m.status == nil {
		m.status = make(map[string]string)
	}
	m.status[key] = status
	return nil
}

type mockPublisher struct {
	called []struct {
		orderID, userID, restID string
		items                   []string
	}
	errToReturn error
}

func (m *mockPublisher) PublishOrderCreated(orderID, userID, restID string, items []string) error {
	if m.errToReturn != nil {
		return m.errToReturn
	}
	m.called = append(m.called, struct {
		orderID, userID, restID string
		items                   []string
	}{orderID, userID, restID, items})
	return nil
}

// --- Тесты ---

func TestCreateOrder_Success(t *testing.T) {
	redis := &mockRedis{}
	publisher := &mockPublisher{}

	service := app.NewService(redis, publisher)

	req := &orderpb.CreateOrderRequest{
		UserId: "user123",
		RestId: "rest456",
		Items: []*orderpb.OrderItem{
			{Name: "burger"},
			{Name: "fries"},
		},
	}

	resp, err := service.CreateOrder(context.Background(), req)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.OrderId)

	// Проверяем, что статус в Redis установлен
	key := "order:" + resp.OrderId + ":status"
	assert.Equal(t, "CREATED", redis.status[key])

	// Проверяем, что Publisher был вызван
	assert.Len(t, publisher.called, 1)
	call := publisher.called[0]
	assert.Equal(t, resp.OrderId, call.orderID)
	assert.Equal(t, "user123", call.userID)
	assert.Equal(t, "rest456", call.restID)
	assert.Equal(t, []string{"burger", "fries"}, call.items)
}

func TestCreateOrder_PublisherError(t *testing.T) {
	redis := &mockRedis{}
	publisher := &mockPublisher{errToReturn: assert.AnError}

	service := app.NewService(redis, publisher)

	req := &orderpb.CreateOrderRequest{
		UserId: "user123",
		RestId: "rest456",
		Items:  []*orderpb.OrderItem{{Name: "pizza"}},
	}

	resp, err := service.CreateOrder(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}
