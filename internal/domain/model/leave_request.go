package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LeaveRequestStatus string

const (
	LeaveRequestStatusPending   LeaveRequestStatus = "pending"
	LeaveRequestStatusApproved  LeaveRequestStatus = "approved"
	LeaveRequestStatusRejected  LeaveRequestStatus = "rejected"
	LeaveRequestStatusCancelled LeaveRequestStatus = "cancelled"
)

type LeaveRequest struct {
	ID            uuid.UUID          `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	EmployeeID    uuid.UUID          `gorm:"type:uuid;not null;index" json:"employeeId"`
	Employee      Employee           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	LeaveType     LeaveType          `gorm:"type:varchar(20);not null" json:"leaveType"`
	StartDate     time.Time          `gorm:"not null" json:"startDate"`
	EndDate       time.Time          `gorm:"not null" json:"endDate"`
	Reason        string             `gorm:"type:text" json:"reason"`
	Status        LeaveRequestStatus `gorm:"type:varchar(20);not null" json:"status"`
	ApprovedBy    *uuid.UUID         `gorm:"type:uuid" json:"approvedBy"`
	ApprovedAt    *time.Time         `json:"approvedAt"`
	ApproverNotes string             `gorm:"type:text" json:"approverNotes"`
	CreatedAt     time.Time          `json:"createdAt"`
	UpdatedAt     time.Time          `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt     `gorm:"index" json:"-"`
}

func (LeaveRequest) TableName() string {
	return "leave_requests"
}
