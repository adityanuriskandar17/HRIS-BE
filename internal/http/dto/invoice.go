package dto

import "time"

// InvoiceResponse represents the response body for invoice
type InvoiceResponse struct {
	ID            string    `json:"id"`
	SubscriptionID string    `json:"subscriptionId"`
	TenantID      string    `json:"tenantId"`
	InvoiceNumber string    `json:"invoiceNumber"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`
	DueDate       time.Time `json:"dueDate"`
	SentAt        *time.Time `json:"sentAt,omitempty"`
	PaidAt        *time.Time `json:"paidAt,omitempty"`
	PaymentID     string    `json:"paymentId,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// SendInvoiceRequest represents the request body for sending an invoice
type SendInvoiceRequest struct {
	Email string `json:"email"`
}

// PayInvoiceRequest represents the request body for paying an invoice
type PayInvoiceRequest struct {
	PaymentMethod string `json:"paymentMethod"`
	TransactionID string `json:"transactionId"`
	PaymentID     string `json:"paymentId"`
}