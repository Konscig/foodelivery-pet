package app

import (
	"context"
	"log"

	"github.com/Konscig/foodelivery-pet/internal/bootstrap"
	eventspb "github.com/Konscig/foodelivery-pet/internal/pb/eventspb"
	"github.com/Konscig/foodelivery-pet/internal/storage"
	models "github.com/Konscig/foodelivery-pet/internal/storage/models"
	"google.golang.org/protobuf/proto"
)

type Consumer struct {
	consumer  bootstrap.EventConsumer
	service   *Service
	db        storage.Storage
	publisher OrderRatedPublisher
}

func NewConsumer(
	c bootstrap.EventConsumer,
	db storage.Storage,
	p OrderRatedPublisher,
) *Consumer {
	repo := db.(storage.ReviewRepository)
	service := NewService(repo)

	return &Consumer{
		consumer:  c,
		service:   service,
		db:        db,
		publisher: p,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	log.Println("rating consumer started")

	for {
		msg, err := c.consumer.ReadMessage(ctx)
		if err != nil {
			log.Println("rating read error:", err)
			continue
		}

		log.Println("rating received message for order:", string(msg.Key))

		var event eventspb.OrderEvent
		if err := proto.Unmarshal(msg.Value, &event); err != nil {
			log.Println("event unmarshal error:", err)
			continue
		}

		log.Println("rating event status:", event.Status)

		if event.Status != eventspb.OrderStatus_DONE {
			continue
		}

		var payload eventspb.OrderDonePayload
		if err := proto.Unmarshal(event.Payload, &payload); err != nil {
			continue
		}

		rating := uint32(3 + (event.EventId[0] % 3)) // 3–5 stars
		comment := "Good food!"

		order, err := c.db.GetOrder(event.OrderId)
		if err != nil {
			log.Println("get order error:", err)
			continue
		}
		restaurantID := order.RestID

		if err := c.service.AddReview(
			event.OrderId,
			restaurantID,
			rating,
			comment,
		); err != nil {
			log.Println("add review error:", err)
			continue
		}

		if err := c.publisher.PublishOrderRated(
			event.OrderId,
			uint8(rating),
			comment,
			restaurantID,
		); err != nil {
			log.Println("publish rated error:", err)
		}

		if err := c.db.UpdateOrderStatus(
			event.OrderId,
			models.StatusRated,
		); err != nil {
			log.Println("update order status error:", err)
		}
	}
}
