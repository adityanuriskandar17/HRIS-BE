package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InvoiceStatus string

const (
	InvoiceStatusDraft    InvoiceStatus = "draft"
	InvoiceStatusSent     InvoiceStatus = "sent"
	InvoiceStatusPaid     InvoiceStatus = "paid"
	InvoiceStatusOverdue  InvoiceStatus = "overdue"
	InvoiceStatusCanceled InvoiceStatus = "canceled"
)

type Invoice struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	TenantID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenantId"`
	Tenant         Tenant         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	SubscriptionID uuid.UUID      `gorm:"type:uuid;not null;index" json:"subscriptionId"`
	Subscription   Subscription   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	InvoiceNumber  string         `gorm:"size:50;not null;uniqueIndex" json:"invoiceNumber"`
	Amount         float64        `gorm:"type:decimal(10,2);not null" json:"amount"`
	DueDate        time.Time      `gorm:"not null" json:"dueDate"`
	Status         InvoiceStatus  `gorm:"type:varchar(20);not null" json:"status"`
	SentAt         *time.Time     `json:"sentAt,omitempty"`
	PaidAt         *time.Time     `json:"paidAt,omitempty"`
	PaymentID      string         `gorm:"size:100" json:"paymentId,omitempty"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Invoice) TableName() string {
	return "invoices"
}
