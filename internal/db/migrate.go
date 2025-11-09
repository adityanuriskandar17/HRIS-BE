package db

import (
	"fmt"

	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/model"
	"github.com/google/uuid"
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
	
	// Handle positions table with null company_id values
	if gdb.Migrator().HasTable(&model.Position{}) {
		// Drop the specific foreign key constraint created by the Company struct relationship
		gdb.Exec("ALTER TABLE positions DROP CONSTRAINT IF EXISTS fk_positions_company")
		
		// Debug: Print all company_id values in positions table
		var companyIDs []string
		gdb.Raw("SELECT DISTINCT company_id FROM positions").Scan(&companyIDs)
		fmt.Printf("Company IDs in positions table: %v\n", companyIDs)
		
		// Debug: Print all id values in companies table
		var companyIDsInCompanies []string
		gdb.Raw("SELECT id FROM companies").Scan(&companyIDsInCompanies)
		fmt.Printf("Company IDs in companies table: %v\n", companyIDsInCompanies)
		
		// Handle case where positions have company_id that doesn't exist in companies table
		if len(companyIDs) > 0 && len(companyIDsInCompanies) == 0 {
			// First, create a default tenant if it doesn't exist
			var tenantCount int
			gdb.Raw("SELECT COUNT(*) FROM tenants").Scan(&tenantCount)
			
			var defaultTenantID uuid.UUID
			if tenantCount == 0 {
				defaultTenantID = uuid.New()
				gdb.Exec("INSERT INTO tenants (id, name, email, company_name, domain, created_at, updated_at) VALUES (?, ?, ?, ?, ?, NOW(), NOW())", 
					defaultTenantID, "Default Tenant", "default@example.com", "Default Company", "example.com")
				fmt.Printf("Created default tenant with ID: %s\n", defaultTenantID.String())
			} else {
				// Get the first tenant ID
				gdb.Raw("SELECT id FROM tenants LIMIT 1").Scan(&defaultTenantID)
			}
			
			// Create a default company
			defaultCompanyID := uuid.New()
			gdb.Exec("INSERT INTO companies (id, tenant_id, name, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW())", 
				defaultCompanyID, defaultTenantID, "Default Company")
			
			// Update all positions with the new company_id
			gdb.Exec("UPDATE positions SET company_id = ?", defaultCompanyID)
			
			fmt.Printf("Created default company with ID: %s and updated positions\n", defaultCompanyID.String())
		}
		
		// Check if company_id column exists
		var columnExists bool
		gdb.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'positions' AND column_name = 'company_id')").Scan(&columnExists)
		
		if !columnExists {
			// First ensure companies table exists
			if !gdb.Migrator().HasTable("companies") {
				// Create companies table if it doesn't exist
				gdb.Exec(`CREATE TABLE IF NOT EXISTS companies (
					id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
					name VARCHAR(255) NOT NULL,
					created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP WITH TIME ZONE
				)`)
			}
			
			// Add company_id column as nullable first
			gdb.Exec("ALTER TABLE positions ADD COLUMN company_id uuid")
			
			// Get a default company ID to use
			var defaultCompanyID uuid.UUID
			if err := gdb.Raw("SELECT id FROM companies LIMIT 1").Scan(&defaultCompanyID).Error; err != nil {
				// If no companies exist, create a default one
				defaultCompanyID = uuid.New()
				gdb.Exec("INSERT INTO companies (id, name, created_at, updated_at) VALUES (?, ?, NOW(), NOW())", defaultCompanyID, "Default Company")
				// Verify the company was created
				var verifyCount int64
				gdb.Raw("SELECT COUNT(*) FROM companies WHERE id = ?", defaultCompanyID).Scan(&verifyCount)
				if verifyCount == 0 {
					return fmt.Errorf("failed to create default company")
				}
			} else {
				// Verify the company exists
				var verifyCount int64
				gdb.Raw("SELECT COUNT(*) FROM companies WHERE id = ?", defaultCompanyID).Scan(&verifyCount)
				if verifyCount == 0 {
					return fmt.Errorf("company ID not found in companies table")
				}
			}
			
			// Check if there are any existing positions with company_id values
			var existingPositionsCount int64
			gdb.Raw("SELECT COUNT(*) FROM positions WHERE company_id IS NOT NULL").Scan(&existingPositionsCount)
			
			if existingPositionsCount > 0 {
				// If there are existing positions with company_id, check if they're valid
				var invalidPositionsCount int64
				gdb.Raw("SELECT COUNT(*) FROM positions p LEFT JOIN companies c ON p.company_id = c.id WHERE p.company_id IS NOT NULL AND c.id IS NULL").Scan(&invalidPositionsCount)
				
				if invalidPositionsCount > 0 {
					// Update invalid company_id values
					gdb.Exec("UPDATE positions SET company_id = ? WHERE company_id IS NOT NULL AND company_id NOT IN (SELECT id FROM companies)", defaultCompanyID)
				}
			} else {
				// If no positions have company_id, update all of them
				gdb.Exec("UPDATE positions SET company_id = ?", defaultCompanyID)
			}
			
			// Debug: Check what company_id values exist in positions
			var positionCompanyIDs []uuid.UUID
			gdb.Raw("SELECT DISTINCT company_id FROM positions WHERE company_id IS NOT NULL").Scan(&positionCompanyIDs)
			fmt.Printf("Position company IDs: %v\n", positionCompanyIDs)
			
			// Debug: Check what companies exist
			var companyIDs []uuid.UUID
			gdb.Raw("SELECT id FROM companies").Scan(&companyIDs)
			fmt.Printf("Company IDs: %v\n", companyIDs)
			
			// Now alter the column to NOT NULL
			gdb.Exec("ALTER TABLE positions ALTER COLUMN company_id SET NOT NULL")
			
			// Let GORM handle the foreign key constraint through AutoMigrate
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
			
			// Create a default company if it doesn't exist
			var companyCount int64
			gdb.Model(&model.Company{}).Count(&companyCount)
			if companyCount == 0 {
				// Create a default tenant if it doesn't exist
				var tenantCount int64
				gdb.Model(&model.Tenant{}).Count(&tenantCount)
				if tenantCount == 0 {
					defaultTenant := model.Tenant{
						Name:        "Default Tenant",
						Email:       "default@example.com",
						CompanyName: "Default Company",
						Domain:      "default.example.com",
					}
					if err := gdb.Create(&defaultTenant).Error; err != nil {
						fmt.Printf("Failed to create default tenant: %v\n", err)
					} else {
						fmt.Printf("Created default tenant with ID: %s\n", defaultTenant.ID)
					}
				}
				
				// Create a default company
				defaultCompany := model.Company{
					Name: "Default Company",
				}
				if err := gdb.Create(&defaultCompany).Error; err != nil {
					fmt.Printf("Failed to create default company: %v\n", err)
				} else {
					fmt.Printf("Created default company with ID: %s\n", defaultCompany.ID)
				}
			}
		}
	}
	
	// Handle employees table type conversion from BIGINT to UUID
	if gdb.Migrator().HasTable("employees") {
		// Check if the employees table uses BIGINT for id
		var isBigIntID bool
		gdb.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'employees' AND column_name = 'id' AND data_type = 'bigint')").Scan(&isBigIntID)

		if isBigIntID {
			// Drop the existing employees table to avoid conflicts
			gdb.Exec("DROP TABLE IF EXISTS employees")
			fmt.Printf("Dropped existing employees table with BIGINT columns\n")
		}
	}
	
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
