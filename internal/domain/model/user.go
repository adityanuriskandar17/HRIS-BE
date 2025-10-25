package model

import "time"

type UserAccount struct {
	ID           uint64   `gorm:"primaryKey" json:"id"`
	Email        string   `gorm:"uniqueIndex:user_accounts_email_key;size:255;not null" json:"email"`
	PasswordHash string   `gorm:"size:255;not null" json:"-"`
	Role         UserRole `gorm:"type:varchar(16);not null" json:"role"`
	Status       uint8    `gorm:"default:1" json:"status"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
