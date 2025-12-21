package bootstrap

import (
	"context"

	"github.com/Konscig/foodelivery-pet/config"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

type Consumer struct {
	reader *kafka.Reader
}

func NewProducer(broker *config.Config) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(broker.Kafka.Broker),
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func NewConsumer(broker *config.Config, groupID string, topic string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{broker.Kafka.Broker},
			Topic:   topic,
			GroupID: groupID,
		}),
	}
}

func (p *Producer) SendProtoMessage(topic string, value []byte) error {
	p.writer.Topic = topic
	msg := kafka.Message{
		Value: value,
	}
	return p.writer.WriteMessages(context.Background(), msg)
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

func (c *Consumer) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return c.reader.ReadMessage(ctx)
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
