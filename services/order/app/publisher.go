package app

import (
	"log"
	"time"

	"github.com/Konscig/foodelivery-pet/api/kafka"
	eventspb "github.com/Konscig/foodelivery-pet/generated/eventspb"
	"github.com/Konscig/foodelivery-pet/services/order/models"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

func PublishOrderCreated(p *kafka.Producer, order *models.Order) error {
	items := make([]*eventspb.OrderItem, len(order.Items))
	for i, it := range order.Items {
		items[i] = &eventspb.OrderItem{
			Name:     it.Name,
			Quantity: int32(it.Quantity),
		}
	}

	payload := &eventspb.OrderCreatedPayload{
		UserId: order.UserID,
		RestId: order.RestID,
		Items:  items,
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		log.Printf("failed to marshal payload: %v", err)
		return err
	}

	event := &eventspb.OrderEvent{
		EventId:   uuid.NewString(),
		OrderId:   order.ID,
		Status:    eventspb.OrderStatus_CREATED,
		Timestamp: time.Now().Unix(),
		Payload:   payloadBytes,
	}

	eventBytes, err := proto.Marshal(event)
	if err != nil {
		log.Printf("failed to marshal event: %v", err)
		return err
	}

	return p.SendProtoMessage(kafka.TopicOrderCreated, eventBytes)
}
