.PHONY: run tidy up down fmt lint test install-hooks dev prepare commit push

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

# Prepare: Run all checks and formatting before commit
prepare: tidy fmt
	@echo "✅ Dependencies synced and code formatted"
	@echo "📦 Checking for changes in go.mod/go.sum..."
	@if ! git diff --quiet go.mod go.sum 2>/dev/null; then \
		echo "⚠️  Changes detected in go.mod or go.sum"; \
		echo "💡 Run: git add go.mod go.sum"; \
	fi
	@echo "✅ Ready to commit!"

# Commit helper: Prepare + show status
commit: prepare
	@echo ""
	@echo "📋 Current git status:"
	@git status --short
	@echo ""
	@echo "💡 To commit, run: git commit -m 'your message'"
	@echo "💡 Or use: make commit-msg MSG='your message'"

# Commit with message: Prepare + commit
# Format: type(scope): description
# Allowed types: ci, chore, docs, ticket, feat, fix, perf, refactor, revert, style
# Examples:
#   make commit-msg MSG='feat: add CORS configuration'
#   make commit-msg MSG='fix: resolve authentication issue'
#   make commit-msg MSG='docs: update API documentation'
commit-msg:
	@if [ -z "$(MSG)" ]; then \
		echo "❌ Error: MSG is required"; \
		echo ""; \
		echo "Usage: make commit-msg MSG='type(scope): description'"; \
		echo ""; \
		echo "Format: type(scope): description"; \
		echo "Allowed types: ci, chore, docs, ticket, feat, fix, perf, refactor, revert, style"; \
		echo ""; \
		echo "Examples:"; \
		echo "  make commit-msg MSG='feat: add CORS configuration'"; \
		echo "  make commit-msg MSG='fix: resolve authentication issue'"; \
		echo "  make commit-msg MSG='docs: update API documentation'"; \
		echo "  make commit-msg MSG='chore: update Makefile'"; \
		exit 1; \
	fi
	@$(MAKE) prepare
	@echo ""
	@echo "📝 Committing with message: $(MSG)"
	@git add -A
	@git commit -m "$(MSG)" || ( \
		echo ""; \
		echo "❌ Commit failed!"; \
		echo ""; \
		echo "💡 Commit message must follow format: type(scope): description"; \
		echo "   Allowed types: ci, chore, docs, ticket, feat, fix, perf, refactor, revert, style"; \
		echo ""; \
		echo "   Examples:"; \
		echo "     make commit-msg MSG='feat: add new feature'"; \
		echo "     make commit-msg MSG='fix: resolve bug'"; \
		echo "     make commit-msg MSG='chore: update Makefile'"; \
		exit 1; \
	)
	@echo "✅ Commit successful!"

# Push: Prepare + commit + push (requires MSG)
# Format: type(scope): description
# Allowed types: ci, chore, docs, ticket, feat, fix, perf, refactor, revert, style
# Examples:
#   make push MSG='feat: add CORS configuration'
#   make push MSG='fix: resolve authentication issue'
push:
	@if [ -z "$(MSG)" ]; then \
		echo "❌ Error: MSG is required"; \
		echo ""; \
		echo "Usage: make push MSG='type(scope): description'"; \
		echo ""; \
		echo "Format: type(scope): description"; \
		echo "Allowed types: ci, chore, docs, ticket, feat, fix, perf, refactor, revert, style"; \
		echo ""; \
		echo "Examples:"; \
		echo "  make push MSG='feat: add CORS configuration'"; \
		echo "  make push MSG='fix: resolve authentication issue'"; \
		echo "  make push MSG='chore: update Makefile'"; \
		exit 1; \
	fi
	@$(MAKE) commit-msg MSG="$(MSG)"
	@echo ""
	@echo "🚀 Pushing to remote..."
	@git push || (echo "❌ Push failed. Check the errors above." && exit 1)
	@echo "✅ Push successful!"
