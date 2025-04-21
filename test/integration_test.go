package test

import (
	"fmt"
	"testing"

	"github.com/prasaria/go-multistorey-parking-lot/internal/cli"
	"github.com/prasaria/go-multistorey-parking-lot/internal/model"
)

// TestIntegrationScenarios tests various integrated scenarios
func TestIntegrationScenarios(t *testing.T) {
	t.Run("BasicWorkflow", func(t *testing.T) {
		// Create command registry
		registry := cli.NewCommandRegistry()
		registry.RegisterAllCommands()

		// Initialize parking lot
		err := registry.ExecuteCommand("init", []string{"2", "3", "4"})
		if err != nil {
			t.Fatalf("Failed to initialize parking lot: %v", err)
		}

		// Park vehicles
		err = registry.ExecuteCommand("park", []string{"bicycle", "B-0001"})
		if err != nil {
			t.Errorf("Failed to park bicycle: %v", err)
		}

		err = registry.ExecuteCommand("park", []string{"motorcycle", "M-0001"})
		if err != nil {
			t.Errorf("Failed to park motorcycle: %v", err)
		}

		err = registry.ExecuteCommand("park", []string{"automobile", "A-0001"})
		if err != nil {
			t.Errorf("Failed to park automobile: %v", err)
		}

		// Check available spots
		err = registry.ExecuteCommand("available", []string{"bicycle"})
		if err != nil {
			t.Errorf("Failed to get available bicycle spots: %v", err)
		}

		err = registry.ExecuteCommand("available", []string{"motorcycle"})
		if err != nil {
			t.Errorf("Failed to get available motorcycle spots: %v", err)
		}

		err = registry.ExecuteCommand("available", []string{"automobile"})
		if err != nil {
			t.Errorf("Failed to get available automobile spots: %v", err)
		}

		// Search for vehicles
		err = registry.ExecuteCommand("search", []string{"B-0001"})
		if err != nil {
			t.Errorf("Failed to search for bicycle: %v", err)
		}

		err = registry.ExecuteCommand("search", []string{"M-0001"})
		if err != nil {
			t.Errorf("Failed to search for motorcycle: %v", err)
		}

		err = registry.ExecuteCommand("search", []string{"A-0001"})
		if err != nil {
			t.Errorf("Failed to search for automobile: %v", err)
		}

		// Check status
		err = registry.ExecuteCommand("status", []string{})
		if err != nil {
			t.Errorf("Failed to check status: %v", err)
		}

		// Find where vehicles are parked to unpark them
		lot := registry.GetParkingLot()

		bSpot, err := lot.FindVehicle("B-0001")
		if err != nil {
			t.Errorf("Failed to find bicycle spot: %v", err)
		} else {
			// Unpark bicycle
			err = registry.ExecuteCommand("unpark", []string{bSpot.GetSpotID(), "B-0001"})
			if err != nil {
				t.Errorf("Failed to unpark bicycle: %v", err)
			}
		}

		mSpot, err := lot.FindVehicle("M-0001")
		if err != nil {
			t.Errorf("Failed to find motorcycle spot: %v", err)
		} else {
			// Unpark motorcycle
			err = registry.ExecuteCommand("unpark", []string{mSpot.GetSpotID(), "M-0001"})
			if err != nil {
				t.Errorf("Failed to unpark motorcycle: %v", err)
			}
		}

		aSpot, err := lot.FindVehicle("A-0001")
		if err != nil {
			t.Errorf("Failed to find automobile spot: %v", err)
		} else {
			// Unpark automobile
			err = registry.ExecuteCommand("unpark", []string{aSpot.GetSpotID(), "A-0001"})
			if err != nil {
				t.Errorf("Failed to unpark automobile: %v", err)
			}
		}

		// Final status check
		err = registry.ExecuteCommand("status", []string{})
		if err != nil {
			t.Errorf("Failed to check final status: %v", err)
		}
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		// Create command registry
		registry := cli.NewCommandRegistry()
		registry.RegisterAllCommands()

		// Initialize parking lot
		err := registry.ExecuteCommand("init", []string{"1", "2", "2"})
		if err != nil {
			t.Fatalf("Failed to initialize parking lot: %v", err)
		}

		// Test error scenarios

		// 1. Try to park with invalid vehicle type
		err = registry.ExecuteCommand("park", []string{"invalid", "V-0001"})
		if err == nil {
			t.Errorf("Expected error for invalid vehicle type")
		}

		// 2. Try to unpark non-existent vehicle
		err = registry.ExecuteCommand("unpark", []string{"0-0-0", "NONE"})
		if err == nil {
			t.Errorf("Expected error for unparking non-existent vehicle")
		}

		// 3. Try to search for non-existent vehicle
		err = registry.ExecuteCommand("search", []string{"NONE"})
		if err == nil {
			t.Errorf("Expected error for searching non-existent vehicle")
		}

		// 4. Fill the lot and test overflow
		// First, park all available spots
		lot := registry.GetParkingLot()
		availableCounts := lot.GetAvailableSpotCountByType()

		var vehicleType model.VehicleType
		var maxCount int

		for vt, count := range availableCounts {
			if count > maxCount {
				maxCount = count
				vehicleType = vt
			}
		}

		// Park vehicles up to capacity
		for i := 0; i < maxCount; i++ {
			err = registry.ExecuteCommand("park",
				[]string{string(vehicleType), fmt.Sprintf("FILL-%d", i)})
			if err != nil {
				t.Logf("Failed to park vehicle %d: %v", i, err)
				break
			}
		}

		// Try to park one more
		err = registry.ExecuteCommand("park",
			[]string{string(vehicleType), "OVERFLOW"})
		if err == nil {
			t.Errorf("Expected error when parking beyond capacity")
		}
	})
}

// TestLargeParkingLot tests scenarios with a large parking lot
func TestLargeParkingLot(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large parking lot test in short mode")
	}

	// Create command registry
	registry := cli.NewCommandRegistry()
	registry.RegisterAllCommands()

	// Initialize a large parking lot
	err := registry.ExecuteCommand("init", []string{"5", "50", "50"})
	if err != nil {
		t.Fatalf("Failed to initialize large parking lot: %v", err)
	}

	// Check dimensions and spot counts
	lot := registry.GetParkingLot()
	if lot.GetNumFloors() != 5 {
		t.Errorf("Expected 5 floors, got %d", lot.GetNumFloors())
	}

	if lot.GetTotalSpotCount() != 5*50*50 {
		t.Errorf("Expected %d total spots, got %d", 5*50*50, lot.GetTotalSpotCount())
	}

	// Park and unpark a large number of vehicles
	const numVehicles = 100
	parkedSpotIDs := make(map[string]string) // vehicleNumber -> spotID

	// Park vehicles
	for i := 0; i < numVehicles; i++ {
		vehicleType := model.VehicleTypeBicycle
		if i%3 == 1 {
			vehicleType = model.VehicleTypeMotorcycle
		} else if i%3 == 2 {
			vehicleType = model.VehicleTypeAutomobile
		}

		vehicleNumber := fmt.Sprintf("LARGE-%04d", i)
		spotID, err := lot.Park(vehicleType, vehicleNumber)
		if err != nil {
			t.Errorf("Failed to park vehicle %s: %v", vehicleNumber, err)
			continue
		}

		parkedSpotIDs[vehicleNumber] = spotID
	}

	// Verify all vehicles are parked
	for vehicleNumber, spotID := range parkedSpotIDs {
		foundSpotID, isParked, err := lot.SearchVehicle(vehicleNumber)
		if err != nil {
			t.Errorf("Failed to search for vehicle %s: %v", vehicleNumber, err)
			continue
		}

		if !isParked {
			t.Errorf("Vehicle %s should be parked", vehicleNumber)
		}

		if foundSpotID != spotID {
			t.Errorf("Expected spot ID %s, got %s", spotID, foundSpotID)
		}
	}

	// Unpark vehicles
	for vehicleNumber, spotID := range parkedSpotIDs {
		err := lot.Unpark(spotID, vehicleNumber)
		if err != nil {
			t.Errorf("Failed to unpark vehicle %s: %v", vehicleNumber, err)
		}
	}

	// Verify all vehicles are unparked
	if lot.GetOccupiedSpotCount() != 0 {
		t.Errorf("Expected 0 occupied spots after unparking, got %d",
			lot.GetOccupiedSpotCount())
	}
}
