package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/model"
	"github.com/adityanuriskandar17/HRIS-BE/internal/http/dto"
	"github.com/adityanuriskandar17/HRIS-BE/internal/repository"
	"github.com/gorilla/mux"
	"github.com/google/uuid"
)

type MasterDataHandler struct {
	unitRepo     repository.UnitRepository
	positionRepo repository.PositionRepository
	employeeRepo repository.EmployeeRepository
}

func NewMasterDataHandler(
	unitRepo repository.UnitRepository,
	positionRepo repository.PositionRepository,
	employeeRepo repository.EmployeeRepository,
) *MasterDataHandler {
	return &MasterDataHandler{
		unitRepo:     unitRepo,
		positionRepo: positionRepo,
		employeeRepo: employeeRepo,
	}
}

// Unit handlers
// ListUnits handles listing all units
// @Summary List all units
// @Description Retrieve a list of all units
// @Tags master-data
// @Accept json
// @Produce json
// @Success 200 {array} model.Unit
// @Failure 500 {object} string
// @Router /master/units [get]
func (h *MasterDataHandler) ListUnits(w http.ResponseWriter, r *http.Request) {
	units, err := h.unitRepo.FindAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(units)
}

// CreateUnit handles creating a new unit
// @Summary Create a new unit
// @Description Create a new unit in the system
// @Tags master-data
// @Accept json
// @Produce json
// @Param request body dto.UnitRequest true "Unit data"
// @Success 201 {object} model.Unit
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /master/units [post]
func (h *MasterDataHandler) CreateUnit(w http.ResponseWriter, r *http.Request) {
	var unitReq dto.UnitRequest
	if err := json.NewDecoder(r.Body).Decode(&unitReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	unit := model.Unit{
		Code: unitReq.Code,
		Name: unitReq.Name,
	}

	createdUnit, err := h.unitRepo.Create(r.Context(), unit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdUnit)
}

// Position handlers
// ListPositions handles listing all positions
// @Summary List all positions
// @Description Retrieve a list of all positions
// @Tags master-data
// @Accept json
// @Produce json
// @Success 200 {array} model.Position
// @Failure 500 {object} string
// @Router /master/positions [get]
func (h *MasterDataHandler) ListPositions(w http.ResponseWriter, r *http.Request) {
	positions, err := h.positionRepo.FindAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(positions)
}

// CreatePosition handles creating a new position
// @Summary Create a new position
// @Description Create a new position in the system
// @Tags master-data
// @Accept json
// @Produce json
// @Param request body dto.PositionRequest true "Position data"
// @Success 201 {object} model.Position
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /master/positions [post]
func (h *MasterDataHandler) CreatePosition(w http.ResponseWriter, r *http.Request) {
	var positionReq dto.PositionRequest
	if err := json.NewDecoder(r.Body).Decode(&positionReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse company ID
	companyID, err := uuid.Parse(positionReq.CompanyID)
	if err != nil {
		http.Error(w, "Invalid company ID", http.StatusBadRequest)
		return
	}

	// Convert Level from string to int
	level, err := strconv.Atoi(positionReq.Level)
	if err != nil {
		http.Error(w, "Invalid level format", http.StatusBadRequest)
		return
	}

	position := model.Position{
		CompanyID:   companyID,
		Title:       positionReq.Title,
		Description: positionReq.Description,
		Level:       level,
	}

	createdPosition, err := h.positionRepo.Create(r.Context(), position)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdPosition)
}

// Employee handlers
// ListEmployees handles listing all employees
// @Summary List all employees
// @Description Retrieve a list of all employees
// @Tags master-data
// @Accept json
// @Produce json
// @Success 200 {array} model.Employee
// @Failure 500 {object} string
// @Router /master/employees [get]
func (h *MasterDataHandler) ListEmployees(w http.ResponseWriter, r *http.Request) {
	employees, err := h.employeeRepo.FindAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employees)
}

// CreateEmployee handles creating a new employee
// @Summary Create a new employee
// @Description Create a new employee in the system
// @Tags master-data
// @Accept json
// @Produce json
// @Param request body dto.EmployeeRequest true "Employee data"
// @Success 201 {object} model.Employee
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /master/employees [post]
func (h *MasterDataHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var employeeReq dto.EmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&employeeReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse unit ID
	unitUUID, err := uuid.Parse(employeeReq.UnitID)
	if err != nil {
		http.Error(w, "Invalid unit ID", http.StatusBadRequest)
		return
	}

	// Parse position ID
	positionUUID, err := uuid.Parse(employeeReq.PositionID)
	if err != nil {
		http.Error(w, "Invalid position ID", http.StatusBadRequest)
		return
	}

	// Convert string to EmploymentStatus
	var employmentStatus model.EmploymentStatus
	switch employeeReq.EmploymentStatus {
	case "FULLTIME":
		employmentStatus = model.EmploymentFullTime
	case "PARTTIME":
		employmentStatus = model.EmploymentPartTime
	case "CONTRACT":
		employmentStatus = model.EmploymentContract
	case "INTERN":
		employmentStatus = model.EmploymentIntern
	default:
		employmentStatus = model.EmploymentFullTime // default value
	}

	employee := model.Employee{
		EmployeeCode:     employeeReq.EmployeeCode,
		FullName:         employeeReq.FullName,
		Email:            employeeReq.Email,
		UnitID:           unitUUID,
		PositionID:       positionUUID,
		EmploymentStatus: employmentStatus,
	}

	createdEmployee, err := h.employeeRepo.Create(r.Context(), employee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdEmployee)
}

// GetEmployee handles retrieving an employee by ID
// @Summary Get employee by ID
// @Description Retrieve a specific employee by their ID
// @Tags master-data
// @Accept json
// @Produce json
// @Param id path string true "Employee ID"
// @Success 200 {object} model.Employee
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Router /master/employees/{id} [get]
func (h *MasterDataHandler) GetEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	employee, err := h.employeeRepo.FindByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employee)
}

// UpdateEmployee handles updating an employee
// @Summary Update an employee
// @Description Update an existing employee's information
// @Tags master-data
// @Accept json
// @Produce json
// @Param id path string true "Employee ID"
// @Param request body dto.EmployeeRequest true "Updated employee data"
// @Success 200 {object} model.Employee
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /master/employees/{id} [put]
func (h *MasterDataHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	var employeeReq dto.EmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&employeeReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse unit ID
	unitID, err := uuid.Parse(employeeReq.UnitID)
	if err != nil {
		http.Error(w, "Invalid unit ID", http.StatusBadRequest)
		return
	}

	// Parse position ID
	positionID, err := uuid.Parse(employeeReq.PositionID)
	if err != nil {
		http.Error(w, "Invalid position ID", http.StatusBadRequest)
		return
	}

	employee := model.Employee{
		ID:           id,
		EmployeeCode: employeeReq.EmployeeCode,
		FullName:     employeeReq.FullName,
		Email:        employeeReq.Email,
		UnitID:       unitID,
		PositionID:   positionID,
	}

	updatedEmployee, err := h.employeeRepo.Update(r.Context(), employee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedEmployee)
}

// DeleteEmployee handles deleting an employee
// @Summary Delete an employee
// @Description Delete an employee from the system
// @Tags master-data
// @Accept json
// @Produce json
// @Param id path string true "Employee ID"
// @Success 204 {object} string
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Router /master/employees/{id} [delete]
func (h *MasterDataHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	if err := h.employeeRepo.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
