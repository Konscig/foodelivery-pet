package storage

import (
	"context"

	models "github.com/Konscig/foodelivery-pet/internal/storage/models"
)

func (s *PGstorage) GetOrder(id string) (*models.Order, error) {
	var order models.Order
	err := s.db.QueryRow(context.Background(), `
		SELECT id, user_id, rest_id, status, items FROM orders WHERE id = $1
	`, id).Scan(&order.ID, &order.UserID, &order.RestID, &order.Status, &order.Items)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (s *PGstorage) GetOrdersByUser(userID string) ([]*models.Order, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT id, user_id, rest_id, status, items FROM orders WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.UserID, &order.RestID, &order.Status, &order.Items)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}
	return orders, nil
}

func (s *PGstorage) GetDelivery(id string) (*models.Delivery, error) {
	var delivery models.Delivery
	err := s.db.QueryRow(context.Background(), `
		SELECT id, order_id, courier_id, status FROM deliveries WHERE id = $1
	`, id).Scan(&delivery.ID, &delivery.OrderID, &delivery.CourierID, &delivery.Status)
	if err != nil {
		return nil, err
	}
	return &delivery, nil
}
