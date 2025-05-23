package repository

import (
	"context"
	"time"
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

type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment Payment) (Payment, error)
	GetPaymentByID(ctx context.Context, id int64) (Payment, error)
	GetPayments(ctx context.Context) ([]Payment, error)
	UpdatePayment(ctx context.Context, payment Payment) (Payment, error)
	DeletePayment(ctx context.Context, id int64) error
}
