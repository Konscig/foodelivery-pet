package app

type RestaurantService interface {
	Start() error
}

type OrderReadyPublisher interface {
	PublishOrderReady(orderID string, restID string) error
}

type OrderStatusStore interface {
	SetOrderStatus(key string, status string) error
}
