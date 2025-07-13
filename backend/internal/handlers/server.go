package handlers

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/flexGURU/zeiba-glam/backend/internal/postgres"
	"github.com/flexGURU/zeiba-glam/backend/pkg"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	ln     net.Listener
	srv    *http.Server

	config     pkg.Config
	tokenMaker pkg.JWTMaker
	repo       *postgres.PostgresRepo
}

func NewServer(config pkg.Config, tokenMaker pkg.JWTMaker, repo *postgres.PostgresRepo) *Server {
	if config.ENVIRONMENT == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	s := &Server{
		router: router,
		ln:     nil,
		srv:    nil,

		config:     config,
		tokenMaker: tokenMaker,
		repo:       repo,
	}

	s.setUpRoutes()

	return s
}

func (s *Server) setUpRoutes() {
	s.router.Use(CORSmiddleware(s.config.FRONTEND_URL))
	v1 := s.router.Group("/api/v1")

	v1Auth := s.router.Group("/api/v1")
	authRoute := v1Auth.Use(authMiddleware(s.tokenMaker))
	v1.GET("/health", s.healthCheckHandler)

	// auth routes
	v1.POST("/auth/login", s.loginUserHandler)
	v1.POST("/auth/refresh", s.refreshTokenHandler)
	v1.GET("/auth/logout", s.logoutUserHandler)

	// users routes
	authRoute.POST("/users", s.createUserHandler)
	authRoute.GET("/users/:id", s.getUserHandler)
	authRoute.GET("/users", s.listUsersHandler)
	authRoute.PATCH("/users/:id", s.updateUserHandler)

	// products routes
	authRoute.POST("/products", s.createProductHandler)
	v1.GET("/products/:id", s.getProductHandler)
	v1.GET("/products", s.listProductsHandler)
	authRoute.PATCH("/products/:id", s.updateProductHandler)
	authRoute.DELETE("/products/:id", s.deleteProductHandler)

	// categories route
	authRoute.POST("/categories", s.createCategoryHandler)
	v1.GET("/categories/:id", s.getCategoryHandler)
	v1.GET("/categories", s.listCategoriesHandler)
	authRoute.PATCH("/categories/:id", s.updateCategoryHandler)

	// sub-categories route
	authRoute.POST("/sub-categories", s.createSubCategoryHandler)
	v1.GET("/sub-categories/:id", s.getSubCategoryHandler)
	v1.GET("/sub-categories", s.listSubcategoriesHandler)
	v1.GET("/sub-categories/category/:category_id", s.listSubcategoriesByCategoryHandler)
	authRoute.PATCH("/sub-categories/:id", s.updateSubCategoryHandler)

	// orders routes
	v1.POST("/orders", s.createOrderHandler)
	v1.GET("/orders/:id", s.getOrderHandler)
	v1.GET("/orders", s.listOrdersHandler)
	v1.PATCH("/orders/:id", s.updateOrderHandler)
	v1.DELETE("/orders/:id", s.deleteOrderHandler)

	// payments routes
	v1.POST("/payments", s.createPaymentHandler)
	v1.GET("/payments/:id", s.getPaymentByIDHandler)
	v1.GET("/payments/order/:id", s.getPaymentByOrderIDHandler)
	v1.GET("/payments", s.listPaymentsHandler)
	v1.PATCH("/payments/:id", s.updatePaymentHandler)

	// helpers route
	authRoute.GET("/helpers/dashboard-stats", s.getDashboardStatsHandler)

	s.srv = &http.Server{
		Addr:         s.config.SERVER_ADDRESS,
		Handler:      s.router.Handler(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func (s *Server) Start() error {
	var err error
	if s.ln, err = net.Listen("tcp", s.config.SERVER_ADDRESS); err != nil {
		return err
	}

	go func(s *Server) {
		err := s.srv.Serve(s.ln)
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}(s)

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	log.Println("Shutting down http server...")

	return s.srv.Shutdown(ctx)
}

// GetPort returns the port the server is listening on for testing purposes
func (s *Server) GetPort() int {
	if s.ln == nil {
		return 0
	}

	return s.ln.Addr().(*net.TCPAddr).Port
}

func (s *Server) healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func errorResponse(err error) gin.H {
	return gin.H{
		"status_code": pkg.ErrorCode(err),
		"message":     pkg.ErrorMessage(err),
	}
}
