package middlewares

import (
	"strings"

	"adarel-api/internal/services"
	"adarel-api/pkg/response"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.Error(c, 401, "unauthorized")
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := authService.ParseToken(token)
		if err != nil {
			response.Error(c, 401, "unauthorized")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("tenant_id", claims.TenantID)
		c.Next()
	}
}
