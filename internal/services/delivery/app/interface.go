package app

type DeliveryService interface {
	Start() error
}

type OrderEventPublisher interface {
	PublishOrderComing(orderID, courierID string) error
	PublishOrderDone(orderID, courierID string) error
}

type OrderStatusStore interface {
	SetOrderStatus(key string, status string) error
}
