package storage

import (
	"context"

	models "github.com/Konscig/foodelivery-pet/internal/storage/models"
)

func (s *PGstorage) AddOrder(order *models.Order) error {
	_, err := s.db.Exec(context.Background(), `
		INSERT INTO orders (id, user_id, rest_id, status, items)
		VALUES ($1, $2, $3, $4, $5)
	`, order.ID, order.UserID, order.RestID, order.Status, order.Items)
	return err
}

func (s *PGstorage) AddDelivery(delivery *models.Delivery) error {
	_, err := s.db.Exec(context.Background(), `
		INSERT INTO deliveries (id, order_id, courier_id, status)
		VALUES ($1, $2, $3, $4)
	`, delivery.ID, delivery.OrderID, delivery.CourierID, delivery.Status)
	return err
}
