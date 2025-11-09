# Go HRIS MVP


## Quick Start
```bash
cp .env .env.local || true # opsional
make tidy
make up # start db & api
# or: go run ./cmd/api
```

## Git Hooks Setup

This project includes Git hooks to ensure code quality and consistency:

1. **commit-msg**: Validates commit messages follow conventional commit format
2. **pre-commit**: Automatically runs `go mod tidy` before each commit

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
