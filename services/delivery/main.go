package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Konscig/foodelivery-pet/services/delivery/kafka"
)

func main() {
	brokers := []string{getEnv("KAFKA_BROKER", "localhost:9092")}
	groupID := getEnv("KAFKA_GROUP", "delivery-group")
	topic := getEnv("KAFKA_TOPIC_ORDER_READY", "order.ready")
	comingTopic := getEnv("KAFKA_TOPIC_ORDER_COMING", "order.coming")
	doneTopic := getEnv("KAFKA_TOPIC_ORDER_DONE", "order.done")

	consumer := kafka.NewConsumer(brokers, groupID, topic, comingTopic, doneTopic)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() { <-sigs; log.Println("shutting down delivery..."); cancel() }()

	if err := consumer.Start(ctx); err != nil {
		log.Printf("delivery consumer finished with error: %v", err)
	}

	time.Sleep(200 * time.Millisecond)
	if err := consumer.Close(); err != nil {
		log.Printf("error closing delivery consumer: %v", err)
	}
}

func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
