package storage

import (
	"context"

	storage "github.com/Konscig/foodelivery-pet/internal/storage/models"
)

func (s *PGstorage) UpdateOrderStatus(id string, status storage.OrderStatus) error {
	_, err := s.db.Exec(context.Background(), `
		UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2
	`, status, id)
	return err
}

func (s *PGstorage) UpdateDeliveryStatus(id string, status string) error {
	_, err := s.db.Exec(context.Background(), `
		UPDATE deliveries SET status = $1, updated_at = NOW() WHERE id = $2
	`, status, id)
	return err
}
