package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/model"
	httpx "github.com/adityanuriskandar17/HRIS-BE/internal/http"
	"gorm.io/gorm"
)

type MasterDataHandler struct {
	DB *gorm.DB
}

// Units

func (h *MasterDataHandler) ListUnits(w http.ResponseWriter, r *http.Request) {
	var units []model.Unit
	if err := h.DB.Order("code ASC").Find(&units).Error; err != nil {
		http.Error(w, "failed to fetch units", http.StatusInternalServerError)
		return
	}
	httpx.OK(w, units)
}

func (h *MasterDataHandler) CreateUnit(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Code string `json:"code"`
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	req.Code = strings.ToUpper(strings.TrimSpace(req.Code))
	req.Name = strings.TrimSpace(req.Name)
	if req.Code == "" || req.Name == "" {
		http.Error(w, "code and name required", http.StatusBadRequest)
		return
	}
	u := model.Unit{Code: req.Code, Name: req.Name}
	if err := h.DB.Create(&u).Error; err != nil {
		http.Error(w, "failed to create unit", http.StatusBadRequest)
		return
	}
	httpx.Created(w, u)
}

// Positions

func (h *MasterDataHandler) ListPositions(w http.ResponseWriter, r *http.Request) {
	var positions []model.Position
	if err := h.DB.Preload("Unit").Order("title ASC").Find(&positions).Error; err != nil {
		http.Error(w, "failed to fetch positions", http.StatusInternalServerError)
		return
	}
	httpx.OK(w, positions)
}

func (h *MasterDataHandler) CreatePosition(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title  string  `json:"title"`
		UnitID *uint64 `json:"unitId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		http.Error(w, "title required", http.StatusBadRequest)
		return
	}
	var unit *model.Unit
	if req.UnitID != nil {
		unit = &model.Unit{}
		if err := h.DB.First(unit, *req.UnitID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "unit not found", http.StatusBadRequest)
				return
			}
			http.Error(w, "failed to lookup unit", http.StatusInternalServerError)
			return
		}
	}
	pos := model.Position{Title: req.Title}
	if unit != nil {
		pos.UnitID = &unit.ID
	}
	if err := h.DB.Create(&pos).Error; err != nil {
		http.Error(w, "failed to create position", http.StatusBadRequest)
		return
	}
	if unit != nil {
		pos.Unit = unit
	}
	httpx.Created(w, pos)
}

// Employees

func (h *MasterDataHandler) ListEmployees(w http.ResponseWriter, r *http.Request) {
	var employees []model.Employee
	if err := h.DB.Preload("Unit").Preload("Position").Order("employee_code ASC").Find(&employees).Error; err != nil {
		http.Error(w, "failed to fetch employees", http.StatusInternalServerError)
		return
	}
	httpx.OK(w, employees)
}

func (h *MasterDataHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var req struct {
		EmployeeCode     string `json:"employeeCode"`
		FullName         string `json:"fullName"`
		Email            string `json:"email"`
		Phone            string `json:"phone"`
		UnitID           uint64 `json:"unitId"`
		PositionID       uint64 `json:"positionId"`
		EmploymentStatus string `json:"employmentStatus"`
		StartDate        string `json:"startDate"`
		EndDate          string `json:"endDate"`
		DateOfBirth      string `json:"dateOfBirth"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	req.EmployeeCode = strings.ToUpper(strings.TrimSpace(req.EmployeeCode))
	req.FullName = strings.TrimSpace(req.FullName)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Phone = strings.TrimSpace(req.Phone)

	if req.EmployeeCode == "" || req.FullName == "" || req.Email == "" || req.UnitID == 0 || req.PositionID == 0 {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	var startDate time.Time
	if req.StartDate != "" {
		parsed, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			http.Error(w, "invalid startDate format (expected YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		startDate = parsed
	} else {
		startDate = time.Now()
	}

	var endDate *time.Time
	if req.EndDate != "" {
		parsed, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			http.Error(w, "invalid endDate format (expected YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		endDate = &parsed
	}

	var dob *time.Time
	if req.DateOfBirth != "" {
		parsed, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err != nil {
			http.Error(w, "invalid dateOfBirth format (expected YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		dob = &parsed
	}

	status := model.EmploymentStatus(strings.ToUpper(req.EmploymentStatus))
	if status == "" {
		status = model.EmploymentFullTime
	}

	if !isValidEmploymentStatus(status) {
		http.Error(w, "invalid employmentStatus value", http.StatusBadRequest)
		return
	}

	var unit model.Unit
	if err := h.DB.First(&unit, req.UnitID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "unit not found", http.StatusBadRequest)
			return
		}
		http.Error(w, "failed to lookup unit", http.StatusInternalServerError)
		return
	}

	var position model.Position
	if err := h.DB.Preload("Unit").First(&position, req.PositionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "position not found", http.StatusBadRequest)
			return
		}
		http.Error(w, "failed to lookup position", http.StatusInternalServerError)
		return
	}

	if position.UnitID != nil && *position.UnitID != unit.ID {
		http.Error(w, "position does not belong to the selected unit", http.StatusBadRequest)
		return
	}

	emp := model.Employee{
		EmployeeCode:     req.EmployeeCode,
		FullName:         req.FullName,
		Email:            req.Email,
		Phone:            req.Phone,
		UnitID:           req.UnitID,
		PositionID:       req.PositionID,
		EmploymentStatus: status,
		StartDate:        startDate,
		EndDate:          endDate,
		DateOfBirth:      dob,
	}

	if err := h.DB.Create(&emp).Error; err != nil {
		http.Error(w, "failed to create employee", http.StatusBadRequest)
		return
	}

	emp.Unit = unit
	emp.Position = position

	httpx.Created(w, emp)
}

func isValidEmploymentStatus(status model.EmploymentStatus) bool {
	switch status {
	case model.EmploymentFullTime,
		model.EmploymentContract,
		model.EmploymentIntern,
		model.EmploymentPartTime:
		return true
	default:
		return false
	}
}
