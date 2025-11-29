package utils

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fatih/color"
)

// Color functions for different types of messages
var (
	// Error colors
	Red    = color.New(color.FgRed).SprintFunc()
	RedBold = color.New(color.FgRed, color.Bold).SprintFunc()
	
	// Success colors
	Green    = color.New(color.FgGreen).SprintFunc()
	GreenBold = color.New(color.FgGreen, color.Bold).SprintFunc()
	
	// Warning colors
	Yellow    = color.New(color.FgYellow).SprintFunc()
	YellowBold = color.New(color.FgYellow, color.Bold).SprintFunc()
	
	// Info colors
	Blue    = color.New(color.FgBlue).SprintFunc()
	BlueBold = color.New(color.FgBlue, color.Bold).SprintFunc()
	
	// Debug colors
	Cyan    = color.New(color.FgCyan).SprintFunc()
	CyanBold = color.New(color.FgCyan, color.Bold).SprintFunc()
	
	// Neutral colors
	White    = color.New(color.FgWhite).SprintFunc()
	WhiteBold = color.New(color.FgWhite, color.Bold).SprintFunc()
	Gray     = color.New(color.FgHiBlack).SprintFunc()
	
	// Background colors for emphasis
	BgRed    = color.New(color.BgRed, color.FgWhite).SprintFunc()
	BgGreen  = color.New(color.BgGreen, color.FgWhite).SprintFunc()
	BgYellow = color.New(color.BgYellow, color.FgBlack).SprintFunc()
)

// Style provides fmt-like functions for colorful output
type Style struct{}

// Global instance for easy access
var S = &Style{}

// Error logging functions
func (s *Style) Errorf(format string, args ...interface{}) {
	log.Printf("%s %s", Red("ERROR:"), fmt.Sprintf(format, args...))
}

func (s *Style) Error(args ...interface{}) {
	log.Printf("%s %s", Red("ERROR:"), fmt.Sprint(args...))
}

func (s *Style) Fatalf(format string, args ...interface{}) {
	log.Fatalf("%s %s", RedBold("FATAL:"), fmt.Sprintf(format, args...))
}

func (s *Style) Fatal(args ...interface{}) {
	log.Fatalf("%s %s", RedBold("FATAL:"), fmt.Sprint(args...))
}

// Warning logging functions
func (s *Style) Warnf(format string, args ...interface{}) {
	log.Printf("%s %s", Yellow("WARN:"), fmt.Sprintf(format, args...))
}

func (s *Style) Warn(args ...interface{}) {
	log.Printf("%s %s", Yellow("WARN:"), fmt.Sprint(args...))
}

// Info logging functions
func (s *Style) Infof(format string, args ...interface{}) {
	log.Printf("%s %s", Blue("INFO:"), fmt.Sprintf(format, args...))
}

func (s *Style) Info(args ...interface{}) {
	log.Printf("%s %s", Blue("INFO:"), fmt.Sprint(args...))
}

// Success logging functions
func (s *Style) Successf(format string, args ...interface{}) {
	log.Printf("%s %s", Green("SUCCESS:"), fmt.Sprintf(format, args...))
}

func (s *Style) Success(args ...interface{}) {
	log.Printf("%s %s", Green("SUCCESS:"), fmt.Sprint(args...))
}

// Debug logging functions
func (s *Style) Debugf(format string, args ...interface{}) {
	log.Printf("%s %s", Cyan("DEBUG:"), fmt.Sprintf(format, args...))
}

func (s *Style) Debug(args ...interface{}) {
	log.Printf("%s %s", Cyan("DEBUG:"), fmt.Sprint(args...))
}

// HTTP Error response with colors in logs
func (s *Style) WriteError(res http.ResponseWriter, status int, err error) {
	// Log with colors
	switch {
	case status >= 500:
		s.Errorf("HTTP %d: %v", status, err)
	case status >= 400:
		s.Warnf("HTTP %d: %v", status, err)
	default:
		s.Infof("HTTP %d: %v", status, err)
	}
	
	// Send JSON response (no colors in HTTP response)
	WriteJSON(res, status, map[string]string{"error": err.Error()})
}

// Colorful printf functions for direct use
func (s *Style) Printf(color func(...interface{}) string, format string, args ...interface{}) {
	fmt.Printf("%s\n", color(fmt.Sprintf(format, args...)))
}

func (s *Style) Print(color func(...interface{}) string, args ...interface{}) {
	fmt.Printf("%s\n", color(fmt.Sprint(args...)))
}

// Convenience functions for common patterns
func (s *Style) ErrorBox(title, message string) {
	fmt.Printf("┌─ %s ─┐\n", RedBold(title))
	fmt.Printf("│ %s │\n", Red(message))
	fmt.Printf("└%s┘\n", Red("─────────────────────────"))
}

func (s *Style) SuccessBox(title, message string) {
	fmt.Printf("┌─ %s ─┐\n", GreenBold(title))
	fmt.Printf("│ %s │\n", Green(message))
	fmt.Printf("└%s┘\n", Green("─────────────────────────"))
}

func (s *Style) InfoBox(title, message string) {
	fmt.Printf("┌─ %s ─┐\n", BlueBold(title))
	fmt.Printf("│ %s │\n", Blue(message))
	fmt.Printf("└%s┘\n", Blue("─────────────────────────"))
}

// Server startup banner
func (s *Style) ServerBanner(addr, version string) {
	fmt.Printf("\n")
	fmt.Printf("🚀 %s\n", GreenBold("SERVER STARTING"))
	fmt.Printf("📍 %s %s\n", Blue("Address:"), White(addr))
	fmt.Printf("📊 %s %s\n", Blue("Version:"), White(version))
	fmt.Printf("🕐 %s %s\n", Blue("Status:"), Green("Running"))
	fmt.Printf("\n")
}

// Database connection status
func (s *Style) DBStatus(status string, details string) {
	switch status {
	case "connected":
		fmt.Printf("🗄️  %s %s\n", Green("Database:"), Green(details))
	case "error":
		fmt.Printf("🗄️  %s %s\n", Red("Database:"), Red(details))
	case "connecting":
		fmt.Printf("🗄️  %s %s\n", Yellow("Database:"), Yellow(details))
	}
}