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

	enumStatements := []string{
		`DO $$ BEGIN
  CREATE TYPE employment_status_enum AS ENUM ('FULLTIME','PARTTIME','CONTRACT','INTERN');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;`,
		`ALTER TABLE public.employees
  ALTER COLUMN employment_status TYPE employment_status_enum USING employment_status::employment_status_enum`,
	}

	for _, stmt := range enumStatements {
		if err := gdb.Exec(stmt).Error; err != nil {
			return fmt.Errorf("auto migrate enum: executing %q: %w", stmt, err)
		}
	}

	indexStatements := []string{
		"CREATE EXTENSION IF NOT EXISTS pg_trgm",
		"DROP INDEX IF EXISTS employees_email_key",
		"CREATE UNIQUE INDEX IF NOT EXISTS employees_email_key_ci ON public.employees (LOWER(email))",
		"DROP INDEX IF EXISTS user_accounts_email_key",
		"CREATE UNIQUE INDEX IF NOT EXISTS user_accounts_email_key_ci ON public.user_accounts (LOWER(email))",
		"CREATE INDEX IF NOT EXISTS idx_employees_unit ON public.employees (unit_id)",
		"CREATE INDEX IF NOT EXISTS idx_employees_position ON public.employees (position_id)",
		"CREATE INDEX IF NOT EXISTS idx_positions_unit ON public.positions (unit_id)",
		"CREATE INDEX IF NOT EXISTS idx_employees_name_trgm ON public.employees USING gin (full_name gin_trgm_ops)",
	}

	for _, stmt := range indexStatements {
		if err := gdb.Exec(stmt).Error; err != nil {
			return fmt.Errorf("auto migrate index: executing %q: %w", stmt, err)
		}
	}

	return nil
}
