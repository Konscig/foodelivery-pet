package main

import (
	"context"
	"log"
	"os"

	"github.com/Konscig/foodelivery-pet/api/kafka"
	"github.com/Konscig/foodelivery-pet/services/restaurant/app"
	"github.com/Konscig/foodelivery-pet/services/restaurant/internal"
	"github.com/Konscig/foodelivery-pet/services/restaurant/internal/models"
	redisClient "github.com/Konscig/foodelivery-pet/services/restaurant/redis"
)

func main() {
	db, err := internal.InitDB()
	if err != nil {
		log.Fatal("postgres error:", err)
	}
	db.AutoMigrate(&models.Order{})

	redis := redisClient.New(os.Getenv("REDIS_ADDR"))

	kafkaConsumer := kafka.NewConsumer(
		[]string{os.Getenv("KAFKA_BROKER")},
		kafka.TopicOrderCreated,
		"restaurant-group",
	)

	kafkaProducer := kafka.NewProducer(
		[]string{os.Getenv("KAFKA_BROKER")},
		kafka.TopicOrderReady,
	)

	publisher := app.NewPublisher(kafkaProducer)

	restaurantConsumer := app.NewConsumer(
		kafkaConsumer,
		db,
		redis,
		publisher,
	)

	log.Println("🍽 restaurant service started")

	restaurantConsumer.Start(context.Background())
}
