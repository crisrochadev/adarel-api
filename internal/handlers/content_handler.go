package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"adarel-api/internal/services"
	"adarel-api/pkg/response"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type ContentHandler struct {
	contentService services.ContentService
}

func NewContentHandler(contentService services.ContentService) *ContentHandler {
	return &ContentHandler{contentService: contentService}
}

type upsertContentRequest struct {
	Page string         `json:"page" binding:"required,max=100"`
	Data map[string]any `json:"data" binding:"required"`
}

func (h *ContentHandler) Upsert(c *gin.Context) {
	var req upsertContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request")
		return
	}

	tenantID := c.GetUint("tenant_id")
	content, err := h.contentService.Upsert(tenantID, req.Page, req.Data)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "could not save content")
		return
	}

	response.Success(c, http.StatusOK, content)
}

func (h *ContentHandler) GetByPage(c *gin.Context) {
	tenantID := c.GetUint("tenant_id")
	page := c.Query("page")
	if page == "" {
		items, err := h.contentService.List(tenantID)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "could not list contents")
			return
		}
		response.Success(c, http.StatusOK, items)
		return
	}

	content, err := h.contentService.GetByPage(tenantID, page)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "content not found")
			return
		}
		response.Error(c, http.StatusBadRequest, "could not fetch content")
		return
	}

	response.Success(c, http.StatusOK, content)
}

func (h *ContentHandler) Delete(c *gin.Context) {
	tenantID := c.GetUint("tenant_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.contentService.Delete(tenantID, uint(id)); err != nil {
		response.Error(c, http.StatusBadRequest, "could not delete content")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "content deleted"})
}
