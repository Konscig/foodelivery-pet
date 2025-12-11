package mockify

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/Konscig/foodelivery-pet/api/kafka"
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

	prod := kafka.NewProducer(brokers, *topic)
	defer prod.Close()

	// Простая генерация мок-заказов прямо в этом файле, чтобы не держать отдельный internal-пакет.
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < *count; i++ {
		order := generateOrder(true)
		b, err := json.Marshal(order)
		if err != nil {
			log.Fatalf("marshal error: %v", err)
		}
		if err := prod.PublishOrderCreated(order.ID, b); err != nil {
			log.Printf("failed to publish: %v", err)
		} else {
			fmt.Printf("published mock order %s\n", order.ID)
		}
		time.Sleep(time.Duration(*interval) * time.Millisecond)
	}
}

// generateOrder — локальная функция генерации мок-объекта заказа.
func generateOrder(useStatus bool) map[string]interface{} {
	id := fmt.Sprintf("order-%d", time.Now().UnixNano())
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
	return map[string]interface{}{
		"id":          id,
		"customer_id": customer,
		"items":       items,
		"status":      status,
		"created_at":  time.Now().Unix(),
	}
}
