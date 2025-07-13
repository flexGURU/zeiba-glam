package repository

import (
	"context"
	"time"

	"github.com/flexGURU/zeiba-glam/backend/pkg"
)

type Product struct {
	ID            uint32     `json:"id"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	Price         float64    `json:"price"`
	Category      string     `json:"category"`
	SubCategory   string     `json:"sub_category,omitempty"` // Optional, for convenience
	ImageURL      []string   `json:"image_url"`
	Size          []string   `json:"size"`
	Color         []string   `json:"color"`
	StockQuantity int64      `json:"stock_quantity"`
	DeletedAt     *time.Time `json:"deleted_at"`
	UpdatedBy     uint32     `json:"updated_by"`
	CreatedAt     time.Time  `json:"created_at"`
}

type UpdateProduct struct {
	ID            uint32    `json:"id"`
	UpdatedBy     uint32    `json:"updated_by"`
	Name          *string   `json:"name"`
	Description   *string   `json:"description"`
	Price         *float64  `json:"price"`
	Category      *string   `json:"category"`
	SubCategory   *string   `json:"sub_category,omitempty"` // Optional, for convenience
	ImageURL      *[]string `json:"image_url"`
	Size          *[]string `json:"size"`
	Color         *[]string `json:"color"`
	StockQuantity *int64    `json:"stock_quantity"`
}

type ProductFilter struct {
	Pagination  *pkg.Pagination
	Search      *string
	PriceFrom   *float64
	PriceTo     *float64
	Category    *[]string
	SubCategory *[]string
	Size        *[]string
	Color       *[]string
}

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *Product) (*Product, error)
	GetProductByID(ctx context.Context, id uint32) (*Product, error)
	ListProducts(ctx context.Context, filter *ProductFilter) ([]*Product, *pkg.Pagination, error)
	UpdateProduct(ctx context.Context, product *UpdateProduct) (*Product, error)
	DeleteProduct(ctx context.Context, id uint32) error
}
