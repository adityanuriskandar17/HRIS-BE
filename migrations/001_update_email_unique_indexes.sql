-- Normalize email uniqueness to be case-insensitive on employees and user accounts.
DROP INDEX IF EXISTS employees_email_key;
CREATE UNIQUE INDEX employees_email_key_ci ON public.employees (LOWER(email));

DROP INDEX IF EXISTS user_accounts_email_key;
CREATE UNIQUE INDEX user_accounts_email_key_ci ON public.user_accounts (LOWER(email));
