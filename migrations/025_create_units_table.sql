-- +migrate Up
CREATE TABLE IF NOT EXISTS units (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(32) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    created_by BIGINT REFERENCES user_accounts(id) ON DELETE SET NULL ON UPDATE CASCADE,
    updated_by BIGINT REFERENCES user_accounts(id) ON DELETE SET NULL ON UPDATE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create index for deleted_at for soft delete functionality
CREATE INDEX IF NOT EXISTS idx_units_deleted_at ON units(deleted_at);

-- +migrate Down
DROP TABLE IF EXISTS units;