# HRIS API – Reference

> Semua endpoint merespons dalam format amplop:
>
> ```json
> {
>   "success": true,
>   "data": {},
>   "error": null,
>   "meta": { "trace_id": "…" }
> }
> ```
>
> Saat terjadi error, `success` bernilai `false` dan `error` berisi detail:
>
> ```json
> {
>   "success": false,
>   "error": {
>     "code": "VALIDATION_ERROR",
>     "message": "Email required",
>     "fields": { "email": "required" }
>   },
>   "meta": { "trace_id": "…" }
> }
> ```

## Base URL

- Lokal (make run): `http://localhost:8081/api/v1`
- Docker compose: `http://localhost:8080/api/v1`

Tambahkan header `Authorization: Bearer <access_token>` untuk semua endpoint protected.

---

## Health Check

```
GET /healthz
```

**Response 200**

```json
{
  "success": true,
  "data": { "status": "ok" },
  "error": null,
  "meta": { "trace_id": "…" }
}
```

---

## Auth

### Login

```
POST /auth/login
Content-Type: application/json
```

**Body**

```json
{
  "email": "admin@example.com",
  "password": "admin12345"
}
```

**Response 200**

```json
{
  "success": true,
  "data": {
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9…",
    "refreshToken": "vQX2e1R…",
    "expiresIn": 900,
    "refreshExpiresIn": 604800,
    "role": "ADMIN"
  },
  "error": null,
  "meta": { "trace_id": "…" }
}
```

### Refresh Token

```
POST /auth/refresh
Content-Type: application/json
```

**Body**

```json
{ "refreshToken": "vQX2e1R…" }
```

**Response 200** – sama seperti login namun dengan pasangan token baru.

**Kemungkinan error**

| HTTP | code                    | Keterangan                        |
|------|-------------------------|-----------------------------------|
| 400  | `INVALID_JSON`          | Payload tidak bisa di-decode      |
| 401  | `INVALID_REFRESH_TOKEN` | Refresh token hilang/invalid      |
| 500  | `INTERNAL_SERVER_ERROR` | Gagal membuat pasangan token baru |

---

## Master Data

Role minimum ditandai untuk setiap endpoint.

### List Units (ADMIN/HR)

```
GET /master/units
```

**Response 200**

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "code": "HRD",
      "name": "Human Resources",
      "createdById": null,
      "updatedById": null,
      "createdAt": "2024-10-26T15:00:00Z",
      "updatedAt": "2024-10-26T15:00:00Z"
    }
  ],
  "error": null,
  "meta": { "trace_id": "…" }
}
```

### Create Unit (ADMIN/HR)

```
POST /master/units
Content-Type: application/json
```

**Body**

```json
{
  "code": "UI/UX",
  "name": "Designer UI/UX"
}
```

**Response 201** – data unit baru.
```json
{
    "success": true,
    "data": {
        "id": 4,
        "code": "UI/UX",
        "name": "Designer UI/UX",
        "createdById": 1,
        "updatedById": 1,
        "createdAt": "2025-10-26T17:41:38.78810323+07:00",
        "updatedAt": "2025-10-26T17:41:38.78810323+07:00"
    },
    "meta": {
        "trace_id": "dc6714eee904a4de4f0c790f7c3c02f5"
    }
}
```

**Error umum**: `VALIDATION_ERROR` (field kosong), `DATABASE_ERROR`.

---

### List Positions (ADMIN/HR)

```
GET /master/positions
```

**Response 200** – daftar posisi dengan nested unit (jika ada).

### Create Position (ADMIN/HR)

```
POST /master/positions
Content-Type: application/json
```

**Body**

```json
{
  "title": "Finance Analyst",
  "unitId": 2
}
```

`unitId` opsional; jika diisi harus merujuk ke unit valid.

---

### List Employees (ADMIN/HR/MANAGER)

```
GET /master/employees
```

Mengembalikan array employee beserta relasi unit dan posisi.

### Create Employee (ADMIN/HR/MANAGER)

```
POST /master/employees
Content-Type: application/json
```

**Body contoh**

```json
{
  "employeeCode": "EMP999",
  "fullName": "John Doe",
  "email": "john.doe@example.com",
  "phone": "+628123456789",
  "unitId": 1,
  "positionId": 3,
  "employmentStatus": "FULLTIME",
  "startDate": "2024-01-01",
  "endDate": "",
  "dateOfBirth": "1990-05-15"
}
```

**Validasi penting**

| Field              | Catatan                                                 |
|--------------------|---------------------------------------------------------|
| `employeeCode`     | wajib, akan di-uppercase                                |
| `email`            | wajib, disimpan lowercase                               |
| `employmentStatus` | salah satu dari `FULLTIME`, `PARTTIME`, `CONTRACT`, `INTERN` |
| `startDate`        | format `YYYY-MM-DD`, default hari ini jika kosong       |
| `endDate`, `dateOfBirth` | opsional, format `YYYY-MM-DD`                     |

Error akan menggunakan `VALIDATION_ERROR` dengan `fields` yang menunjukkan kolom bermasalah.

---

## Error Codes Ringkas

| Code                    | HTTP | Deskripsi                                                |
|-------------------------|------|----------------------------------------------------------|
| `INVALID_JSON`          | 400  | Body request bukan JSON valid                            |
| `VALIDATION_ERROR`      | 400  | Field wajib kosong/format salah                          |
| `INVALID_CREDENTIALS`   | 401  | Login gagal (email / password)                           |
| `INVALID_REFRESH_TOKEN` | 401  | Refresh token tidak valid / kedaluwarsa / dicabut        |
| `UNAUTHENTICATED`       | 401  | Header Authorization hilang atau token salah            |
| `FORBIDDEN`             | 403  | Role pengguna tidak berhak                               |
| `ACCOUNT_DISABLED`      | 403  | Akun user non-aktif                                      |
| `DATABASE_ERROR`        | 400  | Gagal saat proses penyimpanan                            |
| `INTERNAL_SERVER_ERROR` | 500  | Kesalahan tak terduga di server                          |

---

## Postman Tips

1. Buat environment variable `base_url` (`http://localhost:8081/api/v1`) dan `access_token`.
2. Step login, simpan `data.accessToken` ke `access_token`, `data.refreshToken` ke `refresh_token`.
3. Tambahkan header `Authorization: Bearer {{access_token}}` pada request protected.
4. Jika menerima 401 karena akses kedaluwarsa, hit `/auth/refresh` lalu ulangi request.

