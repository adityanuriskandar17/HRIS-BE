package model

import "time"

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
	EmploymentStatus EmploymentStatus `gorm:"type:varchar(16);not null;default:'FULLTIME'" json:"employmentStatus"`
	StartDate        time.Time        `json:"startDate"`
	EndDate          *time.Time       `json:"endDate,omitempty"`
	DateOfBirth      *time.Time       `json:"dateOfBirth,omitempty"`
	CreatedAt        time.Time        `json:"createdAt"`
	UpdatedAt        time.Time        `json:"updatedAt"`
}
