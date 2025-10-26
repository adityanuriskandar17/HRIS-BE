-- Create supporting extensions and indexes for faster lookups.
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_employees_unit ON public.employees (unit_id);
CREATE INDEX IF NOT EXISTS idx_employees_position ON public.employees (position_id);
CREATE INDEX IF NOT EXISTS idx_positions_unit ON public.positions (unit_id);
CREATE INDEX IF NOT EXISTS idx_employees_name_trgm ON public.employees USING gin (full_name gin_trgm_ops);
