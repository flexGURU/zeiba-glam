package handlers

import (
	"net/http"

	"github.com/flexGURU/zeiba-glam/backend/internal/repository"
	"github.com/flexGURU/zeiba-glam/backend/pkg"
	"github.com/gin-gonic/gin"
)

type createCategoryRequest struct {
	Name        string `json:"name"        binding:"required"`
	Description string `json:"description" binding:"required"`
}

func (s *Server) createCategoryHandler(c *gin.Context) {
	var req createCategoryRequest
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
		c.JSON(
			http.StatusForbidden,
			gin.H{"message": "You are not authorized to update this product"},
		)
		return
	}

	category, err := s.repo.CategoryRepo.CreateCategory(c, &repository.Category{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": category})
}

func (s *Server) getCategoryHandler(c *gin.Context) {
	categoryId, err := pkg.StringToUint32(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	category, err := s.repo.CategoryRepo.GetCategory(c, categoryId)
	if err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": category})
}

func (s *Server) listCategoriesHandler(c *gin.Context) {
	categories, err := s.repo.CategoryRepo.ListCategories(c)
	if err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": categories})
}

type updateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *Server) updateCategoryHandler(c *gin.Context) {
	var req updateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	categoryId, err := pkg.StringToUint32(c.Param("id"))
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
			gin.H{"message": "You are not authorized to update this category"},
		)
		return
	}

	updateCategory := &repository.UpdateCategory{
		ID: categoryId,
	}

	if req.Name != "" {
		updateCategory.Name = &req.Name
	}

	if req.Description != "" {
		updateCategory.Description = &req.Description
	}

	category, err := s.repo.CategoryRepo.UpdateCategory(c, updateCategory)
	if err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": category})
}
