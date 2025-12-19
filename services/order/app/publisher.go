package app

import (
	"time"

	"github.com/Konscig/foodelivery-pet/api/kafka"
	eventspb "github.com/Konscig/foodelivery-pet/generated/eventspb"
	"github.com/Konscig/foodelivery-pet/services/order/models"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

func PublishOrderCreated(producer *kafka.Producer, order *models.Order) error {
	items := make([]string, len(order.Items))
	for i, it := range order.Items {
		items[i] = it.Name // или fmt.Sprintf("%s:%d", it.Name, it.Quantity)
	}

	payload := &eventspb.OrderCreatedPayload{
		UserId: order.UserID,
		RestId: order.RestID,
		Items:  items,
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
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
		return err
	}

	return producer.SendProtoMessage(kafka.TopicOrderCreated, order.ID, eventBytes)
}
