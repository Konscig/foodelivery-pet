package storage

import (
	"github.com/Konscig/foodelivery-pet/internal/pb/orderpb"
	"gorm.io/gorm"
)

type Delivery struct {
	gorm.Model
	ID        string `gorm:"primaryKey"`
	OrderID   string
	CourierID string
	Status    string
}

type OrderStatus string

const (
	StatusCreated OrderStatus = "CREATED"
	StatusReady   OrderStatus = "READY"
	StatusComing  OrderStatus = "COMING"
	StatusDone    OrderStatus = "DONE"
	StatusRated   OrderStatus = "RATED"
)

type Order struct {
	gorm.Model
	ID     string               `gorm:"primaryKey"`
	UserID string               `gorm:"not null"`
	RestID string               `gorm:"not null"`
	Status OrderStatus          `gorm:"not null"`
	Items  []*orderpb.OrderItem `gorm:"not null"`
}

type OrderItem struct {
	ID       string `gorm:"primaryKey"`
	OrderID  string
	Name     string
	Price    string
	Quantity int64
}

type Review struct {
	gorm.Model
	ID           string `gorm:"primaryKey"`
	OrderID      string
	RestaurantID string
	Rating       int32
	Comment      string
}

type ReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) CreateReview(review *Review) error {
	return r.db.Create(&review).Error
}

func (r *ReviewRepository) NewService(repo *ReviewRepository) error {
	return nil
}

func (r *ReviewRepository) UpdateRestaurantStats(restaurantID string) error {
	return nil
}

func (r *ReviewRepository) GetRestaurantStats(restaurantID string) (*RestaurantStats, error) {
	return &RestaurantStats{}, nil
}

type RestaurantStats struct {
	RestaurantID  string
	AvgRating     float64
	ReviewsCount  uint32
	WordCloudJSON string
}

func BuildWordCloud(input string) map[string]int {
	// Placeholder implementation
	return map[string]int{"example": 1}
}
