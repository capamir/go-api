## Development Commands

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

---

## Code Quality

### Format Code
```
go fmt ./...
```

### Lint Code (Requires golangci-lint)
```
# Install golangci-lint first:
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run
```

### Vet Code (Find Potential Bugs)
```
go vet ./...
```

---

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

## PowerShell Scripts (Recommended!)

Create `scripts/run.ps1`:
```
# Run the application
Write-Host "Starting application..." -ForegroundColor Green
go run cmd/api/main.go
```

Create `scripts/build.ps1`:
```
# Build the application
Write-Host "Building application..." -ForegroundColor Blue
New-Item -ItemType Directory -Force -Path bin | Out-Null
go build -o bin/api.exe cmd/api/main.go
Write-Host "Build complete: bin/api.exe" -ForegroundColor Green
```

Create `scripts/test.ps1`:
```
# Run tests
Write-Host "Running tests..." -ForegroundColor Blue
go test -v ./...
```

**Usage:**
```
.\scripts\run.ps1
.\scripts\build.ps1
.\scripts\test.ps1
```

---

## Environment Setup

### Create .env file (First Time)
```
Copy-Item .env.example .env
```

### Edit .env file
```
notepad .env
```

---

## Git Commands

### Initial Setup
```
git init
git add .
git commit -m "Initial commit - Phase 2 complete"
```

### Add Remote
```
git remote add origin https://github.com/capamir/go-api.git
git branch -M main
git push -u origin main
```

### Create Phase Branch
```
git checkout -b phase-2-authentication
git add .
git commit -m "Complete Phase 2: Authentication system"
git push origin phase-2-authentication
```

---

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

## Windows-Specific Notes

1. **Use backticks (`) for line continuation** in PowerShell (not backslash `\`)

2. **Escape quotes in JSON strings:**
   ```
   # Wrong:
   -d '{"key": "value"}'
   
   # Correct:
   -d '{\"key\": \"value\"}'
   ```

3. **Path separators:**
   - PowerShell accepts both `/` and `\`
   - Use `/` for cross-platform compatibility

4. **Installing `make` on Windows (Optional):**
   ```
   # Using Chocolatey
   choco install make
   
   # Using Scoop
   scoop install make
   ```

---

## Recommended Tools for Windows

1. **Windows Terminal** - Modern terminal with tabs
2. **Git Bash** - Unix-like commands on Windows
3. **VS Code** - Best Go IDE for Windows
4. **Postman** - API testing (easier than curl)
5. **MySQL Workbench** - Database GUI

---

## Alternative: Use Task Runner

Install **Task** (Makefile alternative):
```
go install github.com/go-task/task/v3/cmd/task@latest
```

Create `Taskfile.yml`:
```
version: '3'

tasks:
  run:
    desc: Run the application
    cmds:
      - go run cmd/api/main.go

  build:
    desc: Build the application
    cmds:
      - go build -o bin/api.exe cmd/api/main.go

  test:
    desc: Run tests
    cmds:
      - go test -v ./...
```

**Usage:**
```
task run
task build
task test
```

---

**All commands tested on Windows 11 with PowerShell 7** ✅
```

***
