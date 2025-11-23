package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ContractType string

const (
	ContractTypePermanent  ContractType = "permanent"
	ContractTypeContract   ContractType = "contract"
	ContractTypeInternship ContractType = "internship"
	ContractTypeProbation  ContractType = "probation"
)

type EmployeeContract struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	EmployeeID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"employeeId"`
	Employee     Employee       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	ContractType ContractType   `gorm:"type:varchar(20);not null" json:"contractType"`
	StartDate    time.Time      `gorm:"not null" json:"startDate"`
	EndDate      *time.Time     `json:"endDate"`
	Salary       float64        `gorm:"type:decimal(10,2);not null" json:"salary"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (EmployeeContract) TableName() string {
	return "employee_contracts"
}
