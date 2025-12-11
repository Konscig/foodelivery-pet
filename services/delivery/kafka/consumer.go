package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/Konscig/foodelivery-pet/services/delivery/models"
	"github.com/segmentio/kafka-go"
)

// Consumer читает `order.ready`, назначает курьера, публикует `order.coming` и `order.done`.
type Consumer struct {
	reader *kafka.Reader
	writer *kafka.Writer
}

func NewConsumer(brokers []string, groupID, topic, comingTopic, doneTopic string) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 1e3,
		MaxBytes: 10e6,
	})

	// Для простоты используем один writer; при необходимости можно иметь два для разных топиков.
	w := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    comingTopic, // по умолчанию пишем в comingTopic, для done будем менять Topic в сообщении
		Balancer: &kafka.LeastBytes{},
	}

	return &Consumer{reader: r, writer: w}
}

func (c *Consumer) Start(ctx context.Context) error {
	log.Println("delivery consumer started")
	for {
		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("consumer context canceled, stopping")
				return nil
			}
			log.Printf("error fetching message: %v", err)
			time.Sleep(500 * time.Millisecond)
			continue
		}

		log.Printf("delivery: message received key=%s offset=%d", string(m.Key), m.Offset)
		order, err := models.FromJSON(m.Value)
		if err != nil {
			log.Printf("delivery: failed to unmarshal order: %v", err)
			if err := c.reader.CommitMessages(ctx, m); err != nil {
				log.Printf("delivery: failed to commit message: %v", err)
			}
			continue
		}

		// Назначаем курьера (мок)
		courier, err := order.AssignCourier(ctx)
		if err != nil {
			log.Printf("failed to assign courier: %v", err)
		}

		// Обновляем статус и сохраняем
		order.Status = "assigned"
		if err := order.SaveToRedis(ctx); err != nil {
			log.Printf("delivery: save to redis error: %v", err)
		}
		if err := order.SaveToDB(ctx); err != nil {
			log.Printf("delivery: save to db error: %v", err)
		}

		// Публикуем order.coming (курьер в пути)
		coming := map[string]interface{}{
			"order_id": order.ID,
			"courier":  courier,
			"status":   "coming",
		}
		comingVal, _ := json.Marshal(coming)
		if err := c.writer.WriteMessages(ctx, kafka.Message{Key: []byte(order.ID), Value: comingVal}); err != nil {
			log.Printf("failed to publish order.coming: %v", err)
		} else {
			log.Printf("published order.coming: %s", order.ID)
		}

		// Имитация доставки: через небольшую задержку публикуем order.done
		go func(o *models.Order) {
			// небольшая пауза, имитируем время доставки
			time.Sleep(3 * time.Second)
			o.Status = "done"
			doneVal, _ := json.Marshal(o)

			// Пишем напрямую в Kafka через временный writer для топика order.done
			wDone := &kafka.Writer{Addr: kafka.TCP("localhost:9092"), Topic: "order.done", Balancer: &kafka.LeastBytes{}}
			defer wDone.Close()
			if err := wDone.WriteMessages(ctx, kafka.Message{Key: []byte(o.ID), Value: doneVal}); err != nil {
				log.Printf("failed to publish order.done: %v", err)
			} else {
				log.Printf("published order.done: %s", o.ID)
			}
		}(order)

		if err := c.reader.CommitMessages(ctx, m); err != nil {
			log.Printf("delivery: failed to commit message: %v", err)
		}
	}
}

func (c *Consumer) Close() error {
	if err := c.reader.Close(); err != nil {
		return err
	}
	return c.writer.Close()
}
