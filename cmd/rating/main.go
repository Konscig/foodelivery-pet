package main

import (
	"log"

	"github.com/Konscig/foodelivery-pet/config"
	"github.com/Konscig/foodelivery-pet/internal/bootstrap"
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

	kafkaProducer, err := bootstrap.NewProducer(cfg.Kafka)
	if err != nil {
		log.Fatal(err)
	}

	kafkaConsumer, err := bootstrap.NewConsumer(cfg.Kafka)
	if err != nil {
		log.Fatal(err)
	}

	ratingApp := bootstrap.NewRatingApp(redis, db, kafkaProducer, kafkaConsumer, models.NewStat())

	ratingApp.Start()
}
