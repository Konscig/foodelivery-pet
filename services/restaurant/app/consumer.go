package app

import (
	"context"
	"log"
	"time"

	"github.com/Konscig/foodelivery-pet/api/kafka"
	eventspb "github.com/Konscig/foodelivery-pet/generated/eventspb"
	"github.com/Konscig/foodelivery-pet/services/restaurant/models"
	"github.com/Konscig/foodelivery-pet/services/restaurant/redis"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

type Consumer struct {
	consumer  *kafka.Consumer
	db        *gorm.DB
	redis     *redis.Client
	publisher *Publisher
}

func NewConsumer(
	consumer *kafka.Consumer,
	db *gorm.DB,
	redis *redis.Client,
	publisher *Publisher,
) *Consumer {
	return &Consumer{consumer, db, redis, publisher}
}

func (c *Consumer) Start(ctx context.Context) {
	for {
		msg, err := c.consumer.Reader.ReadMessage(ctx)
		if err != nil {
			log.Println("kafka read error:", err)
			continue
		}

		var event eventspb.OrderEvent
		if err := proto.Unmarshal(msg.Value, &event); err != nil {
			log.Println("unmarshal error:", err)
			continue
		}

		if event.Status != eventspb.OrderStatus_CREATED {
			continue
		}

		var payload eventspb.OrderCreatedPayload
		_ = proto.Unmarshal(event.Payload, &payload)

		log.Println("🍳 cooking order:", event.OrderId)

		time.Sleep(1 * time.Second) // ИМИТАЦИЯ ГОТОВКИ

		order := models.Order{
			ID:     event.OrderId,
			UserID: payload.UserId,
			RestID: payload.RestId,
			Status: "READY",
		}

		c.db.Save(&order)
		c.redis.SetOrderStatus(order.ID, "READY")

		_ = c.publisher.PublishOrderReady(order.ID, order.RestID)

		log.Println("✅ order ready:", order.ID)
	}
}
