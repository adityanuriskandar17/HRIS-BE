package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Employee struct {
	ID               uuid.UUID        `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	EmployeeCode     string           `gorm:"uniqueIndex:employees_employee_code_key;size:32;not null" json:"employeeCode"`
	FullName         string           `gorm:"size:150;not null" json:"fullName"`
	Email            string           `gorm:"uniqueIndex:employees_email_key;size:255;not null" json:"email"`
	Phone            string           `gorm:"size:30" json:"phone"`
	UnitID           uuid.UUID        `json:"unitId"`
	Unit             Unit             `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"unit"`
	PositionID       uuid.UUID        `json:"positionId"`
	Position         Position         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"position"`
	EmploymentStatus EmploymentStatus `gorm:"type:employment_status_enum;not null;default:'FULLTIME'" json:"employmentStatus"`
	StartDate        time.Time        `gorm:"type:date" json:"startDate"`
	EndDate          *time.Time       `gorm:"type:date" json:"endDate,omitempty"`
	DateOfBirth      *time.Time       `gorm:"type:date" json:"dateOfBirth,omitempty"`
	CreatedAt        time.Time        `json:"createdAt"`
	UpdatedAt        time.Time        `json:"updatedAt"`
	CreatedByID      *uuid.UUID       `gorm:"column:created_by;type:uuid" json:"createdById,omitempty"`
	CreatedBy        *UserAccount     `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	UpdatedByID      *uuid.UUID       `gorm:"column:updated_by;type:uuid" json:"updatedById,omitempty"`
	UpdatedBy        *UserAccount     `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	DeletedAt        gorm.DeletedAt   `gorm:"index" json:"-"`
}
