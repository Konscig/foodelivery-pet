package bootstrap

import (
	"context"
	"log"
)

type OrderService struct {
	kafkaConsumer *Consumer
	kafkaProducer *Producer
}

func (s *OrderService) StartOrderService() error {
	log.Println("Запуск сервиса Order")
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
