package bootstrap

import (
	"context"
	"log"

	"github.com/Konscig/foodelivery-pet/internal/common"
)

type OrderService struct {
	kafkaConsumer *common.Consumer
	kafkaProducer *common.Producer
}

func NewOrderService() *OrderService {
	return &OrderService{
		kafkaConsumer: common.NewConsumer([]string{"localhost:9092"}, "order.created", "order-group"),
		kafkaProducer: common.NewProducer([]string{"localhost:9092"}),
	}
}

func (s *OrderService) Start() error {
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
