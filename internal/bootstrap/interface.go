package bootstrap

import "context"

type Message struct {
	Key   []byte
	Value []byte
}

type EventConsumer interface {
	ReadMessage(ctx context.Context) (Message, error)
}

type EventProducer interface {
	SendProtoMessage(topic string, key []byte, value []byte) error
	Close() error
} 
