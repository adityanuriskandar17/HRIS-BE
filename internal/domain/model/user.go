package model

import "time"

type UserAccount struct {
	ID           uint64   `gorm:"primaryKey" json:"id"`
	Email        string   `gorm:"uniqueIndex;size:255;not null" json:"email"`
	PasswordHash string   `gorm:"size:255;not null" json:"-"`
	Role         UserRole `gorm:"type:enum('EMPLOYEE','MANAGER','HR','FINANCE','ADMIN');not null" json:"role"`
	Status       uint8    `gorm:"default:1" json:"status"`
	CreatedAt    time.Time
}
