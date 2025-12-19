package app

import (
	"context"
	"log"
	"time"

	"github.com/Konscig/foodelivery-pet/api/kafka"
	eventspb "github.com/Konscig/foodelivery-pet/generated/eventspb"
	redisClient "github.com/Konscig/foodelivery-pet/services/delivery/redis"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

type Consumer struct {
	kafkaConsumer *kafka.Consumer
	redis         *redisClient.Client
	publisher     *Publisher
}

func NewConsumer(
	kafkaConsumer *kafka.Consumer,
	redis *redisClient.Client,
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
		msg, err := c.kafkaConsumer.Reader.ReadMessage(ctx)
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

		// Сохраняем статус в Redis
		_ = c.redis.SetOrderStatus(event.OrderId, "COMING")

		// Публикуем order.coming
		if err := c.publisher.PublishOrderComing(event.OrderId, courierID); err != nil {
			log.Println("publish coming error:", err)
			continue
		}

		// Имитируем доставку
		time.Sleep(3 * time.Second)

		// Обновляем статус
		_ = c.redis.SetOrderStatus(event.OrderId, "DONE")

		// Публикуем order.done
		if err := c.publisher.PublishOrderDone(event.OrderId, courierID); err != nil {
			log.Println("publish done error:", err)
			continue
		}

		log.Printf("✅ order %s delivered by courier %s\n", event.OrderId, courierID)
	}
}
