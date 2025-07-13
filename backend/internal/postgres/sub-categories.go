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

var _ repository.SubCategoryRepository = (*SubCategoryRepo)(nil)

type SubCategoryRepo struct {
	queries generated.Querier
	db      *Store
}

func NewSubCategoryRepo(db *Store) *SubCategoryRepo {
	return &SubCategoryRepo{
		queries: generated.New(db.pool),
		db:      db,
	}
}

func (r *SubCategoryRepo) CreateSubCategory(ctx context.Context, subCategory *repository.SubCategory) (*repository.SubCategory, error) {
	subCategoryCreated, err := r.queries.CreateSubCategory(ctx, generated.CreateSubCategoryParams{
		CategoryID:  int64(subCategory.CategoryID),
		Name:        subCategory.Name,
		Description: subCategory.Description,
	})
	if err != nil {
		if pkg.PgxErrorCode(err) == pkg.UNIQUE_VIOLATION {
			return nil, pkg.Errorf(pkg.ALREADY_EXISTS_ERROR, "%s", err.Error())
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error creating sub-category: %s", err.Error())
	}
	subCategory.ID = uint32(subCategoryCreated.ID)
	subCategory.CreatedAt = subCategoryCreated.CreatedAt

	return subCategory, nil
}

func (r *SubCategoryRepo) GetSubCategory(ctx context.Context, id uint32) (*repository.SubCategory, error) {
	subcategory, err := r.queries.GetSubCategoryByID(ctx, int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "sub-category not found")
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting sub-category: %s", err.Error())
	}

	return &repository.SubCategory{
		ID:                  uint32(subcategory.ID),
		CategoryID:          uint32(subcategory.CategoryID),
		Name:                subcategory.Name,
		Description:         subcategory.Description,
		CreatedAt:           subcategory.CreatedAt,
		CategoryName:        &subcategory.CategoryName,
		CategoryDescription: &subcategory.CategoryDescription,
	}, nil
}

func (r *SubCategoryRepo) GetSubCategoryByCategoryIDAndID(ctx context.Context, categoryID, id uint32) (*repository.SubCategory, error) {
	subCategory, err := r.queries.GetSubCategoryByCategoryIDAndID(ctx, generated.GetSubCategoryByCategoryIDAndIDParams{
		CategoryID: int64(categoryID),
		ID:         int64(id),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "sub-category not found for category ID %d and subcategoryID %d", categoryID, id)
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting sub-category by category ID and ID: %s", err.Error())
	}

	return &repository.SubCategory{
		ID:          uint32(subCategory.ID),
		CategoryID:  uint32(subCategory.CategoryID),
		Name:        subCategory.Name,
		Description: subCategory.Description,
		CreatedAt:   subCategory.CreatedAt,
	}, nil
}

func (r *SubCategoryRepo) ListSubCategoriesByCategoryID(ctx context.Context, categoryID uint32) ([]*repository.SubCategory, error) {
	result, err := r.queries.ListSubCategoriesByCategory(ctx, int64(categoryID))
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error listing sub-categories by category ID: %s", err.Error())
	}
	subCategories := make([]*repository.SubCategory, 0, len(result))
	for _, subCategory := range result {
		subCategories = append(subCategories, &repository.SubCategory{
			ID:          uint32(subCategory.ID),
			CategoryID:  uint32(subCategory.CategoryID),
			Name:        subCategory.Name,
			Description: subCategory.Description,
			CreatedAt:   subCategory.CreatedAt,
		})
	}
	return subCategories, nil
}

func (r *SubCategoryRepo) ListSubCategories(ctx context.Context) ([]*repository.SubCategory, error) {
	result, err := r.queries.ListSubCategories(ctx)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error listing sub-categories: %s", err.Error())
	}
	subCategories := make([]*repository.SubCategory, 0, len(result))
	for _, subCategory := range result {
		subCategories = append(subCategories, &repository.SubCategory{
			ID:                  uint32(subCategory.ID),
			CategoryID:          uint32(subCategory.CategoryID),
			Name:                subCategory.Name,
			Description:         subCategory.Description,
			CreatedAt:           subCategory.CreatedAt,
			CategoryName:        &subCategory.CategoryName,
			CategoryDescription: &subCategory.CategoryDescription,
		})
	}
	return subCategories, nil
}

func (r *SubCategoryRepo) UpdateSubCategory(ctx context.Context, subCategory *repository.UpdateSubCategory) (*repository.SubCategory, error) {
	var createdSubCategory generated.SubCategory
	var err error

	oldSubCategory, err := r.queries.GetSubCategoryByID(ctx, int64(subCategory.ID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "sub-category not found")
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting sub-category by id: %s", err.Error())
	}

	err = r.db.ExecTx(ctx, func(q *generated.Queries) error {
		params := generated.UpdateSubCategoryParams{
			ID: int64(subCategory.ID),
		}

		if subCategory.Name != nil {
			if err := q.UpdateProductSubCategory(ctx, generated.UpdateProductSubCategoryParams{
				OldSubCategory: oldSubCategory.Name,
				NewSubCategory: *subCategory.Name,
			}); err != nil {
				return pkg.Errorf(pkg.INTERNAL_ERROR, "error updating product sub-category: %s", err.Error())
			}

			params.Name = pgtype.Text{String: *subCategory.Name, Valid: true}
		}

		if subCategory.Description != nil {
			params.Description = pgtype.Text{String: *subCategory.Description, Valid: true}
		}

		createdSubCategory, err = q.UpdateSubCategory(ctx, params)
		if err != nil {
			return pkg.Errorf(pkg.INTERNAL_ERROR, "error updating sub-category: %s", err.Error())
		}
		return nil
	})

	return &repository.SubCategory{
		ID:          uint32(createdSubCategory.ID),
		CategoryID:  uint32(createdSubCategory.CategoryID),
		Name:        createdSubCategory.Name,
		Description: createdSubCategory.Description,
		CreatedAt:   createdSubCategory.CreatedAt,
	}, err
}

func (r *SubCategoryRepo) DeleteSubCategory(ctx context.Context, id uint32) error {
	if err := r.queries.DeleteSubCategory(ctx, int64(id)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pkg.Errorf(pkg.NOT_FOUND_ERROR, "sub-category not found")
		}
		return pkg.Errorf(pkg.INTERNAL_ERROR, "error deleting sub-category: %s", err.Error())
	}
	return nil
}
