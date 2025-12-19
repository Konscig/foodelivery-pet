package main

import (
	"context"
	"log"
	"os"

	"github.com/Konscig/foodelivery-pet/api/kafka"
	deliveryApp "github.com/Konscig/foodelivery-pet/services/delivery/app"
	redisClient "github.com/Konscig/foodelivery-pet/services/delivery/redis"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	redis := redisClient.New(os.Getenv("REDIS_ADDR"))

	kafkaConsumer := kafka.NewConsumer(
		[]string{os.Getenv("KAFKA_BROKER")},
		kafka.TopicOrderReady,
		"delivery-group",
	)

	kafkaProducer := kafka.NewProducer(
		[]string{os.Getenv("KAFKA_BROKER")},
		"", // topic задаём при SendProtoMessage
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
