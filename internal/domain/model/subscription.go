package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionStatus string

const (
	SubscriptionStatusActive   SubscriptionStatus = "active"
	SubscriptionStatusExpired  SubscriptionStatus = "expired"
	SubscriptionStatusCanceled SubscriptionStatus = "canceled"
)

type Subscription struct {
	ID        uuid.UUID          `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	TenantID  uuid.UUID          `gorm:"type:uuid;not null;index" json:"tenantId"`
	Tenant    Tenant             `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	PlanID    uuid.UUID          `gorm:"type:uuid;not null;index" json:"planId"`
	Plan      Plan               `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	StartDate time.Time          `gorm:"not null" json:"startDate"`
	EndDate   time.Time          `gorm:"not null" json:"endDate"`
	Status    SubscriptionStatus `gorm:"type:varchar(20);not null" json:"status"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
	DeletedAt gorm.DeletedAt     `gorm:"index" json:"-"`
}

func (Subscription) TableName() string {
	return "subscriptions"
}
