package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BillingCycle string

const (
	BillingCycleMonthly BillingCycle = "monthly"
	BillingCycleYearly  BillingCycle = "yearly"
)

type Plan struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name         string         `gorm:"size:100;not null" json:"name"`
	Description  string         `gorm:"type:text" json:"description"`
	Price        float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	BillingCycle BillingCycle   `gorm:"type:varchar(10);not null" json:"billingCycle"`
	MaxUsers     int            `gorm:"default:10" json:"maxUsers"`
	IsActive     bool           `gorm:"default:true" json:"isActive"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Plan) TableName() string {
	return "plans"
}
