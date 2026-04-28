package domain

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusInProgress OrderStatus = "in_progress"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

type Order struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	Status          OrderStatus
	Total           float64
	Comment         string
	ShippingAddress string
	ContactPhone    string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type OrderItem struct {
	ID        uuid.UUID
	OrderID   uuid.UUID
	ProductID uuid.UUID
	Title     string
	Price     float64
	Quantity  int
}
