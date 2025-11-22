// Package dto contains data transfer objects for HTTP requests and responses.
package dto

// LoginRequest represents the request body for login
type LoginRequest struct {
	Email    string `json:"email" example:"john.doe@example.com"`
	Password string `json:"password" example:"SecurePass123!"`
}

// LoginResponse represents the response body for login
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// RegisterRequest represents the request body for registration
type RegisterRequest struct {
	Email     string              `json:"email" example:"john.doe@example.com"`
	Password  string              `json:"password" example:"SecurePass123!"`
	FirstName string              `json:"firstName" example:"John"`
	LastName  string              `json:"lastName" example:"Doe"`
	TenantID  string              `json:"tenantId,omitempty" example:""`
	Tenant    *TenantRegistration `json:"tenant,omitempty"`
}

// TenantRegistration captures tenant details when an admin signs up and
// provisions a new tenant.
type TenantRegistration struct {
	Name        string `json:"name" example:"John Doe"`
	Email       string `json:"email" example:"admin@acmecorp.com"`
	CompanyName string `json:"companyName" example:"Acme Corporation"`
	Domain      string `json:"domain" example:"acmecorp.com"`
}

// UserResponse represents the user information in responses
type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	IsActive  bool   `json:"isActive"`
}

// TokenResponse represents the token response for login and refresh
type TokenResponse struct {
	AccessToken      string `json:"accessToken"`
	RefreshToken     string `json:"refreshToken"`
	ExpiresIn        int64  `json:"expiresIn"`
	RefreshExpiresIn int64  `json:"refreshExpiresIn"`
	Role             string `json:"role"`
}

// CreateEmployeeAccountRequest represents admin-triggered user creation payload.
type CreateEmployeeAccountRequest struct {
	Email     string `json:"email" example:"employee@acmecorp.com"`
	Password  string `json:"password" example:"EmployeePass123!"`
	FirstName string `json:"firstName" example:"Jane"`
	LastName  string `json:"lastName" example:"Smith"`
}
