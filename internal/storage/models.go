package storage

import (
	"github.com/Konscig/foodelivery-pet/generated/orderpb"
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
