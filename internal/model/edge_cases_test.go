package model

import (
	"fmt"
	"testing"
)

// TestEdgeCases tests various edge cases in the parking system
func TestEdgeCases(t *testing.T) {
	t.Run("MinimalParkingLot", func(t *testing.T) {
		// Create a minimal parking lot (1 floor, 1 row, 1 column)
		lot, err := CreateParkingLot("Minimal", 1, 1, 1)
		if err != nil {
			t.Fatalf("Failed to create minimal parking lot: %v", err)
		}

		if lot.GetTotalSpotCount() != 1 {
			t.Errorf("Expected 1 spot, got %d", lot.GetTotalSpotCount())
		}

		// Verify it works for basic operations
		spotType := lot.GetSpotCountByType()

		// Try to park a vehicle of the right type
		var vehicleType VehicleType
		for st, count := range spotType {
			if count > 0 && st.IsActive() {
				switch st {
				case SpotTypeBicycle:
					vehicleType = VehicleTypeBicycle
				case SpotTypeMotorcycle:
					vehicleType = VehicleTypeMotorcycle
				case SpotTypeAutomobile:
					vehicleType = VehicleTypeAutomobile
				}
				break
			}
		}

		if vehicleType == "" {
			t.Fatalf("No active spot types in minimal lot")
		}

		spotID, err := lot.Park(vehicleType, "TEST-1")
		if err != nil {
			t.Errorf("Failed to park in minimal lot: %v", err)
		}

		// Try to park another vehicle (should fail)
		_, err = lot.Park(vehicleType, "TEST-2")
		if err == nil {
			t.Errorf("Expected error when parking in full lot")
		}

		// Unpark and verify we can park again
		err = lot.Unpark(spotID, "TEST-1")
		if err != nil {
			t.Errorf("Failed to unpark from minimal lot: %v", err)
		}

		_, err = lot.Park(vehicleType, "TEST-2")
		if err != nil {
			t.Errorf("Failed to park after unparking: %v", err)
		}
	})

	t.Run("MaximalParkingLot", func(t *testing.T) {
		// Create a large parking lot (maximum allowed dimensions)
		lot, err := CreateParkingLot("Maximal", 8, 100, 100)
		if err != nil {
			t.Fatalf("Failed to create maximal parking lot: %v", err)
		}

		expectedSpots := 8 * 100 * 100
		if lot.GetTotalSpotCount() != expectedSpots {
			t.Errorf("Expected %d spots, got %d", expectedSpots, lot.GetTotalSpotCount())
		}

		// Verify basic operations work
		spotCounts := lot.GetSpotCountByType()
		for spotType, count := range spotCounts {
			t.Logf("Spot type %s: %d", spotType, count)
		}

		// Park a few vehicles
		for i := 0; i < 5; i++ {
			vType := VehicleTypeAutomobile
			if i%3 == 1 {
				vType = VehicleTypeMotorcycle
			} else if i%3 == 2 {
				vType = VehicleTypeBicycle
			}

			_, err := lot.Park(vType, fmt.Sprintf("BIG-%d", i))
			if err != nil {
				t.Errorf("Failed to park in maximal lot: %v", err)
			}
		}

		if lot.GetOccupiedSpotCount() != 5 {
			t.Errorf("Expected 5 occupied spots, got %d", lot.GetOccupiedSpotCount())
		}
	})

	t.Run("BoundaryParkingLot", func(t *testing.T) {
		// Test boundary cases for dimensions

		// Just below maximum
		_, err := CreateParkingLot("AlmostMax", 8, 1000, 1000)
		if err != nil {
			t.Errorf("Failed to create parking lot at maximum dimensions: %v", err)
		}

		// Exceed maximum floors
		_, err = CreateParkingLot("TooManyFloors", 9, 10, 10)
		if err == nil {
			t.Errorf("Expected error for too many floors")
		}

		// Exceed maximum rows
		_, err = CreateParkingLot("TooManyRows", 3, 1001, 10)
		if err == nil {
			t.Errorf("Expected error for too many rows")
		}

		// Exceed maximum columns
		_, err = CreateParkingLot("TooManyColumns", 3, 10, 1001)
		if err == nil {
			t.Errorf("Expected error for too many columns")
		}

		// Zero dimensions
		_, err = CreateParkingLot("ZeroFloors", 0, 10, 10)
		if err == nil {
			t.Errorf("Expected error for zero floors")
		}

		_, err = CreateParkingLot("ZeroRows", 3, 0, 10)
		if err == nil {
			t.Errorf("Expected error for zero rows")
		}

		_, err = CreateParkingLot("ZeroColumns", 3, 10, 0)
		if err == nil {
			t.Errorf("Expected error for zero columns")
		}

		// Negative dimensions
		_, err = CreateParkingLot("NegativeFloors", -1, 10, 10)
		if err == nil {
			t.Errorf("Expected error for negative floors")
		}

		_, err = CreateParkingLot("NegativeRows", 3, -1, 10)
		if err == nil {
			t.Errorf("Expected error for negative rows")
		}

		_, err = CreateParkingLot("NegativeColumns", 3, 10, -1)
		if err == nil {
			t.Errorf("Expected error for negative columns")
		}
	})

	t.Run("EdgeCaseParkingOperations", func(t *testing.T) {
		lot, _ := CreateParkingLot("EdgeOps", 2, 3, 3)

		// 1. Park with empty vehicle number
		_, err := lot.Park(VehicleTypeBicycle, "")
		if err == nil {
			t.Errorf("Expected error for empty vehicle number")
		}

		// 2. Park with invalid vehicle type
		_, err = lot.Park("INVALID", "ABC-123")
		if err == nil {
			t.Errorf("Expected error for invalid vehicle type")
		}

		// 3. Unpark non-existent vehicle
		err = lot.Unpark("0-0-0", "NOT-HERE")
		if err == nil {
			t.Errorf("Expected error for unparking non-existent vehicle")
		}

		// 4. Unpark with invalid spot ID
		err = lot.Unpark("invalid", "ABC-123")
		if err == nil {
			t.Errorf("Expected error for invalid spot ID")
		}

		// 5. Search for non-existent vehicle
		_, _, err = lot.SearchVehicle("NOT-HERE")
		if err == nil {
			t.Errorf("Expected error for searching non-existent vehicle")
		}

		// 6. Park with invalid characters in vehicle number
		_, err = lot.Park(VehicleTypeBicycle, "!@#$%^")
		if err == nil {
			t.Errorf("Expected error for invalid characters in vehicle number")
		}

		// 7. Double parking attempt
		spotID, _ := lot.Park(VehicleTypeBicycle, "DOUBLE-1")
		_, err = lot.Park(VehicleTypeBicycle, "DOUBLE-1")
		if err == nil {
			t.Errorf("Expected error for double parking")
		}

		// 8. Unpark with wrong vehicle number
		err = lot.Unpark(spotID, "WRONG-NUM")
		if err == nil {
			t.Errorf("Expected error for unparking with wrong vehicle number")
		}

		// 9. Try to park more vehicles than spots of a certain type
		// First, count how many spots of each type we have
		counts := lot.GetSpotCountByType()
		bicycleCount := counts[SpotTypeBicycle]

		// Park bicycles up to capacity
		parkedIDs := make([]string, 0)
		for i := 0; i < bicycleCount; i++ {
			id, err := lot.Park(VehicleTypeBicycle, fmt.Sprintf("BIKE-%d", i))
			if err != nil {
				t.Logf("Failed to park bicycle %d: %v", i, err)
				break
			}
			parkedIDs = append(parkedIDs, id)
		}

		// Now try to park one more bicycle
		_, err = lot.Park(VehicleTypeBicycle, "ONE-TOO-MANY")
		if err == nil {
			t.Errorf("Expected error when parking too many bicycles")
		}

		// 10. Verify all our parked bicycles are actually parked
		for i, id := range parkedIDs {
			vehicleNumber := fmt.Sprintf("BIKE-%d", i)
			foundID, isParked, err := lot.SearchVehicle(vehicleNumber)
			if err != nil {
				t.Errorf("Failed to find parked bicycle %s: %v", vehicleNumber, err)
			}
			if !isParked {
				t.Errorf("Expected bicycle %s to be parked", vehicleNumber)
			}
			if foundID != id {
				t.Errorf("Expected spot ID %s, got %s", id, foundID)
			}
		}
	})
}
