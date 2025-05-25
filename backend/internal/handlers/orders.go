package handlers

import (
	"net/http"
	"strconv"

	"github.com/flexGURU/zeiba-glam/backend/internal/repository"
	"github.com/flexGURU/zeiba-glam/backend/pkg"
	"github.com/gin-gonic/gin"
)

func (s *Server) createOrderHandler(c *gin.Context) {}

func (s *Server) getOrderHandler(c *gin.Context) {
	orderId, err := pkg.StringToUint32(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	order, err := s.repo.OrderRepo.GetOrderByID(c, orderId)
	if err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}

func (s *Server) listOrdersHandler(c *gin.Context) {
	filter := &repository.OrderFilter{
		Pagination:    &pkg.Pagination{},
		Status:        nil,
		PaymentStatus: nil,
	}

	orderStatus := c.Query("order_status")
	if orderStatus != "" {
		filter.Status = &orderStatus
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

	orders, pagination, err := s.repo.OrderRepo.ListOrders(c, filter)
	if err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": orders, "pagination": pagination})
}

type updateOrderRequest struct {
	Status        string `json:"status"        `
	PaymentStatus string `json:"payment_status" `
}

func (s *Server) updateOrderHandler(c *gin.Context) {
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
		c.JSON(http.StatusForbidden, gin.H{"message": "You are not authorized to update a loan"})
		return
	}

	var req updateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	orderId, err := pkg.StringToUint32(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	updateParams := &repository.UpdateOrder{
		ID: orderId,
	}

	if req.Status != "" {
		updateParams.Status = &req.Status
	}

	if req.PaymentStatus != "" {
		paymentStatus, err := strconv.ParseBool(req.PaymentStatus)
		if err != nil {
			c.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		updateParams.PaymentStatus = &paymentStatus
	}

	order, err := s.repo.OrderRepo.UpdateOrder(c, updateParams)
	if err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}

func (s *Server) deleteOrderHandler(c *gin.Context) {
	orderId, err := pkg.StringToUint32(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

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
		c.JSON(
			http.StatusForbidden,
			gin.H{"message": "You are not authorized to delete this order"},
		)
		return
	}

	err = s.repo.OrderRepo.DeleteOrder(c, orderId)
	if err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order deleted successfully"})
}
