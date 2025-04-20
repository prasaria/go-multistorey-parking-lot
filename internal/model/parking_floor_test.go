package model

import (
	"testing"
)

func TestNewParkingFloor(t *testing.T) {
	// Create a simple 2x3 grid of spots
	spots := make([][]*ParkingSpot, 2)
	for r := 0; r < 2; r++ {
		spots[r] = make([]*ParkingSpot, 3)
		for c := 0; c < 3; c++ {
			// Alternate between spot types
			var spotType SpotType
			switch (r + c) % 3 {
			case 0:
				spotType = SpotTypeBicycle
			case 1:
				spotType = SpotTypeMotorcycle
			case 2:
				spotType = SpotTypeAutomobile
			}
			spot, _ := NewParkingSpot(spotType, 1, r, c)
			spots[r][c] = spot
		}
	}

	floor, err := NewParkingFloor(1, spots)
	if err != nil {
		t.Fatalf("Failed to create floor: %v", err)
	}

	// Verify floor properties
	if floor.FloorNumber != 1 {
		t.Errorf("Expected floor number 1, got %d", floor.FloorNumber)
	}

	if floor.GetNumRows() != 2 {
		t.Errorf("Expected 2 rows, got %d", floor.GetNumRows())
	}

	if floor.GetNumColumns() != 3 {
		t.Errorf("Expected 3 columns, got %d", floor.GetNumColumns())
	}

	// Test error cases
	_, err = NewParkingFloor(-1, spots)
	if err == nil {
		t.Errorf("Expected error for negative floor number")
	}

	_, err = NewParkingFloor(1, nil)
	if err == nil {
		t.Errorf("Expected error for nil spots")
	}

	_, err = NewParkingFloor(1, [][]*ParkingSpot{})
	if err == nil {
		t.Errorf("Expected error for empty spots")
	}

	// Test uneven rows
	unevenSpots := make([][]*ParkingSpot, 2)
	unevenSpots[0] = make([]*ParkingSpot, 3)
	unevenSpots[1] = make([]*ParkingSpot, 2)
	_, err = NewParkingFloor(1, unevenSpots)
	if err == nil {
		t.Errorf("Expected error for uneven rows")
	}
}

func TestCreateParkingFloor(t *testing.T) {
	// Create a floor with default spot distribution
	floor, err := CreateParkingFloor(1, 5, 10, nil)
	if err != nil {
		t.Fatalf("Failed to create floor: %v", err)
	}

	// Verify dimensions
	if floor.GetNumRows() != 5 {
		t.Errorf("Expected 5 rows, got %d", floor.GetNumRows())
	}

	if floor.GetNumColumns() != 10 {
		t.Errorf("Expected 10 columns, got %d", floor.GetNumColumns())
	}

	// Create a floor with custom spot types
	customTypes := [][]SpotType{
		{SpotTypeBicycle, SpotTypeMotorcycle},
		{SpotTypeAutomobile, SpotTypeInactive},
	}

	floor, err = CreateParkingFloor(2, 2, 2, customTypes)
	if err != nil {
		t.Fatalf("Failed to create floor with custom types: %v", err)
	}

	// Verify spot types
	spot, _ := floor.GetSpot(0, 0)
	if spot.Type != SpotTypeBicycle {
		t.Errorf("Expected spot type %s, got %s", SpotTypeBicycle, spot.Type)
	}

	spot, _ = floor.GetSpot(0, 1)
	if spot.Type != SpotTypeMotorcycle {
		t.Errorf("Expected spot type %s, got %s", SpotTypeMotorcycle, spot.Type)
	}

	spot, _ = floor.GetSpot(1, 0)
	if spot.Type != SpotTypeAutomobile {
		t.Errorf("Expected spot type %s, got %s", SpotTypeAutomobile, spot.Type)
	}

	spot, _ = floor.GetSpot(1, 1)
	if spot.Type != SpotTypeInactive {
		t.Errorf("Expected spot type %s, got %s", SpotTypeInactive, spot.Type)
	}

	// Test error cases
	_, err = CreateParkingFloor(1, 0, 10, nil)
	if err == nil {
		t.Errorf("Expected error for zero rows")
	}

	_, err = CreateParkingFloor(1, 5, 0, nil)
	if err == nil {
		t.Errorf("Expected error for zero columns")
	}

	_, err = CreateParkingFloor(1, 1001, 10, nil)
	if err == nil {
		t.Errorf("Expected error for too many rows")
	}

	_, err = CreateParkingFloor(1, 5, 1001, nil)
	if err == nil {
		t.Errorf("Expected error for too many columns")
	}

	// Test mismatched spot types dimensions
	badTypes1 := [][]SpotType{
		{SpotTypeBicycle, SpotTypeMotorcycle},
		{SpotTypeAutomobile, SpotTypeInactive},
		{SpotTypeBicycle, SpotTypeMotorcycle},
	}

	_, err = CreateParkingFloor(1, 2, 2, badTypes1)
	if err == nil {
		t.Errorf("Expected error for mismatched rows in spot types")
	}

	badTypes2 := [][]SpotType{
		{SpotTypeBicycle, SpotTypeMotorcycle, SpotTypeAutomobile},
		{SpotTypeAutomobile, SpotTypeInactive},
	}

	_, err = CreateParkingFloor(1, 2, 2, badTypes2)
	if err == nil {
		t.Errorf("Expected error for mismatched columns in spot types")
	}
}

func TestGetSpot(t *testing.T) {
	floor, _ := CreateParkingFloor(1, 3, 4, nil)

	// Test valid spot
	spot, err := floor.GetSpot(1, 2)
	if err != nil {
		t.Errorf("Failed to get valid spot: %v", err)
	}

	if spot == nil {
		t.Fatalf("Expected non-nil spot")
	}

	if spot.Floor != 1 || spot.Row != 1 || spot.Column != 2 {
		t.Errorf("Spot has incorrect coordinates: expected (1,1,2), got (%d,%d,%d)",
			spot.Floor, spot.Row, spot.Column)
	}

	// Test out of range
	_, err = floor.GetSpot(-1, 0)
	if err == nil {
		t.Errorf("Expected error for negative row")
	}

	_, err = floor.GetSpot(0, -1)
	if err == nil {
		t.Errorf("Expected error for negative column")
	}

	_, err = floor.GetSpot(3, 0)
	if err == nil {
		t.Errorf("Expected error for row out of range")
	}

	_, err = floor.GetSpot(0, 4)
	if err == nil {
		t.Errorf("Expected error for column out of range")
	}
}

func TestGetAvailableSpots(t *testing.T) {
	// Create a floor with a mix of spot types
	floor, _ := CreateParkingFloor(1, 3, 3, [][]SpotType{
		{SpotTypeBicycle, SpotTypeMotorcycle, SpotTypeAutomobile},
		{SpotTypeBicycle, SpotTypeMotorcycle, SpotTypeAutomobile},
		{SpotTypeInactive, SpotTypeInactive, SpotTypeInactive},
	})

	// Test finding available spots for each vehicle type
	bicycleSpots := floor.GetAvailableSpots(VehicleTypeBicycle)
	motorcycleSpots := floor.GetAvailableSpots(VehicleTypeMotorcycle)
	autoSpots := floor.GetAvailableSpots(VehicleTypeAutomobile)

	// Bicycle can park only in bicycle spots
	if len(bicycleSpots) != 2 {
		t.Errorf("Expected 2 available spots for bicycle, got %d", len(bicycleSpots))
	}

	// Motorcycle can park only in motorcycle spots
	if len(motorcycleSpots) != 2 {
		t.Errorf("Expected 2 available spots for motorcycle, got %d", len(motorcycleSpots))
	}

	// Automobile can park only in automobile spots
	if len(autoSpots) != 2 {
		t.Errorf("Expected 2 available spots for automobile, got %d", len(autoSpots))
	}

	// Verify the spot types
	for _, spot := range bicycleSpots {
		if spot.Type != SpotTypeBicycle {
			t.Errorf("Expected bicycle spot, got %s", spot.Type)
		}
	}

	for _, spot := range motorcycleSpots {
		if spot.Type != SpotTypeMotorcycle {
			t.Errorf("Expected motorcycle spot, got %s", spot.Type)
		}
	}

	for _, spot := range autoSpots {
		if spot.Type != SpotTypeAutomobile {
			t.Errorf("Expected automobile spot, got %s", spot.Type)
		}
	}

	// Occupy a spot and verify it's no longer available
	bicycleSpot := bicycleSpots[0]
	err := bicycleSpot.Occupy("TEST-1234")
	if err != nil {
		t.Errorf("Failed to occupy spot: %v", err)
	}

	newBicycleSpots := floor.GetAvailableSpots(VehicleTypeBicycle)
	if len(newBicycleSpots) != len(bicycleSpots)-1 {
		t.Errorf("Expected one less available spot after occupying")
	}
}

func TestCountingMethods(t *testing.T) {
	// Create a floor with a mix of spot types
	floor, _ := CreateParkingFloor(1, 2, 3, [][]SpotType{
		{SpotTypeBicycle, SpotTypeMotorcycle, SpotTypeAutomobile},
		{SpotTypeBicycle, SpotTypeInactive, SpotTypeAutomobile},
	})

	// Test total spot count
	if floor.GetSpotCount() != 6 {
		t.Errorf("Expected 6 total spots, got %d", floor.GetSpotCount())
	}

	// Test active spot count
	if floor.GetActiveSpotCount() != 5 {
		t.Errorf("Expected 5 active spots, got %d", floor.GetActiveSpotCount())
	}

	// Test spot type counts
	counts := floor.GetSpotCountByType()
	if counts[SpotTypeBicycle] != 2 {
		t.Errorf("Expected 2 bicycle spots, got %d", counts[SpotTypeBicycle])
	}

	if counts[SpotTypeMotorcycle] != 1 {
		t.Errorf("Expected 1 motorcycle spot, got %d", counts[SpotTypeMotorcycle])
	}

	if counts[SpotTypeAutomobile] != 2 {
		t.Errorf("Expected 2 automobile spots, got %d", counts[SpotTypeAutomobile])
	}

	if counts[SpotTypeInactive] != 1 {
		t.Errorf("Expected 1 inactive spot, got %d", counts[SpotTypeInactive])
	}

	// Test occupied count (initially 0)
	if floor.GetOccupiedSpotCount() != 0 {
		t.Errorf("Expected 0 occupied spots initially, got %d", floor.GetOccupiedSpotCount())
	}

	// Occupy a spot
	spot, _ := floor.GetSpot(0, 0)
	_ = spot.Occupy("TEST-1234")

	// Verify occupied count increased
	if floor.GetOccupiedSpotCount() != 1 {
		t.Errorf("Expected 1 occupied spot after parking, got %d", floor.GetOccupiedSpotCount())
	}
}

func TestFindVehicle(t *testing.T) {
	// Create a floor with known spot types
	floor, _ := CreateParkingFloor(1, 2, 2, [][]SpotType{
		{SpotTypeBicycle, SpotTypeMotorcycle},
		{SpotTypeAutomobile, SpotTypeInactive},
	})

	// Initially no vehicles are parked
	spot := floor.FindVehicle("TEST-1234")
	if spot != nil {
		t.Errorf("Expected nil when searching for non-existent vehicle")
	}

	// Park a vehicle in an appropriate spot
	vehicleNumber := "KA-01-HH-1234"

	// For testing, create a vehicle of the right type for the spot
	// First get the spot at (0,0) which is a bicycle spot
	spot1, _ := floor.GetSpot(0, 0)

	// Only bicycles can park in bicycle spots
	vehicle, _ := NewVehicle(VehicleTypeBicycle, vehicleNumber)

	// Occupy the spot with the vehicle
	_ = spot1.Occupy(vehicle.Number)

	// Find the vehicle
	foundSpot := floor.FindVehicle(vehicleNumber)
	if foundSpot == nil {
		t.Fatalf("Failed to find parked vehicle")
	}

	if foundSpot.GetVehicleNumber() != NormalizeVehicleNumber(vehicleNumber) {
		t.Errorf("Found incorrect vehicle: expected %s, got %s",
			NormalizeVehicleNumber(vehicleNumber), foundSpot.GetVehicleNumber())
	}

	// Test normalization - should find vehicle regardless of case/spacing
	foundSpot = floor.FindVehicle("ka-01-hh-1234")
	if foundSpot == nil {
		t.Errorf("Failed to find vehicle with different case")
	}

	foundSpot = floor.FindVehicle("  KA-01-HH-1234  ")
	if foundSpot == nil {
		t.Errorf("Failed to find vehicle with extra spacing")
	}
}

func TestGetLayoutAndDisplayState(t *testing.T) {
	// Create a floor with known layout
	floor, _ := CreateParkingFloor(1, 2, 2, [][]SpotType{
		{SpotTypeBicycle, SpotTypeMotorcycle},
		{SpotTypeAutomobile, SpotTypeInactive},
	})

	// Get layout
	layout := floor.GetLayout()

	// Verify layout dimensions
	if len(layout) != 2 {
		t.Fatalf("Expected 2 rows in layout, got %d", len(layout))
	}

	if len(layout[0]) != 2 {
		t.Fatalf("Expected 2 columns in layout, got %d", len(layout[0]))
	}

	// Verify spot types
	if layout[0][0] != SpotTypeBicycle {
		t.Errorf("Expected spot type %s at (0,0), got %s", SpotTypeBicycle, layout[0][0])
	}

	if layout[0][1] != SpotTypeMotorcycle {
		t.Errorf("Expected spot type %s at (0,1), got %s", SpotTypeMotorcycle, layout[0][1])
	}

	if layout[1][0] != SpotTypeAutomobile {
		t.Errorf("Expected spot type %s at (1,0), got %s", SpotTypeAutomobile, layout[1][0])
	}

	if layout[1][1] != SpotTypeInactive {
		t.Errorf("Expected spot type %s at (1,1), got %s", SpotTypeInactive, layout[1][1])
	}

	// Get display state (all spots unoccupied)
	display := floor.GetDisplayState()

	// Verify display state
	expectedDisplay := [][]string{
		{"B", "M"},
		{"A", "X"},
	}

	for r := 0; r < 2; r++ {
		for c := 0; c < 2; c++ {
			if display[r][c] != expectedDisplay[r][c] {
				t.Errorf("Expected display state %s at (%d,%d), got %s",
					expectedDisplay[r][c], r, c, display[r][c])
			}
		}
	}

	// Occupy a spot
	spot, _ := floor.GetSpot(0, 0)
	_ = spot.Occupy("TEST-1234")

	// Get updated display state
	display = floor.GetDisplayState()

	// Verify occupied spot shows as lowercase
	if display[0][0] != "b" {
		t.Errorf("Expected lowercase 'b' for occupied bicycle spot, got %s", display[0][0])
	}
}

func TestGetDimensions(t *testing.T) {
	floor, _ := CreateParkingFloor(1, 3, 4, nil)

	rows, cols := floor.GetDimensions()
	if rows != 3 || cols != 4 {
		t.Errorf("Expected dimensions (3,4), got (%d,%d)", rows, cols)
	}
}

func TestParkingFloorString(t *testing.T) {
	floor, _ := CreateParkingFloor(1, 2, 3, nil)

	str := floor.String()

	// Verify floor number, dimensions, and counts are included
	expectedParts := []string{
		"Floor 1",
		"2x3",
		"6 total spots",
	}

	for _, part := range expectedParts {
		if !contains(str, part) {
			t.Errorf("Expected string to contain '%s', but got: %s", part, str)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr || len(s) > len(substr) && contains(s[1:], substr)
}
