package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttendanceStatus string

const (
	AttendanceStatusPresent    AttendanceStatus = "present"
	AttendanceStatusAbsent     AttendanceStatus = "absent"
	AttendanceStatusLate       AttendanceStatus = "late"
	AttendanceStatusHalfDay    AttendanceStatus = "half_day"
	AttendanceStatusLeave      AttendanceStatus = "leave"
)

type Attendance struct {
	ID        uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	EmployeeID uuid.UUID       `gorm:"type:uuid;not null;index" json:"employeeId"`
	Employee   Employee        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Date      time.Time       `gorm:"not null;index" json:"date"`
	CheckIn   *time.Time      `json:"checkIn"`
	CheckOut  *time.Time      `json:"checkOut"`
	Status    AttendanceStatus `gorm:"type:varchar(20);not null" json:"status"`
	Notes     string          `gorm:"type:text" json:"notes"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (Attendance) TableName() string {
	return "attendances"
}