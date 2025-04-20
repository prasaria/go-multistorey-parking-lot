// internal/model/vehicle_test.go
package model

import (
	"testing"
)

func TestNewVehicle(t *testing.T) {
	tests := []struct {
		name        string
		vehicleType VehicleType
		number      string
		expectError bool
	}{
		{
			name:        "Valid bicycle",
			vehicleType: VehicleTypeBicycle,
			number:      "B-12345",
			expectError: false,
		},
		{
			name:        "Valid motorcycle",
			vehicleType: VehicleTypeMotorcycle,
			number:      "M-98765",
			expectError: false,
		},
		{
			name:        "Valid automobile",
			vehicleType: VehicleTypeAutomobile,
			number:      "KA-01-HH-1234",
			expectError: false,
		},
		{
			name:        "Invalid vehicle type",
			vehicleType: "INVALID",
			number:      "KA-01-HH-1234",
			expectError: true,
		},
		{
			name:        "Empty vehicle number",
			vehicleType: VehicleTypeAutomobile,
			number:      "",
			expectError: true,
		},
		{
			name:        "Invalid vehicle number with special chars",
			vehicleType: VehicleTypeAutomobile,
			number:      "KA*01#HH@1234",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vehicle, err := NewVehicle(tt.vehicleType, tt.number)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if vehicle == nil {
					t.Fatalf("Expected vehicle, got nil")
				}

				if vehicle.Type != tt.vehicleType {
					t.Errorf("Expected vehicle type %v, got %v", tt.vehicleType, vehicle.Type)
				}

				if vehicle.Number != NormalizeVehicleNumber(tt.number) {
					t.Errorf("Expected vehicle number %v, got %v",
						NormalizeVehicleNumber(tt.number), vehicle.Number)
				}
			}
		})
	}
}

func TestVehicleNumberNormalization(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"KA-01-HH-1234", "KA-01-HH-1234"},
		{"ka-01-hh-1234", "KA-01-HH-1234"},
		{"  KA-01-HH-1234  ", "KA-01-HH-1234"},
		{"KA 01 HH 1234", "KA 01 HH 1234"},
		{"KA  01  HH  1234", "KA 01 HH 1234"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := NormalizeVehicleNumber(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeVehicleNumber(%q) = %q, want %q",
					tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidateVehicleNumber(t *testing.T) {
	tests := []struct {
		number    string
		expectErr bool
	}{
		{"KA-01-HH-1234", false},
		{"B-12345", false},
		{"M-98765", false},
		{"123456", false},
		{"AB-1234", false},
		{"", true},
		{"    ", true},
		{"*INVALID*", true},
		{"!@#$%", true},
	}

	for _, tt := range tests {
		t.Run(tt.number, func(t *testing.T) {
			err := ValidateVehicleNumber(tt.number)

			if tt.expectErr && err == nil {
				t.Errorf("ValidateVehicleNumber(%q): expected error, got nil", tt.number)
			}

			if !tt.expectErr && err != nil {
				t.Errorf("ValidateVehicleNumber(%q): unexpected error: %v", tt.number, err)
			}
		})
	}
}

func TestVehicleEqual(t *testing.T) {
	v1, _ := NewVehicle(VehicleTypeAutomobile, "KA-01-HH-1234")
	v2, _ := NewVehicle(VehicleTypeAutomobile, "KA-01-HH-1234")
	v3, _ := NewVehicle(VehicleTypeAutomobile, "KA-01-HH-5678")
	v4, _ := NewVehicle(VehicleTypeMotorcycle, "KA-01-HH-1234")

	tests := []struct {
		name     string
		vehicle1 *Vehicle
		vehicle2 *Vehicle
		equal    bool
	}{
		{
			name:     "Same vehicle",
			vehicle1: v1,
			vehicle2: v1,
			equal:    true,
		},
		{
			name:     "Equal vehicles",
			vehicle1: v1,
			vehicle2: v2,
			equal:    true,
		},
		{
			name:     "Different numbers",
			vehicle1: v1,
			vehicle2: v3,
			equal:    false,
		},
		{
			name:     "Different types",
			vehicle1: v1,
			vehicle2: v4,
			equal:    false,
		},
		{
			name:     "Nil comparison 1",
			vehicle1: nil,
			vehicle2: v1,
			equal:    false,
		},
		{
			name:     "Nil comparison 2",
			vehicle1: v1,
			vehicle2: nil,
			equal:    false,
		},
		{
			name:     "Both nil",
			vehicle1: nil,
			vehicle2: nil,
			equal:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			equal := tt.vehicle1.Equal(tt.vehicle2)
			if equal != tt.equal {
				t.Errorf("Expected Equal() = %v, got %v", tt.equal, equal)
			}
		})
	}
}

func TestVehicleString(t *testing.T) {
	v1, _ := NewVehicle(VehicleTypeAutomobile, "KA-01-HH-1234")
	v2, _ := NewVehicle(VehicleTypeBicycle, "B-12345")
	v3, _ := NewVehicle(VehicleTypeMotorcycle, "M-98765")

	tests := []struct {
		vehicle  *Vehicle
		expected string
	}{
		{v1, "Automobile [KA-01-HH-1234]"},
		{v2, "Bicycle [B-12345]"},
		{v3, "Motorcycle [M-98765]"},
	}

	for i, tt := range tests {
		t.Run(tt.vehicle.Number, func(t *testing.T) {
			str := tt.vehicle.String()
			if str != tt.expected {
				t.Errorf("Test %d: expected %q, got %q", i, tt.expected, str)
			}
		})
	}
}
