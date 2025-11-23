package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PayrollStatus string

const (
	PayrollStatusDraft     PayrollStatus = "draft"
	PayrollStatusProcessed PayrollStatus = "processed"
	PayrollStatusPaid      PayrollStatus = "paid"
)

type Payroll struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	EmployeeID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"employeeId"`
	Employee    Employee       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	PayPeriod   string         `gorm:"size:20;not null" json:"payPeriod"`
	GrossSalary float64        `gorm:"type:decimal(10,2);not null" json:"grossSalary"`
	Deductions  float64        `gorm:"type:decimal(10,2);default:0" json:"deductions"`
	NetSalary   float64        `gorm:"type:decimal(10,2);not null" json:"netSalary"`
	Status      PayrollStatus  `gorm:"type:varchar(20);not null" json:"status"`
	PaymentDate *time.Time     `json:"paymentDate"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Payroll) TableName() string {
	return "payrolls"
}
