package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/model"
	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/services"
	"github.com/adityanuriskandar17/HRIS-BE/internal/http/dto"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type InvoiceHandler struct {
	invoiceService services.InvoiceService
}

func NewInvoiceHandler(invoiceService services.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{invoiceService: invoiceService}
}

// GetAll godoc
// @Summary Get all invoices
// @Description Get all invoices
// @Tags invoices
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "Invoice status"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Success 200 {array} dto.InvoiceResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /invoices [get]
func (h *InvoiceHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	invoices, err := h.invoiceService.GetAll(r.Context(), status, page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response DTOs
	responses := make([]dto.InvoiceResponse, len(invoices))
	for i, inv := range invoices {
		responses[i] = dto.InvoiceResponse{
			ID:             inv.ID.String(),
			SubscriptionID: inv.SubscriptionID.String(),
			TenantID:       inv.TenantID.String(),
			InvoiceNumber:  inv.InvoiceNumber,
			Amount:         inv.Amount,
			DueDate:        inv.DueDate,
			Status:         string(inv.Status),
			SentAt:         inv.SentAt,
			PaidAt:         inv.PaidAt,
			PaymentID:      inv.PaymentID,
			CreatedAt:      inv.CreatedAt,
			UpdatedAt:      inv.UpdatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

// GetByID godoc
// @Summary Get invoice by ID
// @Description Get invoice by ID
// @Tags invoices
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Invoice ID"
// @Success 200 {object} dto.InvoiceResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /invoices/{id} [get]
func (h *InvoiceHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	invoiceID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "invalid invoice ID", http.StatusBadRequest)
		return
	}

	invoice, err := h.invoiceService.GetByID(r.Context(), invoiceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.InvoiceResponse{
		ID:             invoice.ID.String(),
		SubscriptionID: invoice.SubscriptionID.String(),
		TenantID:       invoice.TenantID.String(),
		InvoiceNumber:  invoice.InvoiceNumber,
		Amount:         invoice.Amount,
		DueDate:        invoice.DueDate,
		Status:         string(invoice.Status),
		SentAt:         invoice.SentAt,
		PaidAt:         invoice.PaidAt,
		PaymentID:      invoice.PaymentID,
		CreatedAt:      invoice.CreatedAt,
		UpdatedAt:      invoice.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetBySubscriptionID godoc
// @Summary Get invoices by subscription ID
// @Description Get invoices by subscription ID
// @Tags invoices
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param subscriptionId path string true "Subscription ID"
// @Success 200 {array} dto.InvoiceResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /invoices/subscription/{subscriptionId} [get]
func (h *InvoiceHandler) GetBySubscriptionID(w http.ResponseWriter, r *http.Request) {
	subscriptionID := chi.URLParam(r, "subscriptionId")
	id, err := uuid.Parse(subscriptionID)
	if err != nil {
		http.Error(w, "invalid subscription ID", http.StatusBadRequest)
		return
	}

	invoices, err := h.invoiceService.GetBySubscriptionID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response DTOs
	responses := make([]dto.InvoiceResponse, len(invoices))
	for i, inv := range invoices {
		responses[i] = dto.InvoiceResponse{
			ID:             inv.ID.String(),
			SubscriptionID: inv.SubscriptionID.String(),
			Amount:         inv.Amount,
			DueDate:        inv.DueDate,
			Status:         string(inv.Status),
			SentAt:         inv.SentAt,
			PaidAt:         inv.PaidAt,
			CreatedAt:      inv.CreatedAt,
			UpdatedAt:      inv.UpdatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

// Create godoc
// @Summary Create a new invoice
// @Description Create a new invoice (Admin only)
// @Tags invoices
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param invoice body model.Invoice true "Invoice data"
// @Success 201 {object} dto.InvoiceResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /invoices [post]
func (h *InvoiceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var invoice model.Invoice
	if err := json.NewDecoder(r.Body).Decode(&invoice); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.invoiceService.Create(r.Context(), &invoice); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.InvoiceResponse{
		ID:             invoice.ID.String(),
		SubscriptionID: invoice.SubscriptionID.String(),
		TenantID:       invoice.TenantID.String(),
		InvoiceNumber:  invoice.InvoiceNumber,
		Amount:         invoice.Amount,
		DueDate:        invoice.DueDate,
		Status:         string(invoice.Status),
		SentAt:         invoice.SentAt,
		PaidAt:         invoice.PaidAt,
		PaymentID:      invoice.PaymentID,
		CreatedAt:      invoice.CreatedAt,
		UpdatedAt:      invoice.UpdatedAt,
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Update godoc
// @Summary Update an invoice
// @Description Update an invoice (Admin only)
// @Tags invoices
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Invoice ID"
// @Param invoice body model.Invoice true "Invoice data"
// @Success 200 {object} dto.InvoiceResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /invoices/{id} [put]
func (h *InvoiceHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	invoiceID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "invalid invoice ID", http.StatusBadRequest)
		return
	}

	var invoice model.Invoice
	if err := json.NewDecoder(r.Body).Decode(&invoice); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.invoiceService.Update(r.Context(), invoiceID, &invoice); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.InvoiceResponse{
		ID:             invoice.ID.String(),
		SubscriptionID: invoice.SubscriptionID.String(),
		Amount:         invoice.Amount,
		DueDate:        invoice.DueDate,
		Status:         string(invoice.Status),
		SentAt:         invoice.SentAt,
		PaidAt:         invoice.PaidAt,
		CreatedAt:      invoice.CreatedAt,
		UpdatedAt:      invoice.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Send godoc
// @Summary Send an invoice
// @Description Send an invoice (Admin only)
// @Tags invoices
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Invoice ID"
// @Param send body dto.SendInvoiceRequest true "Send details"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /invoices/{id}/send [post]
func (h *InvoiceHandler) Send(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	invoiceID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "invalid invoice ID", http.StatusBadRequest)
		return
	}

	var req dto.SendInvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.invoiceService.Send(r.Context(), invoiceID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "Invoice sent successfully",
		"email":   req.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Pay godoc
// @Summary Pay an invoice
// @Description Pay an invoice
// @Tags invoices
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Invoice ID"
// @Param pay body dto.PayInvoiceRequest true "Payment details"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /invoices/{id}/pay [post]
func (h *InvoiceHandler) Pay(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	invoiceID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "invalid invoice ID", http.StatusBadRequest)
		return
	}

	var req dto.PayInvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.invoiceService.Pay(r.Context(), invoiceID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message":       "Invoice paid successfully",
		"paymentId":     req.PaymentID,
		"paymentMethod": req.PaymentMethod,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
