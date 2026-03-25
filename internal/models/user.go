package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"size:120;not null" json:"name"`
	Email        string    `gorm:"size:160;uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	TenantID     uint      `gorm:"index;not null" json:"tenant_id"`
	Tenant       Tenant    `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}
