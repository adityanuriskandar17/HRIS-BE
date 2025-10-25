# Go HRIS MVP


## Quick Start
```bash
cp .env .env.local || true # opsional
make tidy
make up # start db & api
# or: go run ./cmd/api