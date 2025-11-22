package model

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID         uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID     uuid.UUID    `gorm:"type:uuid;not null;index" json:"userId"`
	User       UserAccount  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	TokenHash  string       `gorm:"size:64;uniqueIndex;not null" json:"-"`
	ExpiresAt  time.Time    `gorm:"not null;index" json:"expiresAt"`
	RevokedAt  *time.Time   `json:"revokedAt,omitempty"`
	CreatedAt  time.Time    `json:"createdAt"`
	UpdatedAt  time.Time    `json:"updatedAt"`
}
