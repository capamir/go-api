User submits login form
    ↓
Service.Login()
    ├─ Find user by email
    ├─ User exists? → No: Error "invalid email/password"
    ├─ Email verified? → No: Error "verify email first"
    ├─ Can login? (active + verified + not banned)
    │    └─ No: Error "account not active"
    ├─ Password correct? → No: Error "invalid email/password"
    ├─ Generate JWT token ✅
    └─ Return token + user

Response:
{
  "success": true,
  "token": "eyJhbGc...",  ← NOW user gets token!
  "user": {...}
}
