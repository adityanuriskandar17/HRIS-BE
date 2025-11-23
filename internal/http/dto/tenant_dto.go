package dto

import (
	"time"

	"github.com/google/uuid"
)

// TenantDTO represents a tenant in the system
type TenantDTO struct {
	ID          uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name        string    `json:"name" example:"John Doe"`
	Email       string    `json:"email" example:"john.doe@example.com"`
	CompanyName string    `json:"companyName" example:"Acme Corp"`
	Domain      string    `json:"domain" example:"acme.com"`
	Active      bool      `json:"active" example:"true"`
	CreatedAt   time.Time `json:"createdAt" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updatedAt" example:"2023-01-01T00:00:00Z"`
}

// CreateTenantRequest represents a request to create a tenant
type CreateTenantRequest struct {
	Name        string `json:"name" binding:"required" example:"John Doe"`
	Email       string `json:"email" binding:"required,email" example:"john.doe@example.com"`
	CompanyName string `json:"companyName" binding:"required" example:"Acme Corp"`
	Domain      string `json:"domain" binding:"required" example:"acme.com"`
}

// UpdateTenantRequest represents a request to update a tenant
type UpdateTenantRequest struct {
	Name        *string `json:"name,omitempty" example:"John Smith"`
	Email       *string `json:"email,omitempty" example:"john.smith@example.com"`
	CompanyName *string `json:"companyName,omitempty" example:"Acme Inc"`
	Domain      *string `json:"domain,omitempty" example:"acme-inc.com"`
	Active      *bool   `json:"active,omitempty" example:"false"`
}
