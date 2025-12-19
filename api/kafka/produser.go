package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// Producer обертка над kafka.Writer
type Producer struct {
	writer *kafka.Writer
}

// NewProducer создаёт продюсер
func NewProducer(brokers []string) *Producer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Balancer: &kafka.LeastBytes{},
	}
	return &Producer{writer: w}
}

// SendProtoMessage отправляет сообщение в Kafka в указанный топик
func (p *Producer) SendProtoMessage(topic string, value []byte) error {
	p.writer.Topic = topic
	msg := kafka.Message{
		Value: value,
	}
	return p.writer.WriteMessages(context.Background(), msg)
}

// Close закрывает продюсер
func (p *Producer) Close() error {
	return p.writer.Close()
}
