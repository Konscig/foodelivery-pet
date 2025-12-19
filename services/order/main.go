package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	orderpb "github.com/Konscig/foodelivery-pet/generated/orderpb"
	"github.com/Konscig/foodelivery-pet/services/order/app"
	"github.com/Konscig/foodelivery-pet/services/order/models"
	"github.com/Konscig/foodelivery-pet/services/order/redis"

	"github.com/Konscig/foodelivery-pet/api/kafka"
)

type server struct {
	orderpb.UnimplementedOrderServiceServer
	Producer *kafka.Producer
	DB       *gorm.DB
	Redis    *redis.RedisClient
}

func (s *server) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	orderID := uuid.NewString()

	newOrder := &models.Order{
		ID:     orderID,
		UserID: req.UserId,
		RestID: req.RestId,
		Status: models.StatusCreated,
	}

	if err := s.DB.Create(newOrder).Error; err != nil {
		return nil, err
	}

	for _, item := range req.Items {
		it := models.OrderItem{
			OrderID:  orderID,
			Name:     item.Name,
			Quantity: item.Quantity,
		}
		s.DB.Create(&it)
	}

	if err := app.PublishOrderCreated(s.Producer, newOrder); err != nil {
		log.Printf("failed to publish order.created: %v", err)
	}

	s.Redis.SetOrderStatus(orderID, string(models.StatusCreated))

	return &orderpb.CreateOrderResponse{
		OrderId: orderID,
	}, nil
}

func main() {
	godotenv.Load(".env")

	producer := kafka.NewProducer([]string{os.Getenv("KAFKA_BROKER")}, "order-events")
	defer producer.Close()

	redisClient := redis.NewRedis(os.Getenv("REDIS_ADDR"))

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	orderpb.RegisterOrderServiceServer(s, &server{
		Producer: producer,
		Redis:    redisClient,
	})

	log.Println("Order service running on :8081")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
