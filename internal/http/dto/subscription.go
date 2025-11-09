package dto

import "time"

// SubscriptionRequest represents the request body for creating/updating a subscription
type SubscriptionRequest struct {
	TenantID      string    `json:"tenantId"`
	PlanID        string    `json:"planId"`
	BillingCycle  string    `json:"billingCycle"` // monthly, yearly
	PaymentMethod string    `json:"paymentMethod"`
	StartDate     time.Time `json:"startDate"`
	EndDate       time.Time `json:"endDate"`
	Status        string    `json:"status"`
}

// SubscriptionResponse represents the response body for subscription
type SubscriptionResponse struct {
	ID            string    `json:"id"`
	TenantID      string    `json:"tenantId"`
	PlanID        string    `json:"planId"`
	PlanName      string    `json:"planName"`
	BillingCycle  string    `json:"billingCycle"`
	Status        string    `json:"status"`
	StartDate     time.Time `json:"startDate"`
	EndDate       time.Time `json:"endDate"`
	PaymentMethod string    `json:"paymentMethod"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// CancelSubscriptionRequest represents the request body for canceling a subscription
type CancelSubscriptionRequest struct {
	Reason string `json:"reason"`
}

// RenewSubscriptionRequest represents the request body for renewing a subscription
type RenewSubscriptionRequest struct {
	BillingCycle  string `json:"billingCycle"`
	PaymentMethod string `json:"paymentMethod"`
	Months        int    `json:"months"`
}