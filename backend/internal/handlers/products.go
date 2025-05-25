package handlers

import (
	"net/http"

	"github.com/flexGURU/zeiba-glam/backend/internal/repository"
	"github.com/flexGURU/zeiba-glam/backend/pkg"
	"github.com/gin-gonic/gin"
)

type productRequest struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Price         float64  `json:"price"`
	Category      []string `json:"category"`
	ImageURL      []string `json:"image_url"`
	Size          []string `json:"size"`
	Color         []string `json:"color"`
	StockQuantity int64    `json:"stock_quantity"`
}

func (s *Server) createProductHandler(c *gin.Context) {
	var req productRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))
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
		c.JSON(http.StatusForbidden, gin.H{"message": "You are not authorized to create a product"})
		return
	}

	product, err := s.repo.ProductRepo.CreateProduct(c, &repository.Product{
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		Category:      req.Category,
		ImageURL:      req.ImageURL,
		Size:          req.Size,
		Color:         req.Color,
		StockQuantity: req.StockQuantity,
		UpdatedBy:     payload.UserID,
	})
	if err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}

func (s *Server) getProductHandler(c *gin.Context) {
	productId, err := pkg.StringToUint32(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	product, err := s.repo.ProductRepo.GetProductByID(c, productId)
	if err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}

func (s *Server) listProductsHandler(c *gin.Context) {
	filter := &repository.ProductFilter{
		Pagination: &pkg.Pagination{},
		Search:     nil,
		PriceFrom:  nil,
		PriceTo:    nil,
		Category:   nil,
		Size:       nil,
		Color:      nil,
	}

	if search := c.Query("search"); search != "" {
		filter.Search = &search
	}

	if priceFrom := c.Query("price_from"); priceFrom != "" {
		priceFromFloat, err := pkg.StringToFloat64(priceFrom)
		if err != nil {
			c.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		filter.PriceFrom = &priceFromFloat
	}

	if priceTo := c.Query("price_to"); priceTo != "" {
		priceToFloat, err := pkg.StringToFloat64(priceTo)
		if err != nil {
			c.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		filter.PriceTo = &priceToFloat
	}

	if category := c.Query("category"); category != "" {
		filter.Category = &[]string{category}
	}

	if size := c.Query("size"); size != "" {
		filter.Size = &[]string{size}
	}

	if color := c.Query("color"); color != "" {
		filter.Color = &[]string{color}
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

	products, pagination, err := s.repo.ProductRepo.ListProducts(c, filter)
	if err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": products, "pagination": pagination})
}

func (s *Server) updateProductHandler(c *gin.Context) {
	var req productRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))
		return
	}

	productId, err := pkg.StringToUint32(c.Param("id"))
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
			gin.H{"message": "You are not authorized to update this product"},
		)
		return
	}

	product, err := s.repo.ProductRepo.UpdateProduct(c, &repository.UpdateProduct{
		ID:            productId,
		UpdatedBy:     payload.UserID,
		Name:          &req.Name,
		Description:   &req.Description,
		Price:         &req.Price,
		Category:      &req.Category,
		ImageURL:      &req.ImageURL,
		Size:          &req.Size,
		Color:         &req.Color,
		StockQuantity: &req.StockQuantity,
	})
	if err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}

func (s *Server) deleteProductHandler(c *gin.Context) {
	productId, err := pkg.StringToUint32(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))
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
			gin.H{"message": "You are not authorized to delete this product"},
		)
		return
	}

	if err = s.repo.ProductRepo.DeleteProduct(c, productId); err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "success"})
}
