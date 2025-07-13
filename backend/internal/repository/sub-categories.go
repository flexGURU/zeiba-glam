package repository

import (
	"context"
	"time"
)

type SubCategory struct {
	ID                  uint32    `json:"id"`
	CategoryID          uint32    `json:"category_id"`
	Name                string    `json:"name"`
	Description         string    `json:"description"`
	CreatedAt           time.Time `json:"created_at"`
	CategoryName        *string   `json:"category_name,omitempty"`        // Optional, for convenience
	CategoryDescription *string   `json:"category_description,omitempty"` // Optional, for convenience
}

type UpdateSubCategory struct {
	ID          uint32  `json:"id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type SubCategoryRepository interface {
	CreateSubCategory(ctx context.Context, subCategory *SubCategory) (*SubCategory, error)
	GetSubCategory(ctx context.Context, id uint32) (*SubCategory, error)
	GetSubCategoryByCategoryIDAndID(ctx context.Context, categoryID, id uint32) (*SubCategory, error)
	ListSubCategoriesByCategoryID(ctx context.Context, categoryID uint32) ([]*SubCategory, error)
	ListSubCategories(ctx context.Context) ([]*SubCategory, error)
	UpdateSubCategory(ctx context.Context, subCategory *UpdateSubCategory) (*SubCategory, error)
	DeleteSubCategory(ctx context.Context, id uint32) error
}
