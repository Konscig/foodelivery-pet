package storage

import (
	models "github.com/Konscig/foodelivery-pet/internal/storage/models"
)

type Storage interface {
	AddOrder(order *models.Order) error
	GetOrder(id string) (*models.Order, error)
	UpdateOrderStatus(id string, status models.OrderStatus) error
	GetOrdersByUser(userID string) ([]*models.Order, error)
	AddDelivery(delivery *models.Delivery) error
	GetDelivery(id string) (*models.Delivery, error)
	UpdateDeliveryStatus(id string, status string) error
}

type ReviewRepository interface {
	CreateReview(review *models.Review) error
	UpdateRestaurantStats(restaurantID string) error
	GetRestaurantStats(restaurantID string) (*models.RestaurantStats, error)
}
