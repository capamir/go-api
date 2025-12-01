HTTP Request
    ↓
Handler (receives RegisterRequest DTO)
    ↓
Service (validates, transforms to models.User)
    ↓
Repository (saves models.User to DB)
    ↓
Database (stores as table row)

[Return flow]

Database (returns table row)
    ↓
Repository (returns models.User)
    ↓
Service (may add business logic)
    ↓
Handler (calls user.ToResponse())
    ↓
HTTP Response (sends UserResponse as JSON)
