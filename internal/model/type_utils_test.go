package model

import (
	"testing"
)

func TestIsValidSpotCode(t *testing.T) {
	tests := []struct {
		code  string
		valid bool
	}{
		{"B-1", true},
		{"M-1", true},
		{"A-1", true},
		{"X-0", true},

		{"B-0", true},  // Valid format, but semantically a bicycle spot should be active
		{"Z-1", false}, // Invalid spot type
		{"A-2", false}, // Invalid status
		{"A1", false},  // Missing separator
		{"A-", false},  // Missing status
		{"-1", false},  // Missing type
		{"", false},    // Empty string
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			if IsValidSpotCode(tt.code) != tt.valid {
				t.Errorf("IsValidSpotCode(%s) = %v, want %v",
					tt.code, IsValidSpotCode(tt.code), tt.valid)
			}
		})
	}
}

func TestSpotCodeToSpotType(t *testing.T) {
	tests := []struct {
		code     string
		spotType SpotType
		hasError bool
	}{
		{"B-1", SpotTypeBicycle, false},
		{"M-1", SpotTypeMotorcycle, false},
		{"A-1", SpotTypeAutomobile, false},
		{"X-0", SpotTypeInactive, false},

		{"invalid", "", true},
		{"B-2", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			result, err := SpotCodeToSpotType(tt.code)

			if tt.hasError && err == nil {
				t.Errorf("SpotCodeToSpotType(%s): expected error, got nil", tt.code)
			}

			if !tt.hasError && err != nil {
				t.Errorf("SpotCodeToSpotType(%s): unexpected error: %v", tt.code, err)
			}

			if !tt.hasError && result != tt.spotType {
				t.Errorf("SpotCodeToSpotType(%s) = %s, want %s",
					tt.code, result, tt.spotType)
			}
		})
	}
}

func TestGetVehicleTypeDisplay(t *testing.T) {
	tests := []struct {
		vehicleType VehicleType
		display     string
	}{
		{VehicleTypeBicycle, "Bicycle"},
		{VehicleTypeMotorcycle, "Motorcycle"},
		{VehicleTypeAutomobile, "Automobile"},
		{"UNKNOWN", "Unknown"},
	}

	for _, tt := range tests {
		t.Run(string(tt.vehicleType), func(t *testing.T) {
			if GetVehicleTypeDisplay(tt.vehicleType) != tt.display {
				t.Errorf("GetVehicleTypeDisplay(%s) = %s, want %s",
					tt.vehicleType, GetVehicleTypeDisplay(tt.vehicleType), tt.display)
			}
		})
	}
}

func TestGetSpotTypeDisplay(t *testing.T) {
	tests := []struct {
		spotType SpotType
		display  string
	}{
		{SpotTypeBicycle, "Bicycle Spot"},
		{SpotTypeMotorcycle, "Motorcycle Spot"},
		{SpotTypeAutomobile, "Automobile Spot"},
		{SpotTypeInactive, "Inactive Spot"},
		{"UNKNOWN", "Unknown Spot"},
	}

	for _, tt := range tests {
		t.Run(string(tt.spotType), func(t *testing.T) {
			if GetSpotTypeDisplay(tt.spotType) != tt.display {
				t.Errorf("GetSpotTypeDisplay(%s) = %s, want %s",
					tt.spotType, GetSpotTypeDisplay(tt.spotType), tt.display)
			}
		})
	}
}
