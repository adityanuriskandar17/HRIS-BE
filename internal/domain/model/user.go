package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserAccount struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	TenantID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenantId"`
	Tenant       Tenant         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Email        string         `gorm:"uniqueIndex:user_accounts_email_key;size:255;not null" json:"email"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	FirstName    string         `gorm:"size:100;not null" json:"firstName"`
	LastName     string         `gorm:"size:100;not null" json:"lastName"`
	IsActive     bool           `gorm:"default:true" json:"isActive"`
	LastLoginAt  *time.Time     `json:"lastLoginAt"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (UserAccount) TableName() string {
	return "user_accounts"
}
