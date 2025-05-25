package repository

import (
	"context"
	"time"

	"github.com/flexGURU/zeiba-glam/backend/pkg"
)

type Order struct {
	ID              uint32     `json:"id"`
	UserName        string     `json:"user_name"`
	UserPhoneNumber string     `json:"user_phone_number"`
	TotalAmount     float64    `json:"total_amount"`
	Status          string     `json:"status"`
	ShippingAddress string     `json:"shipping_address"`
	PaymentStatus   bool       `json:"payment_status"`
	DeletedAt       *time.Time `json:"deleted_at"`
	CreatedAt       time.Time  `json:"created_at"`

	OrderItems []*OrderItem `json:"order_items,omitempty"`
	Payments   []*Payment   `json:"payments,omitempty"`
}

type OrderItem struct {
	ID        uint32  `json:"id"`
	OrderID   uint32  `json:"order_id"`
	ProductID uint32  `json:"product_id"`
	Quantity  int64   `json:"quantity"`
	Amount    float64 `json:"amount"`
	Size      string  `json:"size"`
	Color     string  `json:"color"`
}

type UpdateOrder struct {
	ID            uint32  `json:"id"`
	Status        *string `json:"status"`
	PaymentStatus *bool   `json:"payment_status"`
}

type OrderFilter struct {
	Pagination    *pkg.Pagination
	Status        *string
	PaymentStatus *bool
}
type OrderRepository interface {
	CreateOrder(ctx context.Context, order *Order) (*Order, error)
	GetOrderByID(ctx context.Context, id uint32) (*Order, error)
	ListOrders(ctx context.Context, filter *OrderFilter) ([]*Order, *pkg.Pagination, error)
	UpdateOrder(ctx context.Context, order *UpdateOrder) (*Order, error)
	DeleteOrder(ctx context.Context, id uint32) error
}
