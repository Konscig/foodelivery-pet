package bootstrap

import (
	"context"

	"github.com/Konscig/foodelivery-pet/config"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(cfg *config.Config) EventProducer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(cfg.Kafka.Broker),
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) SendProtoMessage(topic string, key []byte, value []byte) error {
	p.writer.Topic = topic
	msg := kafka.Message{
		Key:   key,
		Value: value,
	}
	return p.writer.WriteMessages(context.Background(), msg)
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(cfg *config.Config, groupID, topic string) EventConsumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{cfg.Kafka.Broker},
			Topic:   topic,
			GroupID: groupID,
		}),
	}
}

func (c *Consumer) ReadMessage(ctx context.Context) (Message, error) {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return Message{}, err
	}
	return Message{
		Key:   msg.Key,
		Value: msg.Value,
	}, nil
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
