package model

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/prasaria/go-multistorey-parking-lot/internal/errors"
)

func TestPark(t *testing.T) {
	lot, _ := CreateParkingLot("Park Test Lot", 2, 3, 4)

	// Park a vehicle
	vehicleNumber := "KA-01-HH-1234"
	spotID, err := lot.Park(VehicleTypeAutomobile, vehicleNumber)
	if err != nil {
		t.Fatalf("Failed to park vehicle: %v", err)
	}

	if spotID == "" {
		t.Errorf("Expected non-empty spot ID")
	}

	// Verify the vehicle is parked
	if !lot.IsVehicleParked(vehicleNumber) {
		t.Errorf("IsVehicleParked should return true for parked vehicle")
	}

	// Verify occupied count increased
	if lot.GetOccupiedSpotCount() != 1 {
		t.Errorf("Expected 1 occupied spot, got %d", lot.GetOccupiedSpotCount())
	}

	// Try to park the same vehicle again (should fail)
	_, err = lot.Park(VehicleTypeAutomobile, vehicleNumber)
	if err == nil {
		t.Errorf("Expected error when parking the same vehicle again")
	}

	// Park a different vehicle
	spotID2, err := lot.Park(VehicleTypeMotorcycle, "KA-01-HH-5678")
	if err != nil {
		t.Errorf("Failed to park second vehicle: %v", err)
	}

	if spotID2 == "" || spotID2 == spotID {
		t.Errorf("Expected different spot ID for second vehicle")
	}

	// Try to park with invalid inputs
	_, err = lot.Park("INVALID", "KA-01-HH-9999")
	if err == nil {
		t.Errorf("Expected error for invalid vehicle type")
	}

	_, err = lot.Park(VehicleTypeAutomobile, "")
	if err == nil {
		t.Errorf("Expected error for empty vehicle number")
	}

	// Park many vehicles to fill up available spots of a type
	maxSpots := 20 // Assume we have at most 20 spots
	var lastErr error

	for i := 0; i < maxSpots; i++ {
		_, lastErr = lot.Park(VehicleTypeBicycle, fmt.Sprintf("B-%04d", i))
		if lastErr != nil {
			break
		}
	}

	// Eventually we should run out of spots
	if lastErr == nil {
		t.Errorf("Expected to run out of spots")
	}

	// Check the specific error type
	if _, ok := lastErr.(*errors.NoSpaceError); !ok {
		t.Errorf("Expected NoSpaceError, got %T: %v", lastErr, lastErr)
	}
}

func TestUnpark(t *testing.T) {
	lot, _ := CreateParkingLot("Unpark Test Lot", 2, 3, 4)

	// Park a vehicle
	vehicleNumber := "KA-01-HH-1234"
	spotID, _ := lot.Park(VehicleTypeAutomobile, vehicleNumber)

	// Unpark the vehicle
	err := lot.Unpark(spotID, vehicleNumber)
	if err != nil {
		t.Fatalf("Failed to unpark vehicle: %v", err)
	}

	// Verify the vehicle is not parked anymore
	if lot.IsVehicleParked(vehicleNumber) {
		t.Errorf("IsVehicleParked should return false after unparking")
	}

	// Verify occupied count decreased
	if lot.GetOccupiedSpotCount() != 0 {
		t.Errorf("Expected 0 occupied spots after unparking, got %d",
			lot.GetOccupiedSpotCount())
	}

	// Try to unpark again (should fail)
	err = lot.Unpark(spotID, vehicleNumber)
	if err == nil {
		t.Errorf("Expected error when unparking non-parked vehicle")
	}

	// Park again
	spotID, _ = lot.Park(VehicleTypeAutomobile, vehicleNumber)

	// Try to unpark with wrong spot ID
	err = lot.Unpark("wrong-spot-id", vehicleNumber)
	if err == nil {
		t.Errorf("Expected error when unparking with wrong spot ID")
	}

	// Try to unpark with wrong vehicle number
	err = lot.Unpark(spotID, "wrong-vehicle")
	if err == nil {
		t.Errorf("Expected error when unparking with wrong vehicle number")
	}

	// Now unpark correctly
	err = lot.Unpark(spotID, vehicleNumber)
	if err != nil {
		t.Errorf("Failed to unpark vehicle: %v", err)
	}

	// Check available spots count increased back
	availableBefore := lot.GetAvailableSpotCount()
	spotID, _ = lot.Park(VehicleTypeAutomobile, vehicleNumber)
	if lot.GetAvailableSpotCount() != availableBefore-1 {
		t.Errorf("Available spot count didn't decrease after parking")
	}

	_ = lot.Unpark(spotID, vehicleNumber)
	if lot.GetAvailableSpotCount() != availableBefore {
		t.Errorf("Available spot count didn't increase after unparking")
	}
}

func TestAvailableSpot(t *testing.T) {
	lot, _ := CreateParkingLot("Available Test Lot", 2, 3, 4)

	// Check available spots for each vehicle type
	bicycleSpots, err := lot.AvailableSpot(VehicleTypeBicycle)
	if err != nil {
		t.Fatalf("Failed to get available bicycle spots: %v", err)
	}

	motorcycleSpots, err := lot.AvailableSpot(VehicleTypeMotorcycle)
	if err != nil {
		t.Fatalf("Failed to get available motorcycle spots: %v", err)
	}

	automobileSpots, err := lot.AvailableSpot(VehicleTypeAutomobile)
	if err != nil {
		t.Fatalf("Failed to get available automobile spots: %v", err)
	}

	// Verify we have at least some spots of each type
	if len(bicycleSpots) == 0 {
		t.Errorf("Expected at least one available bicycle spot")
	}

	if len(motorcycleSpots) == 0 {
		t.Errorf("Expected at least one available motorcycle spot")
	}

	if len(automobileSpots) == 0 {
		t.Errorf("Expected at least one available automobile spot")
	}

	// Park a vehicle and verify available spots decreases
	bicycleSpotsBefore := len(bicycleSpots)
	spotID, _ := lot.Park(VehicleTypeBicycle, "B-0001")

	newBicycleSpots, _ := lot.AvailableSpot(VehicleTypeBicycle)
	if len(newBicycleSpots) != bicycleSpotsBefore-1 {
		t.Errorf("Expected one less bicycle spot after parking")
	}

	// Verify the specific spot is no longer in the available list
	for _, s := range newBicycleSpots {
		if s == spotID {
			t.Errorf("Occupied spot should not be in available list")
		}
	}

	// Unpark and verify available spots increases
	_ = lot.Unpark(spotID, "B-0001")
	newBicycleSpots, _ = lot.AvailableSpot(VehicleTypeBicycle)
	if len(newBicycleSpots) != bicycleSpotsBefore {
		t.Errorf("Expected original number of bicycle spots after unparking")
	}

	// Check for invalid vehicle type
	_, err = lot.AvailableSpot("INVALID")
	if err == nil {
		t.Errorf("Expected error for invalid vehicle type")
	}
}

func TestSearchVehicle(t *testing.T) {
	lot, _ := CreateParkingLot("Search Test Lot", 2, 3, 4)

	// Search for a non-existent vehicle
	_, _, err := lot.SearchVehicle("KA-01-HH-1234")
	if err == nil {
		t.Errorf("Expected error when searching for non-existent vehicle")
	}

	// Park a vehicle
	vehicleNumber := "KA-01-HH-1234"
	spotID, _ := lot.Park(VehicleTypeAutomobile, vehicleNumber)

	// Search for the parked vehicle
	foundSpotID, isParked, err := lot.SearchVehicle(vehicleNumber)
	if err != nil {
		t.Fatalf("Failed to search for vehicle: %v", err)
	}

	if !isParked {
		t.Errorf("Expected isParked to be true for parked vehicle")
	}

	if foundSpotID != spotID {
		t.Errorf("Expected spot ID %s, got %s", spotID, foundSpotID)
	}

	// Unpark the vehicle
	_ = lot.Unpark(spotID, vehicleNumber)

	// Search for the unparked vehicle (should find last spot)
	foundSpotID, isParked, err = lot.SearchVehicle(vehicleNumber)
	if err != nil {
		t.Fatalf("Failed to search for unparked vehicle: %v", err)
	}

	if isParked {
		t.Errorf("Expected isParked to be false for unparked vehicle")
	}

	if foundSpotID != spotID {
		t.Errorf("Expected last spot ID %s, got %s", spotID, foundSpotID)
	}

	// Park and unpark in a different spot
	newSpotID, _ := lot.Park(VehicleTypeAutomobile, vehicleNumber)
	_ = lot.Unpark(newSpotID, vehicleNumber)

	// Search again, should find the most recent spot
	foundSpotID, _, err = lot.SearchVehicle(vehicleNumber)
	if err != nil {
		t.Fatalf("Failed to search for vehicle: %v", err)
	}

	if foundSpotID != newSpotID {
		t.Errorf("Expected most recent spot ID %s, got %s", newSpotID, foundSpotID)
	}
}

func TestReset(t *testing.T) {
	lot, _ := CreateParkingLot("Reset Test Lot", 2, 3, 4)

	// Park some vehicles - Check errors
	_, err := lot.Park(VehicleTypeBicycle, "B-0001")
	if err != nil {
		t.Fatalf("Failed to park bicycle: %v", err)
	}

	_, err = lot.Park(VehicleTypeMotorcycle, "M-0001")
	if err != nil {
		t.Fatalf("Failed to park motorcycle: %v", err)
	}

	_, err = lot.Park(VehicleTypeAutomobile, "A-0001")
	if err != nil {
		t.Fatalf("Failed to park automobile: %v", err)
	}

	if lot.GetOccupiedSpotCount() != 3 {
		t.Errorf("Expected 3 occupied spots before reset, got %d",
			lot.GetOccupiedSpotCount())
	}

	// Reset the lot
	lot.Reset()

	// Verify all spots are empty
	if lot.GetOccupiedSpotCount() != 0 {
		t.Errorf("Expected 0 occupied spots after reset, got %d",
			lot.GetOccupiedSpotCount())
	}

	// Verify no vehicles are parked
	if lot.GetParkedVehicleCount() != 0 {
		t.Errorf("Expected 0 parked vehicles after reset, got %d",
			lot.GetParkedVehicleCount())
	}

	// Try to search for a previously parked vehicle
	_, _, err = lot.SearchVehicle("B-0001")
	if err == nil {
		t.Errorf("Expected error when searching for vehicle after reset")
	}
}

func TestConcurrentOperations(t *testing.T) {
	lot, _ := CreateParkingLot("Concurrent Test Lot", 3, 10, 10)

	// Test concurrent parking and unparking
	const numOperations = 50
	var wg sync.WaitGroup
	wg.Add(numOperations * 2) // For both park and unpark operations

	// Track successfully parked vehicles
	var parkedVehicles sync.Map

	// Park goroutines
	for i := 0; i < numOperations; i++ {
		go func(index int) {
			defer wg.Done()

			// Try to park vehicles of different types
			var vehicleType VehicleType
			switch index % 3 {
			case 0:
				vehicleType = VehicleTypeBicycle
			case 1:
				vehicleType = VehicleTypeMotorcycle
			case 2:
				vehicleType = VehicleTypeAutomobile
			}

			vehicleNumber := fmt.Sprintf("%s-%04d", string(vehicleType)[0:1], index)
			spotID, err := lot.Park(vehicleType, vehicleNumber)

			if err == nil {
				// Successfully parked
				parkedVehicles.Store(vehicleNumber, spotID)
			}
		}(i)
	}

	// Wait a moment to give parking operations a head start
	time.Sleep(100 * time.Millisecond)

	// Unpark goroutines
	for i := 0; i < numOperations; i++ {
		go func(index int) {
			defer wg.Done()

			vehicleNumber := fmt.Sprintf("%s-%04d",
				string([]VehicleType{VehicleTypeBicycle, VehicleTypeMotorcycle, VehicleTypeAutomobile}[index%3])[0:1],
				index)

			// Try to unpark if this vehicle was parked
			if spotIDObj, found := parkedVehicles.Load(vehicleNumber); found {
				spotID := spotIDObj.(string)
				err := lot.Unpark(spotID, vehicleNumber)
				if err != nil {
					t.Logf("Failed to unpark vehicle %s: %v", vehicleNumber, err)
				}
			}
		}(i)
	}

	// Wait for all operations to complete
	wg.Wait()

	// Verify the final state is consistent
	occupiedCount := lot.GetOccupiedSpotCount()
	parkedCount := lot.GetParkedVehicleCount()

	if occupiedCount != parkedCount {
		t.Errorf("Inconsistent state: %d occupied spots but %d parked vehicles",
			occupiedCount, parkedCount)
	}

	// Verify each parked vehicle can be found at its spot
	lot.parkedVehicles.Range(func(k, v interface{}) bool {
		vehicleNumber := k.(string)
		spotID := v.(string)

		spot, err := lot.GetSpotByID(spotID)
		if err != nil {
			t.Errorf("Failed to get spot %s: %v", spotID, err)
			return true
		}

		if !spot.IsOccupied() {
			t.Errorf("Spot %s should be occupied", spotID)
		}

		if spot.GetVehicleNumber() != vehicleNumber {
			t.Errorf("Spot %s contains vehicle %s, expected %s",
				spotID, spot.GetVehicleNumber(), vehicleNumber)
		}

		return true
	})
}

// Test edge case where spots fill up and concurrent requests compete
func TestConcurrentFullParking(t *testing.T) {
	// Create a small lot with limited spots
	lot, _ := CreateParkingLot("Small Lot", 1, 2, 2)

	// Get the number of each spot type
	counts := lot.GetSpotCountByType()
	bicycleCount := counts[SpotTypeBicycle]
	motorcycleCount := counts[SpotTypeMotorcycle]
	automobileCount := counts[SpotTypeAutomobile]

	t.Logf("Spot counts - Bicycle: %d, Motorcycle: %d, Automobile: %d",
		bicycleCount, motorcycleCount, automobileCount)

	// Track success and failure counts
	var successCount, failCount int
	var countMutex sync.Mutex

	// Launch more goroutines than spots
	numVehicles := 10 // This should be more than the total available spots
	var wg sync.WaitGroup
	wg.Add(numVehicles)

	for i := 0; i < numVehicles; i++ {
		go func(index int) {
			defer wg.Done()

			// Choose a vehicle type
			var vehicleType VehicleType
			switch index % 3 {
			case 0:
				vehicleType = VehicleTypeBicycle
			case 1:
				vehicleType = VehicleTypeMotorcycle
			case 2:
				vehicleType = VehicleTypeAutomobile
			}

			vehicleNumber := fmt.Sprintf("%s-%04d", string(vehicleType)[0:1], index)
			_, err := lot.Park(vehicleType, vehicleNumber)

			countMutex.Lock()
			defer countMutex.Unlock()

			if err == nil {
				successCount++
			} else {
				failCount++
			}
		}(i)
	}

	wg.Wait()

	// Verify the total number of parked vehicles is correct
	if lot.GetOccupiedSpotCount() != successCount {
		t.Errorf("Occupied spots (%d) doesn't match success count (%d)",
			lot.GetOccupiedSpotCount(), successCount)
	}

	// Verify we didn't park more vehicles than we have spots
	activeSpots := lot.GetActiveSpotCount()
	if successCount > activeSpots {
		t.Errorf("Parked %d vehicles, but only have %d active spots",
			successCount, activeSpots)
	}

	t.Logf("Concurrent parking results - Success: %d, Failed: %d",
		successCount, failCount)
}
