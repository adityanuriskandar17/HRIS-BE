-- Create payroll_lines table
CREATE TABLE IF NOT EXISTS payroll_lines (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    payroll_id UUID NOT NULL REFERENCES payrolls(id) ON DELETE CASCADE,
    description VARCHAR(255) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('earning', 'deduction', 'tax', 'contribution')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create index on payroll_id
CREATE INDEX IF NOT EXISTS idx_payroll_lines_payroll_id ON payroll_lines(payroll_id);

-- Create index on type
CREATE INDEX IF NOT EXISTS idx_payroll_lines_type ON payroll_lines(type);

-- Create index on deleted_at for soft delete
CREATE INDEX IF NOT EXISTS idx_payroll_lines_deleted_at ON payroll_lines(deleted_at);

-- Create trigger to update updated_at column
CREATE OR REPLACE FUNCTION update_payroll_lines_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER payroll_lines_updated_at
    BEFORE UPDATE ON payroll_lines
    FOR EACH ROW
    EXECUTE FUNCTION update_payroll_lines_updated_at();