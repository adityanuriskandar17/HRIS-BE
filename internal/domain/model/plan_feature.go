package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlanFeature struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	PlanID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"planId"`
	Plan      Plan           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	FeatureID uuid.UUID      `gorm:"type:uuid;not null;index" json:"featureId"`
	Feature   Feature        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Enabled   bool           `gorm:"default:true" json:"enabled"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (PlanFeature) TableName() string {
	return "plan_features"
}
