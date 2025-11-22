# Go HRIS MVP


## Quick Start
```bash
cp .env .env.local || true # opsional
make tidy
make up # start db & api
# or: go run ./cmd/api
```

## Development Mode (Auto-Reload)

Untuk development dengan auto-reload dan auto-regenerate Swagger docs:

```bash
make dev
```

Ini akan:
- ✅ Auto-regenerate Swagger docs ketika ada perubahan di handler/DTO
- ✅ Auto-reload server ketika ada perubahan kode
- ✅ Watch file `.go` dan rebuild otomatis

**Catatan:** Setelah perubahan kode, Swagger docs akan ter-update otomatis. Anda hanya perlu **refresh browser** di halaman Swagger UI untuk melihat perubahan.

**Alternatif tanpa auto-reload:**
```bash
make run  # Manual run, perlu restart manual
```

## Git Hooks Setup

This project includes Git hooks to ensure code quality and consistency:

1. **commit-msg**: Validates commit messages follow conventional commit format
2. **pre-commit**: Automatically runs `go mod tidy` before each commit

The hooks are located in the `.husky` directory and are installed to `.git/hooks` using the `install-hooks.sh` script.

### Commit Message Format

The commit-msg hook enforces the conventional commit format:
```
type(scope): description
```

Allowed types:
- `ci`: Changes to CI configuration files and scripts
- `chore`: Maintenance tasks, dependency updates, etc.
- `docs`: Documentation only changes
- `ticket`: Ticket-related changes
- `feat`: New features
- `fix`: Bug fixes
- `perf`: Performance improvements
- `refactor`: Code refactoring without functional changes
- `revert`: Reverting previous changes
- `style`: Code style changes (formatting, missing semicolons, etc.)

### Installation

To install these hooks after cloning the repository, you can use either method:

**Method 1: Using make (recommended)**
```bash
make install-hooks
```

**Method 2: Direct script execution**
```bash
chmod +x scripts/install-hooks.sh
./scripts/install-hooks.sh
```

The hooks will be automatically activated for all subsequent commits.

## Observability

- Tracing diaktifkan via OpenTelemetry; set `OTEL_EXPORTER_JAEGER_ENDPOINT` (contoh: `http://jaeger:14268/api/traces`) dan `OTEL_SERVICE_NAME` sesuai kebutuhan.
- Saat menggunakan `make up`, Jaeger UI tersedia di http://localhost:16686 untuk menelusuri span yang direkam.
