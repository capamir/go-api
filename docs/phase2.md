# Phase 2: Authentication System with Email Verification

## Overview

Built a complete authentication system with user registration, login, JWT token generation, email verification, and protected routes.

## What We Built

### 1. User Model (`internal/models/user.go`)

- **GORM Model** with proper struct tags
- Fields: ID, FirstName, LastName, Email, Password, Phone, Role, Status, Timestamps
- **🆕 Email Verification Fields:**
  - `email_verified` (boolean) - Verification status
  - `email_verified_at` (timestamp) - When verified
  - `verification_token` (string) - Unique verification token
- **Soft Delete** support with `DeletedAt`
- **BeforeCreate Hook** for default values
- **ToResponse()** method to hide sensitive data (password, token)
- **UserResponse** struct for safe API responses
- **🆕 CanLogin()** method - Checks if user can authenticate

**User Status Constants:**
- `pending` - User registered but email not verified
- `active` - Email verified, can login
- `inactive` - Account disabled by admin
- `banned` - Account banned

**Key Features:**

- Password never returned in JSON (using `json:"-"` tag)
- Verification token never exposed in API responses
- Email has unique index
- Automatic timestamp management
- Soft delete capability

---

### 2. Password Utility (`internal/utils/password.go`)

- **HashPassword()** - Bcrypt hashing with cost 10
- **ComparePassword()** - Secure password comparison

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

```
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

### 4. 🆕 Token Utility (`internal/utils/token.go`)

Cryptographically secure token generation for email verification:

**Functions:**
- `GenerateSecureToken(length int)` - Generate random token
- `GenerateVerificationToken()` - 64-character hex token (32 bytes)

**Security:**
- Uses `crypto/rand` (not `math/rand`)
- 256 bits of entropy
- Unpredictable, secure tokens

---

### 5. 🆕 Email Service (`internal/utils/email.go`)

Abstraction layer for sending emails:

**Interface:**
```
type EmailService interface {
    SendVerificationEmail(to, token string) error
}
```

**Implementations:**

1. **ConsoleEmailService** (Development)
   - Logs verification links to console
   - Perfect for testing without email setup
   
2. **SMTPEmailService** (Production - Ready to implement)
   - Sends real emails via SMTP
   - Supports Gmail, SendGrid, custom SMTP
   - Easy to configure via environment variables

**Current Setup:**
- Development: Uses console logging
- Production: Switch to SMTP by updating `.env` and `main.go`

---

### 6. Response Utility (`internal/utils/response.go`)

Standardized JSON responses across the API:
 
**Success Response:**

```
{
  "success": true,
  "message": "Operation successful",
  "data": { ... }
}
```

**Error Response:**

```
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

### 7. User Repository (`internal/repository/user.go`)

Database access layer (GORM operations):

**Methods:**

- `Create(user)` - Insert new user
- `GetByEmail(email)` - Find user by email
- `GetByID(id)` - Find user by ID
- `🆕 GetByVerificationToken(token)` - Find user by verification token
- `Update(user)` - Update user data
- `Delete(id)` - Soft delete user
- `GetAll(limit, offset)` - Paginated user list
- `EmailExists(email)` - Check email availability

**Pattern:**

- Returns `nil` for "not found" (not an error)
- Distinguishes between "not found" and "database error"
- GORM error handling
- Efficient queries with proper indexes

---

### 8. Auth Service (`internal/service/auth.go`)

Business logic layer with email verification:

#### 🆕 Register Flow (Updated):

1. Check if email exists
2. Hash password
3. **Generate verification token**
4. Create user with `status: pending`
5. **Send verification email** (via EmailService)
6. **Do NOT generate JWT token** (user must verify first)
7. Return user data without token

#### Login Flow (Updated):

1. Find user by email
2. **✅ Check if email is verified**
3. Check user status (active/inactive/banned)
4. Verify password
5. Generate JWT token
6. Return token + user data

#### 🆕 Email Verification Flow:

1. Find user by verification token
2. Check if already verified
3. Update user:
   - `email_verified = true`
   - `email_verified_at = now()`
   - `status = active`
   - `verification_token = ""` (clear token)
4. Save to database
5. Return success message

#### 🆕 Resend Verification Flow:

1. Find user by email
2. Check if already verified
3. Check cooldown period (1 minute)
4. Generate new verification token
5. Update user with new token
6. Send new verification email
7. Return success

**Security:**

- Never reveals if email exists (prevents enumeration)
- Generic error messages for failed login
- Password hashing before storage
- Status and verification checks before login
- Tokens are single-use (cleared after verification)
- Cooldown prevents spam

---

### 9. Auth Middleware (`internal/middleware/auth.go`)

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

### 10. Auth Handler (`internal/handler/auth.go`)

HTTP request handlers:

#### Endpoints:

**POST /api/v1/auth/register** (Updated)

- Creates new user account with `status: pending`
- **Does NOT return JWT token**
- Sends verification email
- Status: 201 Created

**POST /api/v1/auth/login** (Updated)

- Authenticates user
- **Checks email verification before allowing login**
- Returns JWT token only if verified
- Status: 200 OK

**🆕 GET /api/v1/auth/verify?token=xxx**

- Verifies user email via token
- Activates user account
- Public endpoint (no auth required)
- Status: 200 OK

**🆕 POST /api/v1/auth/resend-verification**

- Sends new verification email
- Public endpoint (no auth required)
- Generic response (security)
- Status: 200 OK

**GET /api/v1/auth/me**

- Protected route (requires JWT)
- Returns current user profile
- Status: 200 OK

---

## Complete User Journey

### 1️⃣ Registration

```
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "password": "securepass123"
  }'
```

**Response (201):**
```
{
  "success": true,
  "message": "Registration successful! Please check your email to verify your account.",
  "data": {
    "user": {
      "id": 1,
      "first_name": "John",
      "last_name": "Doe",
      "email": "john@example.com",
      "email_verified": false,
      "role": "customer",
      "status": "pending",
      "created_at": "2025-12-01T10:00:00Z"
    }
  }
}
```

**📧 Check server logs for verification link:**
```
📧 [DEV] Verification email to john@example.com: http://localhost:8080/api/v1/auth/verify?token=abc123...
```

---

### 2️⃣ Try Login (Will Fail - Email Not Verified)

```
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepass123"
  }'
```

**Response (401):**
```
{
  "success": false,
  "error": "please verify your email before logging in"
}
```

---

### 3️⃣ Verify Email

```
# Copy token from server logs
curl "http://localhost:8080/api/v1/auth/verify?token=abc123..."
```

**Response (200):**
```
{
  "success": true,
  "message": "Email verified successfully! You can now login.",
  "data": {
    "message": "Email verified successfully! You can now login.",
    "user": {
      "id": 1,
      "first_name": "John",
      "last_name": "Doe",
      "email": "john@example.com",
      "email_verified": true,
      "email_verified_at": "2025-12-01T10:05:00Z",
      "role": "customer",
      "status": "active",
      "created_at": "2025-12-01T10:00:00Z"
    }
  }
}
```

---

### 4️⃣ Login (Now Works!)

```
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepass123"
  }'
```

**Response (200):**
```
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "first_name": "John",
      "last_name": "Doe",
      "email": "john@example.com",
      "email_verified": true,
      "role": "customer",
      "status": "active"
    }
  }
}
```

---

### 5️⃣ Resend Verification (If Needed)

```
curl -X POST http://localhost:8080/api/v1/auth/resend-verification \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com"
  }'
```

**Response (200):**
```
{
  "success": true,
  "message": "If the email exists and is not verified, a verification link has been sent.",
  "data": null
}
```

---

### 6️⃣ Get Profile (Protected Route)

```
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer $TOKEN"
```

**Response (200):**
```
{
  "success": true,
  "message": "Profile retrieved",
  "data": {
    "id": 1,
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "email_verified": true,
    "role": "customer",
    "status": "active"
  }
}
```

---

## Architecture Pattern

We used **Clean Architecture** / **Layered Architecture**:

```
Handler (HTTP)
    ↓
Service (Business Logic) ← Uses EmailService interface
    ↓
Repository (Database)
    ↓
Database (GORM/MySQL)
```

**Benefits:**

- **Separation of Concerns** - Each layer has one responsibility
- **Testable** - Can mock each layer (including email service)
- **Maintainable** - Easy to modify one layer without affecting others
- **Scalable** - Can swap implementations (e.g., change email provider)
- **Dependency Injection** - EmailService injected at runtime

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

3. **Email Verification Security**

   - Cryptographically secure tokens (256 bits)
   - Single-use tokens (cleared after verification)
   - Tokens never exposed in API responses
   - User must verify before login

4. **API Security**

   - Generic error messages (no info leakage)
   - Protected routes with middleware
   - Token validation on every request
   - User status checking
   - Email enumeration prevention

5. **Input Validation**
   - Gin binding validation
   - Email format validation
   - Password length validation
   - Required field checking

---

## Database Schema

**users table:**

```
CREATE TABLE users (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  first_name varchar(100) NOT NULL,
  last_name varchar(100) NOT NULL,
  email varchar(255) NOT NULL,
  password varchar(255) NOT NULL,
  phone varchar(20),
  role varchar(20) DEFAULT 'customer',
  status varchar(20) DEFAULT 'pending',
  email_verified tinyint(1) DEFAULT 0,
  email_verified_at datetime(3),
  verification_token varchar(255),
  created_at datetime(3),
  updated_at datetime(3),
  deleted_at datetime(3),
  PRIMARY KEY (id),
  UNIQUE KEY idx_users_email (email),
  KEY idx_users_deleted_at (deleted_at),
  KEY idx_verification_token (verification_token)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

---

## Files Created/Updated in Phase 2

```
internal/
├── models/
│   └── user.go                  # ✅ Updated with verification fields
├── repository/
│   └── user.go                  # ✅ Updated with GetByVerificationToken()
├── service/
│   └── auth.go                  # ✅ Updated with verification logic
├── handler/
│   └── auth.go                  # ✅ Updated with verification handlers
├── middleware/
│   └── auth.go                  # JWT authentication middleware
├── utils/
│   ├── password.go              # Password hashing utilities
│   ├── jwt.go                   # JWT token utilities
│   ├── 🆕 token.go              # Verification token generation
│   ├── 🆕 email.go              # Email service abstraction
│   └── response.go              # Standard API responses
└── database/
    └── migrate.go               # Auto-migration setup

cmd/api/
└── main.go                      # ✅ Updated with EmailService injection
```

---

## Colorful Terminal Output

**Enhanced logging with email verification:**

```
📧 [DEV] Verification email to john@example.com: http://localhost:8080/api/v1/auth/verify?token=abc123...
✅ SUCCESS: User registered (pending verification): john@example.com
⚠️ WARN: Login attempt with unverified email: john@example.com
✅ SUCCESS: Email verified for user: john@example.com
✅ SUCCESS: User logged in successfully: john@example.com
```

---

## Testing Checklist

- [ ] Register user → Status should be `pending`, no JWT token
- [ ] Try login before verification → Should fail with error message
- [ ] Verify email with token → Status should become `active`
- [ ] Try login after verification → Should succeed with JWT token
- [ ] Get profile with JWT → Should return user data
- [ ] Resend verification → Should generate new token
- [ ] Try verify with invalid token → Should fail
- [ ] Try verify already verified user → Should show "already verified"

---

## Common Issues & Solutions

**Issue:** "Please verify your email before logging in"
- **Solution:** Check server logs for verification link and click it.

**Issue:** "Invalid or expired verification token"
- **Solution:** Token may have been used already or is incorrect. Request resend.

**Issue:** "Email already registered"
- **Solution:** Each email must be unique. Try different email.

**Issue:** "Account is not active"
- **Solution:** User status is not "active". Check if email is verified.

**Issue:** Missing verification link in logs
- **Solution:** Make sure EmailService is properly initialized in `main.go`.

---

## Environment Variables Required

```
# Application
APP_ENV=development
APP_PORT=8080
APP_URL=http://localhost:8080

# Database
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=ecommerce

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRATION=24h

# Server
SERVER_READ_TIMEOUT=15s
SERVER_WRITE_TIMEOUT=15s

# 🆕 Email (Optional - For Production SMTP)
# EMAIL_FROM=your-email@gmail.com
# EMAIL_PASSWORD=your-app-password
# EMAIL_HOST=smtp.gmail.com
# EMAIL_PORT=587
```

---

## Future Enhancements (Optional)

- [ ] Password reset via email
- [ ] Email templates (HTML emails)
- [ ] Production SMTP setup (Gmail/SendGrid)
- [ ] Token expiration (verification links expire after 24 hours)
- [ ] Rate limiting on registration
- [ ] 2FA (Two-Factor Authentication)
- [ ] OAuth2 (Google, GitHub login)
- [ ] Email change verification

---

## Key Learnings

1. **Email Verification Flow** - Standard practice for user registration
2. **Token Generation** - Using crypto/rand for security
3. **Dependency Injection** - EmailService interface pattern
4. **Status Management** - pending → active user states
5. **Single-Use Tokens** - Clear tokens after use
6. **Console Development** - Test without email setup
7. **Clean Architecture** - Easy to swap email providers
8. **Security Best Practices** - No email enumeration, generic errors

---

## Phase 2: Complete! ✅

**What We Built:**
- ✅ User registration with email verification
- ✅ Email verification via secure tokens
- ✅ Login with verification check
- ✅ Resend verification functionality
- ✅ JWT authentication
- ✅ Protected routes
- ✅ Clean architecture with dependency injection
- ✅ Console email service for development
- ✅ Production-ready email abstraction

**Ready for Phase 3!** 🚀

---

**Next Steps:**
- Phase 3: Product Management
- Phase 4: Shopping Cart
- Phase 5: Order Management
- Phase 6: Payment Integration
```
