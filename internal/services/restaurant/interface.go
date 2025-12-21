package restaurant

type RestaurantService interface {
	Start() error
	ProcessOrder(orderID string) error
}
