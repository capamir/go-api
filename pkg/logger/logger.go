package logger

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
)

// Logger provides colorful logging methods
type Logger struct {
	infoColor    *color.Color
	successColor *color.Color
	warnColor    *color.Color
	errorColor   *color.Color
	debugColor   *color.Color
}

// New creates a new Logger instance
func New() *Logger {
	return &Logger{
		infoColor:    color.New(color.FgBlue, color.Bold),
		successColor: color.New(color.FgGreen, color.Bold),
		warnColor:    color.New(color.FgYellow, color.Bold),
		errorColor:   color.New(color.FgRed, color.Bold),
		debugColor:   color.New(color.FgCyan),
	}
}

// Global logger instance
var Log = New()

// Info logs an informational message
func (l *Logger) Info(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	log.Printf("%s %s", l.infoColor.Sprint("INFO:"), message)
}

// Success logs a success message
func (l *Logger) Success(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	log.Printf("%s %s", l.successColor.Sprint("SUCCESS:"), message)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	log.Printf("%s %s", l.warnColor.Sprint("WARN:"), message)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	log.Printf("%s %s", l.errorColor.Sprint("ERROR:"), message)
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	log.Printf("%s %s", l.debugColor.Sprint("DEBUG:"), message)
}

// Fatal logs a fatal error and exits
func (l *Logger) Fatal(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	log.Fatalf("%s %s", l.errorColor.Sprint("FATAL:"), message)
}

// Banner prints a startup banner
func (l *Logger) Banner(appName, version, port string) {
	banner := fmt.Sprintf(`
╔══════════════════════════════════════════╗
║                                          ║
║   🚀 %s                      ║
║   📊 Version: %-27s║
║   🌐 Port: %-30s║
║   ✨ Status: Running                    ║
║                                          ║
╚══════════════════════════════════════════╝
`, appName, version, port)

	fmt.Println(l.successColor.Sprint(banner))
}

// DB logs database-related messages with emoji
func (l *Logger) DB(status, message string) {
	var emoji string
	var colorFunc *color.Color

	switch status {
	case "connecting":
		emoji = "🔄"
		colorFunc = l.warnColor
	case "connected":
		emoji = "✅"
		colorFunc = l.successColor
	case "error":
		emoji = "❌"
		colorFunc = l.errorColor
	default:
		emoji = "📊"
		colorFunc = l.infoColor
	}

	log.Printf("%s %s %s", emoji, colorFunc.Sprint("Database:"), message)
}

// HTTP logs HTTP request information
func (l *Logger) HTTP(method, path string, status int, duration string) {
	var statusColor *color.Color

	switch {
	case status >= 500:
		statusColor = l.errorColor
	case status >= 400:
		statusColor = l.warnColor
	case status >= 300:
		statusColor = color.New(color.FgCyan)
	case status >= 200:
		statusColor = l.successColor
	default:
		statusColor = color.New(color.FgWhite)
	}

	log.Printf("➜ %s %s %s %s",
		color.New(color.Bold).Sprint(method),
		path,
		statusColor.Sprintf("[%d]", status),
		l.debugColor.Sprintf("(%s)", duration),
	)
}

// Custom helper functions

// Info is a package-level function
func Info(format string, args ...interface{}) {
	Log.Info(format, args...)
}

// Success is a package-level function
func Success(format string, args ...interface{}) {
	Log.Success(format, args...)
}

// Warn is a package-level function
func Warn(format string, args ...interface{}) {
	Log.Warn(format, args...)
}

// Error is a package-level function
func Error(format string, args ...interface{}) {
	Log.Error(format, args...)
}

// Debug is a package-level function
func Debug(format string, args ...interface{}) {
	Log.Debug(format, args...)
}

// Fatal is a package-level function
func Fatal(format string, args ...interface{}) {
	Log.Fatal(format, args...)
}

func init() {
	// Set log flags
	log.SetFlags(log.Ldate | log.Ltime)
	log.SetOutput(os.Stdout)
}
