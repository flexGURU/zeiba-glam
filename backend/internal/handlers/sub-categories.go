package handlers

import (
	"net/http"

	"github.com/flexGURU/zeiba-glam/backend/internal/repository"
	"github.com/flexGURU/zeiba-glam/backend/pkg"
	"github.com/gin-gonic/gin"
)

type createSubCategoryRequest struct {
	Name        string `json:"name"        binding:"required"`
	Description string `json:"description" binding:"required"`
	CategoryID  uint32 `json:"category_id" binding:"required"`
}

func (s *Server) createSubCategoryHandler(ctx *gin.Context) {
	var req createSubCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))
		return
	}

	payload, ok := ctx.Get(authorizationPayloadKey)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "missing token"})
		return
	}

	payloadData, ok := payload.(*pkg.Payload)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "incorrect token"})
		return
	}

	if !payloadData.IsAdmin {
		ctx.JSON(
			http.StatusForbidden,
			gin.H{"message": "You are not authorized to create a sub-category"},
		)
		return
	}

	subCategory, err := s.repo.SubCategoryRepo.CreateSubCategory(ctx, &repository.SubCategory{
		Name:        req.Name,
		Description: req.Description,
		CategoryID:  req.CategoryID,
	})
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": subCategory})
}

func (s *Server) getSubCategoryHandler(ctx *gin.Context) {
	subCategoryId, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	subCategory, err := s.repo.SubCategoryRepo.GetSubCategory(ctx, subCategoryId)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": subCategory})
}

func (s *Server) listSubcategoriesHandler(ctx *gin.Context) {
	subCategories, err := s.repo.SubCategoryRepo.ListSubCategories(ctx)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": subCategories})
}

func (s *Server) listSubcategoriesByCategoryHandler(ctx *gin.Context) {
	categoryId, err := pkg.StringToUint32(ctx.Param("category_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	subCategories, err := s.repo.SubCategoryRepo.ListSubCategoriesByCategoryID(ctx, categoryId)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": subCategories})
}

type updateSubCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CategoryID  uint32 `json:"category_id"`
}

func (s *Server) updateSubCategoryHandler(ctx *gin.Context) {
	var req updateSubCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	subCategoryId, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, ok := ctx.Get(authorizationPayloadKey)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "missing token"})

		return
	}

	payloadData, ok := payload.(*pkg.Payload)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "incorrect token"})

		return
	}

	if !payloadData.IsAdmin {
		ctx.JSON(
			http.StatusForbidden,
			gin.H{"message": "You are not authorized to update this category"},
		)
		return
	}

	updateSubCategory := &repository.UpdateSubCategory{
		ID: subCategoryId,
	}

	if req.Name != "" {
		updateSubCategory.Name = &req.Name
	}
	if req.Description != "" {
		updateSubCategory.Description = &req.Description
	}

	subCategory, err := s.repo.SubCategoryRepo.UpdateSubCategory(ctx, updateSubCategory)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": subCategory})
}
