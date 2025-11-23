# Commit Message Guide

## Format

```
type(scope): description
```

**Format wajib**: `type: description`  
**Format dengan scope (opsional)**: `type(scope): description`

---

## Allowed Types

| Type | Kapan Digunakan | Contoh |
|------|----------------|--------|
| `feat` | Menambah fitur baru | `feat: add invoice payment endpoint` |
| `fix` | Memperbaiki bug | `fix: resolve CORS configuration issue` |
| `docs` | Update dokumentasi | `docs: update API documentation` |
| `chore` | Maintenance, config, Makefile | `chore: update Makefile` |
| `refactor` | Refactoring code tanpa perubahan fungsi | `refactor: improve error handling` |
| `style` | Formatting, styling code | `style: format code with gofmt` |
| `perf` | Performance improvement | `perf: optimize database query` |
| `test` | Menambah atau update test | `test: add unit test for auth handler` |
| `ci` | CI/CD configuration changes | `ci: update GitHub Actions workflow` |
| `ticket` | Ticket-related changes | `ticket: implement user story #123` |

---

## Quick Reference

### ✅ Contoh yang BENAR

```bash
# Feature baru
feat: add CORS configuration
feat(auth): add refresh token endpoint
feat(invoice): add payment processing

# Bug fix
fix: resolve authentication issue
fix(api): handle null pointer exception
fix(cors): allow credentials in production

# Documentation
docs: update API documentation
docs: add CORS configuration guide

# Chore (config, Makefile, dependencies)
chore: update Makefile
chore: update dependencies
chore: add commit helper commands

# Refactor
refactor: improve error handling
refactor(auth): simplify token validation

# Style
style: format code with gofmt
style: fix linting errors

# Performance
perf: optimize database query
perf(api): reduce response time

# Test
test: add unit test for auth handler
test: increase test coverage

# CI/CD
ci: update GitHub Actions
ci: add pre-commit hooks

# Ticket
ticket: implement user story #123
ticket: fix bug reported in #456
```

### ❌ Contoh yang SALAH

```bash
# ❌ Tidak ada type
update Makefile
edit config

# ❌ Tidak ada colon
feat add new feature
fix resolve bug

# ❌ Type tidak valid
update: edit Makefile
change: modify config
edit: update file

# ❌ Format tidak lengkap
feat
fix bug
```

---

## Scope (Optional)

Scope membantu menjelaskan area kode yang diubah. Scope **opsional** tapi sangat membantu.

### Contoh dengan Scope

```bash
# Auth related
feat(auth): add refresh token endpoint
fix(auth): resolve token expiration issue

# API related
feat(api): add new endpoint
fix(api): handle error response

# Database related
feat(db): add new migration
fix(db): resolve connection pool issue

# Config related
chore(config): update CORS settings
chore(config): add environment variables

# Documentation
docs(api): update endpoint documentation
docs(cors): add CORS configuration guide
```

---

## Cara Menggunakan dengan Makefile

### Commit dan Push

```bash
# Feature
make push MSG='feat: add CORS configuration'

# Bug fix
make push MSG='fix: resolve authentication issue'

# Documentation
make push MSG='docs: update API documentation'

# Chore
make push MSG='chore: update Makefile'

# Dengan scope
make push MSG='feat(auth): add refresh token endpoint'
make push MSG='fix(api): handle null pointer exception'
```

### Commit saja (tanpa push)

```bash
make commit-msg MSG='feat: add new feature'
```

---

## Best Practices

1. **Gunakan present tense**: "add" bukan "added", "fix" bukan "fixed"
2. **Jangan kapitalisasi**: "add" bukan "Add"
3. **Jangan pakai period di akhir**: "add feature" bukan "add feature."
4. **Singkat tapi jelas**: Maksimal 72 karakter untuk description
5. **Gunakan scope jika membantu**: `feat(auth):` lebih jelas dari `feat:`

---

## Contoh Real-World

### Development Workflow

```bash
# 1. Setelah membuat perubahan
make prepare

# 2. Commit dengan format yang benar
make commit-msg MSG='feat: add invoice payment endpoint'

# 3. Atau langsung commit + push
make push MSG='feat: add invoice payment endpoint'
```

### Multiple Changes

Jika ada banyak perubahan, buat commit terpisah:

```bash
# Commit 1: Feature
make push MSG='feat: add CORS configuration'

# Commit 2: Documentation
make push MSG='docs: add CORS configuration guide'

# Commit 3: Chore
make push MSG='chore: update Makefile with commit helpers'
```

---

## Troubleshooting

### Error: "Commit message must follow the conventional commit format"

**Penyebab**: Format commit message tidak sesuai

**Solusi**: Pastikan format: `type: description`

```bash
# ❌ Salah
make push MSG='update Makefile'

# ✅ Benar
make push MSG='chore: update Makefile'
```

### Error: "Commit type 'xxx' is not allowed"

**Penyebab**: Type yang digunakan tidak ada di allowed types

**Solusi**: Gunakan type yang valid (feat, fix, docs, chore, dll)

```bash
# ❌ Salah
make push MSG='update: edit Makefile'  # 'update' bukan allowed type

# ✅ Benar
make push MSG='chore: update Makefile'
```

---

## Quick Reference Card

```
┌─────────────────────────────────────────────────┐
│  COMMIT MESSAGE FORMAT                          │
├─────────────────────────────────────────────────┤
│  type(scope): description                       │
│                                                 │
│  Types:                                         │
│  • feat    - New feature                        │
│  • fix     - Bug fix                            │
│  • docs    - Documentation                      │
│  • chore   - Maintenance/config                 │
│  • refactor - Code refactoring                  │
│  • style   - Formatting                         │
│  • perf    - Performance                        │
│  • test    - Testing                            │
│  • ci      - CI/CD                              │
│  • ticket  - Ticket-related                     │
│                                                 │
│  Examples:                                      │
│  make push MSG='feat: add new feature'          │
│  make push MSG='fix: resolve bug'               │
│  make push MSG='chore: update Makefile'         │
└─────────────────────────────────────────────────┘
```

---

**Last Updated**: 2024

