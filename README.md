# Go HRIS MVP


## Quick Start
```bash
cp .env .env.local || true # opsional
make tidy
make up # start db & api
# or: go run ./cmd/api
```

## Observability

- Tracing diaktifkan via OpenTelemetry; set `OTEL_EXPORTER_JAEGER_ENDPOINT` (contoh: `http://jaeger:14268/api/traces`) dan `OTEL_SERVICE_NAME` sesuai kebutuhan.
- Saat menggunakan `make up`, Jaeger UI tersedia di http://localhost:16686 untuk menelusuri span yang direkam.
