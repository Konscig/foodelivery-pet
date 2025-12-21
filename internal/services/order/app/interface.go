package app

import "context"

type OrderService interface {
	Start()
}

type Message struct {
	Key   []byte
	Value []byte
}

type EventConsumer interface {
	ReadMessage(ctx context.Context) (Message, error)
}

type EventProducer interface {
	SendProtoMessage(topic string, key []byte, value []byte) error
}

type OrderCreatedPublisher interface {
	PublishOrderCreated(orderID, userID, restID string, items []string) error
}

type OrderStatusStore interface {
	SetOrderStatus(key string, status string) error
}
