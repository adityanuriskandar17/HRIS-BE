# Repository Guidelines

## Project Structure & Module Organization
- `cmd/api` contains the Go `main` package; keep only bootstrap wiring and HTTP server startup here.
- Core logic lives in `internal`: `auth` for authentication helpers, `config` for environment loading, `db` for data access, `domain` for entities/value objects, `http` for transport adapters, and `util` for cross-cutting helpers.
- Database assets belong in `migrations`; name new files with incremental prefixes (`001_init.sql`) to keep ordering consistent.
- Container definitions (`docker-compose.yml`, `DockerFile`) provision the API and supporting services; update them when new dependencies are introduced.

## Build, Test, and Development Commands
- `make run` wraps `go run ./cmd/api` for local iteration.
- `make tidy` syncs `go.mod` / `go.sum`; run after adding or upgrading dependencies.
- `make fmt` enforces `gofmt -s -w .`; always run before opening a PR.
- `make up` / `make down` start and stop the dockerized stack; useful for integration checks.
- `make test` (alias for `go test ./...`) executes all package tests.

## Coding Style & Naming Conventions
- Follow idiomatic Go: tabs for indentation, CamelCase for exported identifiers, lower-case package names.
- Keep HTTP request/response shapes in `internal/http` and suffix structs with `Request` or `Response` for clarity.
- Configuration keys in `.env` / `.env.local` should use upper snake case (`DB_DSN`, `JWT_SECRET`); document new keys in the PR description.
- Prefer returning explicit errors wrapped with context using `fmt.Errorf("...: %w", err)`.

## Testing Guidelines
- Place `_test.go` files alongside the code under test and use table-driven tests for coverage of edge cases.
- Run `make up` to provision the database before tests that require persistence; reset state between tests with migrations or fixtures.
- Stub external integrations behind interfaces in `internal/util` to keep tests deterministic.
- Capture both success and failure paths, and note any gaps in the PR if coverage is not feasible.

## Commit & Pull Request Guidelines
- Write imperative, scoped commit subjects (e.g., `Add employee seed migration`), with optional bodies describing rationale or follow-up work.
- Squash WIP commits before merge to keep history clean.
- PRs should describe the change, list test commands run, call out schema or config impacts, and link related issues or tickets.
- Include screenshots or sample payloads when modifying HTTP handlers to help reviewers validate behavior.
