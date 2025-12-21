package bootstrap

import (
	"context"
	"log"
)

type DeliveryService struct {
	kafkaConsumer *Consumer
	kafkaProducer *Producer
}

func NewDeliveryService() *DeliveryService {
	return &DeliveryService{
		kafkaConsumer: NewConsumer([]string{"localhost:9092"}, "delivery.created", "delivery-group"),
		kafkaProducer: NewProducer([]string{"localhost:9092"}),
	}
}

func (s *DeliveryService) Start() error {
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
