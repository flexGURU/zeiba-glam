package repository

import (
	"context"
	"time"

	"github.com/flexGURU/zeiba-glam/backend/pkg"
)

type Payment struct {
	ID            int64     `json:"id"`
	OrderID       int64     `json:"order_id"`
	Amount        float64   `json:"amount"`
	TransactionID string    `json:"transaction_id"`
	PaymentMethod string    `json:"payment_method"`
	PaymentStatus string    `json:"payment_status"`
	CreatedAt     time.Time `json:"created_at"`
}

type UpdatePayment struct {
	ID            int64   `json:"id"`
	PaymentStatus *string `json:"payment_status"`
}

type PaymentFilter struct {
	Pagination    *pkg.Pagination
	PaymentMethod *string
	PaymentStatus *string
}

type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment *Payment) (*Payment, error)
	GetPaymentByID(ctx context.Context, id int64) (*Payment, error)
	GetPaymentByOrderID(ctx context.Context, orderID int64) (*Payment, error)
	ListPayments(ctx context.Context, filter *PaymentFilter) ([]*Payment, error)
	UpdatePayment(ctx context.Context, payment UpdatePayment) (*Payment, error)
}
