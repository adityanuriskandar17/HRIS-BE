package dto

// UnitRequest represents the request body for creating a unit
type UnitRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// PositionRequest represents the request body for creating a position
type PositionRequest struct {
	CompanyID   string `json:"companyId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Level       string `json:"level"`
}

// EmployeeRequest represents the request body for creating an employee
type EmployeeRequest struct {
	EmployeeCode     string `json:"employeeCode"`
	FullName         string `json:"fullName"`
	Email            string `json:"email"`
	UnitID           string `json:"unitId"`
	PositionID       string `json:"positionId"`
	EmploymentStatus string `json:"employmentStatus"`
}