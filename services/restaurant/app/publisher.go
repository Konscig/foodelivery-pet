package app

import (
	"time"

	"github.com/Konscig/foodelivery-pet/api/kafka"
	eventspb "github.com/Konscig/foodelivery-pet/generated/eventspb"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

type Publisher struct {
	producer *kafka.Producer
}

func NewPublisher(p *kafka.Producer) *Publisher {
	return &Publisher{producer: p}
}

// Публикуем order.ready
func (p *Publisher) PublishOrderReady(orderID, restaurantID string) error {
	payload := &eventspb.OrderReadyPayload{
		RestId: restaurantID,
	}

	payloadBytes, _ := proto.Marshal(payload)

	event := &eventspb.OrderEvent{
		EventId:   uuid.NewString(),
		OrderId:   orderID,
		Status:    eventspb.OrderStatus_READY,
		Timestamp: time.Now().Unix(),
		Payload:   payloadBytes,
	}

	eventBytes, _ := proto.Marshal(event)

	return p.producer.SendProtoMessage(kafka.TopicOrderReady, eventBytes)
}
