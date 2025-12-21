package main

import (
	"context"
	"log"
	"os"

	kafka "github.com/Konscig/foodelivery-pet/internal/bootstrap/"
	deliveryApp "github.com/Konscig/foodelivery-pet/internal/services/delivery/app"
	redisClient "github.com/Konscig/foodelivery-pet/internal/services/delivery/redis"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	redis := redisClient.New(os.Getenv("REDIS_ADDR"))

	kafkaConsumer := kafka.NewConsumer(
		[]string{os.Getenv("KAFKA_BROKER")},
		kafka.TopicOrderReady, // слушаем order.ready
		"delivery-group",
	)

	kafkaProducer := kafka.NewProducer(
		[]string{os.Getenv("KAFKA_BROKER")},
	)

	publisher := deliveryApp.NewPublisher(kafkaProducer)

	deliveryConsumer := deliveryApp.NewConsumer(
		kafkaConsumer,
		redis,
		publisher,
	)

	log.Println("delivery service started")

	deliveryConsumer.Start(context.Background())
}
