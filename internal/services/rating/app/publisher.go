package app

import (
	"time"

	"github.com/Konscig/foodelivery-pet/internal/bootstrap"
	eventspb "github.com/Konscig/foodelivery-pet/internal/pb/eventspb"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

type Publisher struct {
	producer bootstrap.EventProducer
}

func NewPublisher(p bootstrap.EventProducer) *Publisher {
	return &Publisher{producer: p}
}

func (p *Publisher) PublishOrderRated(
	orderID string,
	rating uint8,
	comment string,
	restaurantID string,
) error {
	payload := &eventspb.OrderRatedPayload{
		Rating:       uint32(rating),
		Comment:      comment,
		RestaurantId: restaurantID,
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

	return p.producer.SendProtoMessage(
		bootstrap.TopicOrderRated,
		[]byte(orderID),
		eventBytes,
	)
}
