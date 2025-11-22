package model

import "time"

type RefreshToken struct {
	ID         uint64       `gorm:"primaryKey" json:"id"`
	UserID     uint64       `gorm:"not null;index" json:"userId"`
	User       UserAccount  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	TokenHash  string       `gorm:"size:64;uniqueIndex;not null" json:"-"`
	ExpiresAt  time.Time    `gorm:"not null;index" json:"expiresAt"`
	RevokedAt  *time.Time   `json:"revokedAt,omitempty"`
	CreatedAt  time.Time    `json:"createdAt"`
	UpdatedAt  time.Time    `json:"updatedAt"`
}
