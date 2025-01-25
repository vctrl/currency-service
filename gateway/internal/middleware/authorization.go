package middleware

import (
	"net/http"
	"strings"

	"github.com/vctrl/currency-service/gateway/internal/clients/auth"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type Authorization struct {
	authClient *auth.Client //todo interface
	skipper    func(*gin.Context) bool
	logger     *zap.Logger
}

func NewAuthorization(authClient *auth.Client, skipper func(*gin.Context) bool, logger *zap.Logger) Authorization {
	return Authorization{
		authClient: authClient,
		skipper:    skipper,
		logger:     logger,
	}
}

func (auth *Authorization) Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		if auth.skipper(c) {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if err := auth.authClient.ValidateToken(c.Request.Context(), authHeaderParts[1]); err != nil {
			auth.logger.Error(
				"Invalid token",
				zap.String("token", authHeaderParts[1]),
				zap.String("client_ip", c.ClientIP()),
				zap.String("user_agent", c.GetHeader("User-Agent")),
				zap.Error(err),
			)

			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Next()
	}
}

func shouldSkipMiddleware(c *gin.Context) bool {
	// Check if the current route should skip the middleware
	if c.Request.URL.Path == "/login" {
		return true // Skip the middleware
	}
	return false // Execute the middleware
}
