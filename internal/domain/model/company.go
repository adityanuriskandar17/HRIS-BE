package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Company struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	TenantID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenantId"`
	Tenant         Tenant         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Name           string         `gorm:"size:100;not null" json:"name"`
	RegistrationNo string         `gorm:"size:50" json:"registrationNo"`
	Address        string         `gorm:"type:text" json:"address"`
	Timezone       string         `gorm:"size:50;default:'UTC'" json:"timezone"`
	Currency       string         `gorm:"size:3;default:'USD'" json:"currency"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Company) TableName() string {
	return "companies"
}