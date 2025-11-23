package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Setting struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	TenantID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenantId"`
	Tenant    Tenant         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Key       string         `gorm:"size:255;not null;index" json:"key"`
	Value     string         `gorm:"type:text" json:"value"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Setting) TableName() string {
	return "settings"
}
