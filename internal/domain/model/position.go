package model

import (
	"time"

	"gorm.io/gorm"
)

type Position struct {
	ID          uint64         `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"size:100;not null" json:"title"`
	UnitID      *uint64        `json:"unitId"`
	Unit        *Unit          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"unit,omitempty"`
	CreatedByID *uint64        `gorm:"column:created_by" json:"createdById,omitempty"`
	CreatedBy   *UserAccount   `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	UpdatedByID *uint64        `gorm:"column:updated_by" json:"updatedById,omitempty"`
	UpdatedBy   *UserAccount   `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
