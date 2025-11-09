package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Feature struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Code        string         `gorm:"size:50;not null;uniqueIndex" json:"code"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Feature) TableName() string {
	return "features"
}