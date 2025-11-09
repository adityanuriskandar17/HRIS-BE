package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PayrollLineType string

const (
	PayrollLineTypeAllowance PayrollLineType = "allowance"
	PayrollLineTypeDeduction PayrollLineType = "deduction"
	PayrollLineTypeOvertime  PayrollLineType = "overtime"
	PayrollLineTypeBonus     PayrollLineType = "bonus"
)

type PayrollLine struct {
	ID          uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	PayrollID   uuid.UUID       `gorm:"type:uuid;not null;index" json:"payrollId"`
	Payroll     Payroll         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Description string          `gorm:"size:255;not null" json:"description"`
	Amount      float64         `gorm:"type:decimal(10,2);not null" json:"amount"`
	Type        PayrollLineType `gorm:"type:varchar(20);not null" json:"type"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (PayrollLine) TableName() string {
	return "payroll_lines"
}