package main

import (
	"log"

	"github.com/Konscig/foodelivery-pet/config"
	"github.com/Konscig/foodelivery-pet/internal/bootstrap"
	ratingapp "github.com/Konscig/foodelivery-pet/internal/services/rating/app"
	"github.com/Konscig/foodelivery-pet/internal/services/rating/models"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	producer := bootstrap.NewProducer(cfg)
	consumer := bootstrap.NewConsumer(cfg, "rating-group", bootstrap.TopicOrderRated)
	redis := bootstrap.NewRedis(cfg)
	pg := bootstrap.InitPGStorage(cfg)

	stat := models.NewStat()
	publisher := ratingapp.NewPublisher(producer)
	ratingConsumer := ratingapp.NewConsumer(consumer, redis, pg, publisher, stat)

	// Запуск gRPC сервера и регистрация gRPC-сервиса рейтинга
	bootstrap.StartGRPCServer(cfg.GRPC.RatingPort, func(s *grpc.Server) {
		// TODO: Зарегистрировать gRPC-сервис рейтинга, например:
		// orderpb.RegisterRatingServiceServer(s, ratingConsumer)
	})
}
