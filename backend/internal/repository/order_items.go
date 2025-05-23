package repository

import "context"

type OrderItem struct {
	ID        int64   `json:"id"`
	OrderID   int64   `json:"order_id"`
	ProductID int64   `json:"product_id"`
	Quantity  int64   `json:"quantity"`
	Amount    float64 `json:"amount"`
	Size      string  `json:"size"`
	Color     string  `json:"color"`
}

type OrderItemRepository interface {
	GetOrderItemByID(ctx context.Context, id int64) (*OrderItem, error)
}
