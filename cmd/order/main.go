package main

import (
	"log"

	"github.com/Konscig/foodelivery-pet/config"
	"github.com/Konscig/foodelivery-pet/internal/bootstrap"
	orderapp "github.com/Konscig/foodelivery-pet/internal/services/order/app"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	producer := bootstrap.NewProducer(cfg)
	consumer := bootstrap.NewConsumer(cfg, "order-group", bootstrap.TopicOrderCreated)
	redis := bootstrap.NewRedis(cfg)

	publisher := orderapp.NewPublisher(producer)
	orderConsumer := orderapp.NewConsumer(consumer, redis, publisher)

	// Запуск gRPC сервера и регистрация gRPC-сервиса заказов
	bootstrap.StartGRPCServer(cfg.GRPC.OrderPort, func(s *grpc.Server) {
		// TODO: Зарегистрировать gRPC-сервис заказов, например:
		// orderpb.RegisterOrderServiceServer(s, orderConsumer)
	})
}
