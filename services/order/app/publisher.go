package app

import (
	"log"
	"time"

	eventspb "github.com/Konscig/foodelivery-pet/generated/eventspb"

	"github.com/Konscig/foodelivery-pet/api/kafka"
	"github.com/Konscig/foodelivery-pet/services/order/models"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

func PublishOrderCreated(producer *kafka.Producer, order *models.Order) error {
	payload := &eventspb.OrderCreatedPayload{
		UserId: order.UserID,
		RestId: order.RestID,
	}

	payloadBytes, _ := proto.Marshal(payload)

	event := &eventspb.OrderEvent{
		EventId:   uuid.NewString(),
		OrderId:   order.ID,
		Status:    eventspb.OrderStatus_CREATED,
		Timestamp: time.Now().Unix(),
		Payload:   payloadBytes,
	}

	eventBytes, _ := proto.Marshal(event)

	err := producer.SendProtoMessage(kafka.TopicOrderCreated, order.ID, eventBytes)
	if err != nil {
		log.Printf("failed to publish order.created: %v", err)
		return err
	}

	return nil
}
