package app

import (
	"time"

	"github.com/Konscig/foodelivery-pet/internal/bootstrap"
	eventspb "github.com/Konscig/foodelivery-pet/internal/pb/eventspb"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

type Publisher struct {
	producer *bootstrap.Producer
}

func NewPublisher(p *bootstrap.Producer) *Publisher {
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

	return p.producer.SendProtoMessage(bootstrap.TopicOrderReady, eventBytes)
}
