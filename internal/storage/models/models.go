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

type ReviewRepository interface {
	CreateReview(review *Review) error
}

type GormReviewRepository struct {
	db *gorm.DB
}

func NewGormReviewRepository(db *gorm.DB) *GormReviewRepository {
	return &GormReviewRepository{db: db}
}

func (r *GormReviewRepository) CreateReview(review *Review) error {
	return r.db.Create(review).Error
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
