package app

import (
	"log"
	"time"

	"github.com/Konscig/foodelivery-pet/internal/bootstrap"
	eventspb "github.com/Konscig/foodelivery-pet/internal/pb/eventspb"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

type Publisher struct {
	producer EventProducer
}

func NewPublisher(p EventProducer) *Publisher {
	return &Publisher{producer: p}
}

func (p *Publisher) PublishOrderCreated(
	orderID,
	userID,
	restID string,
	items []string,
) error {
	pbItems := make([]*eventspb.OrderItem, len(items))
	for i, name := range items {
		pbItems[i] = &eventspb.OrderItem{
			Name:     name,
			Quantity: 1,
		}
	}

	payload := &eventspb.OrderCreatedPayload{
		UserId: userID,
		RestId: restID,
		Items:  pbItems,
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		log.Printf("failed to marshal payload: %v", err)
		return err
	}

	event := &eventspb.OrderEvent{
		EventId:   uuid.NewString(),
		OrderId:   orderID,
		Status:    eventspb.OrderStatus_CREATED,
		Timestamp: time.Now().Unix(),
		Payload:   payloadBytes,
	}

	eventBytes, err := proto.Marshal(event)
	if err != nil {
		log.Printf("failed to marshal event: %v", err)
		return err
	}

	return p.producer.SendProtoMessage(
		bootstrap.TopicOrderCreated,
		[]byte(orderID),
		eventBytes,
	)
}
