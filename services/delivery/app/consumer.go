package app

import (
	"context"
	"log"
	"time"

	"github.com/Konscig/foodelivery-pet/api/kafka"
	eventspb "github.com/Konscig/foodelivery-pet/generated/eventspb"
	"github.com/Konscig/foodelivery-pet/services/delivery/models"
	redisClient "github.com/Konscig/foodelivery-pet/services/delivery/redis"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

type Consumer struct {
	kafkaConsumer *kafka.Consumer
	db            *gorm.DB
	redis         *redisClient.Client
	publisher     *Publisher
}

func NewConsumer(
	kafkaConsumer *kafka.Consumer,
	db *gorm.DB,
	redis *redisClient.Client,
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

		// Нас интересует только READY
		if event.Status != eventspb.OrderStatus_READY {
			continue
		}

		var payload eventspb.OrderReadyPayload
		if err := proto.Unmarshal(event.Payload, &payload); err != nil {
			log.Println("payload unmarshal error:", err)
			continue
		}

		// 1️⃣ Назначаем курьера
		courierID := uuid.NewString()

		log.Printf("🚴 courier %s assigned to order %s\n", courierID, event.OrderId)

		// 2️⃣ Сохраняем доставку в БД
		delivery := models.Delivery{
			ID:        uuid.NewString(),
			OrderID:   event.OrderId,
			CourierID: courierID,
			Status:    "COMING",
		}

		if err := c.db.Create(&delivery).Error; err != nil {
			log.Println("db error:", err)
			continue
		}

		// 3️⃣ Пишем статус в Redis
		_ = c.redis.SetOrderStatus(event.OrderId, "COMING")

		// 4️⃣ Публикуем order.coming
		if err := c.publisher.PublishOrderComing(event.OrderId, courierID); err != nil {
			log.Println("publish coming error:", err)
			continue
		}

		// 5️⃣ Имитируем доставку
		time.Sleep(3 * time.Second)

		// 6️⃣ Обновляем статус
		delivery.Status = "DONE"
		c.db.Save(&delivery)
		_ = c.redis.SetOrderStatus(event.OrderId, "DONE")

		// 7️⃣ Публикуем order.done
		if err := c.publisher.PublishOrderDone(event.OrderId, courierID); err != nil {
			log.Println("publish done error:", err)
			continue
		}

		log.Printf("✅ order %s delivered by courier %s\n", event.OrderId, courierID)
	}
}
