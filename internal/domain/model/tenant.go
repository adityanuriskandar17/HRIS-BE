package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tenant struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Email       string         `gorm:"size:100;not null" json:"email"`
	CompanyName string         `gorm:"size:100;not null" json:"companyName"`
	Domain      string         `gorm:"size:100;not null" json:"domain"`
	Active      bool           `gorm:"default:true" json:"active"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
