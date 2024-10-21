package handler

import (
	"net/http"

	"github.com/vctrl/currency-service/gateway/internal/service"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type Server struct {
	AuthService     service.AuthService
	CurrencyService service.CurrencyService
	Router          *gin.Engine
	logger          *zap.Logger
}

func NewServer(authSvc service.AuthService,
	currencySvc service.CurrencyService,
	router *gin.Engine,
	logger *zap.Logger) Server {

	return Server{
		AuthService:     authSvc,
		CurrencyService: currencySvc,
		Router:          router,
		logger:          logger,
	}
}

func (s *Server) RegisterRoutes() {
	s.Router.GET(
		"/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		},
	)

	s.Router.GET("/api/v1/rate", s.GetCurrencyRates)
	s.Router.POST("/api/v1/login", s.Login)
	s.Router.POST("/api/v1/register", s.Register)
	s.Router.POST("/api/v1/logout", s.Logout)
}
