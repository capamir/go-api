User clicks link: /verify?token=abc123...xyz
    ↓
Service.VerifyEmail("abc123...xyz")
    ├─ Find user by token
    ├─ User found?
    │    ├─ Yes: Continue
    │    └─ No: Error "invalid token"
    ├─ Already verified?
    │    ├─ Yes: "Already verified"
    │    └─ No: Continue
    ├─ Update user:
    │    EmailVerified: true
    │    EmailVerifiedAt: 2025-12-01T00:46:00Z
    │    Status: "active"
    │    VerificationToken: ""  ← Clear token
    └─ Save to database

Response:
{
  "message": "Email verified! You can now login.",
  "user": {...}
}
