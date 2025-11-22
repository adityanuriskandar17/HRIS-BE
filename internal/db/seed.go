package db

import (
	"fmt"
	"log"

	"github.com/adityanuriskandar17/HRIS-BE/internal/auth"
	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/model"
	"gorm.io/gorm"
)

// SeedReferenceData ensures a minimal set of master data exists for development.
func SeedReferenceData(gdb *gorm.DB, adminEmail, adminPassword string) error {
	if err := seedTenants(gdb); err != nil {
		return err
	}
	if err := seedCompanies(gdb); err != nil {
		return err
	}
	if err := seedUnits(gdb); err != nil {
		return err
	}
	if err := seedPositions(gdb); err != nil {
		return err
	}
	if err := seedEmployees(gdb); err != nil {
		return err
	}
	if err := seedUserAccounts(gdb, adminEmail, adminPassword); err != nil {
		return err
	}
	return nil
}

func seedTenants(db *gorm.DB) error {
	tenants := []model.Tenant{
		{
			Name:        "Default Tenant",
			Email:       "tenant@example.com",
			CompanyName: "Default Company",
			Domain:      "default.local",
		},
	}

	for _, tenant := range tenants {
		var existingTenant model.Tenant
		if err := db.Where("email = ?", tenant.Email).First(&existingTenant).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&tenant).Error; err != nil {
					log.Printf("Failed to seed tenant %s: %v", tenant.Email, err)
					return err
				}
			} else {
				return err
			}
		}
	}
	return nil
}

func seedCompanies(db *gorm.DB) error {
	// Get the first tenant to associate with companies
	var tenant model.Tenant
	if err := db.First(&tenant).Error; err != nil {
		log.Printf("Failed to get tenant for seeding companies: %v", err)
		return err
	}

	companies := []model.Company{
		{Name: "Default Company", TenantID: tenant.ID},
	}

	for _, company := range companies {
		var existingCompany model.Company
		if err := db.Where("name = ? AND tenant_id = ?", company.Name, company.TenantID).First(&existingCompany).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&company).Error; err != nil {
					log.Printf("Failed to seed company %s: %v", company.Name, err)
					return err
				}
			} else {
				return err
			}
		}
	}
	return nil
}

func seedUnits(db *gorm.DB) error {
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
					return err
				}
			} else {
				return err
			}
		}
	}
	return nil
}

func seedPositions(db *gorm.DB) error {
	// Get the first company to associate with positions
	var company model.Company
	if err := db.First(&company).Error; err != nil {
		log.Printf("Failed to get company for seeding positions: %v", err)
		return err
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
					return err
				}
			} else {
				return err
			}
		}
	}
	return nil
}

func seedEmployees(db *gorm.DB) error {
	// Get the first unit and position to associate with employees
	var unit model.Unit
	if err := db.First(&unit).Error; err != nil {
		log.Printf("Failed to get unit for seeding employees: %v", err)
		return err
	}

	var position model.Position
	if err := db.First(&position).Error; err != nil {
		log.Printf("Failed to get position for seeding employees: %v", err)
		return err
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
					return err
				}
			} else {
				return err
			}
		}
	}
	return nil
}

func seedUserAccounts(db *gorm.DB, adminEmail, adminPassword string) error {
	// Get the first employee to associate with user accounts
	var employee model.Employee
	if err := db.First(&employee).Error; err != nil {
		log.Printf("Failed to get employee for seeding user accounts: %v", err)
		return err
	}

	// Get the first tenant to associate with user accounts
	var tenant model.Tenant
	if err := db.First(&tenant).Error; err != nil {
		log.Printf("Failed to get tenant for seeding user accounts: %v", err)
		return err
	}

	adminHash, err := auth.HashPassword(adminPassword)
	if err != nil {
		return fmt.Errorf("hash admin password: %w", err)
	}

	// Default password for other users
	defaultHash, _ := auth.HashPassword("password")

	userAccounts := []model.UserAccount{
		{
			TenantID:     tenant.ID,
			Email:        adminEmail,
			PasswordHash: adminHash,
			FirstName:    "Admin",
			LastName:     "User",
			Role:         model.RoleAdmin,
			IsActive:     true,
		},
		{
			TenantID:     tenant.ID,
			Email:        employee.Email,
			PasswordHash: defaultHash,
			FirstName:    employee.FullName,
			LastName:     "",
			Role:         model.RoleEmployee,
			IsActive:     true,
		},
	}

	for _, userAccount := range userAccounts {
		var existingUserAccount model.UserAccount
		if err := db.Where("email = ?", userAccount.Email).First(&existingUserAccount).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&userAccount).Error; err != nil {
					log.Printf("Failed to seed user account %s: %v", userAccount.Email, err)
					return err
				}
			}
		}
	}
	return nil
}


