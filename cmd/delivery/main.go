package main

import (
	"log"

	"github.com/Konscig/foodelivery-pet/config"
	"github.com/Konscig/foodelivery-pet/internal/bootstrap"
	deliveryapp "github.com/Konscig/foodelivery-pet/internal/services/delivery/app"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	producer := bootstrap.NewProducer(cfg)
	consumer := bootstrap.NewConsumer(cfg, "delivery-group", bootstrap.TopicOrderReady)
	redis := bootstrap.NewRedis(cfg)

	publisher := deliveryapp.NewPublisher(producer)
	deliveryConsumer := deliveryapp.NewConsumer(consumer, redis, publisher)

	// Запуск gRPC сервера и регистрация gRPC-сервиса доставки
	bootstrap.StartGRPCServer(cfg.GRPC.DeliveryPort, func(s *grpc.Server) {
		// TODO: Зарегистрировать gRPC-сервис доставки, например:
		// orderpb.RegisterDeliveryServiceServer(s, deliveryConsumer)
	})
}
