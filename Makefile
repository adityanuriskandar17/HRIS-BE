.PHONY: run tidy up down fmt lint test install-hooks dev

run:
	go run ./cmd/api

dev:
	export PATH=$$PATH:$$(go env GOPATH)/bin && air

tidy:
	go mod tidy

fmt:
	gofmt -s -w .

up:
	docker compose up -d --build

down:
	docker compose down -v

test:
	go test ./...

install-hooks:
	chmod +x scripts/install-hooks.sh
	./scripts/install-hooks.sh
