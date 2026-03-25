package main

import (
	"log"

	"adarel-api/internal/config"
	"adarel-api/internal/database"
	"adarel-api/internal/handlers"
	"adarel-api/internal/repositories"
	"adarel-api/internal/routes"
	"adarel-api/internal/services"
)

func main() {
	cfg := config.Load()
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("database connection failed")
	}

	userRepo := repositories.NewUserRepository(db)
	tenantRepo := repositories.NewTenantRepository(db)
	contentRepo := repositories.NewContentRepository(db)

	authService := services.NewAuthService(userRepo, tenantRepo, cfg.JWTSecret)
	contentService := services.NewContentService(contentRepo)
	uploadService := services.NewUploadService("./uploads")

	authHandler := handlers.NewAuthHandler(authService)
	contentHandler := handlers.NewContentHandler(contentService)
	uploadHandler := handlers.NewUploadHandler(uploadService)

	r := routes.SetupRouter(authHandler, contentHandler, uploadHandler, authService)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server startup failed")
	}
}
