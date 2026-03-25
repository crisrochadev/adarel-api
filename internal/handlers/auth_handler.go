package handlers

import (
	"net/http"

	"adarel-api/internal/services"
	"adarel-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type registerRequest struct {
	Name       string `json:"name" binding:"required,min=2,max=120"`
	Email      string `json:"email" binding:"required,email,max=160"`
	Password   string `json:"password" binding:"required,min=8,max=64"`
	TenantName string `json:"tenant_name" binding:"required,min=2,max=120"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email,max=160"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request")
		return
	}

	if err := h.authService.Register(req.Name, req.Email, req.Password, req.TenantName); err != nil {
		response.Error(c, http.StatusBadRequest, "could not register user")
		return
	}

	response.Success(c, http.StatusCreated, gin.H{"message": "user registered"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request")
		return
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"token": token})
}
