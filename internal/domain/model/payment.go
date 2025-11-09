package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentMethod string

const (
	PaymentMethodCreditCard PaymentMethod = "credit_card"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodPayPal PaymentMethod = "paypal"
	PaymentMethodStripe PaymentMethod = "stripe"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

type Payment struct {
	ID            uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	InvoiceID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"invoiceId"`
	Invoice       Invoice        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	Amount        float64        `gorm:"type:decimal(10,2);not null" json:"amount"`
	PaymentDate   time.Time      `gorm:"not null" json:"paymentDate"`
	PaymentMethod PaymentMethod  `gorm:"type:varchar(20);not null" json:"paymentMethod"`
	Status        PaymentStatus  `gorm:"type:varchar(20);not null" json:"status"`
	TransactionID string         `gorm:"size:100" json:"transactionId"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Payment) TableName() string {
	return "payments"
}