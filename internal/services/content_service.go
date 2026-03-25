package services

import (
	"encoding/json"
	"errors"
	"strings"

	"adarel-api/internal/models"
	"adarel-api/internal/repositories"
	"adarel-api/pkg/sanitize"
	"gorm.io/datatypes"
)

type ContentService interface {
	Upsert(tenantID uint, page string, payload map[string]any) (*models.Content, error)
	GetByPage(tenantID uint, page string) (*models.Content, error)
	List(tenantID uint) ([]models.Content, error)
	Delete(tenantID uint, id uint) error
}

type contentService struct {
	repo repositories.ContentRepository
}

func NewContentService(repo repositories.ContentRepository) ContentService {
	return &contentService{repo: repo}
}

func (s *contentService) Upsert(tenantID uint, page string, payload map[string]any) (*models.Content, error) {
	page = sanitize.Text(page)
	if tenantID == 0 || page == "" || len(payload) == 0 {
		return nil, errors.New("invalid input")
	}

	sanitizeMap(payload)

	blob, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.New("invalid json payload")
	}

	content := &models.Content{
		TenantID: tenantID,
		Page:     page,
		Data:     datatypes.JSON(blob),
	}
	if err := s.repo.Upsert(content); err != nil {
		return nil, err
	}

	return s.repo.GetByPage(tenantID, page)
}

func (s *contentService) GetByPage(tenantID uint, page string) (*models.Content, error) {
	return s.repo.GetByPage(tenantID, sanitize.Text(page))
}

func (s *contentService) List(tenantID uint) ([]models.Content, error) {
	return s.repo.ListByTenant(tenantID)
}

func (s *contentService) Delete(tenantID uint, id uint) error {
	if tenantID == 0 || id == 0 {
		return errors.New("invalid input")
	}
	return s.repo.DeleteByID(tenantID, id)
}

func sanitizeMap(in map[string]any) {
	for k, v := range in {
		switch tv := v.(type) {
		case string:
			in[k] = sanitize.Text(strings.TrimSpace(tv))
		case map[string]any:
			sanitizeMap(tv)
		case []any:
			for i, item := range tv {
				if s, ok := item.(string); ok {
					tv[i] = sanitize.Text(strings.TrimSpace(s))
				}
			}
		}
	}
}
