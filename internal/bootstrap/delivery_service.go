package bootstrap

import (
	"context"
	"log"
)

type DeliveryService struct {
	kafkaConsumer EventConsumer
	kafkaProducer EventProducer
}

func StartDeliveryService(s *DeliveryService) error {
	log.Println("Запуск сервиса Delivery")
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
