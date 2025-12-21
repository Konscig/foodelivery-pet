package bootstrap

import (
	"context"
	"log"
)

type RatingService struct {
	kafkaConsumer EventConsumer
	kafkaProducer EventProducer
}

func StartRatingService(s *RatingService) error {
	log.Println("Запуск сервиса Rating")
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
