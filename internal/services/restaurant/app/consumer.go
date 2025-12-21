package app

import (
	"context"
	"log"
	"time"

	"github.com/Konscig/foodelivery-pet/internal/bootstrap"
	eventspb "github.com/Konscig/foodelivery-pet/internal/pb/eventspb"
	"github.com/Konscig/foodelivery-pet/internal/services/restaurant/internal/models"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

type Consumer struct {
	kafkaConsumer *kafka.Consumer
	db            *gorm.DB
	redis         *bootstrap.RedisClient
	publisher     *Publisher
}

func NewConsumer(
	kafkaConsumer *kafka.Consumer,
	db *gorm.DB,
	redis *bootstrap.RedisClient,
	publisher *Publisher,
) *Consumer {
	return &Consumer{
		kafkaConsumer: kafkaConsumer,
		db:            db,
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

		// Интересует только CREATED
		if event.Status != eventspb.OrderStatus_CREATED {
			continue
		}

		var payload eventspb.OrderCreatedPayload
		if err := proto.Unmarshal(event.Payload, &payload); err != nil {
			log.Println("payload unmarshal error:", err)
			continue
		}

		log.Println("🍳 cooking order:", event.OrderId)

		time.Sleep(1 * time.Second) // имитация готовки

		order := models.Order{
			ID:     event.OrderId,
			UserID: payload.UserId,
			RestID: payload.RestId,
			Status: "READY",
		}

		// Сохраняем в базу с hash partitioning
		if err := c.db.Save(&order).Error; err != nil {
			log.Println("db save error:", err)
		}

		// Сохраняем статус в Redis
		c.redis.Set("order:"+order.ID+":status", "READY")

		// Публикуем order.ready
		if err := c.publisher.PublishOrderReady(order.ID, order.RestID); err != nil {
			log.Println("publish error:", err)
		}

		log.Println("✅ order ready:", order.ID)
	}
}
