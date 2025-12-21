package app

import (
	"context"
	"log"

	eventspb "github.com/Konscig/foodelivery-pet/internal/pb/eventspb"
	"google.golang.org/protobuf/proto"
)

type Consumer struct {
	consumer *kafka.Consumer
	service  *Service
}

func NewConsumer(c *kafka.Consumer, s *Service) *Consumer {
	return &Consumer{
		consumer: c,
		service:  s,
	}
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
