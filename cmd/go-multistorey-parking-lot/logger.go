package main

import (
	"fmt"
	"os"
	"time"
)

// LogLevel represents the level of a log message
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarning
	LogLevelError
)

// Logger is a simple logger with levels
type Logger struct {
	Verbose bool
	Level   LogLevel
}

// NewLogger creates a new logger
func NewLogger(verbose bool) *Logger {
	level := LogLevelInfo
	if verbose {
		level = LogLevelDebug
	}

	return &Logger{
		Verbose: verbose,
		Level:   level,
	}
}

// formatLogMessage formats a log message with timestamp and level
func formatLogMessage(level string, message string) string {
	timestamp := time.Now().Format("15:04:05.000")
	return fmt.Sprintf("[%s] [%s] %s", timestamp, level, message)
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	if l.Verbose && l.Level <= LogLevelDebug {
		message := fmt.Sprintf(format, args...)
		formattedMsg := formatLogMessage("DEBUG", message)
		fmt.Fprintln(os.Stderr, colorCyan+formattedMsg+colorReset)
	}
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	if l.Level <= LogLevelInfo {
		message := fmt.Sprintf(format, args...)
		formattedMsg := formatLogMessage("INFO", message)
		fmt.Fprintln(os.Stderr, colorBlue+formattedMsg+colorReset)
	}
}

// Warning logs a warning message
func (l *Logger) Warning(format string, args ...interface{}) {
	if l.Level <= LogLevelWarning {
		message := fmt.Sprintf(format, args...)
		formattedMsg := formatLogMessage("WARN", message)
		fmt.Fprintln(os.Stderr, colorYellow+formattedMsg+colorReset)
	}
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	if l.Level <= LogLevelError {
		message := fmt.Sprintf(format, args...)
		formattedMsg := formatLogMessage("ERROR", message)
		fmt.Fprintln(os.Stderr, colorRed+formattedMsg+colorReset)
	}
}
