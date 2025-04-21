package main

import (
	"fmt"
	"strings"
	"time"
)

// Colors for terminal output
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

// PrintError prints an error message in red
func PrintError(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Printf("%sError: %s%s\n", colorRed, message, colorReset)
}

// PrintSuccess prints a success message in green
func PrintSuccess(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Printf("%s%s%s\n", colorGreen, message, colorReset)
}

// PrintInfo prints an info message in blue
func PrintInfo(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Printf("%s%s%s\n", colorBlue, message, colorReset)
}

// PrintWarning prints a warning message in yellow
func PrintWarning(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Printf("%s%s%s\n", colorYellow, message, colorReset)
}

// FormatTable formats data as a table with columns
func FormatTable(headers []string, rows [][]string) string {
	if len(rows) == 0 {
		return "No data to display"
	}

	// Calculate column widths
	colWidths := make([]int, len(headers))
	for i, header := range headers {
		colWidths[i] = len(header)
	}

	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// Build the table
	var builder strings.Builder

	// Add headers
	for i, header := range headers {
		format := fmt.Sprintf("%%-%ds", colWidths[i]+2)
		builder.WriteString(fmt.Sprintf(format, header))
	}
	builder.WriteString("\n")

	// Add separator
	for _, width := range colWidths {
		builder.WriteString(strings.Repeat("-", width+2))
	}
	builder.WriteString("\n")

	// Add rows
	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) {
				format := fmt.Sprintf("%%-%ds", colWidths[i]+2)
				builder.WriteString(fmt.Sprintf(format, cell))
			}
		}
		builder.WriteString("\n")
	}

	return builder.String()
}

// FormatDuration formats a duration in a human-readable format
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%d seconds", int(d.Seconds()))
	} else if d < time.Hour {
		return fmt.Sprintf("%d minutes, %d seconds",
			int(d.Minutes()), int(d.Seconds())%60)
	} else if d < 24*time.Hour {
		return fmt.Sprintf("%d hours, %d minutes",
			int(d.Hours()), int(d.Minutes())%60)
	} else {
		days := int(d.Hours()) / 24
		return fmt.Sprintf("%d days, %d hours",
			days, int(d.Hours())%24)
	}
}
