package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/flexGURU/zeiba-glam/backend/internal/postgres/generated"
	"github.com/flexGURU/zeiba-glam/backend/internal/repository"
	"github.com/flexGURU/zeiba-glam/backend/pkg"
	"github.com/jackc/pgx/v5/pgtype"
)

var _ repository.CategoryRepository = (*CategoryRepo)(nil)

type CategoryRepo struct {
	queries generated.Querier
	db      *Store
}

func NewCategoryRepo(db *Store) *CategoryRepo {
	return &CategoryRepo{
		queries: generated.New(db.pool),
		db:      db,
	}
}

func (s *CategoryRepo) CreateCategory(
	ctx context.Context,
	category *repository.Category,
) (*repository.Category, error) {
	categoryCreated, err := s.queries.CreateCategory(ctx, generated.CreateCategoryParams{
		Name:        category.Name,
		Description: category.Description,
	})

	if err != nil {
		if pkg.PgxErrorCode(err) == pkg.UNIQUE_VIOLATION {
			return nil, pkg.Errorf(pkg.ALREADY_EXISTS_ERROR, "%s", err.Error())
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error creating category: %s", err.Error())
	}

	category.ID = uint32(categoryCreated.ID)
	category.CreatedAt = categoryCreated.CreatedAt

	return category, nil
}

func (s *CategoryRepo) GetCategory(ctx context.Context, id uint32) (*repository.Category, error) {
	category, err := s.queries.GetCategory(ctx, int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "category not found")
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting category: %s", err.Error())
	}

	return &repository.Category{
		ID:          uint32(category.ID),
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
	}, nil
}

func (s *CategoryRepo) ListCategories(ctx context.Context) ([]*repository.Category, error) {
	categories, err := s.queries.ListCategories(ctx)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error listing categories: %s", err.Error())
	}

	categoryList := make([]*repository.Category, len(categories))
	for i, category := range categories {
		categoryList[i] = &repository.Category{
			ID:          uint32(category.ID),
			Name:        category.Name,
			Description: category.Description,
			CreatedAt:   category.CreatedAt,
		}
	}

	return categoryList, nil
}

func (s *CategoryRepo) UpdateCategory(
	ctx context.Context,
	category *repository.UpdateCategory,
) (*repository.Category, error) {
	var categoryUpdated generated.Category
	var err error

	oldCategory, err := s.queries.GetCategory(ctx, int64(category.ID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "category not found")
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting category by id: %s", err.Error())
	}

	err = s.db.ExecTx(ctx, func(q *generated.Queries) error {
		params := generated.UpdateCategoryParams{
			ID: int64(category.ID),
		}
		if category.Name != nil {
			if err := q.UpdateProductCategory(ctx, generated.UpdateProductCategoryParams{
				NewCategory: *category.Name,
				OldCategory: oldCategory.Name,
			}); err != nil {
				return pkg.Errorf(
					pkg.INTERNAL_ERROR,
					"error updating product category: %s",
					err.Error(),
				)
			}

			params.Name = pgtype.Text{String: *category.Name, Valid: true}
		}
		if category.Description != nil {
			params.Description = pgtype.Text{String: *category.Description, Valid: true}
		}

		categoryUpdated, err = q.UpdateCategory(ctx, params)
		if err != nil {
			return pkg.Errorf(pkg.INTERNAL_ERROR, "error updating category: %s", err.Error())
		}

		return nil
	})

	// if category.Name != nil {
	// 	if err := s.queries.UpdateProductCategory(ctx, generated.UpdateProductCategoryParams{
	// 		NewCategory: *category.Name,
	// 		OldCategory: oldName,
	// 	}); err != nil {
	// 		return nil, pkg.Errorf(
	// 			pkg.INTERNAL_ERROR,
	// 			"error updating product category: %s",
	// 			err.Error(),
	// 		)
	// 	}

	// 	params.Name = pgtype.Text{String: *category.Name, Valid: true}
	// }
	// if category.Description != nil {
	// 	params.Description = pgtype.Text{String: *category.Description, Valid: true}
	// }

	// categoryUpdated, err := s.queries.UpdateCategory(ctx, params)
	// if err != nil {
	// 	return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error updating category: %s", err.Error())
	// }

	return &repository.Category{
		ID:          uint32(categoryUpdated.ID),
		Name:        categoryUpdated.Name,
		Description: categoryUpdated.Description,
		CreatedAt:   categoryUpdated.CreatedAt,
	}, err
}

func (s *CategoryRepo) DeleteCategory(ctx context.Context, id uint32) error {
	if err := s.queries.DeleteCategory(ctx, int64(id)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pkg.Errorf(pkg.NOT_FOUND_ERROR, "category not found")
		}
		return pkg.Errorf(pkg.INTERNAL_ERROR, "error deleting category: %s", err.Error())
	}

	return nil
}
