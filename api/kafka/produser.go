package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return &Producer{writer: w}
}

func (p *Producer) SendProtoMessage(topic string, key string, value []byte) error {
	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
	}
	err := p.writer.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Printf("failed to publish kafka message: %v", err)
		return err
	}
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
