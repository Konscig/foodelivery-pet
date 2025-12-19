package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// Producer обертка над kafka.Writer
type Producer struct {
	writer *kafka.Writer
}

// NewProducer создаёт продюсер с заранее заданным топиком
func NewProducer(brokers []string, topic string) *Producer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic, // топик задаём один раз здесь
		Balancer: &kafka.LeastBytes{},
	}
	return &Producer{writer: w}
}

// SendProtoMessage отправляет сообщение в Kafka
func (p *Producer) SendProtoMessage(value []byte) error {
	msg := kafka.Message{
		Value: value,
	}
	return p.writer.WriteMessages(context.Background(), msg)
}

// Close закрывает продюсер
func (p *Producer) Close() error {
	return p.writer.Close()
}
