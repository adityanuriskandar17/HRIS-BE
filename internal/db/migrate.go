package db

import (
	"fmt"

	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/model"
	"gorm.io/gorm"
)

// AutoMigrate applies schema changes required for the MVP features.
func AutoMigrate(gdb *gorm.DB) error {
	if err := gdb.AutoMigrate(
		&model.Unit{},
		&model.Position{},
		&model.Employee{},
		&model.UserAccount{},
	); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}
	return nil
}
