package config

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if err := config.Validate(); err != nil {
		t.Fatalf("Default configuration should be valid, got error : %v", err)
	}

	if config.Floors < 1 || config.Floors > 8 {
		t.Errorf("Default floors (%d) outside valid range", config.Floors)
	}

	if config.Rows < 1 || config.Rows > 1000 {
		t.Errorf("Default rows (%d) outside valid range", config.Rows)
	}

	if config.Columns < 1 || config.Columns > 1000 {
		t.Errorf("Default columns (%d) outside valid range", config.Columns)
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name   string
		config ParkingLotConfig
		valid  bool
	}{
		{
			name: "valid config",
			config: ParkingLotConfig{
				Floors:  3,
				Rows:    10,
				Columns: 20,
			},
			valid: true,
		},
		{
			name: "too many floors",
			config: ParkingLotConfig{
				Floors:  9,
				Rows:    10,
				Columns: 20,
			},
			valid: false,
		},
		{
			name: "negative rows",
			config: ParkingLotConfig{
				Floors:  3,
				Rows:    -1,
				Columns: 20,
			},
			valid: false,
		},
		{
			name: "too many columns",
			config: ParkingLotConfig{
				Floors:  3,
				Rows:    10,
				Columns: 1001,
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.valid && err != nil {
				t.Errorf("Expected valid config, got error: %v", err)
			}
			if !tt.valid && err == nil {
				t.Error("Expected error for invalid config, got nil")
			}
		})
	}
}

func TestParseFlagsValidInput(t *testing.T) {
	args := []string{"-floors", "4", "-rows", "15", "-columns", "25"}

	config, err := ParseInitCommand(args)
	if err != nil {
		t.Fatalf("Failed to parse valid flags: %v", err)
	}

	if config.Floors != 4 {
		t.Errorf("Expected floors = 4, got %d", config.Floors)
	}

	if config.Rows != 15 {
		t.Errorf("Expected rows = 15, got %d", config.Rows)
	}

	if config.Columns != 25 {
		t.Errorf("Expected columns = 25, got %d", config.Columns)
	}
}

func TestParseFlagsInvalidInput(t *testing.T) {
	args := []string{"-floors", "9", "-rows", "15", "-columns", "25"}

	_, err := ParseInitCommand(args)
	if err == nil {
		t.Fatal("Expected error for invalid floors, got nil")
	}
}
