package delivery

import (
	"context"
	"log"
	"os"

	"github.com/Konscig/foodelivery-pet/api/kafka"
	deliveryApp "github.com/Konscig/foodelivery-pet/services/delivery/app"
	"github.com/Konscig/foodelivery-pet/services/delivery/models"
	redisClient "github.com/Konscig/foodelivery-pet/services/delivery/redis"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(
		postgres.Open(os.Getenv("POSTGRES_DSN")),
		&gorm.Config{},
	)
	if err != nil {
		log.Fatal("postgres error:", err)
	}
	if err := db.AutoMigrate(&models.Delivery{}); err != nil {
		log.Fatal("migration error:", err)
	}

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
		db,
		redis,
		publisher,
	)

	log.Println("delivery service started")

	deliveryConsumer.Start(context.Background())
}
