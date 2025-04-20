package config

import "flag"

// ParseInitCommand: parses the init command flag
func ParseInitCommand(arg []string) (ParkingLotConfig, error) {
	// Start  with the default configuration
	config := DefaultConfig()

	// Create a FlagSet for parsing init command
	initCmd := flag.NewFlagSet("init", flag.ContinueOnError)

	// Define flags
	initCmd.IntVar(&config.Floors, "floors", config.Floors, "Number of floors (1-8)")
	initCmd.IntVar(&config.Rows, "rows", config.Rows, "Number of rows per floor (1-1000)")
	initCmd.IntVar(&config.Columns, "columns", config.Columns, "Number of columns per row (1-1000)")

	// Parse flags
	if err := config.Validate(); err != nil {
		return config, err
	}

	return config, nil
}

// GetUsage returns the usage string for the init command
func GetUsage() string {
	return `Usage:
  parking-lot init [options]

Options:
  -floors int    Number of floors (1-8) (default 3)
  -rows int      Number of rows per floor (1-1000) (default 5)
  -columns int   Number of columns per row (1-1000) (default 10)

Example:
  parking-lot init -floors 4 -rows 10 -columns 20`
}
