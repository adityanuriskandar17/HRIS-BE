# CORS Configuration Guide

## Overview

CORS (Cross-Origin Resource Sharing) memungkinkan aplikasi web di domain yang berbeda untuk mengakses API. Konfigurasi CORS di HRIS API dapat disesuaikan melalui environment variables.

## Environment Variables

### CORS_ALLOWED_ORIGINS

**Format**: Comma-separated list of origins

**Default**: Jika tidak diisi, akan menggunakan default development origins:
- `http://localhost:3000` (Create React App default)
- `http://localhost:3001`
- `http://localhost:5173` (Vite default)
- `http://localhost:8080`
- `http://localhost:8081`
- `http://127.0.0.1:3000`
- `http://127.0.0.1:3001`
- `http://127.0.0.1:5173`

**Contoh**:
```bash
# Development
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173

# Production
CORS_ALLOWED_ORIGINS=https://app.yourdomain.com,https://admin.yourdomain.com

# Multiple environments
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://staging.yourdomain.com,https://app.yourdomain.com
```

**Catatan**: 
- Harus include protocol (`http://` atau `https://`)
- Harus include port jika bukan default (80 untuk http, 443 untuk https)
- Tidak boleh ada trailing slash

---

### CORS_ALLOWED_METHODS

**Format**: Comma-separated list of HTTP methods

**Default**: `GET,POST,PUT,PATCH,DELETE,OPTIONS`

**Contoh**:
```bash
CORS_ALLOWED_METHODS=GET,POST,PUT,PATCH,DELETE,OPTIONS
```

**Catatan**: `OPTIONS` harus selalu disertakan untuk preflight requests.

---

### CORS_ALLOWED_HEADERS

**Format**: Comma-separated list of header names

**Default**: `Accept,Authorization,Content-Type,X-CSRF-Token`

**Contoh**:
```bash
CORS_ALLOWED_HEADERS=Accept,Authorization,Content-Type,X-CSRF-Token,X-Requested-With
```

**Header yang umum digunakan**:
- `Accept`: Content type yang diterima
- `Authorization`: JWT token untuk authentication
- `Content-Type`: Request body content type
- `X-CSRF-Token`: CSRF protection token
- `X-Requested-With`: Identifikasi AJAX requests

---

### CORS_EXPOSED_HEADERS

**Format**: Comma-separated list of header names

**Default**: `Link`

**Contoh**:
```bash
CORS_EXPOSED_HEADERS=Link,X-Total-Count,X-Page-Count
```

**Catatan**: Headers ini akan bisa diakses oleh JavaScript di browser.

---

### CORS_ALLOW_CREDENTIALS

**Format**: `true` atau `false`

**Default**: `true`

**Contoh**:
```bash
CORS_ALLOW_CREDENTIALS=true
```

**Catatan**: 
- Harus `true` jika menggunakan cookies atau authentication headers
- Jika `true`, `CORS_ALLOWED_ORIGINS` tidak bisa menggunakan wildcard (`*`)

---

### CORS_MAX_AGE

**Format**: Integer (seconds)

**Default**: `300` (5 minutes)

**Contoh**:
```bash
CORS_MAX_AGE=300
```

**Catatan**: Durasi cache untuk preflight requests. Nilai lebih besar mengurangi jumlah preflight requests.

---

## Konfigurasi per Environment

### Development

```bash
# .env.local atau .env.development
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173,http://localhost:8080
CORS_ALLOW_CREDENTIALS=true
CORS_MAX_AGE=300
```

### Staging

```bash
# .env.staging
CORS_ALLOWED_ORIGINS=https://staging.yourdomain.com,https://admin-staging.yourdomain.com
CORS_ALLOW_CREDENTIALS=true
CORS_MAX_AGE=600
```

### Production

```bash
# .env.production
CORS_ALLOWED_ORIGINS=https://app.yourdomain.com,https://admin.yourdomain.com
CORS_ALLOW_CREDENTIALS=true
CORS_MAX_AGE=3600
```

---

## Contoh Konfigurasi Lengkap

### Minimal Configuration (Development)

```bash
# .env
DATABASE_URL=postgres://hris:hris@127.0.0.1:5432/hris?sslmode=disable
JWT_SECRET=your-secret-key
PORT=8081

# CORS - menggunakan defaults
# Tidak perlu set CORS_ALLOWED_ORIGINS, akan menggunakan default development origins
```

### Full Configuration (Production)

```bash
# .env
DATABASE_URL=postgres://user:pass@db.example.com:5432/hris?sslmode=require
JWT_SECRET=super-secret-production-key
PORT=8080

# CORS Configuration
CORS_ALLOWED_ORIGINS=https://app.yourdomain.com,https://admin.yourdomain.com
CORS_ALLOWED_METHODS=GET,POST,PUT,PATCH,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Accept,Authorization,Content-Type,X-CSRF-Token
CORS_EXPOSED_HEADERS=Link,X-Total-Count
CORS_ALLOW_CREDENTIALS=true
CORS_MAX_AGE=3600
```

---

## Testing CORS Configuration

### Test dengan curl

```bash
# Test preflight request
curl -X OPTIONS http://localhost:8081/api/v1/auth/login \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type" \
  -v

# Expected response headers:
# Access-Control-Allow-Origin: http://localhost:3000
# Access-Control-Allow-Methods: GET,POST,PUT,PATCH,DELETE,OPTIONS
# Access-Control-Allow-Headers: Accept,Authorization,Content-Type,X-CSRF-Token
# Access-Control-Allow-Credentials: true
# Access-Control-Max-Age: 300
```

### Test dari Browser Console

```javascript
// Test dari browser console di http://localhost:3000
fetch('http://localhost:8081/api/v1/auth/login', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    email: 'test@example.com',
    password: 'test123'
  })
})
.then(response => {
  console.log('CORS Headers:', {
    'Access-Control-Allow-Origin': response.headers.get('Access-Control-Allow-Origin'),
    'Access-Control-Allow-Credentials': response.headers.get('Access-Control-Allow-Credentials')
  });
  return response.json();
})
.then(data => console.log('Response:', data))
.catch(error => console.error('CORS Error:', error));
```

---

## Troubleshooting

### Error: "Access to fetch has been blocked by CORS policy"

**Penyebab**: Origin tidak ada di `CORS_ALLOWED_ORIGINS`

**Solusi**:
1. Pastikan origin aplikasi React.js ada di `CORS_ALLOWED_ORIGINS`
2. Restart server setelah mengubah environment variables
3. Pastikan format origin benar (dengan protocol dan port)

**Contoh**:
```bash
# ❌ Salah
CORS_ALLOWED_ORIGINS=localhost:3000

# ✅ Benar
CORS_ALLOWED_ORIGINS=http://localhost:3000
```

---

### Error: "Credentials flag is true, but Access-Control-Allow-Origin is *"

**Penyebab**: Menggunakan wildcard (`*`) dengan `CORS_ALLOW_CREDENTIALS=true`

**Solusi**: Gunakan specific origins, bukan wildcard

```bash
# ❌ Salah
CORS_ALLOWED_ORIGINS=*
CORS_ALLOW_CREDENTIALS=true

# ✅ Benar
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://app.yourdomain.com
CORS_ALLOW_CREDENTIALS=true
```

---

### Preflight requests terlalu sering

**Penyebab**: `CORS_MAX_AGE` terlalu kecil

**Solusi**: Tingkatkan nilai `CORS_MAX_AGE`

```bash
# Default: 300 seconds (5 minutes)
CORS_MAX_AGE=300

# Untuk production, bisa lebih besar
CORS_MAX_AGE=3600  # 1 hour
```

---

### Custom headers tidak dikirim

**Penyebab**: Header tidak ada di `CORS_ALLOWED_HEADERS`

**Solusi**: Tambahkan header ke `CORS_ALLOWED_HEADERS`

```bash
# Jika menggunakan custom header X-Custom-Header
CORS_ALLOWED_HEADERS=Accept,Authorization,Content-Type,X-CSRF-Token,X-Custom-Header
```

---

## Best Practices

1. **Development**: Gunakan default origins atau tambahkan localhost dengan berbagai port
2. **Production**: Hanya allow origins yang benar-benar diperlukan
3. **Security**: Jangan gunakan wildcard (`*`) di production
4. **Performance**: Set `CORS_MAX_AGE` yang sesuai untuk mengurangi preflight requests
5. **Testing**: Test CORS configuration di semua environment sebelum deploy

---

## React.js Configuration

### Setup di React App

```javascript
// src/config/api.js
const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8081/api/v1';

// Pastikan origin React app sesuai dengan CORS_ALLOWED_ORIGINS
// Jika React app di http://localhost:3000, pastikan ada di CORS_ALLOWED_ORIGINS
```

### .env di React App

```bash
# .env.development
REACT_APP_API_URL=http://localhost:8081/api/v1

# .env.production
REACT_APP_API_URL=https://api.yourdomain.com/api/v1
```

---

## Security Considerations

1. **Never use wildcard in production**: Selalu specify exact origins
2. **Limit allowed methods**: Hanya allow methods yang benar-benar digunakan
3. **Limit allowed headers**: Hanya allow headers yang diperlukan
4. **Use HTTPS in production**: Selalu gunakan `https://` untuk production origins
5. **Regular review**: Review CORS configuration secara berkala

---

## Quick Reference

| Variable | Default | Required | Example |
|----------|---------|----------|---------|
| `CORS_ALLOWED_ORIGINS` | Development origins | No | `http://localhost:3000,https://app.com` |
| `CORS_ALLOWED_METHODS` | `GET,POST,PUT,PATCH,DELETE,OPTIONS` | No | `GET,POST,PUT,DELETE,OPTIONS` |
| `CORS_ALLOWED_HEADERS` | `Accept,Authorization,Content-Type,X-CSRF-Token` | No | `Accept,Authorization,Content-Type` |
| `CORS_EXPOSED_HEADERS` | `Link` | No | `Link,X-Total-Count` |
| `CORS_ALLOW_CREDENTIALS` | `true` | No | `true` |
| `CORS_MAX_AGE` | `300` | No | `3600` |

---

**Last Updated**: 2024

