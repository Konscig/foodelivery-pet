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

	return p.producer.SendProtoMessage(bootstrap.TopicOrderRated, eventBytes)
}
