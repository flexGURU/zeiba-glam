package repository

import (
	"context"
	"time"
)

type Order struct {
	ID              int64      `json:"id"`
	UserName        string     `json:"user_name"`
	UserPhoneNumber string     `json:"user_phone_number"`
	TotalAmount     float64    `json:"total_amount"`
	Status          string     `json:"status"`
	ShippingAddress string     `json:"shipping_address"`
	PaymentStatus   string     `json:"payment_status"`
	DeletedAt       *time.Time `json:"deleted_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, order Order) (Order, error)
	GetOrderByID(ctx context.Context, id int64) (Order, error)
	GetOrders(ctx context.Context) ([]Order, error)
	UpdateOrder(ctx context.Context, order Order) (Order, error)
	DeleteOrder(ctx context.Context, id int64) error
}
