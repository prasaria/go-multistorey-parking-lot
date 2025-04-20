package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/prasaria/go-multistorey-parking-lot/pkg/config"
)

const (
	version = "0.1.0"
)

func main() {
	fmt.Printf("Multi-Storey Parking Lot System v%s\n", version)
	fmt.Println("----------------------------------")

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := strings.ToLower((os.Args[1]))
	args := os.Args[2:]

	switch command {
	case "init":
		handleInit(args)
	case "park":
		fmt.Println("Park command not yet implemented")
		// TODO: Implement
	case "unpark":
		fmt.Println("Unpark command not yet implemented")
		// TODO : Implement
	case "available":
		fmt.Println("Available command not yet implemented")
		// TODO : Implement
	case "search":
		fmt.Println("Available command not yet implemented")
		// TODO : Implement
	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}

func handleInit(args []string) {
	fmt.Println("Initializing parking lot...")

	// Parse command-line flags
	cfg, err := config.ParseInitCommand(args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println(config.GetUsage())
		return
	}

	// Print configuration summary
	fmt.Println("Parking lot configuration:")
	fmt.Printf("- Floors: %d\n", cfg.Floors)
	fmt.Printf("- Rows per floor: %d\n", cfg.Rows)
	fmt.Printf("- Columns per row: %d\n", cfg.Columns)

	// TODO: Initialize the parking lot system with this configuration
	fmt.Println("Parking lot initialized successfully")
}

func printUsage() {
	fmt.Println("Usage: parking-lot <command> [arguments]")
	fmt.Println("\nAvailable commands:")
	fmt.Println("  init       Initialize a new parking lot")
	fmt.Println("  park       Park a vehicle")
	fmt.Println("  unpark     Remove a vehicle from its spot")
	fmt.Println("  available  List available spots for a vehicle type")
	fmt.Println("  search     Find a vehicle in the parking lot")
	fmt.Println("  help       Show this help message")
	fmt.Println("\nRun 'parking-lot <command> --help' for more information on a command.")
}
