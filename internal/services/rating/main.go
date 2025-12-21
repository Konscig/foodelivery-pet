package main

import (
	"context"
	"log"
	"os"

	"github.com/Konscig/foodelivery-pet/api/kafka"
	"github.com/Konscig/foodelivery-pet/internal/services/rating/app"
	"github.com/Konscig/foodelivery-pet/internal/services/rating/internal"
	"github.com/Konscig/foodelivery-pet/internal/services/rating/internal/models"
)

func main() {
	db, err := internal.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(
		&models.Review{},
		&models.RestaurantStats{},
	)

	consumer := kafka.NewConsumer(
		[]string{os.Getenv("KAFKA_BROKER")},
		"order.rated",
		"rating-group",
	)

	service := app.NewService(db)
	ratingConsumer := app.NewConsumer(consumer, service)

	log.Println("⭐ rating service started")
	ratingConsumer.Start(context.Background())
}
