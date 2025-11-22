package dto

// LoginRequest represents the request body for login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents the response body for login
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// RegisterRequest represents the request body for registration
type RegisterRequest struct {
	Email     string              `json:"email"`
	Password  string              `json:"password"`
	FirstName string              `json:"firstName"`
	LastName  string              `json:"lastName"`
	TenantID  string              `json:"tenantId,omitempty"`
	Tenant    *TenantRegistration `json:"tenant,omitempty"`
}

// TenantRegistration captures tenant details when an admin signs up and
// provisions a new tenant.
type TenantRegistration struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	CompanyName string `json:"companyName"`
	Domain      string `json:"domain"`
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
