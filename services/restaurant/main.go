package restaurant

import (
	"context"
	"log"
	"os"

	"github.com/Konscig/foodelivery-pet/api/kafka"
	"github.com/Konscig/foodelivery-pet/services/restaurant/app"
	"github.com/Konscig/foodelivery-pet/services/restaurant/models"
	redisClient "github.com/Konscig/foodelivery-pet/services/restaurant/redis"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(postgres.Open(os.Getenv("POSTGRES_DSN")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
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
