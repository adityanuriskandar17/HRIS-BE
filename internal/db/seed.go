package db

import (
	"fmt"
	"time"

	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/model"
	"gorm.io/gorm"
)

// SeedReferenceData ensures a minimal set of master data exists for development.
func SeedReferenceData(gdb *gorm.DB) error {
	if err := seedUnits(gdb); err != nil {
		return err
	}
	if err := seedPositions(gdb); err != nil {
		return err
	}
	if err := seedEmployees(gdb); err != nil {
		return err
	}
	return nil
}

func seedUnits(gdb *gorm.DB) error {
	units := []model.Unit{
		{Code: "HRD", Name: "Human Resources"},
		{Code: "FIN", Name: "Finance"},
		{Code: "ENG", Name: "Engineering"},
	}
	for _, u := range units {
		var existing model.Unit
		if err := gdb.Where("code = ?", u.Code).First(&existing).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return fmt.Errorf("seed unit %s: %w", u.Code, err)
			}
		}
		if existing.ID == 0 {
			if err := gdb.Create(&u).Error; err != nil {
				return fmt.Errorf("create unit %s: %w", u.Code, err)
			}
		}
	}
	return nil
}

func seedPositions(gdb *gorm.DB) error {
	type seed struct {
		Title string
		Unit  string
	}
	seeds := []seed{
		{Title: "HR Manager", Unit: "HRD"},
		{Title: "Finance Analyst", Unit: "FIN"},
		{Title: "Software Engineer", Unit: "ENG"},
	}
	for _, s := range seeds {
		var unit model.Unit
		if err := gdb.Where("code = ?", s.Unit).First(&unit).Error; err != nil {
			return fmt.Errorf("lookup unit %s for position seed: %w", s.Unit, err)
		}
		var existing model.Position
		if err := gdb.Where("title = ?", s.Title).First(&existing).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return fmt.Errorf("seed position %s: %w", s.Title, err)
			}
		}
		if existing.ID == 0 {
			p := model.Position{Title: s.Title, UnitID: &unit.ID}
			if err := gdb.Create(&p).Error; err != nil {
				return fmt.Errorf("create position %s: %w", s.Title, err)
			}
		}
	}
	return nil
}

func seedEmployees(gdb *gorm.DB) error {
	type seed struct {
		Code string
		Name string
		Mail string
		Unit string
		Role string
	}
	seeds := []seed{
		{Code: "EMP001", Name: "Dita HR", Mail: "dita.hr@example.com", Unit: "HRD", Role: "HR Manager"},
		{Code: "EMP002", Name: "Andi Finance", Mail: "andi.finance@example.com", Unit: "FIN", Role: "Finance Analyst"},
		{Code: "EMP003", Name: "Budi Engineer", Mail: "budi.engineer@example.com", Unit: "ENG", Role: "Software Engineer"},
	}
	for _, s := range seeds {
		var existing model.Employee
		if err := gdb.Where("employee_code = ?", s.Code).First(&existing).Error; err == nil {
			continue
		} else if err != gorm.ErrRecordNotFound {
			return fmt.Errorf("seed employee %s: %w", s.Code, err)
		}

		var unit model.Unit
		if err := gdb.Where("code = ?", s.Unit).First(&unit).Error; err != nil {
			return fmt.Errorf("lookup unit for employee seed %s: %w", s.Code, err)
		}

		var position model.Position
		if err := gdb.Where("title = ?", s.Role).First(&position).Error; err != nil {
			return fmt.Errorf("lookup position for employee seed %s: %w", s.Code, err)
		}

		emp := model.Employee{
			EmployeeCode:     s.Code,
			FullName:         s.Name,
			Email:            s.Mail,
			UnitID:           unit.ID,
			PositionID:       position.ID,
			EmploymentStatus: model.EmploymentFullTime,
			StartDate:        time.Now().AddDate(-1, 0, 0),
		}
		if err := gdb.Create(&emp).Error; err != nil {
			return fmt.Errorf("create employee seed %s: %w", s.Code, err)
		}
	}
	return nil
}
