package handler

import (
	"encoding/json"
	"net/http"

	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/model"
	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/services"
	"github.com/adityanuriskandar17/HRIS-BE/internal/http/dto"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type SubscriptionHandler struct {
	subscriptionService services.SubscriptionService
}

func NewSubscriptionHandler(subscriptionService services.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{subscriptionService: subscriptionService}
}

// GetAll godoc
// @Summary Get all subscriptions
// @Description Get all subscriptions
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.SubscriptionResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [get]
func (h *SubscriptionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	subscriptions, err := h.subscriptionService.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response DTOs
	responses := make([]dto.SubscriptionResponse, len(subscriptions))
	for i, sub := range subscriptions {
		responses[i] = dto.SubscriptionResponse{
			ID:        sub.ID.String(),
			TenantID:  sub.TenantID.String(),
			PlanID:    sub.PlanID.String(),
			StartDate: sub.StartDate,
			EndDate:   sub.EndDate,
			Status:    string(sub.Status),
			CreatedAt: sub.CreatedAt,
			UpdatedAt: sub.UpdatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

// GetByID godoc
// @Summary Get subscription by ID
// @Description Get subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Subscription ID"
// @Success 200 {object} dto.SubscriptionResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	subscriptionID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "invalid subscription ID", http.StatusBadRequest)
		return
	}

	subscription, err := h.subscriptionService.GetByID(r.Context(), subscriptionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.SubscriptionResponse{
		ID:        subscription.ID.String(),
		TenantID:  subscription.TenantID.String(),
		PlanID:    subscription.PlanID.String(),
		StartDate: subscription.StartDate,
		EndDate:   subscription.EndDate,
		Status:    string(subscription.Status),
		CreatedAt: subscription.CreatedAt,
		UpdatedAt: subscription.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetByTenantID godoc
// @Summary Get subscriptions by tenant ID
// @Description Get subscriptions by tenant ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param tenantId path string true "Tenant ID"
// @Success 200 {array} dto.SubscriptionResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/tenant/{tenantId} [get]
func (h *SubscriptionHandler) GetByTenantID(w http.ResponseWriter, r *http.Request) {
	tenantID := chi.URLParam(r, "tenantId")
	id, err := uuid.Parse(tenantID)
	if err != nil {
		http.Error(w, "invalid tenant ID", http.StatusBadRequest)
		return
	}

	subscriptions, err := h.subscriptionService.GetByTenantID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response DTOs
	responses := make([]dto.SubscriptionResponse, len(subscriptions))
	for i, sub := range subscriptions {
		responses[i] = dto.SubscriptionResponse{
			ID:        sub.ID.String(),
			TenantID:  sub.TenantID.String(),
			PlanID:    sub.PlanID.String(),
			StartDate: sub.StartDate,
			EndDate:   sub.EndDate,
			Status:    string(sub.Status),
			CreatedAt: sub.CreatedAt,
			UpdatedAt: sub.UpdatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

// Create godoc
// @Summary Create a new subscription
// @Description Create a new subscription (Admin only)
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param subscription body dto.SubscriptionRequest true "Subscription data"
// @Success 201 {object} dto.SubscriptionResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [post]
func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.SubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tenantID, err := uuid.Parse(req.TenantID)
	if err != nil {
		http.Error(w, "Invalid tenant ID", http.StatusBadRequest)
		return
	}

	planID, err := uuid.Parse(req.PlanID)
	if err != nil {
		http.Error(w, "Invalid plan ID", http.StatusBadRequest)
		return
	}

	subscription := model.Subscription{
		TenantID:  tenantID,
		PlanID:    planID,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Status:    model.SubscriptionStatus(req.Status),
	}

	if err := h.subscriptionService.Create(r.Context(), &subscription); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.SubscriptionResponse{
		ID:        subscription.ID.String(),
		TenantID:  subscription.TenantID.String(),
		PlanID:    subscription.PlanID.String(),
		StartDate: subscription.StartDate,
		EndDate:   subscription.EndDate,
		Status:    string(subscription.Status),
		CreatedAt: subscription.CreatedAt,
		UpdatedAt: subscription.UpdatedAt,
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Update godoc
// @Summary Update a subscription
// @Description Update a subscription (Admin only)
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Subscription ID"
// @Param subscription body dto.SubscriptionRequest true "Subscription data"
// @Success 200 {object} dto.SubscriptionResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	subscriptionID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "invalid subscription ID", http.StatusBadRequest)
		return
	}

	var req dto.SubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tenantID, err := uuid.Parse(req.TenantID)
	if err != nil {
		http.Error(w, "Invalid tenant ID", http.StatusBadRequest)
		return
	}

	planID, err := uuid.Parse(req.PlanID)
	if err != nil {
		http.Error(w, "Invalid plan ID", http.StatusBadRequest)
		return
	}

	subscription := model.Subscription{
		TenantID:  tenantID,
		PlanID:    planID,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Status:    model.SubscriptionStatus(req.Status),
	}

	if err := h.subscriptionService.Update(r.Context(), subscriptionID, &subscription); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.SubscriptionResponse{
		ID:        subscription.ID.String(),
		TenantID:  subscription.TenantID.String(),
		PlanID:    subscription.PlanID.String(),
		StartDate: subscription.StartDate,
		EndDate:   subscription.EndDate,
		Status:    string(subscription.Status),
		CreatedAt: subscription.CreatedAt,
		UpdatedAt: subscription.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Cancel godoc
// @Summary Cancel a subscription
// @Description Cancel a subscription (Admin only)
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Subscription ID"
// @Param cancel body dto.CancelSubscriptionRequest true "Cancel reason"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id}/cancel [post]
func (h *SubscriptionHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	subscriptionID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "invalid subscription ID", http.StatusBadRequest)
		return
	}

	var req dto.CancelSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.subscriptionService.Cancel(r.Context(), subscriptionID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "Subscription cancelled successfully",
		"reason":  req.Reason,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Renew godoc
// @Summary Renew a subscription
// @Description Renew a subscription (Admin only)
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Subscription ID"
// @Param renew body dto.RenewSubscriptionRequest true "Renewal details"
// @Success 200 {object} dto.SubscriptionResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id}/renew [post]
func (h *SubscriptionHandler) Renew(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	subscriptionID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "invalid subscription ID", http.StatusBadRequest)
		return
	}

	var req dto.RenewSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Calculate months based on billing cycle
	months := 1 // Default to 1 month
	if req.BillingCycle == "yearly" {
		months = 12
	}
	// Use the Months field if provided
	if req.Months > 0 {
		months = req.Months
	}

	if err := h.subscriptionService.Renew(r.Context(), subscriptionID, months); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	subscription, err := h.subscriptionService.GetByID(r.Context(), subscriptionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.SubscriptionResponse{
		ID:        subscription.ID.String(),
		TenantID:  subscription.TenantID.String(),
		PlanID:    subscription.PlanID.String(),
		StartDate: subscription.StartDate,
		EndDate:   subscription.EndDate,
		Status:    string(subscription.Status),
		CreatedAt: subscription.CreatedAt,
		UpdatedAt: subscription.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
