package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Position struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	CompanyID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"companyId"`
	Company     Company        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Title       string         `gorm:"size:100;not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Level       int            `gorm:"default:1" json:"level"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Position) TableName() string {
	return "positions"
}
