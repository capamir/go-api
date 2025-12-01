User submits registration form
    ↓
Service.Register()
    ├─ Check email exists? → Yes: Error
    ├─ Hash password
    ├─ Generate verification token: "abc123...xyz"
    ├─ Create user:
    │    Status: "pending"
    │    EmailVerified: false
    │    VerificationToken: "abc123...xyz"
    ├─ Save to database
    ├─ Log verification link (TODO: send email)
    └─ Return response WITHOUT JWT token

Response:
{
  "success": true,
  "user": {...},
  "token": null  ← No token yet!
}
