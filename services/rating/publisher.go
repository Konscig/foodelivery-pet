package api

import (
	"time"

	"github.com/Konscig/foodelivery-pet/api/kafka"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

func publishOrderCreated(producer *kafka.Producer, orderID string, userID string, restID string) error {
	payload := &eventspb.OrderCreatedPayload{
		UserId:       userID,
		RestaurantId: restID,
	}

	payloadBytes, _ := proto.Marshal(payload)

	event := &eventspb.OrderEvent{
		EventId:   uuid.NewString(),
		OrderId:   orderID,
		Status:    eventspb.OrderStatus_ORDER_STATUS_CREATED,
		Timestamp: time.Now().Unix(),
		Payload:   payloadBytes,
	}

	eventBytes, _ := proto.Marshal(event)

	return producer.SendProtoMessage(kafka.TopicOrderCreated, orderID, eventBytes)
}
