package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	perrors "github.com/prasaria/go-multistorey-parking-lot/internal/errors"
	"github.com/prasaria/go-multistorey-parking-lot/internal/model"
)

// Command represents a CLI command
type Command struct {
	Name        string
	Usage       string
	Description string
	MinArgs     int
	MaxArgs     int
	Handler     func(args []string) error
}

// OutputFormat represents the format of command output
type OutputFormat int

const (
	OutputFormatText OutputFormat = iota
	OutputFormatJSON
)

// CommandOptions contains options for command execution
type CommandOptions struct {
	Format  OutputFormat
	Verbose bool
}

// Update CommandRegistry to include options
type CommandRegistry struct {
	Commands   map[string]*Command
	parkingLot *model.ParkingLot
	Options    CommandOptions
}

// NewCommandRegistry creates a new command registry
func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		Commands: make(map[string]*Command),
		Options: CommandOptions{
			Format:  OutputFormatText,
			Verbose: false,
		},
	}
}

// RegisterCommand adds a command to the registry
func (r *CommandRegistry) RegisterCommand(cmd *Command) {
	r.Commands[cmd.Name] = cmd
}

// GetCommand returns a command by name
func (r *CommandRegistry) GetCommand(name string) (*Command, bool) {
	cmd, found := r.Commands[name]
	return cmd, found
}

// Add option parsing to ExecuteCommand
func (r *CommandRegistry) ExecuteCommand(name string, args []string) error {
	// Parse options first
	filteredArgs := make([]string, 0)
	for _, arg := range args {
		if arg == "--json" {
			r.Options.Format = OutputFormatJSON
		} else if arg == "--verbose" || arg == "-v" {
			r.Options.Verbose = true
		} else {
			filteredArgs = append(filteredArgs, arg)
		}
	}

	// Look up command
	cmd, found := r.GetCommand(name)
	if !found {
		return fmt.Errorf("unknown command: %s\nType 'help' to see available commands", name)
	}

	// Validate argument count (with filtered args now)
	if len(filteredArgs) < cmd.MinArgs {
		return fmt.Errorf("too few arguments for command '%s'\nUsage: %s", name, cmd.Usage)
	}

	if cmd.MaxArgs >= 0 && len(filteredArgs) > cmd.MaxArgs {
		return fmt.Errorf("too many arguments for command '%s'\nUsage: %s", name, cmd.Usage)
	}

	// Execute the command with filtered args
	err := cmd.Handler(filteredArgs)

	// Reset options after command execution
	r.Options = CommandOptions{
		Format:  OutputFormatText,
		Verbose: false,
	}

	return err
}

// SetParkingLot sets the parking lot instance for the command registry
func (r *CommandRegistry) SetParkingLot(lot *model.ParkingLot) {
	r.parkingLot = lot
}

// GetParkingLot returns the current parking lot instance
func (r *CommandRegistry) GetParkingLot() *model.ParkingLot {
	return r.parkingLot
}

// RegisterAllCommands registers all available commands
func (r *CommandRegistry) RegisterAllCommands() {
	// Help command
	r.RegisterCommand(&Command{
		Name:        "help",
		Usage:       "help [command]",
		Description: "Show help for all commands or a specific command",
		MinArgs:     0,
		MaxArgs:     1,
		Handler:     r.handleHelp,
	})

	// Init command
	r.RegisterCommand(&Command{
		Name:        "init",
		Usage:       "init <floors> <rows> <columns>",
		Description: "Initialize a new parking lot",
		MinArgs:     3,
		MaxArgs:     3,
		Handler:     r.handleInit,
	})

	// Park command
	r.RegisterCommand(&Command{
		Name:        "park",
		Usage:       "park <vehicle_type> <vehicle_number>",
		Description: "Park a vehicle in the lot",
		MinArgs:     2,
		MaxArgs:     2,
		Handler:     r.handlePark,
	})

	// Unpark command
	r.RegisterCommand(&Command{
		Name:        "unpark",
		Usage:       "unpark <spot_id> <vehicle_number>",
		Description: "Remove a vehicle from the lot",
		MinArgs:     2,
		MaxArgs:     2,
		Handler:     r.handleUnpark,
	})

	// Available command
	r.RegisterCommand(&Command{
		Name:        "available",
		Usage:       "available <vehicle_type>",
		Description: "Display available spots for a vehicle type",
		MinArgs:     1,
		MaxArgs:     1,
		Handler:     r.handleAvailable,
	})

	// Search command
	r.RegisterCommand(&Command{
		Name:        "search",
		Usage:       "search <vehicle_number>",
		Description: "Search for a vehicle in the lot",
		MinArgs:     1,
		MaxArgs:     1,
		Handler:     r.handleSearch,
	})

	// Status command
	r.RegisterCommand(&Command{
		Name:        "status",
		Usage:       "status",
		Description: "Show the current status of the parking lot",
		MinArgs:     0,
		MaxArgs:     0,
		Handler:     r.handleStatus,
	})

	// Exit command
	r.RegisterCommand(&Command{
		Name:        "exit",
		Usage:       "exit",
		Description: "Exit the application",
		MinArgs:     0,
		MaxArgs:     0,
		Handler:     r.handleExit,
	})
}

// Command handlers

// handleHelp handles the help command
func (r *CommandRegistry) handleHelp(args []string) error {
	if len(args) == 0 {
		// Show help for all commands
		fmt.Println("Available commands:")
		fmt.Println()

		// Get and sort command names
		commandNames := make([]string, 0, len(r.Commands))
		for name := range r.Commands {
			commandNames = append(commandNames, name)
		}

		// Sort commands alphabetically
		// sort.Strings(commandNames)

		// Display all commands
		for _, name := range commandNames {
			cmd := r.Commands[name]
			fmt.Printf("  %-12s %s\n", name, cmd.Description)
		}

		fmt.Println()
		fmt.Println("Type 'help <command>' for more information about a specific command.")
	} else {
		// Show help for a specific command
		cmdName := args[0]
		cmd, found := r.GetCommand(cmdName)
		if !found {
			return fmt.Errorf("unknown command: %s", cmdName)
		}

		fmt.Printf("Command: %s\n", cmd.Name)
		fmt.Printf("Description: %s\n", cmd.Description)
		fmt.Printf("Usage: %s\n", cmd.Usage)
	}

	return nil
}

// handleInit handles the init command
func (r *CommandRegistry) handleInit(args []string) error {
	// Parse arguments
	floors, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid floors value: %s", args[0])
	}

	rows, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid rows value: %s", args[1])
	}

	columns, err := strconv.Atoi(args[2])
	if err != nil {
		return fmt.Errorf("invalid columns value: %s", args[2])
	}

	// Create the parking lot
	parkingLot, err := model.CreateParkingLot("Parking Lot", floors, rows, columns)
	if err != nil {
		return fmt.Errorf("failed to create parking lot: %v", err)
	}

	// Store the parking lot in the registry
	r.SetParkingLot(parkingLot)

	// Use formatted output
	PrintSuccess("Created parking lot with %d floors, %d rows, and %d columns",
		floors, rows, columns)
	PrintInfo("Total spots: %d", parkingLot.GetTotalSpotCount())

	// Show counts by type in a table
	counts := parkingLot.GetSpotCountByType()
	tableRows := [][]string{
		{"Bicycle", fmt.Sprintf("%d", counts[model.SpotTypeBicycle])},
		{"Motorcycle", fmt.Sprintf("%d", counts[model.SpotTypeMotorcycle])},
		{"Automobile", fmt.Sprintf("%d", counts[model.SpotTypeAutomobile])},
		{"Inactive", fmt.Sprintf("%d", counts[model.SpotTypeInactive])},
	}

	fmt.Println("Spot types:")
	fmt.Println(FormatTable([]string{"Type", "Count"}, tableRows))

	return nil
}

// handlePark handles the park command
func (r *CommandRegistry) handlePark(args []string) error {
	// Check if parking lot is initialized
	if r.parkingLot == nil {
		return fmt.Errorf("parking lot not initialized, use 'init' command first")
	}

	// Parse arguments
	vehicleTypeStr := strings.ToUpper(args[0])
	vehicleNumber := args[1]

	// Convert vehicle type
	vehicleType, err := model.ParseVehicleType(vehicleTypeStr)
	if err != nil {
		return fmt.Errorf("invalid vehicle type: %s", vehicleTypeStr)
	}

	// Try to park the vehicle
	spotID, err := r.parkingLot.Park(vehicleType, vehicleNumber)
	if err != nil {
		return fmt.Errorf("failed to park vehicle: %v", err)
	}

	PrintSuccess("Vehicle %s parked successfully at spot %s", vehicleNumber, spotID)
	return nil
}

// handleUnpark handles the unpark command
func (r *CommandRegistry) handleUnpark(args []string) error {
	// Check if parking lot is initialized
	if r.parkingLot == nil {
		return fmt.Errorf("parking lot not initialized, use 'init' command first")
	}

	// Parse arguments
	spotID := args[0]
	vehicleNumber := args[1]

	// Try to unpark the vehicle
	err := r.parkingLot.Unpark(spotID, vehicleNumber)
	if err != nil {
		return fmt.Errorf("failed to unpark vehicle: %v", err)
	}

	fmt.Printf("Vehicle %s successfully removed from spot %s\n", vehicleNumber, spotID)
	return nil
}

// handleAvailable handles the available command
func (r *CommandRegistry) handleAvailable(args []string) error {
	// Check if parking lot is initialized
	if r.parkingLot == nil {
		return fmt.Errorf("parking lot not initialized, use 'init' command first")
	}

	// Parse arguments
	vehicleTypeStr := strings.ToUpper(args[0])

	// Convert vehicle type
	vehicleType, err := model.ParseVehicleType(vehicleTypeStr)
	if err != nil {
		return fmt.Errorf("invalid vehicle type: %s", vehicleTypeStr)
	}

	// Get available spots
	spots, err := r.parkingLot.AvailableSpot(vehicleType)
	if err != nil {
		return fmt.Errorf("failed to get available spots: %v", err)
	}

	if len(spots) == 0 {
		fmt.Printf("No available spots for vehicle type %s\n", vehicleTypeStr)
		return nil
	}

	fmt.Printf("Available spots for %s:\n", model.GetVehicleTypeDisplay(vehicleType))

	// Display spots in a grid-like format
	for i, spotID := range spots {
		fmt.Printf("%s", spotID)
		if i < len(spots)-1 {
			fmt.Print(", ")
		}

		// Add a line break every 5 spots for readability
		if (i+1)%5 == 0 {
			fmt.Println()
		}
	}
	fmt.Println()

	fmt.Printf("Total available: %d\n", len(spots))
	return nil
}

// handleSearch handles the search command
func (r *CommandRegistry) handleSearch(args []string) error {
	// Check if parking lot is initialized
	if r.parkingLot == nil {
		return fmt.Errorf("parking lot not initialized, use 'init' command first")
	}

	// Parse arguments
	vehicleNumber := args[0]

	// Search for the vehicle
	spotID, isParked, err := r.parkingLot.SearchVehicle(vehicleNumber)
	if err != nil {
		// Special case for "not found" errors
		var notFoundErr *perrors.VehicleNotFoundError
		if errors.As(err, &notFoundErr) {
			fmt.Printf("Vehicle %s not found in the parking lot\n", vehicleNumber)
			return nil
		}

		return fmt.Errorf("failed to search for vehicle: %v", err)
	}

	if isParked {
		fmt.Printf("Vehicle %s is currently parked at spot %s\n", vehicleNumber, spotID)
	} else {
		fmt.Printf("Vehicle %s is not currently parked, but was last seen at spot %s\n",
			vehicleNumber, spotID)
	}

	return nil
}

// handleStatus handles the status command
func (r *CommandRegistry) handleStatus(args []string) error {
	// Check if parking lot is initialized
	if r.parkingLot == nil {
		return fmt.Errorf("parking lot not initialized, use 'init' command first")
	}

	fmt.Println(r.parkingLot.String())

	// Show counts by type
	counts := r.parkingLot.GetSpotCountByType()
	fmt.Printf("  Bicycle spots: %d\n", counts[model.SpotTypeBicycle])
	fmt.Printf("  Motorcycle spots: %d\n", counts[model.SpotTypeMotorcycle])
	fmt.Printf("  Automobile spots: %d\n", counts[model.SpotTypeAutomobile])
	fmt.Printf("  Inactive spots: %d\n", counts[model.SpotTypeInactive])

	// Show available spots by vehicle type
	availableCounts := r.parkingLot.GetAvailableSpotCountByType()
	fmt.Printf("Available spots by vehicle type:\n")
	fmt.Printf("  Bicycle: %d\n", availableCounts[model.VehicleTypeBicycle])
	fmt.Printf("  Motorcycle: %d\n", availableCounts[model.VehicleTypeMotorcycle])
	fmt.Printf("  Automobile: %d\n", availableCounts[model.VehicleTypeAutomobile])

	// Show parked vehicles
	parkedVehicles := r.parkingLot.GetAllParkedVehicles()
	fmt.Printf("Currently parked vehicles: %d\n", len(parkedVehicles))
	if len(parkedVehicles) > 0 {
		for vehicleNumber, spotID := range parkedVehicles {
			fmt.Printf("  %s at spot %s\n", vehicleNumber, spotID)
		}
	}

	return nil
}

// handleExit handles the exit command
func (r *CommandRegistry) handleExit(args []string) error {
	fmt.Println("Exiting...")
	return nil
}
