package repository

import (
	"context"
	"time"

	"github.com/flexGURU/zeiba-glam/backend/pkg"
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

type UpdateOrder struct {
	ID            int64   `json:"id"`
	Status        *string `json:"status"`
	PaymentStatus *string `json:"payment_status"`
}

type OrderFilter struct {
	Pagination    *pkg.Pagination
	Status        *string
	PaymentStatus *string
}
type OrderRepository interface {
	CreateOrder(ctx context.Context, order *Order) (*Order, error)
	GetOrderByID(ctx context.Context, id int64) (*Order, error)
	ListOrders(ctx context.Context, filter *OrderFilter) ([]*Order, error)
	UpdateOrder(ctx context.Context, order *UpdateOrder) (*Order, error)
	DeleteOrder(ctx context.Context, id int64) error
}
