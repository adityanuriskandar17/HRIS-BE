-- Create leave_requests table
CREATE TABLE IF NOT EXISTS leave_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    leave_type VARCHAR(20) NOT NULL CHECK (leave_type IN ('annual', 'sick', 'personal', 'maternity', 'paternity', 'unpaid')),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    reason TEXT,
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'approved', 'rejected', 'cancelled')),
    approved_by UUID REFERENCES user_accounts(id),
    approved_at TIMESTAMP WITH TIME ZONE,
    approver_notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create index on employee_id
CREATE INDEX IF NOT EXISTS idx_leave_requests_employee_id ON leave_requests(employee_id);

-- Create index on leave_type
CREATE INDEX IF NOT EXISTS idx_leave_requests_leave_type ON leave_requests(leave_type);

-- Create index on status
CREATE INDEX IF NOT EXISTS idx_leave_requests_status ON leave_requests(status);

-- Create index on start_date
CREATE INDEX IF NOT EXISTS idx_leave_requests_start_date ON leave_requests(start_date);

-- Create index on approved_by
CREATE INDEX IF NOT EXISTS idx_leave_requests_approved_by ON leave_requests(approved_by);

-- Create index on deleted_at for soft delete
CREATE INDEX IF NOT EXISTS idx_leave_requests_deleted_at ON leave_requests(deleted_at);

-- Create trigger to update updated_at column
CREATE OR REPLACE FUNCTION update_leave_requests_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER leave_requests_updated_at
    BEFORE UPDATE ON leave_requests
    FOR EACH ROW
    EXECUTE FUNCTION update_leave_requests_updated_at();