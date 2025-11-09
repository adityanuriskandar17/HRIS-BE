-- Create payrolls table
CREATE TABLE IF NOT EXISTS payrolls (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    pay_period VARCHAR(20) NOT NULL,
    gross_salary DECIMAL(10,2) NOT NULL,
    deductions DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    net_salary DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('draft', 'calculated', 'approved', 'paid', 'cancelled')),
    payment_date DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create index on employee_id
CREATE INDEX IF NOT EXISTS idx_payrolls_employee_id ON payrolls(employee_id);

-- Create index on pay_period
CREATE INDEX IF NOT EXISTS idx_payrolls_pay_period ON payrolls(pay_period);

-- Create index on status
CREATE INDEX IF NOT EXISTS idx_payrolls_status ON payrolls(status);

-- Create index on payment_date
CREATE INDEX IF NOT EXISTS idx_payrolls_payment_date ON payrolls(payment_date);

-- Create index on deleted_at for soft delete
CREATE INDEX IF NOT EXISTS idx_payrolls_deleted_at ON payrolls(deleted_at);

-- Create trigger to update updated_at column
CREATE OR REPLACE FUNCTION update_payrolls_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER payrolls_updated_at
    BEFORE UPDATE ON payrolls
    FOR EACH ROW
    EXECUTE FUNCTION update_payrolls_updated_at();