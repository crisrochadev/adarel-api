package middlewares

import (
	"adarel-api/pkg/response"

	"github.com/gin-gonic/gin"
)

func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, ok := c.Get("tenant_id")
		if !ok || tenantID == nil {
			response.Error(c, 403, "forbidden")
			c.Abort()
			return
		}
		c.Next()
	}
}
