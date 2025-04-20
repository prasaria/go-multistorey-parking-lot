package main

import (
	"fmt"
	"os"
)

const (
	version = "0.1.0"
)

func main() {
	fmt.Println("Multi-Storey Parking Lot System", version)
	fmt.Println("----------------------------------")

	// Command-line arguments available
	if len(os.Args) > 1 {
		fmt.Println("Command:", os.Args[1])
		fmt.Println("Arguments:", os.Args[2:])
		// TODO: Implement command parsing
	} else {
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage: parking-lot <command> [arguments]")
	fmt.Println("\nAvailable commands:")
	fmt.Println("  init      Initialize a new parking lot")
	fmt.Println("  park      Park a vehicle")
	fmt.Println("  unpark    Remove a vehicle from its spot")
	fmt.Println("  available List available spots for a vehicle type")
	fmt.Println("  search    Find a vehicle in the parking lot")
	fmt.Println("\nRun 'parking-lot <command> --help' for more information on a command.")
}
