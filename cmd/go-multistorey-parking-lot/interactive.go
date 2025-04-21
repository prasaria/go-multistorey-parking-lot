package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/prasaria/go-multistorey-parking-lot/internal/cli"
	"github.com/prasaria/go-multistorey-parking-lot/internal/model"
)

// CommandRegistryInterface defines the interface for the command registry
type CommandRegistryInterface interface {
	ExecuteCommand(name string, args []string) error
	GetParkingLot() *model.ParkingLot
	GetCommands() map[string]*cli.Command
}

// InteractiveMode contains enhancements for interactive command-line mode
type InteractiveMode struct {
	Registry CommandRegistryInterface
	History  []string
}

// NewInteractiveMode creates a new interactive mode
func NewInteractiveMode(registry CommandRegistryInterface) *InteractiveMode {
	return &InteractiveMode{
		Registry: registry,
		History:  make([]string, 0),
	}
}

// AddToHistory adds a command to the history
func (i *InteractiveMode) AddToHistory(command string) {
	i.History = append(i.History, command)
}

// ProcessCommand processes a command line
func (i *InteractiveMode) ProcessCommand(line string) bool {
	line = strings.TrimSpace(line)

	// Skip empty lines
	if line == "" {
		return true
	}

	// Add to history
	i.AddToHistory(line)

	// Parse command and arguments
	parts := splitCommandLine(line)
	if len(parts) == 0 {
		return true
	}

	command := parts[0]
	args := parts[1:]

	// Handle exit command directly
	if command == "exit" {
		return false
	}

	// Execute the command
	err := i.Registry.ExecuteCommand(command, args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	return true
}

// AutoComplete provides auto-completion for commands
func (i *InteractiveMode) AutoComplete(partial string) []string {
	var completions []string

	// Complete commands
	for name := range i.Registry.GetCommands() {
		if strings.HasPrefix(name, partial) {
			completions = append(completions, name)
		}
	}

	// If it's the "park" command, suggest vehicle types
	parts := splitCommandLine(partial)
	if len(parts) > 0 && len(parts) < 3 {
		command := parts[0]

		if command == "park" || command == "available" {
			// Suggest vehicle types for the first arg
			if len(parts) == 1 {
				completions = append(completions,
					command+" bicycle",
					command+" motorcycle",
					command+" automobile")
			}
		}
	}

	return completions
}

// GetExampleCommands returns example commands for an empty prompt
func (i *InteractiveMode) GetExampleCommands() []string {
	examples := []string{
		"init 3 5 10",
		"park automobile KA-01-HH-1234",
		"unpark 0-1-2 KA-01-HH-1234",
		"available bicycle",
		"search KA-01-HH-1234",
		"status",
		"help",
		"exit",
	}

	// Add more tailored examples if the parking lot is initialized
	parkingLot := i.Registry.GetParkingLot()
	if parkingLot != nil {
		// Get counts by vehicle type to suggest available types
		availableCounts := parkingLot.GetAvailableSpotCountByType()

		var vehicleType model.VehicleType
		if availableCounts[model.VehicleTypeBicycle] > 0 {
			vehicleType = model.VehicleTypeBicycle
		} else if availableCounts[model.VehicleTypeMotorcycle] > 0 {
			vehicleType = model.VehicleTypeMotorcycle
		} else if availableCounts[model.VehicleTypeAutomobile] > 0 {
			vehicleType = model.VehicleTypeAutomobile
		}

		if vehicleType != "" {
			examples = append(examples,
				fmt.Sprintf("park %s ABC-123", strings.ToLower(string(vehicleType))))
		}
	}

	return examples
}
