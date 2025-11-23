package db

import (
	"fmt"

	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/model"

	"gorm.io/gorm"
)

// AutoMigrate applies schema changes required for the MVP features.
func AutoMigrate(gdb *gorm.DB) error {

	// Check if units table exists with BIGSERIAL ID and drop it if needed
	if gdb.Migrator().HasTable(&model.Unit{}) {
		var columnType string
		gdb.Raw("SELECT data_type FROM information_schema.columns WHERE table_name = 'units' AND column_name = 'id'").Scan(&columnType)
		if columnType == "bigint" {
			if err := gdb.Migrator().DropTable(&model.Unit{}); err != nil {
				return fmt.Errorf("drop units table: %w", err)
			}
		}
	}

	// Check if employees table exists with BIGSERIAL ID and drop it if needed
	if gdb.Migrator().HasTable(&model.Employee{}) {
		var columnType string
		gdb.Raw("SELECT data_type FROM information_schema.columns WHERE table_name = 'employees' AND column_name = 'id'").Scan(&columnType)
		if columnType == "bigint" {
			if err := gdb.Migrator().DropTable(&model.Employee{}); err != nil {
				return fmt.Errorf("drop employees table: %w", err)
			}
		}
	}

	// Check if user_accounts table exists with BIGINT ID and drop it if needed
	if gdb.Migrator().HasTable(&model.UserAccount{}) {
		var columnType string
		gdb.Raw("SELECT data_type FROM information_schema.columns WHERE table_name = 'user_accounts' AND column_name = 'id'").Scan(&columnType)
		if columnType == "bigint" {
			if err := gdb.Migrator().DropTable(&model.UserAccount{}); err != nil {
				return fmt.Errorf("drop user_accounts table: %w", err)
			}
		}
	}

	// Check if refresh_tokens table exists with BIGINT user_id and drop it if needed
	if gdb.Migrator().HasTable(&model.RefreshToken{}) {
		var columnType string
		gdb.Raw("SELECT data_type FROM information_schema.columns WHERE table_name = 'refresh_tokens' AND column_name = 'user_id'").Scan(&columnType)
		if columnType == "bigint" {
			if err := gdb.Migrator().DropTable(&model.RefreshToken{}); err != nil {
				return fmt.Errorf("drop refresh_tokens table: %w", err)
			}
		}
	}

	// Handle units table type conversion from BIGINT to UUID
	if gdb.Migrator().HasTable("units") {
		// Check if the units table uses BIGINT for id
		var isBigIntID bool
		gdb.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'units' AND column_name = 'id' AND data_type = 'bigint')").Scan(&isBigIntID)

		if isBigIntID {
			// Drop the existing units table to avoid conflicts
			gdb.Exec("DROP TABLE IF EXISTS units")
			fmt.Printf("Dropped existing units table with BIGINT columns\n")
		}

		// Check if the units table still exists and print its structure
		if gdb.Migrator().HasTable("units") {
			var columnInfo []struct {
				ColumnName string
				DataType   string
			}
			gdb.Raw("SELECT column_name, data_type FROM information_schema.columns WHERE table_name = 'units' ORDER BY ordinal_position").Scan(&columnInfo)
			fmt.Printf("Units table structure:\n")
			for _, col := range columnInfo {
				fmt.Printf("  %s: %s\n", col.ColumnName, col.DataType)
			}
		}
	}

	// Handle positions table type conversion from BIGINT to UUID
	if gdb.Migrator().HasTable("positions") {
		// Check if the positions table uses BIGINT for id
		var isBigIntID bool
		gdb.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'positions' AND column_name = 'id' AND data_type = 'bigint')").Scan(&isBigIntID)

		if isBigIntID {
			// Drop the existing positions table to avoid conflicts
			gdb.Exec("DROP TABLE IF EXISTS positions")
			fmt.Printf("Dropped existing positions table with BIGINT columns\n")

		}
	}

	// Pre-migration statements (Enums)
	preStatements := []string{
		`DO $$ BEGIN
  CREATE TYPE employment_status_enum AS ENUM ('FULLTIME','PARTTIME','CONTRACT','INTERN');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;`,
	}

	for _, stmt := range preStatements {
		if err := gdb.Exec(stmt).Error; err != nil {
			return fmt.Errorf("auto migrate pre: executing %q: %w", stmt, err)
		}
	}

	// Perform AutoMigrate
	if err := gdb.AutoMigrate(
		&model.Unit{},
		&model.Position{},
		&model.Employee{},
		&model.UserAccount{},
		&model.Tenant{},
		&model.Subscription{},
		&model.Plan{},
		&model.Invoice{},
		&model.Company{},
		&model.RefreshToken{},
	); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}

	enumStatements := []string{
		`ALTER TABLE public.employees
  ALTER COLUMN employment_status TYPE employment_status_enum USING employment_status::employment_status_enum`,
	}

	for _, stmt := range enumStatements {
		if err := gdb.Exec(stmt).Error; err != nil {
			return fmt.Errorf("auto migrate enum: executing %q: %w", stmt, err)
		}
	}

	// Enable UUID extension before creating indexes
	uuidStatements := []string{
		"CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"",
	}

	for _, stmt := range uuidStatements {
		if err := gdb.Exec(stmt).Error; err != nil {
			return fmt.Errorf("auto migrate uuid: executing %q: %w", stmt, err)
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
		"CREATE INDEX IF NOT EXISTS idx_employees_name_trgm ON public.employees USING gin (full_name gin_trgm_ops)",
	}

	for _, stmt := range indexStatements {
		if err := gdb.Exec(stmt).Error; err != nil {
			return fmt.Errorf("auto migrate index: executing %q: %w", stmt, err)
		}
	}

	return nil
}
