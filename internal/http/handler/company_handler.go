package handler

import (
	"net/http"
	"strings"

	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/model"
	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/services"
	"github.com/adityanuriskandar17/HRIS-BE/internal/http/dto"
	res "github.com/adityanuriskandar17/HRIS-BE/internal/http/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CompanyHandler struct {
	companyService services.CompanyService
}

func NewCompanyHandler(companyService services.CompanyService) *CompanyHandler {
	return &CompanyHandler{companyService: companyService}
}

// GetAllCompanies godoc
// @Summary Get all companies
// @Description Get a list of all companies
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]dto.CompanyDTO}
// @Failure 500 {object} response.Response
// @Router /companies [get]
func (h *CompanyHandler) GetAllCompanies(w http.ResponseWriter, r *http.Request) {
	companies, err := h.companyService.GetAll()
	if err != nil {
		res.Error(w, http.StatusInternalServerError, "Failed to retrieve companies")
		return
	}

	companyDTOs := make([]*dto.CompanyDTO, len(companies))
	for i, company := range companies {
		companyDTOs[i] = &dto.CompanyDTO{
			ID:             company.ID,
			TenantID:       company.TenantID,
			Name:           company.Name,
			RegistrationNo: company.RegistrationNo,
			Address:        company.Address,
			Timezone:       company.Timezone,
			Currency:       company.Currency,
			CreatedAt:      company.CreatedAt,
			UpdatedAt:      company.UpdatedAt,
		}
	}

	res.Success(w, http.StatusOK, "Companies retrieved successfully", companyDTOs)
}

// CreateCompany godoc
// @Summary Create a new company
// @Description Create a new company for a tenant
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateCompanyRequest true "Company data"
// @Success 201 {object} response.Response{data=dto.CompanyDTO}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /companies [post]
func (h *CompanyHandler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCompanyRequest
	if err := res.ReadJSON(r, &req); err != nil {
		res.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if strings.TrimSpace(req.Name) == "" {
		res.Error(w, http.StatusBadRequest, "Name is required")
		return
	}

	tenantID, err := uuid.Parse(req.TenantID)
	if err != nil {
		res.Error(w, http.StatusBadRequest, "Invalid tenant ID")
		return
	}

	// Set defaults
	timezone := strings.TrimSpace(req.Timezone)
	if timezone == "" {
		timezone = "UTC"
	}

	currency := strings.TrimSpace(req.Currency)
	if currency == "" {
		currency = "USD"
	}

	company := &model.Company{
		TenantID:       tenantID,
		Name:           strings.TrimSpace(req.Name),
		RegistrationNo: strings.TrimSpace(req.RegistrationNo),
		Address:        strings.TrimSpace(req.Address),
		Timezone:       timezone,
		Currency:       currency,
	}

	if err := h.companyService.Create(company); err != nil {
		res.Error(w, http.StatusInternalServerError, "Failed to create company")
		return
	}

	companyDTO := &dto.CompanyDTO{
		ID:             company.ID,
		TenantID:       company.TenantID,
		Name:           company.Name,
		RegistrationNo: company.RegistrationNo,
		Address:        company.Address,
		Timezone:       company.Timezone,
		Currency:       company.Currency,
		CreatedAt:      company.CreatedAt,
		UpdatedAt:      company.UpdatedAt,
	}

	res.Success(w, http.StatusCreated, "Company created successfully", companyDTO)
}

// GetCompanyProfile godoc
// @Summary Get company profile
// @Description Get company profile by ID
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Company ID"
// @Success 200 {object} response.Response{data=dto.CompanyDTO}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /companies/{id} [get]
func (h *CompanyHandler) GetCompanyProfile(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		res.Error(w, http.StatusBadRequest, "Invalid company ID")
		return
	}

	company, err := h.companyService.GetByID(id)
	if err != nil {
		res.Error(w, http.StatusNotFound, "Company not found")
		return
	}

	companyDTO := &dto.CompanyDTO{
		ID:             company.ID,
		TenantID:       company.TenantID,
		Name:           company.Name,
		RegistrationNo: company.RegistrationNo,
		Address:        company.Address,
		Timezone:       company.Timezone,
		Currency:       company.Currency,
		CreatedAt:      company.CreatedAt,
		UpdatedAt:      company.UpdatedAt,
	}

	res.Success(w, http.StatusOK, "Company profile retrieved successfully", companyDTO)
}

// UpdateCompanyProfile godoc
// @Summary Update company profile
// @Description Update company profile by ID
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Company ID"
// @Param request body dto.UpdateCompanyRequest true "Update company request"
// @Success 200 {object} response.Response{data=dto.CompanyDTO}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /companies/{id} [patch]
func (h *CompanyHandler) UpdateCompanyProfile(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		res.Error(w, http.StatusBadRequest, "Invalid company ID")
		return
	}

	var req dto.UpdateCompanyRequest
	if err = res.ReadJSON(r, &req); err != nil {
		res.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Get existing company
	var company *model.Company
	company, err = h.companyService.GetByID(id)
	if err != nil {
		res.Error(w, http.StatusNotFound, "Company not found")
		return
	}

	// Update fields if provided
	if req.Name != nil {
		company.Name = *req.Name
	}
	if req.RegistrationNo != nil {
		company.RegistrationNo = *req.RegistrationNo
	}
	if req.Address != nil {
		company.Address = *req.Address
	}
	if req.Timezone != nil {
		company.Timezone = *req.Timezone
	}
	if req.Currency != nil {
		company.Currency = *req.Currency
	}

	if err = h.companyService.Update(id, company); err != nil {
		res.Error(w, http.StatusInternalServerError, "Failed to update company")
		return
	}

	// Get updated company
	var updatedCompany *model.Company
	updatedCompany, err = h.companyService.GetByID(id)
	if err != nil {
		res.Error(w, http.StatusInternalServerError, "Failed to retrieve updated company")
		return
	}

	companyDTO := &dto.CompanyDTO{
		ID:             updatedCompany.ID,
		TenantID:       updatedCompany.TenantID,
		Name:           updatedCompany.Name,
		RegistrationNo: updatedCompany.RegistrationNo,
		Address:        updatedCompany.Address,
		Timezone:       updatedCompany.Timezone,
		Currency:       updatedCompany.Currency,
		CreatedAt:      updatedCompany.CreatedAt,
		UpdatedAt:      updatedCompany.UpdatedAt,
	}

	res.Success(w, http.StatusOK, "Company profile updated successfully", companyDTO)
}

// GetCompanySettings godoc
// @Summary Get company settings
// @Description Get company settings by ID
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Company ID"
// @Success 200 {object} response.Response{data=dto.CompanySettingsDTO}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /companies/{id}/settings [get]
func (h *CompanyHandler) GetCompanySettings(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		res.Error(w, http.StatusBadRequest, "Invalid company ID")
		return
	}

	company, err := h.companyService.GetByID(id)
	if err != nil {
		res.Error(w, http.StatusNotFound, "Company not found")
		return
	}

	companySettingsDTO := &dto.CompanySettingsDTO{
		ID:             company.ID,
		TenantID:       company.TenantID,
		Name:           company.Name,
		RegistrationNo: company.RegistrationNo,
		Address:        company.Address,
		Timezone:       company.Timezone,
		Currency:       company.Currency,
	}

	res.Success(w, http.StatusOK, "Company settings retrieved successfully", companySettingsDTO)
}

// UpdateCompanySettings godoc
// @Summary Update company settings
// @Description Update company settings by ID
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Company ID"
// @Param request body dto.UpdateCompanySettingsRequest true "Update company settings request"
// @Success 200 {object} response.Response{data=dto.CompanySettingsDTO}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /companies/{id}/settings [patch]
func (h *CompanyHandler) UpdateCompanySettings(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		res.Error(w, http.StatusBadRequest, "Invalid company ID")
		return
	}

	var req dto.UpdateCompanySettingsRequest
	if err = res.ReadJSON(r, &req); err != nil {
		res.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Get existing company
	company, err := h.companyService.GetByID(id)
	if err != nil {
		res.Error(w, http.StatusNotFound, "Company not found")
		return
	}

	// Update fields if provided
	if req.Name != nil {
		company.Name = *req.Name
	}
	if req.RegistrationNo != nil {
		company.RegistrationNo = *req.RegistrationNo
	}
	if req.Address != nil {
		company.Address = *req.Address
	}
	if req.Timezone != nil {
		company.Timezone = *req.Timezone
	}
	if req.Currency != nil {
		company.Currency = *req.Currency
	}

	if err = h.companyService.Update(id, company); err != nil {
		res.Error(w, http.StatusInternalServerError, "Failed to update company settings")
		return
	}

	// Get updated company
	updatedCompany, err := h.companyService.GetByID(id)
	if err != nil {
		res.Error(w, http.StatusInternalServerError, "Failed to retrieve updated company")
		return
	}

	companySettingsDTO := &dto.CompanySettingsDTO{
		ID:             updatedCompany.ID,
		TenantID:       updatedCompany.TenantID,
		Name:           updatedCompany.Name,
		RegistrationNo: updatedCompany.RegistrationNo,
		Address:        updatedCompany.Address,
		Timezone:       updatedCompany.Timezone,
		Currency:       updatedCompany.Currency,
	}

	res.Success(w, http.StatusOK, "Company settings updated successfully", companySettingsDTO)
}

// GetCompanyLimits godoc
// @Summary Get company limits
// @Description Get company limits based on subscription plan
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Company ID"
// @Success 200 {object} response.Response{data=dto.CompanyLimitsDTO}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /companies/{id}/limits [get]
func (h *CompanyHandler) GetCompanyLimits(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		res.Error(w, http.StatusBadRequest, "Invalid company ID")
		return
	}

	company, err := h.companyService.GetLimits(id)
	if err != nil {
		res.Error(w, http.StatusNotFound, "Company not found")
		return
	}

	// For now, we'll use default values. In a real implementation,
	// these would be fetched from the subscription plan
	companyLimitsDTO := &dto.CompanyLimitsDTO{
		ID:             company.ID,
		TenantID:       company.TenantID,
		Name:           company.Name,
		MaxEmployees:   100, // Default value, should come from subscription plan
		MaxDepartments: 10,  // Default value, should come from subscription plan
		MaxPositions:   50,  // Default value, should come from subscription plan
	}

	res.Success(w, http.StatusOK, "Company limits retrieved successfully", companyLimitsDTO)
}
