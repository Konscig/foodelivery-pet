package bootstrap

import (
	"context"
	"log"

	"gorm.io/gorm"
)

type RestaurantService struct {
	db            *gorm.DB
	RedisClient   *RedisClient
	kafkaConsumer *Consumer
	kafkaProducer *Producer
}

func (s *RestaurantService) StartRestaurantService() error {
	log.Println("Запуск сервиса Restaurant")
	ctx := context.Background()
	for {
		msg, err := s.kafkaConsumer.ReadMessage(ctx)
		if err != nil {
			log.Printf("Ошибка чтения сообщения: %v", err)
			continue
		}
		log.Printf("Получено сообщение: %s", string(msg.Value))
	}
}
