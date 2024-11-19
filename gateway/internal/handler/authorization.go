package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/vctrl/currency-service/gateway/internal/dto"
	"github.com/vctrl/currency-service/gateway/internal/pkg/auth"
	"github.com/vctrl/currency-service/gateway/internal/repository"

	"github.com/gin-gonic/gin"
)

func (s *Server) Register(c *gin.Context) {
	var req dto.RegisterRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = s.AuthService.Register(req)
	if err != nil {
		s.handleError(c, err)
		return
	}

	c.Status(http.StatusCreated)
}

func (s *Server) Login(c *gin.Context) {
	var req dto.LoginRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := s.AuthService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		s.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (s *Server) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization token is required"})
		return
	}

	err := s.AuthService.Logout(token)
	if err != nil {
		s.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}

// todo move errors to separate package
func (s *Server) handleError(c *gin.Context, err error) {
	log.Printf("internal error: %v", err)
	switch {
	case errors.Is(err, repository.ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	case errors.Is(err, repository.ErrUserAlreadyExist):
		c.JSON(http.StatusConflict, gin.H{"error": "User already exist"})
	case errors.Is(err, auth.ErrUnexpectedStatusCode):
		log.Printf("unexpected status code error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected server error"}) // Обычный ответ клиенту
	case errors.Is(err, auth.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	case errors.Is(err, auth.ErrTokenGeneration):
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to generate token"},
		)
	case errors.Is(err, auth.ErrTokenNotFound):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token not found"})
	case errors.Is(err, auth.ErrInvalidOrExpiredToken):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is invalid or expired"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}
}
