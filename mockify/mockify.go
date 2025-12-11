package mockify

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// Простая утилита для генерации мок-данных (orders, ratings и т.п.).
// Название "mockify" — внутренний пакет-генератор, не зависит от внешних библиотек.

type Order struct {
	ID         string   `json:"id"`
	CustomerID string   `json:"customer_id"`
	Items      []string `json:"items"`
	Status     string   `json:"status"`
	CreatedAt  int64    `json:"created_at"`
}

type Rating struct {
	OrderID string `json:"order_id"`
	Score   int    `json:"score"`
	Comment string `json:"comment,omitempty"`
	Created int64  `json:"created_at"`
}

// Инициализация генератора
func init() {
	rand.Seed(time.Now().UnixNano())
}

// GenerateOrder создает случайный заказ.
// useStatus — если true, устанавливает начальный статус (created)
func GenerateOrder(useStatus bool) Order {
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
	return Order{
		ID:         id,
		CustomerID: customer,
		Items:      items,
		Status:     status,
		CreatedAt:  time.Now().Unix(),
	}
}

// GenerateOrders создает n заказов.
func GenerateOrders(n int) []Order {
	out := make([]Order, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, GenerateOrder(true))
		// короткая пауза, чтобы id не совпадали
		time.Sleep(time.Nanosecond)
	}
	return out
}

// GenerateRating создает простую запись оценки для заказа
func GenerateRating(orderID string) Rating {
	score := rand.Intn(5) + 1
	comment := ""
	if rand.Float32() < 0.6 {
		comment = fmt.Sprintf("Good service %d", rand.Intn(100))
	}
	return Rating{
		OrderID: orderID,
		Score:   score,
		Comment: comment,
		Created: time.Now().Unix(),
	}
}

// ToJSON helper
func ToJSON(v any) ([]byte, error) {
	return json.Marshal(v)
}
