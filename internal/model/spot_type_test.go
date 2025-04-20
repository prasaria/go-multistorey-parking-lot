package model

import (
	"testing"
)

func TestSpotTypeIsActive(t *testing.T) {
	tests := []struct {
		spotType SpotType
		active   bool
	}{
		{SpotTypeBicycle, true},
		{SpotTypeMotorcycle, true},
		{SpotTypeAutomobile, true},
		{SpotTypeInactive, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.spotType), func(t *testing.T) {
			if tt.spotType.IsActive() != tt.active {
				t.Errorf("SpotType %s: IsActive() = %v, want %v",
					tt.spotType, tt.spotType.IsActive(), tt.active)
			}
		})
	}
}

func TestSpotTypeCanParkVehicleType(t *testing.T) {
	tests := []struct {
		spotType    SpotType
		vehicleType VehicleType
		canPark     bool
	}{
		// Bicycle spots
		{SpotTypeBicycle, VehicleTypeBicycle, true},
		{SpotTypeBicycle, VehicleTypeMotorcycle, false},
		{SpotTypeBicycle, VehicleTypeAutomobile, false},

		// Motorcycle spots
		{SpotTypeMotorcycle, VehicleTypeBicycle, false},
		{SpotTypeMotorcycle, VehicleTypeMotorcycle, true},
		{SpotTypeMotorcycle, VehicleTypeAutomobile, false},

		// Automobile spots
		{SpotTypeAutomobile, VehicleTypeBicycle, false},
		{SpotTypeAutomobile, VehicleTypeMotorcycle, false},
		{SpotTypeAutomobile, VehicleTypeAutomobile, true},

		// Inactive spots
		{SpotTypeInactive, VehicleTypeBicycle, false},
		{SpotTypeInactive, VehicleTypeMotorcycle, false},
		{SpotTypeInactive, VehicleTypeAutomobile, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.spotType)+"_"+string(tt.vehicleType), func(t *testing.T) {
			if tt.spotType.CanParkVehicleType(tt.vehicleType) != tt.canPark {
				t.Errorf("SpotType %s: CanParkVehicleType(%s) = %v, want %v",
					tt.spotType, tt.vehicleType,
					tt.spotType.CanParkVehicleType(tt.vehicleType), tt.canPark)
			}
		})
	}
}

func TestParseSpotType(t *testing.T) {
	tests := []struct {
		input    string
		expected SpotType
		hasError bool
	}{
		{"B-1", SpotTypeBicycle, false},
		{"M-1", SpotTypeMotorcycle, false},
		{"A-1", SpotTypeAutomobile, false},
		{"X-0", SpotTypeInactive, false},
		{"B-0", "", true},
		{"Z-1", "", true},
		{"invalid", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ParseSpotType(tt.input)

			if tt.hasError && err == nil {
				t.Errorf("ParseSpotType(%s): expected error, got nil", tt.input)
			}

			if !tt.hasError && err != nil {
				t.Errorf("ParseSpotType(%s): unexpected error: %v", tt.input, err)
			}

			if !tt.hasError && result != tt.expected {
				t.Errorf("ParseSpotType(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}
