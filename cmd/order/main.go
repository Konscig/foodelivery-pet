package main

import (
	"log"

	"github.com/Konscig/foodelivery-pet/config"
	"github.com/Konscig/foodelivery-pet/internal/bootstrap"
	"github.com/Konscig/foodelivery-pet/internal/pb/orderpb"
	orderapp "github.com/Konscig/foodelivery-pet/internal/services/order/app"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	producer := bootstrap.NewProducer(cfg)
	redis := bootstrap.NewRedis(cfg)

	publisher := orderapp.NewPublisher(producer)
	orderService := orderapp.NewService(redis, publisher)

	// Запуск gRPC сервера и регистрация gRPC-сервиса заказов
	bootstrap.StartGRPCServer(cfg.GRPC.OrderPort, func(s *grpc.Server) {
		orderpb.RegisterOrderServiceServer(s, orderService)
		reflection.Register(s)
	})
}
