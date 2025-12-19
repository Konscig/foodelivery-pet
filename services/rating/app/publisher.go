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

func (p *Publisher) PublishOrderRated(orderID string, rating uint8, comment string) error {
	payload := &eventspb.OrderRatedPayload{
		Rating:  uint32(rating),
		Comment: comment,
	}

	payloadBytes, _ := proto.Marshal(payload)

	event := &eventspb.OrderEvent{
		EventId:   uuid.NewString(),
		OrderId:   orderID,
		Status:    eventspb.OrderStatus_RATED,
		Timestamp: time.Now().Unix(),
		Payload:   payloadBytes,
	}

	eventBytes, _ := proto.Marshal(event)

	return p.producer.SendProtoMessage(eventBytes)
}
