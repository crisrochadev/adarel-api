package repositories

import (
	"adarel-api/internal/models"

	"gorm.io/gorm"
)

type ContentRepository interface {
	Upsert(content *models.Content) error
	GetByPage(tenantID uint, page string) (*models.Content, error)
	ListByTenant(tenantID uint) ([]models.Content, error)
	DeleteByID(tenantID uint, id uint) error
}

type contentRepository struct {
	db *gorm.DB
}

func NewContentRepository(db *gorm.DB) ContentRepository {
	return &contentRepository{db: db}
}

func (r *contentRepository) Upsert(content *models.Content) error {
	return r.db.Where("tenant_id = ? AND page = ?", content.TenantID, content.Page).
		Assign(map[string]any{"data": content.Data}).
		FirstOrCreate(content).Error
}

func (r *contentRepository) GetByPage(tenantID uint, page string) (*models.Content, error) {
	var content models.Content
	if err := r.db.Where("tenant_id = ? AND page = ?", tenantID, page).First(&content).Error; err != nil {
		return nil, err
	}
	return &content, nil
}

func (r *contentRepository) ListByTenant(tenantID uint) ([]models.Content, error) {
	var contents []models.Content
	if err := r.db.Where("tenant_id = ?", tenantID).Order("updated_at desc").Find(&contents).Error; err != nil {
		return nil, err
	}
	return contents, nil
}

func (r *contentRepository) DeleteByID(tenantID uint, id uint) error {
	return r.db.Where("tenant_id = ?", tenantID).Delete(&models.Content{}, id).Error
}
