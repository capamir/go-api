# Phase 2: Authentication System

## Overview

Built a complete authentication system with user registration, login, JWT token generation, and protected routes.

## What We Built

### 1. User Model (`internal/models/user.go`)

- **GORM Model** with proper struct tags
- Fields: ID, FirstName, LastName, Email, Password, Phone, Role, Status, Timestamps
- **Soft Delete** support with `DeletedAt`
- **BeforeCreate Hook** for default values
- **ToResponse()** method to hide sensitive data (password)
- **UserResponse** struct for safe API responses

**Key Features:**

- Password never returned in JSON (using `json:"-"` tag)
- Email has unique index
- Automatic timestamp management
- Soft delete capability

---

### 2. Password Utility (`internal/utils/password.go`)

- **HashPassword()** - Bcrypt hashing with cost 10
- **ComparePassword()** - Secure password comparison
- **ValidatePassword()** - 8-72 character validation

**Security:**

- Uses bcrypt (industry standard)
- Salt automatically generated
- Constant-time comparison prevents timing attacks

---

### 3. JWT Utility (`internal/utils/jwt.go`)

- **GenerateToken()** - Create JWT with custom claims
- **ValidateToken()** - Parse and verify JWT
- **JWTClaims** struct with UserID, Email, Role

**JWT Structure:**

```json
{
  "user_id": 1,
  "email": "user@example.com",
  "role": "customer",
  "exp": 1732934400,
  "iat": 1732848000,
  "iss": "ecommerce-api",
  "sub": "1"
}
```

**Features:**

- HMAC-SHA256 signing
- Configurable expiration (default 24 hours)
- Standard JWT claims (exp, iat, nbf, iss, sub)
- Custom claims (user_id, email, role)

---

### 4. Response Utility (`internal/utils/response.go`)

Standardized JSON responses across the API:

**Success Response:**

```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... }
}
```

**Error Response:**

```json
{
  "success": false,
  "error": "Error message"
}
```

**Helper Functions:**

- `RespondSuccess()` - 200 OK
- `RespondCreated()` - 201 Created
- `RespondBadRequest()` - 400 Bad Request
- `RespondUnauthorized()` - 401 Unauthorized
- `RespondForbidden()` - 403 Forbidden
- `RespondNotFound()` - 404 Not Found
- `RespondInternalError()` - 500 Internal Server Error

---

### 5. User Repository (`internal/repository/user.go`)

Database access layer (GORM operations):

**Methods:**

- `Create(user)` - Insert new user
- `GetByEmail(email)` - Find user by email
- `GetByID(id)` - Find user by ID
- `Update(user)` - Update user data
- `Delete(id)` - Soft delete user
- `GetAll(limit, offset)` - Paginated user list
- `EmailExists(email)` - Check email availability

**Pattern:**

- Returns `nil` for "not found" (not an error)
- Distinguishes between "not found" and "database error"
- GORM error handling

---

### 6. Auth Service (`internal/service/auth_service.go`)

Business logic layer:

#### Register Flow:

1. Check if email exists
2. Validate password strength
3. Hash password
4. Create user in database
5. Generate JWT token
6. Return token + user data

#### Login Flow:

1. Find user by email
2. Check user status (active/inactive)
3. Verify password
4. Generate JWT token
5. Return token + user data

**Security:**

- Never reveals if email exists (prevents enumeration)
- Generic error messages for failed login
- Password hashing before storage
- Status check before login

---

### 7. Auth Middleware (`internal/middleware/auth.go`)

JWT validation middleware:

**Flow:**

1. Extract token from `Authorization: Bearer <token>` header
2. Validate token format
3. Parse and verify JWT signature
4. Check expiration
5. Store user info in Gin context
6. Continue to handler OR abort with 401

**Context Keys:**

- `user_id` - User ID (uint)
- `user_email` - Email (string)
- `user_role` - Role (string)

**Helper Functions:**

- `GetUserID(c)` - Extract user ID from context
- `GetUserEmail(c)` - Extract email from context
- `GetUserRole(c)` - Extract role from context

---

### 8. Auth Handler (`internal/handler/auth_handler.go`)

HTTP request handlers:

#### Endpoints:

**POST /api/v1/auth/register**

- Creates new user account
- Returns JWT token
- Status: 201 Created

**POST /api/v1/auth/login**

- Authenticates user
- Returns JWT token
- Status: 200 OK

**GET /api/v1/auth/me**

- Protected route (requires JWT)
- Returns current user profile
- Status: 200 OK

**Request/Response Examples:**

**Register:**

Request:

```http
POST /api/v1/auth/register
Content-Type: application/json
```

```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john@example.com",
  "password": "securepass123",
  "phone": "+1234567890"
}
```

Response (201):

```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "first_name": "John",
      "last_name": "Doe",
      "email": "john@example.com",
      "phone": "+1234567890",
      "role": "customer",
      "status": "active",
      "created_at": "2025-11-30T01:57:00Z"
    }
  }
}
```

**Login:**
**Login:**

Request:

```http
POST /api/v1/auth/login
Content-Type: application/json
```

```json
{
  "email": "john@example.com",
  "password": "securepass123"
}
```

Response (200):

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": { ... }
  }
}
```

**Get Profile:**

Request:

```http
GET /api/v1/auth/me
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

Response (200):

```json
{
  "success": true,
  "message": "Profile retrieved",
  "data": {
    "id": 1,
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "role": "customer",
    "status": "active"
  }
}
```

---

### 9. Database Auto-Migration (`internal/database/migrate.go`)

- Automatic schema generation from GORM models
- Creates/updates tables based on struct definitions
- Handles indexes, foreign keys, constraints

**Usage:**
database.AutoMigrate() // Creates users table

---

## Architecture Pattern

We used **Clean Architecture** / **Layered Architecture**:
Handler (HTTP)
↓
Service (Business Logic)
↓
Repository (Database)
↓
Database (GORM/MySQL)

**Benefits:**

- **Separation of Concerns** - Each layer has one responsibility
- **Testable** - Can mock each layer
- **Maintainable** - Easy to modify one layer without affecting others
- **Scalable** - Can swap implementations (e.g., change database)

---

## Security Features Implemented

1. **Password Security**

   - Bcrypt hashing (cost 10)
   - Salting automatic
   - Never stored in plain text
   - Never returned in API responses

2. **JWT Security**

   - HMAC-SHA256 signing
   - Expiration validation
   - Signature verification
   - Secure secret from environment

3. **API Security**

   - Generic error messages (no info leakage)
   - Protected routes with middleware
   - Token validation on every request
   - User status checking

4. **Input Validation**
   - Gin binding validation
   - Email format validation
   - Password length validation
   - Required field checking

---

## Colorful Terminal Output

**Yes!** We use our custom logger (`pkg/logger`) with colorful output:

✅ SUCCESS: User registered successfully: john@example.com
⚠️ WARN: Login attempt with non-existent email: fake@example.com
❌ ERROR: Failed to hash password: invalid length
🔄 Database: Connecting to MySQL...
✅ Database: Successfully connected to database 'ecommerce'

**Color Scheme:**

- 🔵 **INFO** (Blue) - General information
- 🟢 **SUCCESS** (Green) - Successful operations
- 🟡 **WARN** (Yellow) - Warnings
- 🔴 **ERROR** (Red) - Errors
- 🔵 **DEBUG** (Cyan) - Debug info

---

## Files Created in Phase 2

internal/
├── models/
│ └── user.go # User GORM model
├── repository/
│ └── user.go # User database operations
├── service/
│ └── auth_service.go # Authentication business logic
├── handler/
│ └── auth_handler.go # HTTP handlers for auth
├── middleware/
│ └── auth.go # JWT authentication middleware
├── utils/
│ ├── password.go # Password hashing utilities
│ ├── jwt.go # JWT token utilities
│ └── response.go # Standard API responses
└── database/
└── migrate.go # Auto-migration setup

---

## Database Schema

**users table:**

```sql
CREATE TABLE users (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  first_name varchar(100) NOT NULL,
  last_name varchar(100) NOT NULL,
  email varchar(255) NOT NULL,
  password varchar(255) NOT NULL,
  phone varchar(20),
  role varchar(20) DEFAULT 'customer',
  status varchar(20) DEFAULT 'active',
  created_at datetime(3),
  updated_at datetime(3),
  deleted_at datetime(3),
  PRIMARY KEY (id),
  UNIQUE KEY idx_users_email (email),
  KEY idx_users_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

---

## Testing the API

### 1. Register a User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "password": "securepass123"
  }'
```

### 2. Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepass123"
  }'
```

### 3. Get Profile (Protected)

TOKEN="your-jwt-token-here"

curl -X GET http://localhost:8080/api/v1/auth/me
-H "Authorization: Bearer $TOKEN"

---

## Next Steps (Phase 3)

- [ ] Product model (GORM)
- [ ] Product repository
- [ ] Product service
- [ ] Product handlers
- [ ] Product CRUD endpoints
- [ ] Image upload support
- [ ] Product search & filtering
- [ ] Pagination

---

## Key Learnings

1. **GORM Hooks** - `BeforeCreate` for default values
2. **Gin Binding** - Automatic request validation
3. **JWT Claims** - Custom + Standard claims
4. **Middleware Pattern** - Reusable authentication
5. **Clean Architecture** - Layer separation
6. **Error Handling** - Generic messages for security
7. **Context Usage** - Passing data between middleware and handlers
8. **Bcrypt** - Secure password hashing

---

## Common Issues & Solutions

**Issue:** "Email already registered"

- **Solution:** Each email must be unique. Try different email.

**Issue:** "Invalid or expired token"

- **Solution:** Token expired (24h). Login again to get new token.

**Issue:** "Missing authentication token"

- **Solution:** Add `Authorization: Bearer <token>` header to request.

**Issue:** "Invalid request data"

- **Solution:** Check JSON format and required fields (first_name, last_name, email, password).

---

## Environment Variables Required

Application
APP_ENV=development
APP_PORT=8080

Database
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=ecommerce

JWT
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRATION_HOURS=24

Server
SERVER_READ_TIMEOUT=15
SERVER_WRITE_TIMEOUT=15
