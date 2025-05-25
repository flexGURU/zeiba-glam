package repository

import (
	"context"
	"time"

	"github.com/flexGURU/zeiba-glam/backend/pkg"
)

type Payment struct {
	ID            uint32    `json:"id"`
	OrderID       uint32    `json:"order_id"`
	Amount        float64   `json:"amount"`
	TransactionID string    `json:"transaction_id"`
	PaymentMethod string    `json:"payment_method"`
	PaymentStatus bool      `json:"payment_status"`
	PaidAt        time.Time `json:"paid_at"`
	CreatedAt     time.Time `json:"created_at"`

	OrderDetails *Order `json:"order_details"`
}

type UpdatePayment struct {
	ID            uint32     `json:"id"`
	PaymentStatus *bool      `json:"payment_status"`
	PaidAt        *time.Time `json:"paid_at"`
}

type PaymentFilter struct {
	Pagination    *pkg.Pagination
	PaymentMethod *string
	PaymentStatus *bool
}

type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment *Payment) (*Payment, error)
	GetPayment(ctx context.Context, id uint32, orderID uint32) (*Payment, error)
	ListPayments(ctx context.Context, filter *PaymentFilter) ([]*Payment, *pkg.Pagination, error)
	UpdatePayment(ctx context.Context, payment *UpdatePayment) (*Payment, error)
}
