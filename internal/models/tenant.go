package models

import "time"

type Tenant struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:120;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
