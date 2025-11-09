package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditLog struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	TenantID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenantId"`
	Tenant      Tenant         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	User        UserAccount    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Action      string         `gorm:"size:50;not null" json:"action"`
	Resource    string         `gorm:"size:100;not null" json:"resource"`
	ResourceID  *uuid.UUID     `gorm:"type:uuid" json:"resourceId"`
	OldValues   string         `gorm:"type:text" json:"oldValues"`
	NewValues   string         `gorm:"type:text" json:"newValues"`
	IPAddress   string         `gorm:"size:45" json:"ipAddress"`
	UserAgent   string         `gorm:"size:255" json:"userAgent"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}