package app

type RatingService interface {
	Start()
}

type OrderRatedPublisher interface {
	PublishOrderRated(
		orderID string,
		rating uint8,
		comment string,
		restaurantID string,
	) error
}

type OrderStatusStore interface {
	SetOrderStatus(key string, status string) error
}
