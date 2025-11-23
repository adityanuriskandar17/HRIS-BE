# HRIS API Documentation for React.js

## Table of Contents
- [Base Configuration](#base-configuration)
- [Authentication](#authentication)
- [API Endpoints](#api-endpoints)
  - [Auth](#auth)
  - [Tenants](#tenants)
  - [Subscriptions](#subscriptions)
  - [Invoices](#invoices)
  - [Master Data](#master-data)
  - [Companies](#companies)
- [Error Handling](#error-handling)
- [React.js Examples](#reactjs-examples)

---

## Base Configuration

### Base URL
```javascript
const API_BASE_URL = 'http://localhost:8081/api/v1';
```

### Headers
```javascript
const headers = {
  'Content-Type': 'application/json',
  'Authorization': `Bearer ${accessToken}` // Required for protected endpoints
};
```

---

## Authentication

### Token Storage
```javascript
// Store tokens after login
localStorage.setItem('accessToken', response.accessToken);
localStorage.setItem('refreshToken', response.refreshToken);

// Retrieve token
const accessToken = localStorage.getItem('accessToken');
```

### Helper Function
```javascript
const getAuthHeaders = () => {
  const token = localStorage.getItem('accessToken');
  return {
    'Content-Type': 'application/json',
    'Authorization': token ? `Bearer ${token}` : ''
  };
};
```

---

## API Endpoints

## Auth

### 1. Login
**POST** `/api/v1/auth/login`

**No Authentication Required**

**Request:**
```javascript
const response = await fetch(`${API_BASE_URL}/auth/login`, {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    email: 'admin@example.com',
    password: 'SecurePass123!'
  })
});
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "admin@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "isActive": true
  }
}
```

---

### 2. Register
**POST** `/api/v1/auth/register`

**No Authentication Required**

**Request (Create New Tenant):**
```javascript
const response = await fetch(`${API_BASE_URL}/auth/register`, {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    email: 'admin@acmecorp.com',
    password: 'SecurePass123!',
    firstName: 'John',
    lastName: 'Doe',
    tenant: {
      name: 'John Doe',
      email: 'admin@acmecorp.com',
      companyName: 'Acme Corporation',
      domain: 'acmecorp.com'
    }
  })
});
```

**Request (Join Existing Tenant):**
```javascript
const response = await fetch(`${API_BASE_URL}/auth/register`, {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    email: 'employee@acmecorp.com',
    password: 'SecurePass123!',
    firstName: 'Jane',
    lastName: 'Smith',
    tenantId: '123e4567-e89b-12d3-a456-426614174000'
  })
});
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "admin@acmecorp.com",
    "firstName": "John",
    "lastName": "Doe"
  }
}
```

---

### 3. Refresh Token
**POST** `/api/v1/auth/refresh`

**No Authentication Required**

**Request:**
```javascript
const response = await fetch(`${API_BASE_URL}/auth/refresh`, {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    refreshToken: localStorage.getItem('refreshToken')
  })
});
```

**Response:**
```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refreshToken": "new_refresh_token...",
  "expiresIn": 900,
  "refreshExpiresIn": 604800,
  "role": "ADMIN"
}
```

---

### 4. Get Profile (Self)
**GET** `/api/v1/auth/profile`

**Authentication Required**

**Request:**
```javascript
const response = await fetch(`${API_BASE_URL}/auth/profile`, {
  method: 'GET',
  headers: getAuthHeaders()
});
```

---

### 5. Get Profile by ID
**GET** `/api/v1/auth/profile/{id}`

**Authentication Required**

**Request:**
```javascript
const userId = '123e4567-e89b-12d3-a456-426614174000';
const response = await fetch(`${API_BASE_URL}/auth/profile/${userId}`, {
  method: 'GET',
  headers: getAuthHeaders()
});
```

---

### 6. Create Employee Account
**POST** `/api/v1/auth/employee`

**Authentication Required | Admin Only**

**Request:**
```javascript
const response = await fetch(`${API_BASE_URL}/auth/employee`, {
  method: 'POST',
  headers: getAuthHeaders(),
  body: JSON.stringify({
    email: 'employee@acmecorp.com',
    password: 'EmployeePass123!',
    firstName: 'Jane',
    lastName: 'Smith'
  })
});
```

---

## Tenants

### 1. Get All Tenants
**GET** `/api/v1/tenants`

**Authentication Required**

**Request:**
```javascript
const response = await fetch(`${API_BASE_URL}/tenants`, {
  method: 'GET',
  headers: getAuthHeaders()
});
```

**Response:**
```json
[
  {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "companyName": "Acme Corp",
    "domain": "acme.com",
    "active": true,
    "createdAt": "2023-01-01T00:00:00Z",
    "updatedAt": "2023-01-01T00:00:00Z"
  }
]
```

---

### 2. Get Tenant by ID
**GET** `/api/v1/tenants/{id}`

**Authentication Required**

**Request:**
```javascript
const tenantId = '123e4567-e89b-12d3-a456-426614174000';
const response = await fetch(`${API_BASE_URL}/tenants/${tenantId}`, {
  method: 'GET',
  headers: getAuthHeaders()
});
```

---

### 3. Create Tenant
**POST** `/api/v1/tenants`

**Authentication Required | Admin Only**

**Request:**
```javascript
const response = await fetch(`${API_BASE_URL}/tenants`, {
  method: 'POST',
  headers: getAuthHeaders(),
  body: JSON.stringify({
    name: 'John Doe',
    email: 'john.doe@example.com',
    companyName: 'Acme Corp',
    domain: 'acme.com'
  })
});
```

---

### 4. Update Tenant
**PUT** `/api/v1/tenants/{id}`

**Authentication Required | Admin Only**

**Request:**
```javascript
const tenantId = '123e4567-e89b-12d3-a456-426614174000';
const response = await fetch(`${API_BASE_URL}/tenants/${tenantId}`, {
  method: 'PUT',
  headers: getAuthHeaders(),
  body: JSON.stringify({
    name: 'John Smith',
    email: 'john.smith@example.com',
    companyName: 'Acme Inc',
    domain: 'acme-inc.com',
    active: true
  })
});
```

---

## Subscriptions

### 1. Get All Subscriptions
**GET** `/api/v1/subscriptions`

**Authentication Required**

**Request:**
```javascript
const response = await fetch(`${API_BASE_URL}/subscriptions`, {
  method: 'GET',
  headers: getAuthHeaders()
});
```

**Response:**
```json
[
  {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "tenantId": "123e4567-e89b-12d3-a456-426614174001",
    "planId": "123e4567-e89b-12d3-a456-426614174002",
    "startDate": "2023-01-01T00:00:00Z",
    "endDate": "2023-12-31T23:59:59Z",
    "status": "active",
    "createdAt": "2023-01-01T00:00:00Z",
    "updatedAt": "2023-01-01T00:00:00Z"
  }
]
```

---

### 2. Get Subscription by ID
**GET** `/api/v1/subscriptions/{id}`

**Authentication Required**

**Request:**
```javascript
const subscriptionId = '123e4567-e89b-12d3-a456-426614174000';
const response = await fetch(`${API_BASE_URL}/subscriptions/${subscriptionId}`, {
  method: 'GET',
  headers: getAuthHeaders()
});
```

---

### 3. Get Subscriptions by Tenant ID
**GET** `/api/v1/subscriptions/tenant/{tenantId}`

**Authentication Required**

**Request:**
```javascript
const tenantId = '123e4567-e89b-12d3-a456-426614174000';
const response = await fetch(`${API_BASE_URL}/subscriptions/tenant/${tenantId}`, {
  method: 'GET',
  headers: getAuthHeaders()
});
```

---

### 4. Create Subscription
**POST** `/api/v1/subscriptions`

**Authentication Required | Admin Only**

**Request:**
```javascript
const response = await fetch(`${API_BASE_URL}/subscriptions`, {
  method: 'POST',
  headers: getAuthHeaders(),
  body: JSON.stringify({
    tenantId: '123e4567-e89b-12d3-a456-426614174000',
    planId: '123e4567-e89b-12d3-a456-426614174001',
    billingCycle: 'monthly',
    paymentMethod: 'credit_card',
    startDate: '2023-01-01T00:00:00Z',
    endDate: '2023-12-31T23:59:59Z',
    status: 'active'
  })
});
```

---

### 5. Update Subscription
**PUT** `/api/v1/subscriptions/{id}`

**Authentication Required | Admin Only**

**Request:**
```javascript
const subscriptionId = '123e4567-e89b-12d3-a456-426614174000';
const response = await fetch(`${API_BASE_URL}/subscriptions/${subscriptionId}`, {
  method: 'PUT',
  headers: getAuthHeaders(),
  body: JSON.stringify({
    tenantId: '123e4567-e89b-12d3-a456-426614174000',
    planId: '123e4567-e89b-12d3-a456-426614174001',
    startDate: '2023-01-01T00:00:00Z',
    endDate: '2023-12-31T23:59:59Z',
    status: 'active'
  })
});
```

---

### 6. Cancel Subscription
**POST** `/api/v1/subscriptions/{id}/cancel`

**Authentication Required | Admin Only**

**Request:**
```javascript
const subscriptionId = '123e4567-e89b-12d3-a456-426614174000';
const response = await fetch(`${API_BASE_URL}/subscriptions/${subscriptionId}/cancel`, {
  method: 'POST',
  headers: getAuthHeaders(),
  body: JSON.stringify({
    reason: 'Customer requested cancellation'
  })
});
```

**Response:**
```json
{
  "message": "Subscription cancelled successfully",
  "reason": "Customer requested cancellation"
}
```

---

### 7. Renew Subscription
**POST** `/api/v1/subscriptions/{id}/renew`

**Authentication Required | Admin Only**

**Request:**
```javascript
const subscriptionId = '123e4567-e89b-12d3-a456-426614174000';
const response = await fetch(`${API_BASE_URL}/subscriptions/${subscriptionId}/renew`, {
  method: 'POST',
  headers: getAuthHeaders(),
  body: JSON.stringify({
    billingCycle: 'yearly',
    paymentMethod: 'credit_card',
    months: 12
  })
});
```

---

## Invoices

### 1. Get All Invoices
**GET** `/api/v1/invoices`

**Authentication Required**

**Query Parameters:**
- `status` (optional): Filter by status (draft, sent, paid, overdue, canceled)
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10)

**Request:**
```javascript
const response = await fetch(`${API_BASE_URL}/invoices?status=paid&page=1&limit=10`, {
  method: 'GET',
  headers: getAuthHeaders()
});
```

**Response:**
```json
[
  {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "subscriptionId": "123e4567-e89b-12d3-a456-426614174001",
    "tenantId": "123e4567-e89b-12d3-a456-426614174002",
    "invoiceNumber": "INV-2023-001",
    "amount": 99.99,
    "status": "paid",
    "dueDate": "2023-01-31T00:00:00Z",
    "sentAt": "2023-01-01T00:00:00Z",
    "paidAt": "2023-01-15T00:00:00Z",
    "paymentId": "pay_1234567890",
    "createdAt": "2023-01-01T00:00:00Z",
    "updatedAt": "2023-01-15T00:00:00Z"
  }
]
```

---

### 2. Get Invoice by ID
**GET** `/api/v1/invoices/{id}`

**Authentication Required**

**Request:**
```javascript
const invoiceId = '123e4567-e89b-12d3-a456-426614174000';
const response = await fetch(`${API_BASE_URL}/invoices/${invoiceId}`, {
  method: 'GET',
  headers: getAuthHeaders()
});
```

---

### 3. Get Invoices by Subscription ID
**GET** `/api/v1/invoices/subscription/{subscriptionId}`

**Authentication Required**

**Request:**
```javascript
const subscriptionId = '123e4567-e89b-12d3-a456-426614174000';
const response = await fetch(`${API_BASE_URL}/invoices/subscription/${subscriptionId}`, {
  method: 'GET',
  headers: getAuthHeaders()
});
```

---

### 4. Create Invoice
**POST** `/api/v1/invoices`

**Authentication Required | Admin Only**

**Request:**
```javascript
const response = await fetch(`${API_BASE_URL}/invoices`, {
  method: 'POST',
  headers: getAuthHeaders(),
  body: JSON.stringify({
    subscriptionId: '123e4567-e89b-12d3-a456-426614174000',
    tenantId: '123e4567-e89b-12d3-a456-426614174001',
    invoiceNumber: 'INV-2023-001',
    amount: 99.99,
    dueDate: '2023-01-31T00:00:00Z',
    status: 'draft'
  })
});
```

---

### 5. Update Invoice
**PUT** `/api/v1/invoices/{id}`

**Authentication Required | Admin Only**

**Request:**
```javascript
const invoiceId = '123e4567-e89b-12d3-a456-426614174000';
const response = await fetch(`${API_BASE_URL}/invoices/${invoiceId}`, {
  method: 'PUT',
  headers: getAuthHeaders(),
  body: JSON.stringify({
    subscriptionId: '123e4567-e89b-12d3-a456-426614174000',
    amount: 149.99,
    dueDate: '2023-02-28T00:00:00Z',
    status: 'sent'
  })
});
```

---

### 6. Send Invoice
**POST** `/api/v1/invoices/{id}/send`

**Authentication Required | Admin Only**

**Request:**
```javascript
const invoiceId = '123e4567-e89b-12d3-a456-426614174000';
const response = await fetch(`${API_BASE_URL}/invoices/${invoiceId}/send`, {
  method: 'POST',
  headers: getAuthHeaders(),
  body: JSON.stringify({
    email: 'customer@example.com'
  })
});
```

**Response:**
```json
{
  "message": "Invoice sent successfully",
  "email": "customer@example.com"
}
```

---

### 7. Pay Invoice
**POST** `/api/v1/invoices/{id}/pay`

**Authentication Required**

**Request:**
```javascript
const invoiceId = '123e4567-e89b-12d3-a456-426614174000';
const response = await fetch(`${API_BASE_URL}/invoices/${invoiceId}/pay`, {
  method: 'POST',
  headers: getAuthHeaders(),
  body: JSON.stringify({
    paymentMethod: 'credit_card',
    transactionId: 'txn_1234567890',
    paymentId: 'pay_1234567890'
  })
});
```

**Response:**
```json
{
  "message": "Invoice paid successfully",
  "paymentId": "pay_1234567890",
  "paymentMethod": "credit_card"
}
```

---

## Master Data

### Units

#### 1. List Units
**GET** `/api/v1/master/units`

**Authentication Required | Admin, HR Only**

**Request:**
```javascript
const response = await fetch(`${API_BASE_URL}/master/units`, {
  method: 'GET',
  headers: getAuthHeaders()
});
```

**Response:**
```json
[
  {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "code": "IT",
    "name": "Information Technology",
    "createdAt": "2023-01-01T00:00:00Z",
    "updatedAt": "2023-01-01T00:00:00Z"
  }
]
```

---

#### 2. Create Unit
**POST** `/api/v1/master/units`

**Authentication Required | Admin, HR Only**

**Request:**
```javascript
const response = await fetch(`${API_BASE_URL}/master/units`, {
  method: 'POST',
  headers: getAuthHeaders(),
  body: JSON.stringify({
    code: 'IT',
    name: 'Information Technology'
  })
});
```

---

### Positions

#### 1. List Positions
**GET** `/api/v1/master/positions`

**Authentication Required | Admin, HR Only**

**Request:**
```javascript
const response = await fetch(`${API_BASE_URL}/master/positions`, {
  method: 'GET',
  headers: getAuthHeaders()
});
```

---

#### 2. Create Position
**POST** `/api/v1/master/positions`

**Authentication Required | Admin, HR Only**

**Request:**
```javascript
const response = await fetch(`${API_BASE_URL}/master/positions`, {
  method: 'POST',
  headers: getAuthHeaders(),
  body: JSON.stringify({
    companyId: '123e4567-e89b-12d3-a456-426614174000',
    title: 'Software Engineer',
    description: 'Develop and maintain software applications',
    level: '3'
  })
});
```

---

### Employees

#### 1. List Employees
**GET** `/api/v1/master/employees`

**Authentication Required | Admin, HR, Manager Only**

**Request:**
```javascript
const response = await fetch(`${API_BASE_URL}/master/employees`, {
  method: 'GET',
  headers: getAuthHeaders()
});
```

---

#### 2. Get Employee by ID
**GET** `/api/v1/master/employees/{id}`

**Authentication Required | Admin, HR, Manager Only**

**Request:**
```javascript
const employeeId = '123e4567-e89b-12d3-a456-426614174000';
const response = await fetch(`${API_BASE_URL}/master/employees/${employeeId}`, {
  method: 'GET',
  headers: getAuthHeaders()
});
```

---

#### 3. Create Employee
**POST** `/api/v1/master/employees`

**Authentication Required | Admin, HR, Manager Only**

**Request:**
```javascript
const response = await fetch(`${API_BASE_URL}/master/employees`, {
  method: 'POST',
  headers: getAuthHeaders(),
  body: JSON.stringify({
    employeeCode: 'EMP001',
    fullName: 'John Doe',
    email: 'john.doe@example.com',
    unitId: '123e4567-e89b-12d3-a456-426614174000',
    positionId: '123e4567-e89b-12d3-a456-426614174001',
    employmentStatus: 'FULLTIME' // FULLTIME, PARTTIME, CONTRACT, INTERN
  })
});
```

---

#### 4. Update Employee
**PUT** `/api/v1/master/employees/{id}`

**Authentication Required | Admin, HR, Manager Only**

**Request:**
```javascript
const employeeId = '123e4567-e89b-12d3-a456-426614174000';
const response = await fetch(`${API_BASE_URL}/master/employees/${employeeId}`, {
  method: 'PUT',
  headers: getAuthHeaders(),
  body: JSON.stringify({
    employeeCode: 'EMP001',
    fullName: 'John Smith',
    email: 'john.smith@example.com',
    unitId: '123e4567-e89b-12d3-a456-426614174000',
    positionId: '123e4567-e89b-12d3-a456-426614174001'
  })
});
```

---

#### 5. Delete Employee
**DELETE** `/api/v1/master/employees/{id}`

**Authentication Required | Admin, HR, Manager Only**

**Request:**
```javascript
const employeeId = '123e4567-e89b-12d3-a456-426614174000';
const response = await fetch(`${API_BASE_URL}/master/employees/${employeeId}`, {
  method: 'DELETE',
  headers: getAuthHeaders()
});
```

**Response:** `204 No Content`

---

## Companies

### 1. Get All Companies
**GET** `/api/v1/companies`

**Authentication Required**

**Request:**
```javascript
const response = await fetch(`${API_BASE_URL}/companies`, {
  method: 'GET',
  headers: getAuthHeaders()
});
```

---

### 2. Create Company
**POST** `/api/v1/companies`

**Authentication Required**

**Request:**
```javascript
const response = await fetch(`${API_BASE_URL}/companies`, {
  method: 'POST',
  headers: getAuthHeaders(),
  body: JSON.stringify({
    tenantId: '123e4567-e89b-12d3-a456-426614174000',
    name: 'Acme Corporation',
    registrationNo: '123456789',
    address: '123 Main St, Anytown, USA',
    timezone: 'UTC',
    currency: 'USD'
  })
});
```

---

### 3. Get Company Profile
**GET** `/api/v1/companies/{id}`

**Authentication Required**

**Request:**
```javascript
const companyId = '123e4567-e89b-12d3-a456-426614174000';
const response = await fetch(`${API_BASE_URL}/companies/${companyId}`, {
  method: 'GET',
  headers: getAuthHeaders()
});
```

---

### 4. Update Company Profile
**PATCH** `/api/v1/companies/{id}`

**Authentication Required**

**Request:**
```javascript
const companyId = '123e4567-e89b-12d3-a456-426614174000';
const response = await fetch(`${API_BASE_URL}/companies/${companyId}`, {
  method: 'PATCH',
  headers: getAuthHeaders(),
  body: JSON.stringify({
    name: 'Acme Inc',
    registrationNo: '987654321',
    address: '456 Oak Ave, Anytown, USA',
    timezone: 'America/New_York',
    currency: 'EUR'
  })
});
```

---

### 5. Get Company Settings
**GET** `/api/v1/companies/{id}/settings`

**Authentication Required**

**Request:**
```javascript
const companyId = '123e4567-e89b-12d3-a456-426614174000';
const response = await fetch(`${API_BASE_URL}/companies/${companyId}/settings`, {
  method: 'GET',
  headers: getAuthHeaders()
});
```

---

### 6. Update Company Settings
**PATCH** `/api/v1/companies/{id}/settings`

**Authentication Required**

**Request:**
```javascript
const companyId = '123e4567-e89b-12d3-a456-426614174000';
const response = await fetch(`${API_BASE_URL}/companies/${companyId}/settings`, {
  method: 'PATCH',
  headers: getAuthHeaders(),
  body: JSON.stringify({
    name: 'Acme Inc',
    registrationNo: '987654321',
    address: '456 Oak Ave, Anytown, USA',
    timezone: 'America/New_York',
    currency: 'EUR'
  })
});
```

---

### 7. Get Company Limits
**GET** `/api/v1/companies/{id}/limits`

**Authentication Required**

**Request:**
```javascript
const companyId = '123e4567-e89b-12d3-a456-426614174000';
const response = await fetch(`${API_BASE_URL}/companies/${companyId}/limits`, {
  method: 'GET',
  headers: getAuthHeaders()
});
```

**Response:**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "tenantId": "123e4567-e89b-12d3-a456-426614174001",
  "name": "Acme Corp",
  "maxEmployees": 100,
  "maxDepartments": 10,
  "maxPositions": 50
}
```

---

## Error Handling

### Error Response Format
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Error message description",
    "fields": {
      "fieldName": "Field specific error"
    }
  },
  "meta": {
    "trace_id": "trace-id-here"
  }
}
```

### Common Error Codes

| HTTP Status | Error Code | Description |
|-------------|------------|-------------|
| 400 | `INVALID_JSON` | Invalid JSON payload |
| 400 | `VALIDATION_ERROR` | Validation failed |
| 401 | `UNAUTHENTICATED` | Missing or invalid token |
| 403 | `FORBIDDEN` | Insufficient permissions |
| 404 | `NOT_FOUND` | Resource not found |
| 500 | `INTERNAL_SERVER_ERROR` | Server error |

### Error Handling Example
```javascript
const handleApiCall = async (url, options) => {
  try {
    const response = await fetch(url, options);
    const data = await response.json();
    
    if (!response.ok) {
      if (response.status === 401) {
        // Token expired, try refresh
        await refreshToken();
        // Retry request
        return await fetch(url, options);
      }
      throw new Error(data.error?.message || 'Request failed');
    }
    
    return data;
  } catch (error) {
    console.error('API Error:', error);
    throw error;
  }
};
```

---

## React.js Examples

### API Service Class
```javascript
class ApiService {
  constructor(baseURL) {
    this.baseURL = baseURL;
  }

  getAuthHeaders() {
    const token = localStorage.getItem('accessToken');
    return {
      'Content-Type': 'application/json',
      'Authorization': token ? `Bearer ${token}` : ''
    };
  }

  async request(endpoint, options = {}) {
    const url = `${this.baseURL}${endpoint}`;
    const config = {
      ...options,
      headers: {
        ...this.getAuthHeaders(),
        ...options.headers
      }
    };

    const response = await fetch(url, config);
    const data = await response.json();

    if (!response.ok) {
      if (response.status === 401) {
        // Handle token refresh
        await this.refreshToken();
        // Retry with new token
        return this.request(endpoint, options);
      }
      throw new Error(data.error?.message || 'Request failed');
    }

    return data;
  }

  async refreshToken() {
    const refreshToken = localStorage.getItem('refreshToken');
    if (!refreshToken) {
      throw new Error('No refresh token available');
    }

    const response = await fetch(`${this.baseURL}/auth/refresh`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refreshToken })
    });

    const data = await response.json();
    if (response.ok) {
      localStorage.setItem('accessToken', data.accessToken);
      localStorage.setItem('refreshToken', data.refreshToken);
    } else {
      // Redirect to login
      window.location.href = '/login';
    }
  }

  // Auth methods
  async login(email, password) {
    return this.request('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password })
    });
  }

  async register(userData) {
    return this.request('/auth/register', {
      method: 'POST',
      body: JSON.stringify(userData)
    });
  }

  // Tenant methods
  async getTenants() {
    return this.request('/tenants');
  }

  async getTenant(id) {
    return this.request(`/tenants/${id}`);
  }

  async createTenant(tenantData) {
    return this.request('/tenants', {
      method: 'POST',
      body: JSON.stringify(tenantData)
    });
  }

  async updateTenant(id, tenantData) {
    return this.request(`/tenants/${id}`, {
      method: 'PUT',
      body: JSON.stringify(tenantData)
    });
  }

  // Subscription methods
  async getSubscriptions() {
    return this.request('/subscriptions');
  }

  async getSubscription(id) {
    return this.request(`/subscriptions/${id}`);
  }

  async createSubscription(subscriptionData) {
    return this.request('/subscriptions', {
      method: 'POST',
      body: JSON.stringify(subscriptionData)
    });
  }

  // Invoice methods
  async getInvoices(params = {}) {
    const queryString = new URLSearchParams(params).toString();
    return this.request(`/invoices${queryString ? `?${queryString}` : ''}`);
  }

  async getInvoice(id) {
    return this.request(`/invoices/${id}`);
  }

  async payInvoice(id, paymentData) {
    return this.request(`/invoices/${id}/pay`, {
      method: 'POST',
      body: JSON.stringify(paymentData)
    });
  }

  // Master Data methods
  async getUnits() {
    return this.request('/master/units');
  }

  async createUnit(unitData) {
    return this.request('/master/units', {
      method: 'POST',
      body: JSON.stringify(unitData)
    });
  }

  async getEmployees() {
    return this.request('/master/employees');
  }

  async createEmployee(employeeData) {
    return this.request('/master/employees', {
      method: 'POST',
      body: JSON.stringify(employeeData)
    });
  }

  // Company methods
  async getCompanies() {
    return this.request('/companies');
  }

  async getCompany(id) {
    return this.request(`/companies/${id}`);
  }

  async createCompany(companyData) {
    return this.request('/companies', {
      method: 'POST',
      body: JSON.stringify(companyData)
    });
  }
}

// Export singleton instance
export const apiService = new ApiService('http://localhost:8081/api/v1');
```

### React Hook Example
```javascript
import { useState, useEffect } from 'react';
import { apiService } from './services/api';

function useTenants() {
  const [tenants, setTenants] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchTenants = async () => {
      try {
        setLoading(true);
        const data = await apiService.getTenants();
        setTenants(data);
        setError(null);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchTenants();
  }, []);

  return { tenants, loading, error };
}

// Usage in component
function TenantsList() {
  const { tenants, loading, error } = useTenants();

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <ul>
      {tenants.map(tenant => (
        <li key={tenant.id}>{tenant.name}</li>
      ))}
    </ul>
  );
}
```

### Login Component Example
```javascript
import { useState } from 'react';
import { apiService } from './services/api';
import { useNavigate } from 'react-router-dom';

function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    try {
      const response = await apiService.login(email, password);
      localStorage.setItem('accessToken', response.accessToken);
      localStorage.setItem('refreshToken', response.refreshToken);
      navigate('/dashboard');
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        placeholder="Email"
        required
      />
      <input
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        placeholder="Password"
        required
      />
      {error && <div className="error">{error}</div>}
      <button type="submit">Login</button>
    </form>
  );
}
```

---

## Role-Based Access Control

### Roles
- **ADMIN**: Full access to all endpoints
- **HR**: Access to master data (units, positions, employees)
- **MANAGER**: Access to employees (read/write)
- **EMPLOYEE**: Limited access (read-only for most resources)

### Endpoint Access Matrix

| Endpoint | ADMIN | HR | MANAGER | EMPLOYEE |
|----------|-------|----|---------|----------|
| Auth endpoints | ✅ | ✅ | ✅ | ✅ |
| Tenants (GET) | ✅ | ✅ | ✅ | ✅ |
| Tenants (POST/PUT) | ✅ | ❌ | ❌ | ❌ |
| Subscriptions (GET) | ✅ | ✅ | ✅ | ✅ |
| Subscriptions (POST/PUT) | ✅ | ❌ | ❌ | ❌ |
| Invoices (GET) | ✅ | ✅ | ✅ | ✅ |
| Invoices (POST/PUT) | ✅ | ❌ | ❌ | ❌ |
| Invoices (Pay) | ✅ | ✅ | ✅ | ✅ |
| Master Data (Units/Positions) | ✅ | ✅ | ❌ | ❌ |
| Master Data (Employees) | ✅ | ✅ | ✅ | ❌ |
| Companies | ✅ | ✅ | ✅ | ✅ |

---

## Notes

1. **Base URL**: Default is `http://localhost:8081/api/v1`. Update for production.
2. **Token Expiry**: Access tokens expire in 15 minutes (900 seconds). Use refresh token to get new tokens.
3. **CORS**: Ensure your React app origin is added to allowed origins in backend config.
4. **Date Format**: All dates are in ISO 8601 format (e.g., `2023-01-01T00:00:00Z`).
5. **UUID Format**: All IDs are UUIDs in format `123e4567-e89b-12d3-a456-426614174000`.

---

**Last Updated**: 2024
**API Version**: v1

