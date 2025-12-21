package app

import (
	"context"
	"log"

	"github.com/Konscig/foodelivery-pet/internal/bootstrap"
	eventspb "github.com/Konscig/foodelivery-pet/internal/pb/eventspb"
	models "github.com/Konscig/foodelivery-pet/internal/storage/models"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

type Consumer struct {
	consumer *bootstrap.Consumer
	service  *Service
}

func NewConsumer(c *bootstrap.Consumer, db *gorm.DB) *Consumer {
	repo := models.NewReviewRepository(db)
	service := NewService(repo)
	return &Consumer{
		consumer: c,
		service:  service,
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
			continue
		}

		if event.Status != eventspb.OrderStatus_RATED {
			continue
		}

		var payload eventspb.OrderRatedPayload
		if err := proto.Unmarshal(event.Payload, &payload); err != nil {
			continue
		}

		err = c.service.AddReview(
			event.OrderId,
			payload.RestaurantId,
			payload.Rating,
			payload.Comment,
		)

		if err != nil {
			log.Println("add review error:", err)
		}
	}
}
