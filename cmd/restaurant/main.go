package main

import (
	"log"

	"github.com/Konscig/foodelivery-pet/config"
	"github.com/Konscig/foodelivery-pet/internal/bootstrap"
	restaurantapp "github.com/Konscig/foodelivery-pet/internal/services/restaurant/app"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	producer := bootstrap.NewProducer(cfg)
	consumer := bootstrap.NewConsumer(cfg, "restaurant-group", bootstrap.TopicOrderCreated)
	redis := bootstrap.NewRedis(cfg)
	pg := bootstrap.InitPGStorage(cfg)

	publisher := restaurantapp.NewPublisher(producer)
	restaurantConsumer := restaurantapp.NewConsumer(consumer, pg, redis, publisher)

	// Запуск gRPC сервера и регистрация gRPC-сервиса ресторана
	bootstrap.StartGRPCServer(cfg.GRPC.RestaurantPort, func(s *grpc.Server) {
		// TODO: Зарегистрировать gRPC-сервис ресторана, например:
		// orderpb.RegisterRestaurantServiceServer(s, restaurantConsumer)
	})
}
