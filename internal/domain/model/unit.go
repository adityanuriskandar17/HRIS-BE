package model

import "time"

type Unit struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Code      string    `gorm:"uniqueIndex:units_code_key;size:32;not null" json:"code"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
