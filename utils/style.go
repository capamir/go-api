package utils

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/fatih/color"
)

// LogLevel represents different logging levels
type LogLevel int

const (
	LogLevelError LogLevel = iota
	LogLevelFatal
	LogLevelWarn
	LogLevelInfo
	LogLevelSuccess
	LogLevelDebug
)

// BoxType represents different box styles
type BoxType int

const (
	BoxTypeError BoxType = iota
	BoxTypeSuccess
	BoxTypeInfo
)

// Color functions for different types of messages
var (
	// Error colors
	Red      = color.New(color.FgRed).SprintFunc()
	RedBold  = color.New(color.FgRed, color.Bold).SprintFunc()
	
	// Success colors
	Green     = color.New(color.FgGreen).SprintFunc()
	GreenBold = color.New(color.FgGreen, color.Bold).SprintFunc()
	
	// Warning colors
	Yellow     = color.New(color.FgYellow).SprintFunc()
	YellowBold = color.New(color.FgYellow, color.Bold).SprintFunc()
	
	// Info colors
	Blue     = color.New(color.FgBlue).SprintFunc()
	BlueBold = color.New(color.FgBlue, color.Bold).SprintFunc()
	
	// Debug colors
	Cyan     = color.New(color.FgCyan).SprintFunc()
	CyanBold = color.New(color.FgCyan, color.Bold).SprintFunc()
	
	// Neutral colors
	White     = color.New(color.FgWhite).SprintFunc()
	WhiteBold = color.New(color.FgWhite, color.Bold).SprintFunc()
	Gray      = color.New(color.FgHiBlack).SprintFunc()
	
	// Background colors for emphasis
	BgRed    = color.New(color.BgRed, color.FgWhite).SprintFunc()
	BgGreen  = color.New(color.BgGreen, color.FgWhite).SprintFunc()
	BgYellow = color.New(color.BgYellow, color.FgBlack).SprintFunc()
)

// logConfig holds configuration for each log level
type logConfig struct {
	prefix    string
	colorFunc func(...interface{}) string
	logFunc   func(string, ...interface{})
}

// logConfigs maps log levels to their configurations
var logConfigs = map[LogLevel]logConfig{
	LogLevelError: {
		prefix:    "ERROR:",
		colorFunc: Red,
		logFunc:   log.Printf,
	},
	LogLevelFatal: {
		prefix:    "FATAL:",
		colorFunc: RedBold,
		logFunc:   log.Fatalf,
	},
	LogLevelWarn: {
		prefix:    "WARN:",
		colorFunc: Yellow,
		logFunc:   log.Printf,
	},
	LogLevelInfo: {
		prefix:    "INFO:",
		colorFunc: Blue,
		logFunc:   log.Printf,
	},
	LogLevelSuccess: {
		prefix:    "SUCCESS:",
		colorFunc: Green,
		logFunc:   log.Printf,
	},
	LogLevelDebug: {
		prefix:    "DEBUG:",
		colorFunc: Cyan,
		logFunc:   log.Printf,
	},
}

// boxConfig holds configuration for each box type
type boxConfig struct {
	titleColor   func(...interface{}) string
	messageColor func(...interface{}) string
}

// boxConfigs maps box types to their configurations
var boxConfigs = map[BoxType]boxConfig{
	BoxTypeError: {
		titleColor:   RedBold,
		messageColor: Red,
	},
	BoxTypeSuccess: {
		titleColor:   GreenBold,
		messageColor: Green,
	},
	BoxTypeInfo: {
		titleColor:   BlueBold,
		messageColor: Blue,
	},
}

// Style provides fmt-like functions for colorful output
type Style struct{}

// Global instance for easy access
var S = &Style{}

// logMessage is the core logging function that all others use
func (s *Style) logMessage(level LogLevel, format string, args []interface{}, isFormat bool) {
	config, exists := logConfigs[level]
	if !exists {
		log.Printf("Unknown log level: %v", level)
		return
	}
	
	var message string
	if isFormat {
		message = fmt.Sprintf(format, args...)
	} else {
		message = fmt.Sprint(args...)
	}
	
	formattedMessage := fmt.Sprintf("%s %s", config.colorFunc(config.prefix), message)
	config.logFunc(formattedMessage)
}

// Error logging functions
func (s *Style) Errorf(format string, args ...interface{}) {
	s.logMessage(LogLevelError, format, args, true)
}

func (s *Style) Error(args ...interface{}) {
	s.logMessage(LogLevelError, "", args, false)
}

func (s *Style) Fatalf(format string, args ...interface{}) {
	s.logMessage(LogLevelFatal, format, args, true)
}

func (s *Style) Fatal(args ...interface{}) {
	s.logMessage(LogLevelFatal, "", args, false)
}

// Warning logging functions
func (s *Style) Warnf(format string, args ...interface{}) {
	s.logMessage(LogLevelWarn, format, args, true)
}

func (s *Style) Warn(args ...interface{}) {
	s.logMessage(LogLevelWarn, "", args, false)
}

// Info logging functions
func (s *Style) Infof(format string, args ...interface{}) {
	s.logMessage(LogLevelInfo, format, args, true)
}

func (s *Style) Info(args ...interface{}) {
	s.logMessage(LogLevelInfo, "", args, false)
}

// Success logging functions
func (s *Style) Successf(format string, args ...interface{}) {
	s.logMessage(LogLevelSuccess, format, args, true)
}

func (s *Style) Success(args ...interface{}) {
	s.logMessage(LogLevelSuccess, "", args, false)
}

// Debug logging functions
func (s *Style) Debugf(format string, args ...interface{}) {
	s.logMessage(LogLevelDebug, format, args, true)
}

func (s *Style) Debug(args ...interface{}) {
	s.logMessage(LogLevelDebug, "", args, false)
}

// HTTP Error response with colors in logs
func (s *Style) WriteError(res http.ResponseWriter, status int, err error) {
	// Log with colors based on status code
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
func (s *Style) Printf(colorFunc func(...interface{}) string, format string, args ...interface{}) {
	fmt.Printf("%s\n", colorFunc(fmt.Sprintf(format, args...)))
}

func (s *Style) Print(colorFunc func(...interface{}) string, args ...interface{}) {
	fmt.Printf("%s\n", colorFunc(fmt.Sprint(args...)))
}

// createBox is the core box drawing function
func (s *Style) createBox(boxType BoxType, title, message string) {
	config, exists := boxConfigs[boxType]
	if !exists {
		log.Printf("Unknown box type: %v", boxType)
		return
	}
	
	maxLen := max(len(title), len(message))
	width := maxLen + 4
	border := strings.Repeat("─", width)
	titlePad := strings.Repeat(" ", width-len(title)-2)
	msgPad := strings.Repeat(" ", width-len(message)-2)

	fmt.Printf("┌%s┐\n", border)
	fmt.Printf("│ %s%s │\n", config.titleColor(title), titlePad)
	fmt.Printf("│ %s%s │\n", config.messageColor(message), msgPad)
	fmt.Printf("└%s┘\n", border)
}

// Box functions using the core createBox function
func (s *Style) ErrorBox(title, message string) {
	s.createBox(BoxTypeError, title, message)
}

func (s *Style) SuccessBox(title, message string) {
	s.createBox(BoxTypeSuccess, title, message)
}

func (s *Style) InfoBox(title, message string) {
	s.createBox(BoxTypeInfo, title, message)
}

// max returns the larger of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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