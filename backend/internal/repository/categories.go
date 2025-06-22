package repository

import (
	"context"
	"time"
)

type Category struct {
	ID          uint32    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type UpdateCategory struct {
	ID          uint32  `json:"id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category *Category) (*Category, error)
	GetCategory(ctx context.Context, id uint32) (*Category, error)
	ListCategories(ctx context.Context) ([]*Category, error)
	UpdateCategory(ctx context.Context, category *UpdateCategory) (*Category, error)
	DeleteCategory(ctx context.Context, id uint32) error
}
