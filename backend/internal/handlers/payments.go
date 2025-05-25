package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/flexGURU/zeiba-glam/backend/internal/repository"
	"github.com/flexGURU/zeiba-glam/backend/pkg"
	"github.com/gin-gonic/gin"
)

type createPaymentRequest struct {
	OrderID       uint32  `json:"order_id"       binding:"required"`
	Amount        float64 `json:"amount"         binding:"required"`
	TransactionID string  `json:"transaction_id" binding:"required"`
	PaymentMethod string  `json:"payment_method" binding:"required,oneof=cash on_delivery"`
	PaymentStatus string  `json:"payment_status" binding:"required"`
	PaidAt        string  `json:"paid_at"`
}

func (s *Server) createPaymentHandler(c *gin.Context) {
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Refresh token not found"})
		return
	}

	payload, err := s.tokenMaker.VerifyToken(refreshToken)
	if err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	if !payload.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"message": "You are not authorized to create a payment"})
		return
	}

	var req createPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var paidAt time.Time
	if req.PaidAt != "" {
		paidAt, err = time.Parse(time.RFC3339, req.PaidAt)
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid paid_at: %s", err.Error())),
			)
			return
		}
	} else {
		paidAt = time.Now()
	}

	paymentStatus, err := strconv.ParseBool(req.PaymentStatus)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payment, err := s.repo.PaymentRepo.CreatePayment(c, &repository.Payment{
		OrderID:       req.OrderID,
		Amount:        req.Amount,
		TransactionID: req.TransactionID,
		PaymentMethod: req.PaymentMethod,
		PaymentStatus: paymentStatus,
		PaidAt:        paidAt,
	})
	if err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": payment})
}

func (s *Server) getPaymentByIDHandler(c *gin.Context) {
	paymentId, err := pkg.StringToUint32(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payment, err := s.repo.PaymentRepo.GetPayment(c, paymentId, 0)
	if err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": payment})
}

func (s *Server) getPaymentByOrderIDHandler(c *gin.Context) {
	orderId, err := pkg.StringToUint32(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payment, err := s.repo.PaymentRepo.GetPayment(c, 0, orderId)
	if err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": payment})
}

func (s *Server) listPaymentsHandler(c *gin.Context) {
	filter := &repository.PaymentFilter{
		Pagination:    &pkg.Pagination{},
		PaymentMethod: nil,
		PaymentStatus: nil,
	}

	paymentMethod := c.Query("payment_method")
	if paymentMethod != "" {
		filter.PaymentMethod = &paymentMethod
	}

	paymentStatus := c.Query("payment_status")
	if paymentStatus != "" {
		paymentStatusBool, err := strconv.ParseBool(paymentStatus)
		if err != nil {
			c.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		filter.PaymentStatus = &paymentStatusBool
	}

	pageNo, err := pkg.StringToUint32(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	filter.Pagination.Page = pageNo

	pageSize, err := pkg.StringToUint32(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	filter.Pagination.PageSize = pageSize

	payments, pagination, err := s.repo.PaymentRepo.ListPayments(c, filter)
	if err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": payments, "pagination": pagination})
}

type updatePaymentRequest struct {
	PaymentStatus string `json:"payment_status" binding:"required"`
	PaidAt        string `json:"paid_at"        binding:"required"`
}

// payment for cash and on_delivery payment is updated by admin
func (s *Server) updatePaymentHandler(c *gin.Context) {
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Refresh token not found"})
		return
	}

	payload, err := s.tokenMaker.VerifyToken(refreshToken)
	if err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	if !payload.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"message": "You are not authorized to update a payment"})
		return
	}

	var req updatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	paymentId, err := pkg.StringToUint32(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payment, err := s.repo.PaymentRepo.GetPayment(c, paymentId, 0)
	if err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	if payment.PaymentMethod != "cash" && payment.PaymentMethod != "on_delivery" {
		c.JSON(
			http.StatusBadRequest,
			errorResponse(
				pkg.Errorf(pkg.INVALID_ERROR, "payment method is not cash or on_delivery"),
			),
		)
		return
	}

	var paidAt time.Time
	if req.PaidAt != "" {
		paidAt, err = time.Parse(time.RFC3339, req.PaidAt)
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid paid_at: %s", err.Error())),
			)
			return
		}
	} else {
		paidAt = time.Now()
	}

	paymentStatus, err := strconv.ParseBool(req.PaymentStatus)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payment, err = s.repo.PaymentRepo.UpdatePayment(c, &repository.UpdatePayment{
		ID:            paymentId,
		PaymentStatus: &paymentStatus,
		PaidAt:        &paidAt,
	})
	if err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": payment})
}
