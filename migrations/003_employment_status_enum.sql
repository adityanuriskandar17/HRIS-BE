DO $$ BEGIN
  CREATE TYPE employment_status_enum AS ENUM ('FULLTIME','PARTTIME','CONTRACT','INTERN');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

ALTER TABLE public.employees
  ALTER COLUMN employment_status TYPE employment_status_enum USING employment_status::employment_status_enum;
