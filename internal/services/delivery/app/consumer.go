package app

// Пакет app отвечает за реализацию бизнес-логики конкретного сервиса.
// Здесь находятся обработчики событий, взаимодействие с базами данных и внешними сервисами.

import (
	"context"
	"log"
	"time"

	"github.com/Konscig/foodelivery-pet/internal/bootstrap"
	eventspb "github.com/Konscig/foodelivery-pet/internal/pb/eventspb"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

type Consumer struct {
	kafkaConsumer *bootstrap.Consumer
	redis         *bootstrap.RedisClient
	publisher     *Publisher
}

func NewConsumer(
	kafkaConsumer *bootstrap.Consumer,
	redis *bootstrap.RedisClient,
	publisher *Publisher,
) *Consumer {
	return &Consumer{
		kafkaConsumer: kafkaConsumer,
		redis:         redis,
		publisher:     publisher,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	for {
		msg, err := c.kafkaConsumer.ReadMessage(ctx)
		if err != nil {
			log.Println("kafka read error:", err)
			continue
		}

		var event eventspb.OrderEvent
		if err := proto.Unmarshal(msg.Value, &event); err != nil {
			log.Println("event unmarshal error:", err)
			continue
		}

		if event.Status != eventspb.OrderStatus_READY {
			continue
		}

		var payload eventspb.OrderReadyPayload
		if err := proto.Unmarshal(event.Payload, &payload); err != nil {
			log.Println("payload unmarshal error:", err)
			continue
		}

		courierID := uuid.NewString()
		log.Printf("🚴 courier %s assigned to order %s\n", courierID, event.OrderId)

		_ = c.redis.SetOrderStatus("order:"+event.OrderId+":status", "COMING")

		if err := c.publisher.PublishOrderComing(event.OrderId, courierID); err != nil {
			log.Println("publish coming error:", err)
			continue
		}

		time.Sleep(3 * time.Second)

		_ = c.redis.SetOrderStatus("order:"+event.OrderId+":status", "DONE")

		if err := c.publisher.PublishOrderDone(event.OrderId, courierID); err != nil {
			log.Println("publish done error:", err)
			continue
		}

		log.Printf("✅ order %s delivered by courier %s\n", event.OrderId, courierID)
	}
}
