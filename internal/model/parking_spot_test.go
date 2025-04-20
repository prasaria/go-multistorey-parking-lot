package model

import (
	"strings"
	"testing"
)

func TestNewParkingSpot(t *testing.T) {
	tests := []struct {
		name      string
		spotType  SpotType
		floor     int
		row       int
		column    int
		expectErr bool
	}{
		{
			name:      "Valid bicycle spot",
			spotType:  SpotTypeBicycle,
			floor:     1,
			row:       2,
			column:    3,
			expectErr: false,
		},
		{
			name:      "Valid motorcycle spot",
			spotType:  SpotTypeMotorcycle,
			floor:     2,
			row:       3,
			column:    4,
			expectErr: false,
		},
		{
			name:      "Valid automobile spot",
			spotType:  SpotTypeAutomobile,
			floor:     3,
			row:       4,
			column:    5,
			expectErr: false,
		},
		{
			name:      "Valid inactive spot",
			spotType:  SpotTypeInactive,
			floor:     1,
			row:       1,
			column:    1,
			expectErr: false,
		},
		{
			name:      "Invalid spot type",
			spotType:  "INVALID",
			floor:     1,
			row:       1,
			column:    1,
			expectErr: true,
		},
		{
			name:      "Negative floor",
			spotType:  SpotTypeBicycle,
			floor:     -1,
			row:       1,
			column:    1,
			expectErr: true,
		},
		{
			name:      "Negative row",
			spotType:  SpotTypeBicycle,
			floor:     1,
			row:       -1,
			column:    1,
			expectErr: true,
		},
		{
			name:      "Negative column",
			spotType:  SpotTypeBicycle,
			floor:     1,
			row:       1,
			column:    -1,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spot, err := NewParkingSpot(tt.spotType, tt.floor, tt.row, tt.column)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if spot == nil {
					t.Fatalf("Expected spot, got nil")
				}

				if spot.Type != tt.spotType {
					t.Errorf("Expected spot type %v, got %v", tt.spotType, spot.Type)
				}

				if spot.Floor != tt.floor {
					t.Errorf("Expected floor %v, got %v", tt.floor, spot.Floor)
				}

				if spot.Row != tt.row {
					t.Errorf("Expected row %v, got %v", tt.row, spot.Row)
				}

				if spot.Column != tt.column {
					t.Errorf("Expected column %v, got %v", tt.column, spot.Column)
				}

				if spot.IsOccupied() {
					t.Errorf("New spot should not be occupied")
				}

				if spot.GetVehicleNumber() != "" {
					t.Errorf("New spot should not have a vehicle number")
				}
			}
		})
	}
}

func TestGetSpotID(t *testing.T) {
	spot, _ := NewParkingSpot(SpotTypeBicycle, 1, 2, 3)
	expected := "1-2-3"

	if spot.GetSpotID() != expected {
		t.Errorf("Expected spot ID %s, got %s", expected, spot.GetSpotID())
	}
}

func TestIsActiveAndIsOccupied(t *testing.T) {
	tests := []struct {
		spotType       SpotType
		shouldBeActive bool
	}{
		{SpotTypeBicycle, true},
		{SpotTypeMotorcycle, true},
		{SpotTypeAutomobile, true},
		{SpotTypeInactive, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.spotType), func(t *testing.T) {
			spot, _ := NewParkingSpot(tt.spotType, 1, 1, 1)

			if spot.IsActive() != tt.shouldBeActive {
				t.Errorf("Spot of type %s: IsActive() = %v, want %v",
					tt.spotType, spot.IsActive(), tt.shouldBeActive)
			}

			if spot.IsOccupied() {
				t.Errorf("New spot should not be occupied")
			}
		})
	}
}

func TestCanPark(t *testing.T) {
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
		{SpotTypeMotorcycle, VehicleTypeBicycle, true},
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
			spot, _ := NewParkingSpot(tt.spotType, 1, 1, 1)

			if spot.CanPark(tt.vehicleType) != tt.canPark {
				t.Errorf("Spot of type %s: CanPark(%s) = %v, want %v",
					tt.spotType, tt.vehicleType,
					spot.CanPark(tt.vehicleType), tt.canPark)
			}

			// Try to occupy and verify CanPark returns false
			if tt.canPark {
				err := spot.Occupy("TEST-1234")
				if err != nil {
					t.Errorf("Failed to occupy spot: %v", err)
				}

				if spot.CanPark(tt.vehicleType) {
					t.Errorf("CanPark should return false for occupied spot")
				}

				// Vacate for further tests
				_ = spot.Vacate("TEST-1234")
			}
		})
	}
}

func TestOccupyAndVacate(t *testing.T) {
	spot, _ := NewParkingSpot(SpotTypeAutomobile, 1, 1, 1)
	vehicleNumber := "KA-01-HH-1234"

	// Test occupying
	err := spot.Occupy(vehicleNumber)
	if err != nil {
		t.Errorf("Failed to occupy spot: %v", err)
	}

	if !spot.IsOccupied() {
		t.Errorf("Spot should be occupied after Occupy()")
	}

	if spot.GetVehicleNumber() != NormalizeVehicleNumber(vehicleNumber) {
		t.Errorf("Expected vehicle number %s, got %s",
			NormalizeVehicleNumber(vehicleNumber), spot.GetVehicleNumber())
	}

	// Test occupying again (should fail)
	err = spot.Occupy("Another-1234")
	if err == nil {
		t.Errorf("Expected error when occupying already occupied spot")
	}

	// Test vacating with wrong vehicle number
	err = spot.Vacate("Wrong-5678")
	if err == nil {
		t.Errorf("Expected error when vacating with wrong vehicle number")
	}

	// Test vacating with correct vehicle number
	err = spot.Vacate(vehicleNumber)
	if err != nil {
		t.Errorf("Failed to vacate spot: %v", err)
	}

	if spot.IsOccupied() {
		t.Errorf("Spot should not be occupied after Vacate()")
	}

	if spot.GetVehicleNumber() != "" {
		t.Errorf("Vehicle number should be empty after Vacate()")
	}

	// Test vacating again (should fail)
	err = spot.Vacate(vehicleNumber)
	if err == nil {
		t.Errorf("Expected error when vacating unoccupied spot")
	}
}

func TestInactiveSpotOperations(t *testing.T) {
	spot, _ := NewParkingSpot(SpotTypeInactive, 1, 1, 1)
	vehicleNumber := "KA-01-HH-1234"

	// Verify inactive status
	if spot.IsActive() {
		t.Errorf("Spot should be inactive")
	}

	// Test occupying inactive spot
	err := spot.Occupy(vehicleNumber)
	if err == nil {
		t.Errorf("Expected error when occupying inactive spot")
	}

	// Test vacating inactive spot
	err = spot.Vacate(vehicleNumber)
	if err == nil {
		t.Errorf("Expected error when vacating inactive spot")
	}
}

func TestString(t *testing.T) {
	// Create spots of different types
	bicycleSpot, _ := NewParkingSpot(SpotTypeBicycle, 1, 2, 3)
	autoSpot, _ := NewParkingSpot(SpotTypeAutomobile, 2, 3, 4)
	inactiveSpot, _ := NewParkingSpot(SpotTypeInactive, 3, 4, 5)

	// Test unoccupied spot strings
	if !strings.Contains(bicycleSpot.String(), "Bicycle") ||
		!strings.Contains(bicycleSpot.String(), "1-2-3") ||
		!strings.Contains(bicycleSpot.String(), "Available") {
		t.Errorf("Bicycle spot string doesn't contain expected information: %s",
			bicycleSpot.String())
	}

	// Occupy a spot
	vehicleNumber := "KA-01-HH-1234"
	err := autoSpot.Occupy(vehicleNumber)
	if err != nil {
		t.Errorf("Failed to occupy spot: %v", err)
	}

	// Test occupied spot string
	if !strings.Contains(autoSpot.String(), "Automobile") ||
		!strings.Contains(autoSpot.String(), "2-3-4") ||
		!strings.Contains(autoSpot.String(), "Occupied") ||
		!strings.Contains(autoSpot.String(), vehicleNumber) {
		t.Errorf("Occupied automobile spot string doesn't contain expected information: %s",
			autoSpot.String())
	}

	// Test inactive spot string
	if !strings.Contains(inactiveSpot.String(), "Inactive") ||
		!strings.Contains(inactiveSpot.String(), "3-4-5") {
		t.Errorf("Inactive spot string doesn't contain expected information: %s",
			inactiveSpot.String())
	}
}

func TestParseSpotID(t *testing.T) {
	tests := []struct {
		spotID    string
		floor     int
		row       int
		column    int
		expectErr bool
	}{
		{"1-2-3", 1, 2, 3, false},
		{"0-0-0", 0, 0, 0, false},
		{"10-20-30", 10, 20, 30, false},

		// Invalid formats
		{"1-2", 0, 0, 0, true},
		{"1-2-", 0, 0, 0, true},
		{"-2-3", 0, 0, 0, true},
		{"a-2-3", 0, 0, 0, true},
		{"1-b-3", 0, 0, 0, true},
		{"1-2-c", 0, 0, 0, true},
		{"", 0, 0, 0, true},

		// Negative numbers
		{"-1-2-3", 0, 0, 0, true},
		{"1--2-3", 0, 0, 0, true},
		{"1-2--3", 0, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.spotID, func(t *testing.T) {
			floor, row, column, err := ParseSpotID(tt.spotID)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error for spotID %s, got nil", tt.spotID)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for spotID %s: %v", tt.spotID, err)
				}

				if floor != tt.floor {
					t.Errorf("Expected floor %d, got %d", tt.floor, floor)
				}

				if row != tt.row {
					t.Errorf("Expected row %d, got %d", tt.row, row)
				}

				if column != tt.column {
					t.Errorf("Expected column %d, got %d", tt.column, column)
				}
			}
		})
	}
}
