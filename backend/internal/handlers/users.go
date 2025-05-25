package handlers

import (
	"net/http"
	"time"

	"github.com/flexGURU/zeiba-glam/backend/internal/repository"
	"github.com/flexGURU/zeiba-glam/backend/pkg"
	"github.com/gin-gonic/gin"
)

type createUserRequest struct {
	Name        string `json:"name"         binding:"required"`
	Email       string `json:"email"        binding:"required,email"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password"     binding:"required"`
	IsAdmin     string `json:"is_admin"     binding:"required"`
}

func (s *Server) createUserHandler(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))
		return
	}

	hashPassword, err := pkg.HashPassword(req.Password, s.config.PASSWORD_COST)
	if err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	defaultRefreshToken := "default_refresh_token"
	userPayload := &repository.User{
		Name:         req.Name,
		Email:        req.Email,
		PhoneNumber:  req.PhoneNumber,
		Password:     &hashPassword,
		IsAdmin:      req.IsAdmin == "true",
		RefreshToken: &defaultRefreshToken,
	}

	user, err := s.repo.UserRepo.CreateUser(c, userPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (s *Server) getUserHandler(c *gin.Context) {
	id, err := pkg.StringToUint32(c.Param("id"))
	if err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	user, err := s.repo.UserRepo.GetUser(c, id, "")
	if err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

type loginUserRequest struct {
	Email    string `json:"email"    binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (s *Server) loginUserHandler(c *gin.Context) {
	var req loginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))
		return
	}

	user, err := s.repo.UserRepo.GetUserInternal(c, 0, req.Email)
	if err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	if err = pkg.VerifyPassword(req.Password, *user.Password); err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	accessToken, err := s.tokenMaker.CreateToken(
		user.ID,
		user.Email,
		user.IsAdmin,
		s.config.ACCESS_TOKEN_DURATION,
	)
	if err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	refreshToken, err := s.tokenMaker.CreateToken(
		user.ID,
		user.Email,
		user.IsAdmin,
		s.config.REFRESH_TOKEN_DURATION,
	)
	if err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	c.SetCookie(
		"refreshToken",
		refreshToken,
		int(s.config.REFRESH_TOKEN_DURATION),
		"/",
		"",
		true,
		true,
	)

	_, err = s.repo.UserRepo.UpdateUser(c, &repository.UpdateUser{
		ID:           user.ID,
		RefreshToken: &refreshToken,
	})
	if err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}})
}

func (s *Server) refreshTokenHandler(c *gin.Context) {
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

	accessToken, err := s.tokenMaker.CreateToken(
		payload.UserID,
		payload.UserEmail,
		payload.IsAdmin,
		s.config.ACCESS_TOKEN_DURATION,
	)
	if err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"access_token":            accessToken,
		"access_token_expires_at": time.Now().Add(s.config.ACCESS_TOKEN_DURATION),
	}})
}

func (s *Server) logoutUserHandler(c *gin.Context) {
	c.SetCookie("refreshToken", "", -1, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{"data": "success"})
}

func (s *Server) listUsersHandler(c *gin.Context) {
	filter := &repository.UserFilter{
		Search:     nil,
		IsAdmin:    nil,
		Pagination: &pkg.Pagination{},
	}
	if search := c.Query("search"); search != "" {
		filter.Search = &search
	}
	if isAdmin := c.Query("is_admin"); isAdmin != "" {
		isAdminBool, err := pkg.StringToBool(isAdmin)
		if err != nil {
			c.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		filter.IsAdmin = &isAdminBool
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

	users, pagination, err := s.repo.UserRepo.ListUsers(c, filter)
	if err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users, "pagination": pagination})
}

type updateUserRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
	IsAdmin     string `json:"is_admin"`
}

func (s *Server) updateUserHandler(c *gin.Context) {
	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))
		return
	}

	userId, err := pkg.StringToUint32(c.Param("id"))
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

	if payload.UserID != userId {
		c.JSON(http.StatusForbidden, gin.H{"message": "You are not authorized to update this user"})
		return
	}

	isAdmin, err := pkg.StringToBool(req.IsAdmin)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.repo.UserRepo.UpdateUser(c, &repository.UpdateUser{
		ID:          userId,
		Name:        &req.Name,
		Email:       &req.Email,
		PhoneNumber: &req.PhoneNumber,
		Password:    &req.Password,
		IsAdmin:     &isAdmin,
	})
	if err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}
