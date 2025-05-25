package postgres

import (
	"context"

	"github.com/flexGURU/zeiba-glam/backend/internal/postgres/generated"
	"github.com/flexGURU/zeiba-glam/backend/internal/repository"
	"github.com/flexGURU/zeiba-glam/backend/pkg"
	"github.com/jackc/pgx/v5/pgtype"
)

var _ repository.PaymentRepository = (*PaymentRepo)(nil)

// PaymentRepo is the repository for the payment model
type PaymentRepo struct {
	queries generated.Querier
}

func NewPaymentRepo(db *Store) *PaymentRepo {
	return &PaymentRepo{
		queries: generated.New(db.pool),
	}
}

func (r *PaymentRepo) CreatePayment(
	ctx context.Context,
	payment *repository.Payment,
) (*repository.Payment, error) {
	createdPayment, err := r.queries.CreatePayment(ctx, generated.CreatePaymentParams{
		OrderID:       int64(payment.OrderID),
		Amount:        pkg.Float64ToPgTypeNumeric(payment.Amount),
		TransactionID: payment.TransactionID,
		PaymentMethod: payment.PaymentMethod,
		PaymentStatus: payment.PaymentStatus,
		PaidAt:        payment.PaidAt,
	})
	if err != nil {
		if pkg.PgxErrorCode(err) == pkg.UNIQUE_VIOLATION {
			return nil, pkg.Errorf(pkg.ALREADY_EXISTS_ERROR, "%s", err.Error())
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error creating payment: %s", err.Error())
	}

	payment.ID = uint32(createdPayment.ID)
	payment.CreatedAt = createdPayment.CreatedAt

	return payment, nil
}

func (r *PaymentRepo) GetPayment(
	ctx context.Context,
	id uint32,
	orderID uint32,
) (*repository.Payment, error) {
	getPaymentParams := generated.GetPaymentParams{}
	if id != 0 {
		getPaymentParams.ID = id
	}

	if orderID != 0 {
		getPaymentParams.OrderID = orderID
	}

	if id == 0 && orderID == 0 {
		return nil, pkg.Errorf(pkg.INVALID_ERROR, "id or order_id is required")
	}

	payment, err := r.queries.GetPayment(ctx, getPaymentParams)
	if err != nil {
		if pkg.PgxErrorCode(err) == pkg.NOT_FOUND_ERROR {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "payment not found")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting payment by id: %s", err.Error())
	}

	orderItems, err := r.queries.ListOrderItems(ctx, payment.OrderID)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting order items: %s", err.Error())
	}

	orderItemsResult := make([]*repository.OrderItem, len(orderItems))
	for idx, orderItem := range orderItems {
		orderItemsResult[idx] = marshalOrderItem(orderItem)
	}

	return &repository.Payment{
		ID:            uint32(payment.ID),
		OrderID:       uint32(payment.OrderID),
		Amount:        pkg.PgTypeNumericToFloat64(payment.Amount),
		TransactionID: payment.TransactionID,
		PaymentMethod: payment.PaymentMethod,
		PaymentStatus: payment.PaymentStatus,
		PaidAt:        payment.PaidAt,
		CreatedAt:     payment.CreatedAt,
		OrderDetails: &repository.Order{
			ID:              uint32(payment.OrderID),
			UserName:        payment.OrderUserName,
			UserPhoneNumber: payment.OrderUserPhoneNumber,
			TotalAmount:     pkg.PgTypeNumericToFloat64(payment.OrderTotalAmount),
			Status:          payment.OrderStatus,
			ShippingAddress: payment.OrderShippingAddress,
			PaymentStatus:   payment.OrderPaymentStatus,
			CreatedAt:       payment.OrderCreatedAt,

			OrderItems: orderItemsResult,
		},
	}, nil
}

func (r *PaymentRepo) ListPayments(
	ctx context.Context,
	filter *repository.PaymentFilter,
) ([]*repository.Payment, *pkg.Pagination, error) {
	paramListPayments := generated.ListPaymentsParams{
		Limit:  int32(filter.Pagination.PageSize),
		Offset: pkg.Offset(filter.Pagination.Page, filter.Pagination.PageSize),
	}

	paramListPaymentsCount := generated.ListPaymentsCountParams{}

	if filter.PaymentMethod != nil {
		paramListPayments.PaymentMethod = pgtype.Text{
			String: *filter.PaymentMethod,
			Valid:  true,
		}
		paramListPaymentsCount.PaymentMethod = pgtype.Text{
			String: *filter.PaymentMethod,
			Valid:  true,
		}
	}

	if filter.PaymentStatus != nil {
		paramListPayments.PaymentStatus = pgtype.Bool{
			Bool:  *filter.PaymentStatus,
			Valid: true,
		}
		paramListPaymentsCount.PaymentStatus = pgtype.Bool{
			Bool:  *filter.PaymentStatus,
			Valid: true,
		}
	}

	payments, err := r.queries.ListPayments(ctx, paramListPayments)
	if err != nil {
		return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error listing payments: %s", err.Error())
	}

	totalPayments, err := r.queries.ListPaymentsCount(ctx, paramListPaymentsCount)
	if err != nil {
		return nil, nil, pkg.Errorf(
			pkg.INTERNAL_ERROR,
			"error listing payments count: %s",
			err.Error(),
		)
	}

	paymentsResult := make([]*repository.Payment, len(payments))
	for idx, payment := range payments {
		paymentsResult[idx] = marshalPayment(payment)
	}

	return paymentsResult, pkg.CalculatePagination(
		uint32(totalPayments),
		filter.Pagination.PageSize,
		filter.Pagination.Page,
	), nil
}

func (r *PaymentRepo) UpdatePayment(
	ctx context.Context,
	payment *repository.UpdatePayment,
) (*repository.Payment, error) {
	params := generated.UpdatePaymentParams{
		ID: int64(payment.ID),
	}

	if payment.PaymentStatus != nil {
		params.PaymentStatus = pgtype.Bool{
			Bool:  *payment.PaymentStatus,
			Valid: true,
		}
	}

	if payment.PaidAt != nil {
		params.PaidAt = pgtype.Timestamptz{
			Time:  *payment.PaidAt,
			Valid: true,
		}
	}

	updatedPayment, err := r.queries.UpdatePayment(ctx, params)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error updating payment: %s", err.Error())
	}

	return marshalPayment(updatedPayment), nil
}

func marshalPayment(payment generated.Payment) *repository.Payment {
	return &repository.Payment{
		ID:            uint32(payment.ID),
		OrderID:       uint32(payment.OrderID),
		Amount:        pkg.PgTypeNumericToFloat64(payment.Amount),
		TransactionID: payment.TransactionID,
		PaymentMethod: payment.PaymentMethod,
		PaymentStatus: payment.PaymentStatus,
		CreatedAt:     payment.CreatedAt,
	}
}
