-- Create employee_contracts table
CREATE TABLE IF NOT EXISTS employee_contracts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    contract_type VARCHAR(20) NOT NULL CHECK (contract_type IN ('permanent', 'contract', 'internship', 'probation')),
    start_date DATE NOT NULL,
    end_date DATE,
    salary DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create index on employee_id
CREATE INDEX IF NOT EXISTS idx_employee_contracts_employee_id ON employee_contracts(employee_id);

-- Create index on contract_type
CREATE INDEX IF NOT EXISTS idx_employee_contracts_contract_type ON employee_contracts(contract_type);

-- Create index on start_date
CREATE INDEX IF NOT EXISTS idx_employee_contracts_start_date ON employee_contracts(start_date);

-- Create index on end_date to track expiring contracts
CREATE INDEX IF NOT EXISTS idx_employee_contracts_end_date ON employee_contracts(end_date);

-- Create index on deleted_at for soft delete
CREATE INDEX IF NOT EXISTS idx_employee_contracts_deleted_at ON employee_contracts(deleted_at);

-- Create trigger to update updated_at column
CREATE OR REPLACE FUNCTION update_employee_contracts_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER employee_contracts_updated_at
    BEFORE UPDATE ON employee_contracts
    FOR EACH ROW
    EXECUTE FUNCTION update_employee_contracts_updated_at();