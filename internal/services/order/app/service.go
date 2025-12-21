package app

import (
	"context"
	"log"

	"github.com/Konscig/foodelivery-pet/internal/bootstrap"
	"github.com/Konscig/foodelivery-pet/internal/pb/orderpb"
	"github.com/google/uuid"
)

type Service struct {
	orderpb.UnimplementedOrderServiceServer
	redis     *bootstrap.RedisClient
	publisher *Publisher
}

func NewService(redis *bootstrap.RedisClient, publisher *Publisher) *Service {
	return &Service{
		redis:     redis,
		publisher: publisher,
	}
}

func (s *Service) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	orderID := uuid.NewString()

	// Сохранить в Redis
	s.redis.SetOrderStatus("order:"+orderID+":status", "CREATED")

	// Опубликовать событие
	items := make([]string, len(req.Items))
	for i, item := range req.Items {
		items[i] = item.Name
	}
	if err := s.publisher.PublishOrderCreated(orderID, req.UserId, req.RestId, items); err != nil {
		log.Printf("failed to publish order created: %v", err)
		return nil, err
	}

	return &orderpb.CreateOrderResponse{OrderId: orderID}, nil
}
