package main

import (
	"testing"
)

func TestCommandRegistry(t *testing.T) {
	registry := NewCommandRegistry()
	registry.RegisterAllCommands()

	// Test help command
	err := registry.ExecuteCommand("help", []string{})
	if err != nil {
		t.Errorf("Failed to execute help command: %v", err)
	}

	// Test unknown command
	err = registry.ExecuteCommand("unknown", []string{})
	if err == nil {
		t.Errorf("Expected error for unknown command")
	}

	// Test init command
	err = registry.ExecuteCommand("init", []string{"2", "3", "4"})
	if err != nil {
		t.Errorf("Failed to execute init command: %v", err)
	}

	// Verify parking lot was created
	if registry.GetParkingLot() == nil {
		t.Errorf("Parking lot not initialized after init command")
	}

	// Test park command
	err = registry.ExecuteCommand("park", []string{"automobile", "TEST-1234"})
	if err != nil {
		t.Errorf("Failed to execute park command: %v", err)
	}

	// Test available command
	err = registry.ExecuteCommand("available", []string{"automobile"})
	if err != nil {
		t.Errorf("Failed to execute available command: %v", err)
	}

	// Test search command
	err = registry.ExecuteCommand("search", []string{"TEST-1234"})
	if err != nil {
		t.Errorf("Failed to execute search command: %v", err)
	}

	// Test unpark command
	// First find where the vehicle is parked
	spot, err := registry.GetParkingLot().FindVehicle("TEST-1234")
	if err != nil {
		t.Errorf("Failed to find parked vehicle: %v", err)
	}

	err = registry.ExecuteCommand("unpark", []string{spot.GetSpotID(), "TEST-1234"})
	if err != nil {
		t.Errorf("Failed to execute unpark command: %v", err)
	}

	// Test status command
	err = registry.ExecuteCommand("status", []string{})
	if err != nil {
		t.Errorf("Failed to execute status command: %v", err)
	}
}
