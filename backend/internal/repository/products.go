package repository

import (
	"context"
	"time"
)

type Product struct {
	ID            int64      `json:"id"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	Price         float64    `json:"price"`
	Category      []string   `json:"category"`
	ImageURL      []string   `json:"image_url"`
	Size          []string   `json:"size"`
	Color         []string   `json:"color"`
	StockQuantity int64      `json:"stock_quantity"`
	DeletedAt     *time.Time `json:"deleted_at"`
	CreatedAt     time.Time  `json:"created_at"`
}

type ProductRepository interface {
	CreateProduct(ctx context.Context, product Product) (Product, error)
	GetProductByID(ctx context.Context, id int64) (Product, error)
	GetProducts(ctx context.Context) ([]Product, error)
	UpdateProduct(ctx context.Context, product Product) (Product, error)
	DeleteProduct(ctx context.Context, id int64) error
}
