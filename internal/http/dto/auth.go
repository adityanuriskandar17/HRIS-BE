package dto

// LoginRequest represents the request body for login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents the response body for login
type LoginResponse struct {
	Token string      `json:"token"`
	User  UserResponse `json:"user"`
}

// RegisterRequest represents the request body for registration
type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// UserResponse represents the user information in responses
type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	IsActive  bool   `json:"isActive"`
}