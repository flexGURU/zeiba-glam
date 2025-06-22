package handlers

import (
	"net/http"

	"github.com/flexGURU/zeiba-glam/backend/pkg"
	"github.com/gin-gonic/gin"
)

func (s *Server) getDashboardStatsHandler(c *gin.Context) {
	stats, err := s.repo.HelperRepo.GetDashboardStats(c)
	if err != nil {
		c.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}
