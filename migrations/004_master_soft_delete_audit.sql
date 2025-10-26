-- Soft delete and audit columns for master tables.

-- Units
ALTER TABLE public.units ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;
ALTER TABLE public.units ADD COLUMN IF NOT EXISTS created_by BIGINT;
ALTER TABLE public.units ADD COLUMN IF NOT EXISTS updated_by BIGINT;

DO $$ BEGIN
  ALTER TABLE public.units
    ADD CONSTRAINT units_created_by_fkey FOREIGN KEY (created_by)
    REFERENCES public.user_accounts(id) ON UPDATE CASCADE ON DELETE SET NULL;
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  ALTER TABLE public.units
    ADD CONSTRAINT units_updated_by_fkey FOREIGN KEY (updated_by)
    REFERENCES public.user_accounts(id) ON UPDATE CASCADE ON DELETE SET NULL;
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

-- Positions
ALTER TABLE public.positions ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;
ALTER TABLE public.positions ADD COLUMN IF NOT EXISTS created_by BIGINT;
ALTER TABLE public.positions ADD COLUMN IF NOT EXISTS updated_by BIGINT;

DO $$ BEGIN
  ALTER TABLE public.positions
    ADD CONSTRAINT positions_created_by_fkey FOREIGN KEY (created_by)
    REFERENCES public.user_accounts(id) ON UPDATE CASCADE ON DELETE SET NULL;
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  ALTER TABLE public.positions
    ADD CONSTRAINT positions_updated_by_fkey FOREIGN KEY (updated_by)
    REFERENCES public.user_accounts(id) ON UPDATE CASCADE ON DELETE SET NULL;
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

-- Employees
ALTER TABLE public.employees ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;
ALTER TABLE public.employees ADD COLUMN IF NOT EXISTS created_by BIGINT;
ALTER TABLE public.employees ADD COLUMN IF NOT EXISTS updated_by BIGINT;

DO $$ BEGIN
  ALTER TABLE public.employees
    ADD CONSTRAINT employees_created_by_fkey FOREIGN KEY (created_by)
    REFERENCES public.user_accounts(id) ON UPDATE CASCADE ON DELETE SET NULL;
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  ALTER TABLE public.employees
    ADD CONSTRAINT employees_updated_by_fkey FOREIGN KEY (updated_by)
    REFERENCES public.user_accounts(id) ON UPDATE CASCADE ON DELETE SET NULL;
EXCEPTION WHEN duplicate_object THEN NULL; END $$;
