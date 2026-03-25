package repositories

import (
	"adarel-api/internal/models"

	"gorm.io/gorm"
)

type TenantRepository interface {
	Create(tenant *models.Tenant) error
	FindByID(id uint) (*models.Tenant, error)
}

type tenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) TenantRepository {
	return &tenantRepository{db: db}
}

func (r *tenantRepository) Create(tenant *models.Tenant) error {
	return r.db.Create(tenant).Error
}

func (r *tenantRepository) FindByID(id uint) (*models.Tenant, error) {
	var tenant models.Tenant
	if err := r.db.First(&tenant, id).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}
