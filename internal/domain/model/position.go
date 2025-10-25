package model

import "time"

type Position struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:100;not null" json:"title"`
	UnitID    *uint64   `json:"unitId"`
	Unit      *Unit     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"unit,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
