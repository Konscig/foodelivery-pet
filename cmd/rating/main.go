package main

import (
	"log"

	"github.com/Konscig/foodelivery-pet/api/kafka"
	"github.com/Konscig/foodelivery-pet/config"
	"github.com/Konscig/foodelivery-pet/internal/services/rating/app"
	"github.com/Konscig/foodelivery-pet/internal/services/rating/internal"
	"github.com/Konscig/foodelivery-pet/internal/services/rating/internal/models"
	redisClient "github.com/Konscig/foodelivery-pet/internal/services/rating/redis"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	redis, err := redisClient.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatal(err)
	}

	db, err := internal.NewDB(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}

	kafkaProducer, err := kafka.NewProducer(cfg.Kafka)
	if err != nil {
		log.Fatal(err)
	}

	kafkaConsumer, err := kafka.NewConsumer(cfg.Kafka)
	if err != nil {
		log.Fatal(err)
	}

	ratingApp := app.NewRatingApp(redis, db, kafkaProducer, kafkaConsumer, models.NewStat())

	ratingApp.Start()
}
