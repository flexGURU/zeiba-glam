package postgres

import (
	"context"
	"database/sql"

	"github.com/flexGURU/zeiba-glam/backend/internal/postgres/generated"
	"github.com/flexGURU/zeiba-glam/backend/internal/repository"
	"github.com/flexGURU/zeiba-glam/backend/pkg"
	"github.com/jackc/pgx/v5/pgtype"
)

// var _ repository.OrderRepository = (*OrderRepo)(nil)

type OrderRepo struct {
	db      *Store
	queries *generated.Queries
}

func NewOrderRepo(db *Store) *OrderRepo {
	return &OrderRepo{db: db, queries: generated.New(db.pool)}
}

// func (r *OrderRepo) CreateOrder(
// 	ctx context.Context,
// 	order *repository.Order,
// ) (*repository.Order, error) {
// }

func (r *OrderRepo) GetOrderByID(ctx context.Context, id uint32) (*repository.Order, error) {
	order, err := r.queries.GetOrderByID(ctx, int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "order not found")
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting order: %s", err.Error())
	}

	orderItems, err := r.queries.ListOrderItems(ctx, int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "order items not found")
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting order items: %s", err.Error())
	}

	payments, err := r.queries.GetPaymentsOverviewByOrderID(ctx, int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "payments not found")
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting payments: %s", err.Error())
	}

	orderItemsResult := make([]*repository.OrderItem, len(orderItems))
	for idx, orderItem := range orderItems {
		orderItemsResult[idx] = marshalOrderItem(orderItem)
	}

	paymentsResult := make([]*repository.Payment, len(payments))
	for idx, payment := range payments {
		paymentsResult[idx] = marshalPayment(payment)
	}

	rsp := marshalOrder(order)
	rsp.OrderItems = orderItemsResult
	rsp.Payments = paymentsResult

	return rsp, nil
}

func (r *OrderRepo) ListOrders(
	ctx context.Context,
	filter *repository.OrderFilter,
) ([]*repository.Order, *pkg.Pagination, error) {
	paramListOrders := generated.ListOrdersParams{
		Limit:  int32(filter.Pagination.PageSize),
		Offset: pkg.Offset(filter.Pagination.Page, filter.Pagination.PageSize),
	}

	paramListOrdersCount := generated.ListOrdersCountParams{}

	if filter.Status != nil {
		paramListOrders.Status = pgtype.Text{
			String: *filter.Status,
			Valid:  true,
		}

		paramListOrdersCount.Status = pgtype.Text{
			String: *filter.Status,
			Valid:  true,
		}
	}

	if filter.PaymentStatus != nil {
		paramListOrders.PaymentStatus = pgtype.Bool{
			Bool:  *filter.PaymentStatus,
			Valid: true,
		}

		paramListOrdersCount.PaymentStatus = pgtype.Bool{
			Bool:  *filter.PaymentStatus,
			Valid: true,
		}
	}

	orders, err := r.queries.ListOrders(ctx, paramListOrders)
	if err != nil {
		return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error listing orders: %s", err.Error())
	}

	totalOrders, err := r.queries.ListOrdersCount(ctx, paramListOrdersCount)
	if err != nil {
		return nil, nil, pkg.Errorf(
			pkg.INTERNAL_ERROR,
			"error listing orders count: %s",
			err.Error(),
		)
	}

	ordersResult := make([]*repository.Order, len(orders))
	for idx, order := range orders {
		ordersResult[idx] = marshalOrder(order)
	}

	return ordersResult, pkg.CalculatePagination(
		uint32(totalOrders),
		filter.Pagination.PageSize,
		filter.Pagination.Page,
	), nil
}

func (r *OrderRepo) UpdateOrder(
	ctx context.Context,
	order *repository.UpdateOrder,
) (*repository.Order, error) {
	params := generated.UpdateOrderParams{
		ID: int64(order.ID),
	}

	if order.Status != nil {
		params.Status = pgtype.Text{
			String: *order.Status,
			Valid:  true,
		}
	}

	if order.PaymentStatus != nil {
		params.PaymentStatus = pgtype.Bool{
			Bool:  *order.PaymentStatus,
			Valid: true,
		}
	}

	updatedOrder, err := r.queries.UpdateOrder(ctx, params)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error updating order: %s", err.Error())
	}

	return marshalOrder(updatedOrder), nil
}

func (r *OrderRepo) DeleteOrder(ctx context.Context, id uint32) error {
	if err := r.queries.DeleteOrder(ctx, int64(id)); err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "error deleting order: %s", err.Error())
	}

	return nil
}

func marshalOrder(order generated.Order) *repository.Order {
	return &repository.Order{
		ID:              uint32(order.ID),
		UserName:        order.UserName,
		UserPhoneNumber: order.UserPhoneNumber,
		TotalAmount:     pkg.PgTypeNumericToFloat64(order.TotalAmount),
		Status:          order.Status,
		ShippingAddress: order.ShippingAddress,
		PaymentStatus:   order.PaymentStatus,
		CreatedAt:       order.CreatedAt,
	}
}

func marshalOrderItem(orderItem generated.OrderItem) *repository.OrderItem {
	return &repository.OrderItem{
		ID:        uint32(orderItem.ID),
		OrderID:   uint32(orderItem.OrderID),
		ProductID: uint32(orderItem.ProductID),
		Quantity:  int64(orderItem.Quantity),
		Amount:    pkg.PgTypeNumericToFloat64(orderItem.Amount),
		Size:      orderItem.Size,
		Color:     orderItem.Color,
	}
}
