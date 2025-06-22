package repository

import "context"

type DashboardStats struct {
	TotalProducts int64 `json:"total_products"`
	InStock       int64 `json:"in_stock"`
	LowStock      int64 `json:"low_stock"`
	OutOfStock    int64 `json:"out_of_stock"`
}

type HelperRepository interface {
	GetDashboardStats(ctx context.Context) (*DashboardStats, error)
}
