package dto

import (
	"time"

	"github.com/google/uuid"
)

// CompanyDTO represents a company in the system
type CompanyDTO struct {
	ID             uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	TenantID       uuid.UUID `json:"tenantId" example:"123e4567-e89b-12d3-a456-426614174001"`
	Name           string    `json:"name" example:"Acme Corp"`
	RegistrationNo string    `json:"registrationNo" example:"123456789"`
	Address        string    `json:"address" example:"123 Main St, Anytown, USA"`
	Timezone       string    `json:"timezone" example:"UTC"`
	Currency       string    `json:"currency" example:"USD"`
	CreatedAt      time.Time `json:"createdAt" example:"2023-01-01T00:00:00Z"`
	UpdatedAt      time.Time `json:"updatedAt" example:"2023-01-01T00:00:00Z"`
}

// CreateCompanyRequest represents a request to create a company
type CreateCompanyRequest struct {
	TenantID       string `json:"tenantId" example:"123e4567-e89b-12d3-a456-426614174001"`
	Name           string `json:"name" example:"Acme Corporation"`
	RegistrationNo string `json:"registrationNo,omitempty" example:"123456789"`
	Address        string `json:"address,omitempty" example:"123 Main St, Anytown, USA"`
	Timezone       string `json:"timezone,omitempty" example:"UTC"`
	Currency       string `json:"currency,omitempty" example:"USD"`
}

// UpdateCompanyRequest represents a request to update a company
type UpdateCompanyRequest struct {
	Name           *string `json:"name,omitempty" example:"Acme Inc"`
	RegistrationNo *string `json:"registrationNo,omitempty" example:"987654321"`
	Address        *string `json:"address,omitempty" example:"456 Oak Ave, Anytown, USA"`
	Timezone       *string `json:"timezone,omitempty" example:"America/New_York"`
	Currency       *string `json:"currency,omitempty" example:"EUR"`
}

// CompanySettingsDTO represents company settings
type CompanySettingsDTO struct {
	ID             uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	TenantID       uuid.UUID `json:"tenantId" example:"123e4567-e89b-12d3-a456-426614174001"`
	Name           string    `json:"name" example:"Acme Corp"`
	RegistrationNo string    `json:"registrationNo" example:"123456789"`
	Address        string    `json:"address" example:"123 Main St, Anytown, USA"`
	Timezone       string    `json:"timezone" example:"UTC"`
	Currency       string    `json:"currency" example:"USD"`
}

// UpdateCompanySettingsRequest represents a request to update company settings
type UpdateCompanySettingsRequest struct {
	Name           *string `json:"name,omitempty" example:"Acme Inc"`
	RegistrationNo *string `json:"registrationNo,omitempty" example:"987654321"`
	Address        *string `json:"address,omitempty" example:"456 Oak Ave, Anytown, USA"`
	Timezone       *string `json:"timezone,omitempty" example:"America/New_York"`
	Currency       *string `json:"currency,omitempty" example:"EUR"`
}

// CompanyLimitsDTO represents company limits based on subscription plan
type CompanyLimitsDTO struct {
	ID             uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	TenantID       uuid.UUID `json:"tenantId" example:"123e4567-e89b-12d3-a456-426614174001"`
	Name           string    `json:"name" example:"Acme Corp"`
	MaxEmployees   int       `json:"maxEmployees" example:"100"`
	MaxDepartments int       `json:"maxDepartments" example:"10"`
	MaxPositions   int       `json:"maxPositions" example:"50"`
}
