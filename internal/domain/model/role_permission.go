package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RolePermission struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	RoleID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"roleId"`
	Role         Role           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	PermissionID uuid.UUID      `gorm:"type:uuid;not null;index" json:"permissionId"`
	Permission   Permission     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}