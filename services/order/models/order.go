package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

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
	ID     uuid.UUID   `gorm:"primaryKey"`
	UserID uuid.UUID   `gorm:"not null"`
	RestID uuid.UUID   `gorm:"not null"`
	Status OrderStatus `gorm:"not null"`
	Items  []OrderItem `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID       uuid.UUID `gorm:"primaryKey"`
	OrderID  uuid.UUID
	Name     string
	Price    string
	Quantity int64
}
