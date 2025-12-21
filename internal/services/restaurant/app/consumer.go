package app

import (
	"context"
	"log"
	"time"

	"github.com/Konscig/foodelivery-pet/internal/bootstrap"
	eventspb "github.com/Konscig/foodelivery-pet/internal/pb/eventspb"
	"github.com/Konscig/foodelivery-pet/internal/pb/orderpb"
	"github.com/Konscig/foodelivery-pet/internal/storage"
	models "github.com/Konscig/foodelivery-pet/internal/storage/models"
	"google.golang.org/protobuf/proto"
)

type Consumer struct {
	consumer  bootstrap.EventConsumer
	db        storage.Storage
	redis     OrderStatusStore
	publisher OrderReadyPublisher
}

func NewConsumer(
	consumer bootstrap.EventConsumer,
	db storage.Storage,
	redis OrderStatusStore,
	publisher OrderReadyPublisher,
) *Consumer {
	return &Consumer{
		consumer:  consumer,
		db:        db,
		redis:     redis,
		publisher: publisher,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	for {
		msg, err := c.consumer.ReadMessage(ctx)
		if err != nil {
			log.Println("kafka read error:", err)
			continue
		}

		var event eventspb.OrderEvent
		if err := proto.Unmarshal(msg.Value, &event); err != nil {
			log.Println("event unmarshal error:", err)
			continue
		}

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
			Status: models.StatusReady,
			Items:  convertItems(payload.Items),
		}

		// Сохраняем в базу
		if err := c.db.AddOrder(&order); err != nil {
			log.Println("db save error:", err)
		}

		// Сохраняем статус в Redis
		if err := c.redis.SetOrderStatus("order:"+order.ID+":status", "READY"); err != nil {
			log.Println("redis set status error:", err)
		}

		// Публикуем order.ready
		if err := c.publisher.PublishOrderReady(order.ID, order.RestID); err != nil {
			log.Println("publish error:", err)
		}

		log.Println("✅ order ready:", order.ID)
	}
}

func convertItems(items []*eventspb.OrderItem) []*orderpb.OrderItem {
	result := make([]*orderpb.OrderItem, len(items))
	for i, item := range items {
		result[i] = &orderpb.OrderItem{
			Name:     item.Name,
			Quantity: int64(item.Quantity),
		}
	}
	return result
}
