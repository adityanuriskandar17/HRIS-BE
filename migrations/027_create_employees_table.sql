-- +migrate Up
CREATE TABLE IF NOT EXISTS employees (
    id BIGSERIAL PRIMARY KEY,
    employee_code VARCHAR(32) NOT NULL UNIQUE,
    full_name VARCHAR(150) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(30),
    unit_id BIGINT REFERENCES units(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    position_id BIGINT REFERENCES positions(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    employment_status employment_status_enum NOT NULL DEFAULT 'FULLTIME',
    start_date DATE NOT NULL,
    end_date DATE,
    date_of_birth DATE,
    created_by BIGINT REFERENCES user_accounts(id) ON DELETE SET NULL ON UPDATE CASCADE,
    updated_by BIGINT REFERENCES user_accounts(id) ON DELETE SET NULL ON UPDATE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create index for deleted_at for soft delete functionality
CREATE INDEX IF NOT EXISTS idx_employees_deleted_at ON employees(deleted_at);

-- +migrate Down
DROP TABLE IF EXISTS employees;