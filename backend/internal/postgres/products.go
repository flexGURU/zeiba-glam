package postgres

import (
	"context"

	"github.com/flexGURU/zeiba-glam/backend/internal/postgres/generated"
	"github.com/flexGURU/zeiba-glam/backend/internal/repository"
	"github.com/flexGURU/zeiba-glam/backend/pkg"
	"github.com/jackc/pgx/v5/pgtype"
)

var _ repository.ProductRepository = (*ProductRepo)(nil)

// ProductRepo is the repository for the product model
type ProductRepo struct {
	queries generated.Querier
}

func NewProductRepo(db *Store) *ProductRepo {
	return &ProductRepo{
		queries: generated.New(db.pool),
	}
}

func (r *ProductRepo) CreateProduct(
	ctx context.Context,
	product *repository.Product,
) (*repository.Product, error) {
	productCreated, err := r.queries.CreateProduct(ctx, generated.CreateProductParams{
		Name:          product.Name,
		Description:   product.Description,
		Price:         pkg.Float64ToPgTypeNumeric(product.Price),
		Category:      product.Category,
		ImageUrl:      product.ImageURL,
		Size:          product.Size,
		Color:         product.Color,
		StockQuantity: product.StockQuantity,
		UpdatedBy:     int64(product.UpdatedBy),
	})

	if err != nil {
		if pkg.PgxErrorCode(err) == pkg.UNIQUE_VIOLATION {
			return nil, pkg.Errorf(pkg.ALREADY_EXISTS_ERROR, "%s", err.Error())
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error creating product: %s", err.Error())
	}

	product.ID = uint32(productCreated.ID)
	product.CreatedAt = productCreated.CreatedAt

	return product, nil
}

func (r *ProductRepo) GetProductByID(ctx context.Context, id uint32) (*repository.Product, error) {
	product, err := r.queries.GetProductByID(ctx, int64(id))
	if err != nil {
		if pkg.PgxErrorCode(err) == pkg.NOT_FOUND_ERROR {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "product not found")
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting product by id: %s", err.Error())
	}

	return marshalProduct(product), nil
}

func (r *ProductRepo) ListProducts(
	ctx context.Context,
	filter *repository.ProductFilter,
) ([]*repository.Product, *pkg.Pagination, error) {
	paramListProduct := generated.ListProductsParams{
		Limit:  int32(filter.Pagination.PageSize),
		Offset: pkg.Offset(filter.Pagination.Page, filter.Pagination.PageSize),
	}

	paramListProductCount := generated.ListProductsCountParams{}

	if filter.Search != nil {
		paramListProduct.Search = pgtype.Text{
			String: *filter.Search,
			Valid:  true,
		}
		paramListProductCount.Search = pgtype.Text{
			String: *filter.Search,
			Valid:  true,
		}
	}

	if filter.PriceFrom != nil && filter.PriceTo != nil {
		paramListProduct.PriceFrom = pkg.Float64ToPgTypeNumeric(*filter.PriceFrom)
		paramListProduct.PriceTo = pkg.Float64ToPgTypeNumeric(*filter.PriceTo)
		paramListProductCount.PriceFrom = pkg.Float64ToPgTypeNumeric(*filter.PriceFrom)
		paramListProductCount.PriceTo = pkg.Float64ToPgTypeNumeric(*filter.PriceTo)
	}

	if filter.Category != nil {
		paramListProduct.Category = pgtype.Array[string]{
			Valid:    true,
			Elements: *filter.Category,
		}
		paramListProductCount.Category = pgtype.Array[string]{
			Valid:    true,
			Elements: *filter.Category,
		}
	}

	if filter.Size != nil {
		paramListProduct.Size = pgtype.Array[string]{
			Valid:    true,
			Elements: *filter.Size,
		}
		paramListProductCount.Size = pgtype.Array[string]{
			Valid:    true,
			Elements: *filter.Size,
		}
	}

	if filter.Color != nil {
		paramListProduct.Color = pgtype.Array[string]{
			Valid:    true,
			Elements: *filter.Color,
		}
		paramListProductCount.Color = pgtype.Array[string]{
			Valid:    true,
			Elements: *filter.Color,
		}
	}

	products, err := r.queries.ListProducts(ctx, paramListProduct)
	if err != nil {
		return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error listing products: %s", err.Error())
	}

	productsCount, err := r.queries.ListProductsCount(ctx, paramListProductCount)
	if err != nil {
		return nil, nil, pkg.Errorf(
			pkg.INTERNAL_ERROR,
			"error listing products count: %s",
			err.Error(),
		)
	}

	productsList := make([]*repository.Product, len(products))
	for idx, product := range products {
		productsList[idx] = marshalProduct(product)
	}

	return productsList, pkg.CalculatePagination(
		uint32(productsCount),
		filter.Pagination.PageSize,
		filter.Pagination.Page,
	), nil
}

func (r *ProductRepo) UpdateProduct(
	ctx context.Context,
	product *repository.UpdateProduct,
) (*repository.Product, error) {
	params := generated.UpdateProductParams{
		ID:        int64(product.ID),
		UpdatedBy: int64(product.UpdatedBy),
	}

	if product.Name != nil {
		params.Name = pgtype.Text{
			String: *product.Name,
			Valid:  true,
		}
	}

	if product.Description != nil {
		params.Description = pgtype.Text{
			String: *product.Description,
			Valid:  true,
		}
	}

	if product.Price != nil {
		params.Price = pkg.Float64ToPgTypeNumeric(*product.Price)
	}

	if product.Category != nil {
		params.Category = *product.Category
	}

	if product.ImageURL != nil {
		params.ImageUrl = *product.ImageURL
	}

	if product.Size != nil {
		params.Size = *product.Size
	}

	if product.Color != nil {
		params.Color = *product.Color
	}

	if product.StockQuantity != nil {
		params.StockQuantity = pgtype.Int8{
			Int64: *product.StockQuantity,
			Valid: true,
		}
	}

	productUpdated, err := r.queries.UpdateProduct(ctx, params)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error updating product: %s", err.Error())
	}

	return marshalProduct(productUpdated), nil
}

func (r *ProductRepo) DeleteProduct(ctx context.Context, id uint32) error {
	if err := r.queries.DeleteProduct(ctx, int64(id)); err != nil {
		if pkg.PgxErrorCode(err) == pkg.NOT_FOUND_ERROR {
			return pkg.Errorf(pkg.NOT_FOUND_ERROR, "product not found")
		}

		return pkg.Errorf(pkg.INTERNAL_ERROR, "error deleting product: %s", err.Error())
	}

	return nil
}

func marshalProduct(product generated.Product) *repository.Product {
	p := &repository.Product{
		ID:            uint32(product.ID),
		Name:          product.Name,
		Description:   product.Description,
		Price:         pkg.PgTypeNumericToFloat64(product.Price),
		Category:      product.Category,
		ImageURL:      product.ImageUrl,
		Size:          product.Size,
		Color:         product.Color,
		StockQuantity: product.StockQuantity,
		DeletedAt:     nil,
		UpdatedBy:     uint32(product.UpdatedBy),
		CreatedAt:     product.CreatedAt,
	}

	if product.DeletedAt.Valid {
		p.DeletedAt = &product.DeletedAt.Time
	}

	return p
}
