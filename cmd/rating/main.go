package main

import (
	"context"
	"log"

	"github.com/Konscig/foodelivery-pet/config"
	"github.com/Konscig/foodelivery-pet/internal/bootstrap"
	ratingapp "github.com/Konscig/foodelivery-pet/internal/services/rating/app"
)

func main() {
	log.Println("rating service starting")
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	consumer := bootstrap.NewConsumer(cfg, "rating-group-done", bootstrap.TopicOrderDone)
	log.Println("rating consumer created")
	producer := bootstrap.NewProducer(cfg)
	log.Println("rating producer created")
	pg := bootstrap.InitPGStorage(cfg)

	publisher := ratingapp.NewPublisher(producer)
	ratingConsumer := ratingapp.NewConsumer(consumer, pg, publisher)
	log.Println("rating consumer initialized")

	// Запуск consumer в goroutine
	go ratingConsumer.Start(context.Background())

	// Keep the main goroutine alive
	select {}
}
