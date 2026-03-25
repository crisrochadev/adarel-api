package handlers

import (
	"net/http"

	"adarel-api/internal/services"
	"adarel-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	uploadService services.UploadService
}

func NewUploadHandler(uploadService services.UploadService) *UploadHandler {
	return &UploadHandler{uploadService: uploadService}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "file is required")
		return
	}

	url, err := h.uploadService.SaveImage(file)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid upload")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"url": url})
}
