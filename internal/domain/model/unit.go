package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Unit struct {
	ID          uuid.UUID      `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Code        string         `gorm:"uniqueIndex:units_code_key;size:32;not null" json:"code"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	CreatedByID *uuid.UUID     `gorm:"column:created_by;type:uuid" json:"createdById,omitempty"`
	CreatedBy   *UserAccount   `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	UpdatedByID *uuid.UUID     `gorm:"column:updated_by;type:uuid" json:"updatedById,omitempty"`
	UpdatedBy   *UserAccount   `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
