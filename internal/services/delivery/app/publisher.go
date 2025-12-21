package app

import (
	"time"

	kafka "github.com/Konscig/foodelivery-pet/internal/bootstrap"
	eventspb "github.com/Konscig/foodelivery-pet/internal/pb/eventspb"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

type Publisher struct {
	producer *kafka.Producer
}

func NewPublisher(p *kafka.Producer) *Publisher {
	return &Publisher{producer: p}
}

func (p *Publisher) PublishOrderComing(orderID, courierID string) error {
	payload := &eventspb.OrderComingPayload{
		CourierId: courierID,
	}

	payloadBytes, _ := proto.Marshal(payload)

	event := &eventspb.OrderEvent{
		EventId:   uuid.NewString(),
		OrderId:   orderID,
		Status:    eventspb.OrderStatus_COMING,
		Timestamp: time.Now().Unix(),
		Payload:   payloadBytes,
	}

	eventBytes, _ := proto.Marshal(event)
	return p.producer.SendProtoMessage(kafka.TopicOrderComing, eventBytes)
}

func (p *Publisher) PublishOrderDone(orderID, courierID string) error {
	payload := &eventspb.OrderDonePayload{
		CourierId: courierID,
	}

	payloadBytes, _ := proto.Marshal(payload)

	event := &eventspb.OrderEvent{
		EventId:   uuid.NewString(),
		OrderId:   orderID,
		Status:    eventspb.OrderStatus_DONE,
		Timestamp: time.Now().Unix(),
		Payload:   payloadBytes,
	}

	eventBytes, _ := proto.Marshal(event)

	return p.producer.SendProtoMessage(kafka.TopicOrderDone, eventBytes)
}
