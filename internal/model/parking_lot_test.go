package model

import (
	"testing"
)

func TestNewParkingLot(t *testing.T) {
	// Create a simple parking lot with 2 floors
	floor0, _ := CreateParkingFloor(0, 2, 3, nil)
	floor1, _ := CreateParkingFloor(1, 2, 3, nil)

	// Test successful creation
	lot, err := NewParkingLot("Test Parking Lot", []*ParkingFloor{floor0, floor1})
	if err != nil {
		t.Fatalf("Failed to create parking lot: %v", err)
	}

	if lot.Name != "Test Parking Lot" {
		t.Errorf("Expected name 'Test Parking Lot', got '%s'", lot.Name)
	}

	if lot.GetNumFloors() != 2 {
		t.Errorf("Expected 2 floors, got %d", lot.GetNumFloors())
	}

	// Test error cases
	_, err = NewParkingLot("Empty Lot", []*ParkingFloor{})
	if err == nil {
		t.Errorf("Expected error for empty floor list")
	}

	// Test too many floors
	tooManyFloors := make([]*ParkingFloor, 9)
	for i := 0; i < 9; i++ {
		tooManyFloors[i], _ = CreateParkingFloor(i, 1, 1, nil)
	}

	_, err = NewParkingLot("Too Many Floors", tooManyFloors)
	if err == nil {
		t.Errorf("Expected error for too many floors")
	}

	// Test duplicate floor numbers
	duplicateFloors := []*ParkingFloor{floor0, floor0}
	_, err = NewParkingLot("Duplicate Floors", duplicateFloors)
	if err == nil {
		t.Errorf("Expected error for duplicate floor numbers")
	}

	// Test nil floor
	nilFloors := []*ParkingFloor{floor0, nil}
	_, err = NewParkingLot("Nil Floor", nilFloors)
	if err == nil {
		t.Errorf("Expected error for nil floor")
	}
}

func TestCreateParkingLot(t *testing.T) {
	// Test successful creation
	lot, err := CreateParkingLot("New Lot", 3, 5, 10)
	if err != nil {
		t.Fatalf("Failed to create parking lot: %v", err)
	}

	if lot.Name != "New Lot" {
		t.Errorf("Expected name 'New Lot', got '%s'", lot.Name)
	}

	if lot.GetNumFloors() != 3 {
		t.Errorf("Expected 3 floors, got %d", lot.GetNumFloors())
	}

	// Test validation
	_, err = CreateParkingLot("Invalid Floors", 0, 5, 10)
	if err == nil {
		t.Errorf("Expected error for 0 floors")
	}

	_, err = CreateParkingLot("Too Many Floors", 9, 5, 10)
	if err == nil {
		t.Errorf("Expected error for 9 floors")
	}

	_, err = CreateParkingLot("Invalid Rows", 3, 0, 10)
	if err == nil {
		t.Errorf("Expected error for 0 rows")
	}

	_, err = CreateParkingLot("Too Many Rows", 3, 1001, 10)
	if err == nil {
		t.Errorf("Expected error for 1001 rows")
	}

	_, err = CreateParkingLot("Invalid Columns", 3, 5, 0)
	if err == nil {
		t.Errorf("Expected error for 0 columns")
	}

	_, err = CreateParkingLot("Too Many Columns", 3, 5, 1001)
	if err == nil {
		t.Errorf("Expected error for 1001 columns")
	}
}

func TestGetFloorAndSpot(t *testing.T) {
	lot, _ := CreateParkingLot("Test Lot", 3, 5, 10)

	// Test getting floors
	floor, err := lot.GetFloor(1)
	if err != nil {
		t.Errorf("Failed to get floor 1: %v", err)
	}

	if floor.FloorNumber != 1 {
		t.Errorf("Expected floor number 1, got %d", floor.FloorNumber)
	}

	// Test getting non-existent floor
	_, err = lot.GetFloor(10)
	if err == nil {
		t.Errorf("Expected error for non-existent floor")
	}

	// Test getting spot
	spot, err := lot.GetSpot(1, 2, 3)
	if err != nil {
		t.Errorf("Failed to get spot: %v", err)
	}

	if spot.Floor != 1 || spot.Row != 2 || spot.Column != 3 {
		t.Errorf("Got incorrect spot: floor=%d, row=%d, column=%d",
			spot.Floor, spot.Row, spot.Column)
	}

	// Test getting spot by ID
	spotByID, err := lot.GetSpotByID("1-2-3")
	if err != nil {
		t.Errorf("Failed to get spot by ID: %v", err)
	}

	if spotByID.Floor != 1 || spotByID.Row != 2 || spotByID.Column != 3 {
		t.Errorf("Got incorrect spot by ID: floor=%d, row=%d, column=%d",
			spotByID.Floor, spotByID.Row, spotByID.Column)
	}

	// Test invalid spot ID
	_, err = lot.GetSpotByID("invalid")
	if err == nil {
		t.Errorf("Expected error for invalid spot ID")
	}

	// Test out of range spot
	_, err = lot.GetSpot(1, 10, 3)
	if err == nil {
		t.Errorf("Expected error for out of range row")
	}
}

func TestParkingLotCounts(t *testing.T) {
	lot, _ := CreateParkingLot("Count Test Lot", 2, 3, 4)

	// Test total spot count (2 floors x 3 rows x 4 columns = 24 spots)
	if lot.GetTotalSpotCount() != 24 {
		t.Errorf("Expected 24 total spots, got %d", lot.GetTotalSpotCount())
	}

	// Get active spot count (should be less than total due to inactive spots)
	activeCount := lot.GetActiveSpotCount()
	if activeCount >= 24 || activeCount <= 0 {
		t.Errorf("Active spot count %d seems incorrect", activeCount)
	}

	// Initially no spots are occupied
	if lot.GetOccupiedSpotCount() != 0 {
		t.Errorf("Expected 0 occupied spots initially, got %d", lot.GetOccupiedSpotCount())
	}

	// Available = Active - Occupied
	if lot.GetAvailableSpotCount() != activeCount {
		t.Errorf("Expected %d available spots, got %d",
			activeCount, lot.GetAvailableSpotCount())
	}

	// Check spot counts by type
	counts := lot.GetSpotCountByType()
	totalFromCounts := 0
	for _, count := range counts {
		totalFromCounts += count
	}

	if totalFromCounts != lot.GetTotalSpotCount() {
		t.Errorf("Sum of spot counts by type (%d) doesn't match total count (%d)",
			totalFromCounts, lot.GetTotalSpotCount())
	}

	// Ensure we have spots of each type
	if counts[SpotTypeBicycle] == 0 {
		t.Errorf("Expected at least one bicycle spot")
	}

	if counts[SpotTypeMotorcycle] == 0 {
		t.Errorf("Expected at least one motorcycle spot")
	}

	if counts[SpotTypeAutomobile] == 0 {
		t.Errorf("Expected at least one automobile spot")
	}
}

func TestParkingLotFindVehicle(t *testing.T) {
	lot, _ := CreateParkingLot("Find Test Lot", 2, 3, 4)

	// Initially no vehicles are parked
	if lot.IsVehicleParked("TEST-1234") {
		t.Errorf("Expected no vehicles to be parked initially")
	}

	// Finding non-existent vehicle should return error
	_, err := lot.FindVehicle("TEST-1234")
	if err == nil {
		t.Errorf("Expected error when finding non-existent vehicle")
	}

	// We'll fully test FindVehicle after implementing Park functionality
}

func TestParkingLotString(t *testing.T) {
	lot, _ := CreateParkingLot("String Test Lot", 2, 3, 4)

	str := lot.String()

	// Verify name and counts are included
	expectedParts := []string{
		"String Test Lot",
		"2 floors",
		"24 total spots",
	}

	for _, part := range expectedParts {
		if !contains(str, part) {
			t.Errorf("Expected string to contain '%s', but got: %s", part, str)
		}
	}
}
