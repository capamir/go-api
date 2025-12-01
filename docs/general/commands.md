## Development Commands

## Quick Commands Reference

| Task | Command |
|------|---------|
| Run app | `go run cmd/api/main.go` |
| Build | `go build -o bin/api.exe cmd/api/main.go` |
| Test | `go test -v ./...` |
| Format | `go fmt ./...` |
| Install deps | `go mod download` |
| Tidy deps | `go mod tidy` |
| Clean | `Remove-Item -Recurse -Force bin` |

---

### Run Application
```
go run cmd/api/main.go
```

### Build Application
```
# Create bin directory if it doesn't exist
mkdir bin -Force

# Build executable
go build -o bin/api.exe cmd/api/main.go

# Run the built executable
.\bin\api.exe
```

---

## Testing

### Run All Tests
```
go test -v ./...
```

### Run Specific Package Tests
```
go test -v ./internal/service
```

### Run with Coverage
```
go test -cover ./...
```

### Generate Coverage Report
```
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## Database Operations

### Check Database Tables
```
# Connect to MySQL
mysql -u root -p

# In MySQL shell:
USE ecommerce;
SHOW TABLES;
DESCRIBE users;
```

### Reset Database (Drop and Recreate)
```
# In MySQL shell:
DROP DATABASE IF EXISTS ecommerce;
CREATE DATABASE ecommerce CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

---

## Dependencies

### Install All Dependencies
```
go mod download
```

### Add New Dependency
```
go get -u github.com/package/name
```

### Update Dependencies
```
go get -u ./...
```

### Tidy Dependencies (Remove Unused)
```
go mod tidy
```

### Verify Dependencies
```
go mod verify
```

## Cleaning

### Remove Build Artifacts
```
Remove-Item -Recurse -Force bin
```

### Clean Go Cache
```
go clean -cache
go clean -modcache
```

---

## API Testing with curl (Windows)

### Register User
```
curl -X POST http://localhost:8080/api/v1/auth/register `
  -H "Content-Type: application/json" `
  -d '{\"first_name\":\"John\",\"last_name\":\"Doe\",\"email\":\"john@example.com\",\"password\":\"securepass123\"}'
```

### Login
```
curl -X POST http://localhost:8080/api/v1/auth/login `
  -H "Content-Type: application/json" `
  -d '{\"email\":\"john@example.com\",\"password\":\"securepass123\"}'
```

### Get Profile (Replace TOKEN)
```
$TOKEN = "your-jwt-token-here"

curl -X GET http://localhost:8080/api/v1/auth/me `
  -H "Authorization: Bearer $TOKEN"
```

---

## Recommended Tools for Windows

1. **Windows Terminal** - Modern terminal with tabs
2. **Git Bash** - Unix-like commands on Windows
3. **VS Code** - Best Go IDE for Windows
4. **Postman** - API testing (easier than curl)
5. **MySQL Workbench** - Database GUI

---
