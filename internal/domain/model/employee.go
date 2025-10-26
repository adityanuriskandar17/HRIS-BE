package model

import (
	"time"

	"gorm.io/gorm"
)

type Employee struct {
	ID               uint64           `gorm:"primaryKey" json:"id"`
	EmployeeCode     string           `gorm:"uniqueIndex:employees_employee_code_key;size:32;not null" json:"employeeCode"`
	FullName         string           `gorm:"size:150;not null" json:"fullName"`
	Email            string           `gorm:"uniqueIndex:employees_email_key;size:255;not null" json:"email"`
	Phone            string           `gorm:"size:30" json:"phone"`
	UnitID           uint64           `json:"unitId"`
	Unit             Unit             `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"unit"`
	PositionID       uint64           `json:"positionId"`
	Position         Position         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"position"`
	EmploymentStatus EmploymentStatus `gorm:"type:employment_status_enum;not null;default:'FULLTIME'" json:"employmentStatus"`
	StartDate        time.Time        `gorm:"type:date" json:"startDate"`
	EndDate          *time.Time       `gorm:"type:date" json:"endDate,omitempty"`
	DateOfBirth      *time.Time       `gorm:"type:date" json:"dateOfBirth,omitempty"`
	CreatedAt        time.Time        `json:"createdAt"`
	UpdatedAt        time.Time        `json:"updatedAt"`
	CreatedByID      *uint64          `gorm:"column:created_by" json:"createdById,omitempty"`
	CreatedBy        *UserAccount     `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	UpdatedByID      *uint64          `gorm:"column:updated_by" json:"updatedById,omitempty"`
	UpdatedBy        *UserAccount     `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	DeletedAt        gorm.DeletedAt   `gorm:"index" json:"-"`
}
