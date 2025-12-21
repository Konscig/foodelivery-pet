package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	eventspb "github.com/Konscig/foodelivery-pet/generated/eventspb"
	kafka "github.com/Konscig/foodelivery-pet/internal/bootstrap"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

// cmd/mockify — небольшая утилита для генерации и отправки тестовых событий в Kafka.
// Пример:
// go run ./cmd/mockify --brokers localhost:9092 --topic order.created --count 5 --interval 500

func main() {
	brokersFlag := flag.String("brokers", "localhost:9092", "comma-separated list of kafka brokers")
	topic := flag.String("topic", "order.created", "kafka topic to publish to")
	count := flag.Int("count", 1, "number of messages to send")
	interval := flag.Int("interval", 500, "milliseconds between messages")
	flag.Parse()

	brokers := strings.Split(*brokersFlag, ",")
	log.Printf("mockify: brokers=%v topic=%s count=%d interval=%dms", brokers, *topic, *count, *interval)

	prod := kafka.NewProducer(brokers)
	defer prod.Close()

	// Простая генерация мок-заказов прямо в этом файле, чтобы не держать отдельный internal-пакет.
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < *count; i++ {
		order := GenerateOrder(true)
		items := make([]*eventspb.OrderItem, len(order.Items))
		for j, item := range order.Items {
			items[j] = &eventspb.OrderItem{Name: item, Quantity: 1}
		}
		payload := &eventspb.OrderCreatedPayload{
			UserId: order.CustomerID,
			RestId: uuid.NewString(),
			Items:  items,
		}
		payloadBytes, _ := proto.Marshal(payload)
		event := &eventspb.OrderEvent{
			EventId:   uuid.NewString(),
			OrderId:   order.ID,
			Status:    eventspb.OrderStatus_CREATED,
			Timestamp: time.Now().Unix(),
			Payload:   payloadBytes,
		}
		eventBytes, _ := proto.Marshal(event)
		if err := prod.SendProtoMessage(kafka.TopicOrderCreated, eventBytes); err != nil {
			log.Printf("failed to publish: %v", err)
		} else {
			fmt.Printf("published mock order %s\n", order.ID)
		}
		time.Sleep(time.Duration(*interval) * time.Millisecond)
	}
}

// GenerateOrder создает случайный заказ.
func GenerateOrder(useStatus bool) Order {
	id := uuid.NewString()
	customer := fmt.Sprintf("cust-%d", rand.Intn(1000))
	itemsCount := rand.Intn(4) + 1
	items := make([]string, 0, itemsCount)
	for i := 0; i < itemsCount; i++ {
		items = append(items, fmt.Sprintf("item-%d", rand.Intn(100)))
	}
	status := "created"
	if !useStatus {
		status = ""
	}
	return Order{
		ID:         id,
		CustomerID: customer,
		Items:      items,
		Status:     status,
		CreatedAt:  time.Now().Unix(),
	}
}

type Order struct {
	ID         string   `json:"id"`
	CustomerID string   `json:"customer_id"`
	Items      []string `json:"items"`
	Status     string   `json:"status"`
	CreatedAt  int64    `json:"created_at"`
}
