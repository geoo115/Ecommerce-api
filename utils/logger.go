package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

// LogLevel represents the logging level
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// Logger represents a structured logger
type Logger struct {
	level LogLevel
	*log.Logger
}

// Global logger instance
var AppLogger *Logger

// Initialize logger
func init() {
	AppLogger = NewLogger(INFO)
}

// NewLogger creates a new logger instance
func NewLogger(level LogLevel) *Logger {
	return &Logger{
		level:  level,
		Logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

// SetLogLevel sets the logging level
func (l *Logger) SetLogLevel(level LogLevel) {
	l.level = level
}

// SetLogLevelFromString sets the logging level from a string
func (l *Logger) SetLogLevelFromString(level string) {
	switch level {
	case "debug":
		l.SetLogLevel(DEBUG)
	case "info":
		l.SetLogLevel(INFO)
	case "warn":
		l.SetLogLevel(WARN)
	case "error":
		l.SetLogLevel(ERROR)
	case "fatal":
		l.SetLogLevel(FATAL)
	default:
		l.SetLogLevel(INFO)
	}
}

// logMessage formats and logs a message with the given level
func (l *Logger) logMessage(level LogLevel, format string, args ...interface{}) {
	if level >= l.level {
		levelStr := []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}[level]
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		message := fmt.Sprintf(format, args...)
		l.Printf("[%s] %s: %s", timestamp, levelStr, message)
	}
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.logMessage(DEBUG, format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.logMessage(INFO, format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	l.logMessage(WARN, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.logMessage(ERROR, format, args...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.logMessage(FATAL, format, args...)
	os.Exit(1)
}

// LogRequest logs HTTP request details
func (l *Logger) LogRequest(method, path, clientIP string, statusCode int, duration time.Duration) {
	l.Info("HTTP Request: %s %s from %s - Status: %d - Duration: %v",
		method, path, clientIP, statusCode, duration)
}

// LogError logs error details with context
func (l *Logger) LogError(err error, context string) {
	l.Error("Error in %s: %v", context, err)
}

// LogDatabase logs database operation details
func (l *Logger) LogDatabase(operation, table string, duration time.Duration) {
	l.Debug("Database %s on %s - Duration: %v", operation, table, duration)
}

// LogSecurity logs security-related events
func (l *Logger) LogSecurity(event, clientIP string, details ...interface{}) {
	l.Warn("Security Event: %s from %s - Details: %v", event, clientIP, details)
}

// Convenience functions for global logger
func Debug(format string, args ...interface{}) {
	AppLogger.Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	AppLogger.Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	AppLogger.Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	AppLogger.Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
	AppLogger.Fatal(format, args...)
}
