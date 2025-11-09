package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Permission struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Resource  string         `gorm:"size:100;not null" json:"resource"`
	Action    string         `gorm:"size:50;not null" json:"action"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Permission) TableName() string {
	return "permissions"
}