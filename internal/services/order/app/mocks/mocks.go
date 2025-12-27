package mocks

import (
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	orderpb "github.com/Konscig/foodelivery-pet/internal/pb/orderpb"
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

func GenerateCreateOrderRequest(t *testing.T) *orderpb.CreateOrderRequest {
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
