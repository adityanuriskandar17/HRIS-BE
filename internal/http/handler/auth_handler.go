package handler

import (
	"encoding/binary"
	"encoding/json"
	"net/http"

	"github.com/adityanuriskandar17/HRIS-BE/internal/auth"
	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/model"
	"github.com/adityanuriskandar17/HRIS-BE/internal/http/dto"
	"github.com/adityanuriskandar17/HRIS-BE/internal/repository"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo repository.UserAccountRepository
	secret   string
}

func NewAuthHandler(userRepo repository.UserAccountRepository, secret string) *AuthHandler {
	return &AuthHandler{
		userRepo: userRepo,
		secret:   secret,
	}
}

// Login handles user login
// @Summary Login user
// @Description Authenticate a user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} string
// @Failure 401 {object} string
// @Failure 500 {object} string
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Find user by email
	u, err := h.userRepo.FindByEmail(r.Context(), loginReq.Email)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(loginReq.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Convert UUID to uint64 for JWT signing
	uid := binary.BigEndian.Uint64(u.ID[:8])

	// Generate JWT token
	token, err := auth.SignJWT(uid, "user", h.secret, 24*60*60) // 24 hours
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare response
	loginRes := dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        u.ID.String(),
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginRes)
}

// Register handles user registration
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User registration data"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var registerReq dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&registerReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create user
	user := model.UserAccount{
		Email:        registerReq.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    registerReq.FirstName,
		LastName:     registerReq.LastName,
	}



	createdUser, err := h.userRepo.Create(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert UUID to uint64 for JWT signing
	uid := binary.BigEndian.Uint64(createdUser.ID[:8])

	// Generate JWT token
	token, err := auth.SignJWT(uid, "user", h.secret, 24*60*60) // 24 hours
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare response
	loginRes := dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        createdUser.ID.String(),
			Email:     createdUser.Email,
			FirstName: createdUser.FirstName,
			LastName:  createdUser.LastName,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginRes)
}

// Profile handles user profile retrieval
// @Summary Get user profile
// @Description Retrieve the authenticated user's profile
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} string
// @Failure 500 {object} string
// @Router /auth/profile [get]
func (h *AuthHandler) Profile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	// Parse UUID
	uuid, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Find user by ID
	user, err := h.userRepo.FindByID(r.Context(), uuid)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Prepare response
	userRes := dto.UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userRes)
}
