package models

import (
	"time"

	"gorm.io/datatypes"
)

type Content struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Page      string         `gorm:"size:100;not null;index:idx_tenant_page,unique" json:"page"`
	Data      datatypes.JSON `gorm:"type:jsonb;not null" json:"data"`
	TenantID  uint           `gorm:"not null;index;index:idx_tenant_page,unique" json:"tenant_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
