package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Create and initialize command registry
	registry := NewCommandRegistry()
	registry.RegisterAllCommands()

	// Print welcome message
	fmt.Println("Welcome to Parking Lot CLI")
	fmt.Println("Type 'help' to see available commands or 'exit' to quit")

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

		// Get command line
		line := scanner.Text()
		line = strings.TrimSpace(line)

		// Skip empty lines
		if line == "" {
			continue
		}

		// Parse command and arguments
		parts := splitCommandLine(line)
		if len(parts) == 0 {
			continue
		}

		command := parts[0]
		args := parts[1:]

		// Handle exit command directly
		if command == "exit" {
			fmt.Println("Exiting...")
			break
		}

		// Execute the command
		err := registry.ExecuteCommand(command, args)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
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
