package postgres

import (
	"context"

	"github.com/flexGURU/zeiba-glam/backend/internal/postgres/generated"
	"github.com/flexGURU/zeiba-glam/backend/internal/repository"
	"github.com/flexGURU/zeiba-glam/backend/pkg"
)

type ProductRepo struct {
	db      *Store
	queries generated.Querier
}

func NewProductRepo(db *Store) *ProductRepo {
	return &ProductRepo{
		db:      db,
		queries: generated.New(db.pool),
	}
}

func (r *ProductRepo) CreateProduct(
	ctx context.Context,
	product *repository.Product,
) (*repository.Product, error) {
}
func (r *ProductRepo) GetProductByID(ctx context.Context, id int64) (*repository.Product, error) {}

func (r *ProductRepo) ListProducts(
	ctx context.Context,
	filter *repository.ProductFilter,
) ([]*repository.Product, *pkg.Pagination, error) {
}

func (r *ProductRepo) UpdateProduct(
	ctx context.Context,
	product *repository.UpdateProduct,
) (*repository.Product, error) {
}
func (r *ProductRepo) DeleteProduct(ctx context.Context, id int64) error {}
