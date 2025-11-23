package handler

import (
	"encoding/json"
	"net/http"

	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/model"
	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/services"
	_ "github.com/adityanuriskandar17/HRIS-BE/internal/http/dto"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type TenantHandler struct {
	tenantService services.TenantService
}

func NewTenantHandler(tenantService services.TenantService) *TenantHandler {
	return &TenantHandler{tenantService: tenantService}
}

// GetAll godoc
// @Summary Get all tenants
// @Description Get all tenants
// @Tags tenants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.TenantDTO
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tenants [get]
func (h *TenantHandler) GetAll(c http.ResponseWriter, r *http.Request) {
	tenants, err := h.tenantService.GetAll()

	if err != nil {
		http.Error(c, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Header().Set("Content-Type", "application/json")
	json.NewEncoder(c).Encode(tenants)
}

// GetByID godoc
// @Summary Get tenant by ID
// @Description Get tenant by ID
// @Tags tenants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Tenant ID"
// @Success 200 {object} dto.TenantDTO
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tenants/{id} [get]
func (h *TenantHandler) GetByID(c http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenantID, err := uuid.Parse(id)
	if err != nil {
		http.Error(c, "invalid tenant ID", http.StatusBadRequest)
		return
	}

	tenant, err := h.tenantService.GetByID(tenantID)
	if err != nil {
		http.Error(c, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Header().Set("Content-Type", "application/json")
	json.NewEncoder(c).Encode(tenant)
}

// Create godoc
// @Summary Create a new tenant
// @Description Create a new tenant (Admin only)
// @Tags tenants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param tenant body dto.CreateTenantRequest true "Tenant data"
// @Success 201 {object} dto.TenantDTO
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tenants [post]
func (h *TenantHandler) Create(c http.ResponseWriter, r *http.Request) {
	var tenant model.Tenant
	if err := json.NewDecoder(r.Body).Decode(&tenant); err != nil {
		http.Error(c, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.tenantService.Create(&tenant); err != nil {
		http.Error(c, err.Error(), http.StatusInternalServerError)
		return
	}
	c.WriteHeader(http.StatusCreated)
	c.Header().Set("Content-Type", "application/json")
	json.NewEncoder(c).Encode(tenant)
}

// Update godoc
// @Summary Update a tenant
// @Description Update a tenant (Admin only)
// @Tags tenants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Tenant ID"
// @Param tenant body dto.UpdateTenantRequest true "Tenant data"
// @Success 200 {object} dto.TenantDTO
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tenants/{id} [put]
func (h *TenantHandler) Update(c http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenantID, err := uuid.Parse(id)
	if err != nil {
		http.Error(c, "invalid tenant ID", http.StatusBadRequest)
		return
	}

	var tenant model.Tenant
	if err := json.NewDecoder(r.Body).Decode(&tenant); err != nil {
		http.Error(c, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.tenantService.Update(tenantID, &tenant); err != nil {
		http.Error(c, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Header().Set("Content-Type", "application/json")
	json.NewEncoder(c).Encode(tenant)
}
