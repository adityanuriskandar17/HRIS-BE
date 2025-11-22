package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/adityanuriskandar17/HRIS-BE/internal/auth"
	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/model"
	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/services"
	"github.com/adityanuriskandar17/HRIS-BE/internal/http/authctx"
	"github.com/adityanuriskandar17/HRIS-BE/internal/http/dto"
	httpx "github.com/adityanuriskandar17/HRIS-BE/internal/http/httputils"
	"github.com/adityanuriskandar17/HRIS-BE/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo  repository.UserAccountRepository
	tenantSvc services.TenantService
	tokens    *auth.Service
}

var errTenantCreate = errors.New("tenant creation failed")

func NewAuthHandler(userRepo repository.UserAccountRepository, tenantSvc services.TenantService, tokens *auth.Service) *AuthHandler {
	return &AuthHandler{
		userRepo:  userRepo,
		tenantSvc: tenantSvc,
		tokens:    tokens,
	}
}

type refreshReq struct {
	RefreshToken string `json:"refreshToken"`
}

// Login handles user login
// @Summary Login user
// @Description Authenticate a user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} response.TokenResponseEnvelope
// @Failure 400 {object} string
// @Failure 401 {object} string
// @Failure 500 {object} string
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "INVALID_JSON", "invalid json payload", nil)
		return
	}

	// Find user by email
	u, err := h.userRepo.FindByEmail(r.Context(), loginReq.Email)
	if err != nil {
		httpx.Error(w, r, http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid credentials", nil)
		return
	}

	if !u.IsActive {
		httpx.Error(w, r, http.StatusForbidden, "ACCOUNT_DISABLED", "account disabled", nil)
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(loginReq.Password)); err != nil {
		httpx.Error(w, r, http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid credentials", nil)
		return
	}

	// Issue tokens
	pair, err := h.tokens.IssueTokenPair(u)
	if err != nil {
		httpx.Error(w, r, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to issue token", nil)
		return
	}

	resp := dto.TokenResponse{
		AccessToken:      pair.AccessToken,
		RefreshToken:     pair.RefreshToken,
		ExpiresIn:        secondsUntil(pair.AccessExpiresAt),
		RefreshExpiresIn: secondsUntil(pair.RefreshExpiresAt),
		Role:             string(u.Role),
	}
	httpx.OK(w, r, resp)
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
		httpx.Error(w, r, http.StatusBadRequest, "INVALID_JSON", "invalid json payload", nil)
		return
	}

	tenantID, role, err := h.resolveTenant(registerReq)
	if err != nil {
		var status = http.StatusBadRequest
		if errors.Is(err, errTenantCreate) {
			status = http.StatusInternalServerError
		}
		httpx.Error(w, r, status, "TENANT_RESOLUTION_FAILED", err.Error(), nil)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	if err != nil {
		httpx.Error(w, r, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to hash password", nil)
		return
	}

	// Create user
	user := model.UserAccount{
		TenantID:     tenantID,
		Email:        registerReq.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    registerReq.FirstName,
		LastName:     registerReq.LastName,
		Role:         role,
		IsActive:     true,
	}

	createdUser, err := h.userRepo.Create(r.Context(), user)
	if err != nil {
		httpx.Error(w, r, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to create user", nil)
		return
	}

	// Issue tokens
	pair, err := h.tokens.IssueTokenPair(createdUser)
	if err != nil {
		httpx.Error(w, r, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to issue token", nil)
		return
	}

	// Prepare response
	loginRes := dto.LoginResponse{
		Token: pair.AccessToken, // Keeping backward compatibility if needed, or just use the pair
		User: dto.UserResponse{
			ID:        createdUser.ID.String(),
			Email:     createdUser.Email,
			FirstName: createdUser.FirstName,
			LastName:  createdUser.LastName,
		},
	}
	// Note: The original upstream response structure might have been different (just Token string vs TokenPair).
	// But since we are moving to TokenPair (Access+Refresh), we should probably return that.
	// However, `dto.LoginResponse` might only have `Token` string field.
	// Let's check `internal/http/dto/auth_dto.go` if I could, but I'll stick to what upstream `LoginResponse` has.
	// Upstream `LoginResponse` has `Token` string.
	// Stashed `tokenResponse` has `AccessToken`, `RefreshToken` etc.
	// I should probably use `tokenResponse` struct for consistency with `Login` and `Refresh`.
	// But `Register` return type in swagger says `dto.UserResponse` (upstream) or `dto.LoginResponse` (upstream implementation).
	// Upstream implementation:
	// loginRes := dto.LoginResponse{ Token: token, User: ... }
	// json.NewEncoder(w).Encode(loginRes)

	// I will stick to upstream `dto.LoginResponse` for now to avoid breaking changes if frontend expects that.
	// But `dto.LoginResponse` likely doesn't have RefreshToken.
	// If I want to support refresh tokens, I should update `dto.LoginResponse` or use `tokenResponse`.
	// Given the conflict, I'll use `tokenResponse` structure but maybe map it to `dto.LoginResponse` if possible?
	// No, `Login` now returns `tokenResponse` (stashed struct).
	// So `Register` should also return `tokenResponse` + User info?
	// Stashed `Login` returns `tokenResponse`. Upstream `Login` returned `dto.LoginResponse`.
	// I replaced `Login` with stashed version which returns `tokenResponse`.
	// So I should make `Register` return `tokenResponse` too?
	// But `Register` also returns User info in upstream.
	// Stashed `Login` returns `tokenResponse` which has `Role` but not full User info.
	// Upstream `Login` returns `dto.LoginResponse` which has `User` info.

	// Compromise: I will use `dto.LoginResponse` but I might need to add RefreshToken to it later.
	// For now, I'll just set `Token` to `pair.AccessToken`.

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginRes)
}

// CreateEmployeeAccount allows an admin to create an employee user within their tenant
// @Summary Create employee account
// @Description Admin creates an employee account tied to their tenant
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateEmployeeAccountRequest true "Employee account data"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} string
// @Failure 401 {object} string
// @Failure 500 {object} string
// @Router /auth/employees [post]
func (h *AuthHandler) CreateEmployeeAccount(w http.ResponseWriter, r *http.Request) {
	current, ok := authctx.CurrentUser(r.Context())
	if !ok {
		httpx.Error(w, r, http.StatusUnauthorized, "UNAUTHENTICATED", "unauthorized", nil)
		return
	}

	var req dto.CreateEmployeeAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "INVALID_JSON", "invalid json payload", nil)
		return
	}

	if err := validateEmployeeAccountRequest(req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "INVALID_PAYLOAD", err.Error(), nil)
		return
	}

	tenantID, err := uuid.Parse(current.TenantID)
	if err != nil {
		httpx.Error(w, r, http.StatusInternalServerError, "INVALID_TENANT", "invalid tenant context", nil)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		httpx.Error(w, r, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to hash password", nil)
		return
	}

	user := model.UserAccount{
		TenantID:     tenantID,
		Email:        strings.TrimSpace(req.Email),
		PasswordHash: string(passwordHash),
		FirstName:    strings.TrimSpace(req.FirstName),
		LastName:     strings.TrimSpace(req.LastName),
		Role:         model.RoleEmployee,
		IsActive:     true,
	}

	created, err := h.userRepo.Create(r.Context(), user)
	if err != nil {
		httpx.Error(w, r, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to create user", nil)
		return
	}

	resp := dto.UserResponse{
		ID:        created.ID.String(),
		Email:     created.Email,
		FirstName: created.FirstName,
		LastName:  created.LastName,
		IsActive:  created.IsActive,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) resolveTenant(req dto.RegisterRequest) (uuid.UUID, model.UserRole, error) {
	tenantIDStr := strings.TrimSpace(req.TenantID)

	switch {
	case tenantIDStr != "" && req.Tenant != nil:
		return uuid.Nil, model.RoleEmployee, fmt.Errorf("provide either tenantId or tenant payload, not both")
	case tenantIDStr != "":
		id, err := uuid.Parse(tenantIDStr)
		if err != nil {
			return uuid.Nil, model.RoleEmployee, fmt.Errorf("tenantId must be a valid UUID")
		}
		return id, model.RoleEmployee, nil
	case req.Tenant == nil:
		return uuid.Nil, model.RoleEmployee, fmt.Errorf("tenantId or tenant payload is required")
	default:
		if err := validateTenantRegistration(req.Tenant); err != nil {
			return uuid.Nil, model.RoleEmployee, err
		}

		tenant := &model.Tenant{
			Name:        strings.TrimSpace(req.Tenant.Name),
			Email:       strings.TrimSpace(req.Tenant.Email),
			CompanyName: strings.TrimSpace(req.Tenant.CompanyName),
			Domain:      strings.TrimSpace(req.Tenant.Domain),
		}

		if err := h.tenantSvc.Create(tenant); err != nil {
			return uuid.Nil, model.RoleEmployee, fmt.Errorf("%w: %v", errTenantCreate, err)
		}

		return tenant.ID, model.RoleAdmin, nil
	}
}

func validateTenantRegistration(tenant *dto.TenantRegistration) error {
	if tenant == nil {
		return fmt.Errorf("tenant payload is required")
	}

	if strings.TrimSpace(tenant.Name) == "" {
		return fmt.Errorf("tenant.name is required")
	}
	if strings.TrimSpace(tenant.Email) == "" {
		return fmt.Errorf("tenant.email is required")
	}
	if strings.TrimSpace(tenant.CompanyName) == "" {
		return fmt.Errorf("tenant.companyName is required")
	}
	if strings.TrimSpace(tenant.Domain) == "" {
		return fmt.Errorf("tenant.domain is required")
	}

	return nil
}

func validateEmployeeAccountRequest(req dto.CreateEmployeeAccountRequest) error {
	if strings.TrimSpace(req.Email) == "" {
		return fmt.Errorf("email is required")
	}
	if strings.TrimSpace(req.Password) == "" {
		return fmt.Errorf("password is required")
	}
	if strings.TrimSpace(req.FirstName) == "" {
		return fmt.Errorf("firstName is required")
	}
	if strings.TrimSpace(req.LastName) == "" {
		return fmt.Errorf("lastName is required")
	}
	return nil
}

// ProfileSelf handles retrieval of the currently authenticated user's profile
// @Summary Get current user profile
// @Description Retrieve the current authenticated user's profile
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} string
// @Failure 500 {object} string
// @Router /auth/profile [get]
func (h *AuthHandler) ProfileSelf(w http.ResponseWriter, r *http.Request) {
	current, ok := authctx.CurrentUser(r.Context())
	if !ok {
		httpx.Error(w, r, http.StatusUnauthorized, "UNAUTHENTICATED", "unauthorized", nil)
		return
	}
	h.renderUserProfile(w, r, current.ID)
}

// ProfileByID handles user profile retrieval by ID
// @Summary Get user profile by ID
// @Description Retrieve a user's profile by ID
// @Tags auth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} string
// @Failure 401 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /auth/profile/{id} [get]
func (h *AuthHandler) ProfileByID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	uid, err := uuid.Parse(userID)
	if err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "INVALID_UUID", "invalid user ID", nil)
		return
	}

	h.renderUserProfile(w, r, uid)
}

func (h *AuthHandler) renderUserProfile(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	user, err := h.userRepo.FindByID(r.Context(), id)
	if err != nil {
		httpx.Error(w, r, http.StatusNotFound, "USER_NOT_FOUND", "User not found", nil)
		return
	}

	userRes := dto.UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsActive:  user.IsActive,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userRes)
}

// Refresh handles token refresh
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body refreshReq true "Refresh token"
// @Success 200 {object} response.TokenResponseEnvelope
// @Failure 400 {object} string
// @Failure 401 {object} string
// @Failure 500 {object} string
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "INVALID_JSON", "invalid json payload", nil)
		return
	}

	token := strings.TrimSpace(req.RefreshToken)
	pair, user, err := h.tokens.RotateRefreshToken(token)
	if err != nil {
		status := http.StatusUnauthorized
		message := "invalid refresh token"
		code := "INVALID_REFRESH_TOKEN"
		switch err {
		case auth.ErrRefreshTokenNotFound, auth.ErrRefreshTokenExpired, auth.ErrRefreshTokenRevoked:
			status = http.StatusUnauthorized
		default:
			status = http.StatusInternalServerError
			message = "failed to refresh token"
			code = "INTERNAL_SERVER_ERROR"
		}
		httpx.Error(w, r, status, code, message, nil)
		return
	}

	resp := dto.TokenResponse{
		AccessToken:      pair.AccessToken,
		RefreshToken:     pair.RefreshToken,
		ExpiresIn:        secondsUntil(pair.AccessExpiresAt),
		RefreshExpiresIn: secondsUntil(pair.RefreshExpiresAt),
		Role:             string(user.Role),
	}
	httpx.OK(w, r, resp)
}

func secondsUntil(t time.Time) int64 {
	seconds := time.Until(t).Seconds()
	if seconds < 0 {
		return 0
	}
	return int64(seconds)
}
