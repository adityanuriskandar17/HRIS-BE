package db

import (
	"fmt"
	"log"

	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/model"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	// Seed Units
	seedUnits(db)

	// Seed Positions
	seedPositions(db)

	// Seed Employees
	seedEmployees(db)

	// Seed UserAccounts
	seedUserAccounts(db)
}

func seedUnits(db *gorm.DB) {
	units := []model.Unit{
		{Code: "IT", Name: "Information Technology"},
		{Code: "HR", Name: "Human Resources"},
		{Code: "FIN", Name: "Finance"},
		{Code: "OPS", Name: "Operations"},
	}

	for _, unit := range units {
		var existingUnit model.Unit
		if err := db.Where("code = ?", unit.Code).First(&existingUnit).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&unit).Error; err != nil {
					log.Printf("Failed to seed unit %s: %v", unit.Code, err)
				}
			}
		}
	}
}

func seedPositions(db *gorm.DB) {
	// Get the first company to associate with positions
	var company model.Company
	if err := db.First(&company).Error; err != nil {
		log.Printf("Failed to get company for seeding positions: %v", err)
		return
	}

	positions := []model.Position{
		{Title: "Software Engineer", Description: "Develops software applications"},
		{Title: "HR Manager", Description: "Manages human resources"},
		{Title: "Financial Analyst", Description: "Analyzes financial data"},
		{Title: "Operations Manager", Description: "Manages operations"},
	}

	for _, position := range positions {
		var existingPosition model.Position
		if err := db.Where("title = ?", position.Title).First(&existingPosition).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Set CompanyID to the existing company ID
				position.CompanyID = company.ID
				if err := db.Create(&position).Error; err != nil {
					log.Printf("Failed to seed position %s: %v", position.Title, err)
				}
			}
		}
	}
}

func seedEmployees(db *gorm.DB) {
	// Get the first unit and position to associate with employees
	var unit model.Unit
	if err := db.First(&unit).Error; err != nil {
		log.Printf("Failed to get unit for seeding employees: %v", err)
		return
	}

	var position model.Position
	if err := db.First(&position).Error; err != nil {
		log.Printf("Failed to get position for seeding employees: %v", err)
		return
	}

	// Create employees
	for i := 1; i <= 3; i++ {
		employee := model.Employee{
			EmployeeCode: fmt.Sprintf("EMP%03d", i),
			FullName:     fmt.Sprintf("Employee %d", i),
			Email:        fmt.Sprintf("employee%d@example.com", i),
			UnitID:       unit.ID,
			PositionID:   position.ID,
		}
		var existingEmployee model.Employee
		if err := db.Where("employee_code = ?", employee.EmployeeCode).First(&existingEmployee).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&employee).Error; err != nil {
					log.Printf("Failed to seed employee %s: %v", employee.EmployeeCode, err)
				}
			}
		}
	}
}

func seedUserAccounts(db *gorm.DB) {
	// Get the first employee to associate with user accounts
	var employee model.Employee
	if err := db.First(&employee).Error; err != nil {
		log.Printf("Failed to get employee for seeding user accounts: %v", err)
		return
	}

	// Get the first tenant to associate with user accounts
	var tenant model.Tenant
	if err := db.First(&tenant).Error; err != nil {
		log.Printf("Failed to get tenant for seeding user accounts: %v", err)
		return
	}

	userAccounts := []model.UserAccount{
		{
			TenantID:     tenant.ID,
			Email:        "admin@example.com",
			PasswordHash: "hashed_password_here",
			FirstName:    "Admin",
			LastName:     "User",
		},
		{
			TenantID:     tenant.ID,
			Email:        employee.Email,
			PasswordHash: "hashed_password_here",
			FirstName:    employee.FullName,
			LastName:     "",
		},
	}

	for _, userAccount := range userAccounts {
		var existingUserAccount model.UserAccount
		if err := db.Where("email = ?", userAccount.Email).First(&existingUserAccount).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&userAccount).Error; err != nil {
					log.Printf("Failed to seed user account %s: %v", userAccount.Email, err)
				}
			}
		}
	}
}
