-- Create attendances table
CREATE TABLE IF NOT EXISTS attendances (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    check_in TIME,
    check_out TIME,
    status VARCHAR(20) NOT NULL CHECK (status IN ('present', 'absent', 'late', 'half_day', 'leave')),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(employee_id, date)
);

-- Create index on employee_id
CREATE INDEX IF NOT EXISTS idx_attendances_employee_id ON attendances(employee_id);

-- Create index on date
CREATE INDEX IF NOT EXISTS idx_attendances_date ON attendances(date);

-- Create index on status
CREATE INDEX IF NOT EXISTS idx_attendances_status ON attendances(status);

-- Create composite index on employee_id and date for faster queries
CREATE INDEX IF NOT EXISTS idx_attendances_employee_date ON attendances(employee_id, date);

-- Create index on deleted_at for soft delete
CREATE INDEX IF NOT EXISTS idx_attendances_deleted_at ON attendances(deleted_at);

-- Create trigger to update updated_at column
CREATE OR REPLACE FUNCTION update_attendances_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER attendances_updated_at
    BEFORE UPDATE ON attendances
    FOR EACH ROW
    EXECUTE FUNCTION update_attendances_updated_at();