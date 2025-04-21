package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/prasaria/go-multistorey-parking-lot/internal/cli"
)

func main() {
	// Create and initialize command registry
	registry := cli.NewCommandRegistry()
	registry.RegisterAllCommands()

	// Create interactive mode
	interactive := NewInteractiveMode(registry)

	// Print welcome message
	fmt.Println("Welcome to Parking Lot CLI")
	fmt.Println("Type 'help' to see available commands or 'exit' to quit")
	fmt.Println("Options:")
	fmt.Println("  --json    Output results in JSON format")
	fmt.Println("  --verbose Show detailed operation logs")

	// Create scanner for reading user input
	scanner := bufio.NewScanner(os.Stdin)

	// Main loop
	for {
		// Show prompt
		fmt.Print("> ")

		// Read input
		if !scanner.Scan() {
			break
		}

		// Process command
		line := scanner.Text()
		if !interactive.ProcessCommand(line) {
			break
		}
	}
}

// splitCommandLine splits a command line into parts, handling quotes
func splitCommandLine(line string) []string {
	var parts []string
	var currentPart strings.Builder
	inQuotes := false

	for _, char := range line {
		switch {
		case char == ' ' && !inQuotes:
			// Space outside quotes, end current part
			if currentPart.Len() > 0 {
				parts = append(parts, currentPart.String())
				currentPart.Reset()
			}
		case char == '"':
			// Toggle quote state
			inQuotes = !inQuotes
		default:
			// Add character to current part
			currentPart.WriteRune(char)
		}
	}

	// Add the last part if not empty
	if currentPart.Len() > 0 {
		parts = append(parts, currentPart.String())
	}

	return parts
}
