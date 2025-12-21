package app

import (
	"log"
	"time"

	eventspb "github.com/Konscig/foodelivery-pet/internal/pb/eventspb"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

// PublishOrderCreated публикует событие создания заказа в Kafka.
// TODO: заменить параметры на актуальные для вашего order-объекта.
func PublishOrderCreated(publish func([]byte) error, orderID, userID, restID string, items []string) error {
	pbItems := make([]*eventspb.OrderItem, len(items))
	for i, name := range items {
		pbItems[i] = &eventspb.OrderItem{
			Name:     name,
			Quantity: 1, // TODO: если есть количество, передавать его
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
	// publish — функция, отправляющая байты в Kafka (например, producer.Publish)
	return publish(eventBytes)
}
