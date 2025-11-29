# Colorful Error Messages Usage Guide

This guide shows you how to use the new colorful styling utility in your Go API.

## Basic Usage

Import the utils package and use the global `S` instance like `fmt`:

```go
import "github.com/capamir/go-api/utils"

// Error logging
utils.S.Error("Something went wrong")
utils.S.Errorf("Failed to connect to database: %v", err)

// Warning logging  
utils.S.Warn("User not found")
utils.S.Warnf("Invalid login attempt for user: %s", email)

// Success logging
utils.S.Success("User created successfully")
utils.S.Successf("Order %d processed", orderID)

// Info logging
utils.S.Info("Server starting")
utils.S.Infof("Processing %d items", count)

// Debug logging
utils.S.Debug("Cache hit")
utils.S.Debugf("Request took %dms", duration)
```

## HTTP Error Responses

Replace your current `utils.WriteError()` calls with the enhanced version that adds colors to logs:

```go
// Before (still works, but now with colors in logs!)
utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))

// The HTTP response remains the same JSON format, but logs are now colorful:
// 🟡 WARN: HTTP 400: invalid user ID
```

## Direct Color Usage

Use colors directly for custom formatting:

```go
import "github.com/capamir/go-api/utils"

fmt.Printf("Status: %s\n", utils.Green("Active"))
fmt.Printf("Error: %s\n", utils.RedBold("Connection failed"))
fmt.Printf("Warning: %s\n", utils.Yellow("Rate limit exceeded"))

// Available colors:
// utils.Red(), utils.Green(), utils.Blue(), utils.Yellow(), utils.Cyan()
// utils.RedBold(), utils.GreenBold(), utils.BlueBold(), etc.
// utils.BgRed(), utils.BgGreen(), utils.BgYellow() (background colors)
```

## Special Boxes and Banners

```go
// Error boxes
utils.S.ErrorBox("VALIDATION ERROR", "Email format is invalid")

// Success boxes  
utils.S.SuccessBox("USER CREATED", "Registration completed")

// Info boxes
utils.S.InfoBox("API INFO", "Rate limit: 100 req/hour")

// Server startup banner
utils.S.ServerBanner("localhost:8080", "v1.0.0")

// Database status
utils.S.DBStatus("connected", "MySQL pool active")
utils.S.DBStatus("error", "Connection failed")
utils.S.DBStatus("connecting", "Establishing connection...")
```

## Example Handler with Colors

```go
func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
    var payload types.LoginUserPayload
    if err := utils.ParseJSON(r, &payload); err != nil {
        utils.S.Warnf("Invalid JSON payload from %s", r.RemoteAddr)
        utils.WriteError(w, http.StatusBadRequest, err)
        return
    }

    user, err := h.store.GetUserByEmail(payload.Email)
    if err != nil {
        utils.S.Warnf("Login attempt for non-existent user: %s", payload.Email)
        utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid credentials"))
        return
    }

    if !auth.ComparePasswords(user.Password, []byte(payload.Password)) {
        utils.S.Warnf("Invalid password for user: %s", payload.Email)
        utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid credentials"))
        return
    }

    utils.S.Successf("User %s logged in successfully", user.Email)
    token, _ := auth.CreateJWT(secret, user.ID)
    utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}
```

## Color Scheme

- 🔴 **Red**: Errors, fatal messages
- 🟡 **Yellow**: Warnings, deprecations
- 🔵 **Blue**: Info, general messages  
- 🟢 **Green**: Success, positive outcomes
- 🔷 **Cyan**: Debug, development info
- ⚫ **Gray**: Secondary information

## Benefits

1. **Better Development Experience**: Quickly spot errors and important messages in logs
2. **Maintain Clean JSON APIs**: HTTP responses remain unchanged, only logs get colors
3. **Drop-in Replacement**: Your existing `utils.WriteError()` calls automatically get colors
4. **Flexible**: Use simple logging or fancy boxes depending on context
5. **Professional**: Server banners and status indicators make your app feel polished

## Terminal Compatibility

The `fatih/color` package automatically detects terminal capabilities and disables colors when:
- Output is redirected to a file
- Terminal doesn't support colors
- Running on Windows without color support

This ensures your logs work everywhere without issues.