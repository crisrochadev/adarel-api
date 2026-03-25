package routes

import (
	"time"

	"adarel-api/internal/handlers"
	"adarel-api/internal/middlewares"
	"adarel-api/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(authHandler *handlers.AuthHandler, contentHandler *handlers.ContentHandler, uploadHandler *handlers.UploadHandler, swaggerHandler *handlers.SwaggerHandler, authService services.AuthService) *gin.Engine {
	r := gin.New()
	r.Use(middlewares.RecoveryMiddleware())
	r.Use(middlewares.LoggingMiddleware())
	r.Use(middlewares.SecurityHeaders())
	r.Use(middlewares.RateLimitMiddleware(100, time.Minute))

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	r.Static("/uploads", "./uploads")
	r.StaticFile("/swagger/openapi.json", "./docs/openapi.json")
	r.GET("/swagger", swaggerHandler.UI)

	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	secured := r.Group("/")
	secured.Use(middlewares.AuthMiddleware(authService), middlewares.TenantMiddleware())
	{
		secured.GET("/content", contentHandler.GetByPage)
		secured.POST("/content", contentHandler.Upsert)
		secured.PUT("/content", contentHandler.Upsert)
		secured.DELETE("/content/:id", contentHandler.Delete)
		secured.POST("/upload", uploadHandler.Upload)
	}

	return r
}
