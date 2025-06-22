package postgres

import (
	"context"

	"github.com/flexGURU/zeiba-glam/backend/internal/postgres/generated"
	"github.com/flexGURU/zeiba-glam/backend/internal/repository"
	"github.com/flexGURU/zeiba-glam/backend/pkg"
)

var _ repository.HelperRepository = (*HelperRepo)(nil)

type HelperRepo struct {
	queries generated.Querier
}

func NewHelperRepo(db *Store) *HelperRepo {
	return &HelperRepo{
		queries: generated.New(db.pool),
	}
}

func (r *HelperRepo) GetDashboardStats(ctx context.Context) (*repository.DashboardStats, error) {
	stats, err := r.queries.GetDashboardStats(ctx)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting dashboard stats")
	}

	return &repository.DashboardStats{
		TotalProducts: stats.TotalProducts,
		InStock:       stats.InStock,
		LowStock:      stats.LowStock,
		OutOfStock:    stats.OutOfStock,
	}, nil
}
