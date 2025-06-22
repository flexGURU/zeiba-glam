package handlers

import (
	"net/http"

	"github.com/flexGURU/zeiba-glam/backend/internal/repository"
	"github.com/flexGURU/zeiba-glam/backend/pkg"
	"github.com/gin-gonic/gin"
)

type createProductRequest struct {
	Name          string   `json:"name"           binding:"required"`
	Description   string   `json:"description"    binding:"required"`
	Price         float64  `json:"price"          binding:"required"`
	Category      int64    `json:"category"       binding:"required"`
	ImageURL      []string `json:"image_url"      binding:"required"`
	Size          []string `json:"size"           binding:"required"`
	Color         []string `json:"color"          binding:"required"`
	StockQuantity int64    `json:"stock_quantity"`
}

func (s *Server) createProductHandler(c *gin.Context) {
	var req createProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))
		return
	}

	payload, ok := c.Get(authorizationPayloadKey)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "missing token"})

		return
	}

	payloadData, ok := payload.(*pkg.Payload)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "incorrect token"})

		return
	}

	if !payloadData.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"message": "You are not authorized to create a product"})
		return
	}

	category, err := s.repo.CategoryRepo.GetCategory(c, uint32(req.Category))
	if err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	product, err := s.repo.ProductRepo.CreateProduct(c, &repository.Product{
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		Category:      category.Name,
		ImageURL:      req.ImageURL,
		Size:          req.Size,
		Color:         req.Color,
		StockQuantity: req.StockQuantity,
		UpdatedBy:     payloadData.UserID,
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

	if category := c.QueryArray("category"); len(category) > 0 {
		filter.Category = &category
	}

	if size := c.QueryArray("size"); len(size) > 0 {
		filter.Size = &size
	}

	if color := c.QueryArray("color"); len(color) > 0 {
		filter.Color = &color
	}

	pageNo, err := pkg.StringToUint32(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if pageNo < 1 {
		pageNo = 1
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

type updateProductRequest struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Price         float64  `json:"price"`
	Category      int64    `json:"category"`
	ImageURL      []string `json:"image_url"`
	Size          []string `json:"size"`
	Color         []string `json:"color"`
	StockQuantity int64    `json:"stock_quantity"`
}

func (s *Server) updateProductHandler(c *gin.Context) {
	var req updateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))
		return
	}

	productId, err := pkg.StringToUint32(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, ok := c.Get(authorizationPayloadKey)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "missing token"})

		return
	}

	payloadData, ok := payload.(*pkg.Payload)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "incorrect token"})

		return
	}

	if !payloadData.IsAdmin {
		c.JSON(
			http.StatusForbidden,
			gin.H{"message": "You are not authorized to update this product"},
		)
		return
	}

	updateProduct := &repository.UpdateProduct{
		ID:        productId,
		UpdatedBy: payloadData.UserID,
	}

	if req.Name != "" {
		updateProduct.Name = &req.Name
	}
	if req.Description != "" {
		updateProduct.Description = &req.Description
	}
	if req.Price != 0 {
		updateProduct.Price = &req.Price
	}
	if req.Category != 0 {
		category, err := s.repo.CategoryRepo.GetCategory(c, uint32(req.Category))
		if err != nil {
			c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
			return
		}
		updateProduct.Category = &category.Name
	}
	if len(req.ImageURL) > 0 {
		updateProduct.ImageURL = &req.ImageURL
	}
	if len(req.Size) > 0 {
		updateProduct.Size = &req.Size
	}
	if len(req.Color) > 0 {
		updateProduct.Color = &req.Color
	}
	if req.StockQuantity != 0 {
		updateProduct.StockQuantity = &req.StockQuantity
	}

	product, err := s.repo.ProductRepo.UpdateProduct(c, updateProduct)
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

	payload, ok := c.Get(authorizationPayloadKey)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "missing token"})

		return
	}

	payloadData, ok := payload.(*pkg.Payload)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "incorrect token"})

		return
	}

	if !payloadData.IsAdmin {
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
