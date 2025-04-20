package model

import (
	"reflect"
	"testing"
)

func TestVehicleTypeGetPreferredSpotType(t *testing.T) {
	tests := []struct {
		vehicleType VehicleType
		spotType    SpotType
	}{
		{VehicleTypeBicycle, SpotTypeBicycle},
		{VehicleTypeMotorcycle, SpotTypeMotorcycle},
		{VehicleTypeAutomobile, SpotTypeAutomobile},
	}

	for _, tt := range tests {
		t.Run(string(tt.vehicleType), func(t *testing.T) {
			if tt.vehicleType.GetPreferredSpotType() != tt.spotType {
				t.Errorf("VehicleType %s: GetPreferredSpotType() = %s, want %s",
					tt.vehicleType, tt.vehicleType.GetPreferredSpotType(), tt.spotType)
			}
		})
	}
}

func TestVehicleTypeGetCompatibleSpotTypes(t *testing.T) {
	tests := []struct {
		vehicleType VehicleType
		spotTypes   []SpotType
	}{
		{VehicleTypeBicycle, []SpotType{SpotTypeBicycle, SpotTypeMotorcycle, SpotTypeAutomobile}},
		{VehicleTypeMotorcycle, []SpotType{SpotTypeMotorcycle, SpotTypeAutomobile}},
		{VehicleTypeAutomobile, []SpotType{SpotTypeAutomobile}},
	}

	for _, tt := range tests {
		t.Run(string(tt.vehicleType), func(t *testing.T) {
			result := tt.vehicleType.GetCompatibleSpotTypes()

			if !reflect.DeepEqual(result, tt.spotTypes) {
				t.Errorf("VehicleType %s: GetCompatibleSpotTypes() = %v, want %v",
					tt.vehicleType, result, tt.spotTypes)
			}
		})
	}
}

func TestParseVehicleType(t *testing.T) {
	tests := []struct {
		input    string
		expected VehicleType
		hasError bool
	}{
		{"BICYCLE", VehicleTypeBicycle, false},
		{"bicycle", VehicleTypeBicycle, false},
		{"B", VehicleTypeBicycle, false},
		{"BIKE", VehicleTypeBicycle, false},

		{"MOTORCYCLE", VehicleTypeMotorcycle, false},
		{"motorcycle", VehicleTypeMotorcycle, false},
		{"M", VehicleTypeMotorcycle, false},
		{"MOTORBIKE", VehicleTypeMotorcycle, false},

		{"AUTOMOBILE", VehicleTypeAutomobile, false},
		{"automobile", VehicleTypeAutomobile, false},
		{"A", VehicleTypeAutomobile, false},
		{"CAR", VehicleTypeAutomobile, false},
		{"AUTO", VehicleTypeAutomobile, false},

		{"invalid", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ParseVehicleType(tt.input)

			if tt.hasError && err == nil {
				t.Errorf("ParseVehicleType(%s): expected error, got nil", tt.input)
			}

			if !tt.hasError && err != nil {
				t.Errorf("ParseVehicleType(%s): unexpected error: %v", tt.input, err)
			}

			if !tt.hasError && result != tt.expected {
				t.Errorf("ParseVehicleType(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}
