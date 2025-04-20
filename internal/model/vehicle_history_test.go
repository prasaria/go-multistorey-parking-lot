package model

import (
	"testing"
	"time"
)

func TestParkingRecordIsComplete(t *testing.T) {
	now := time.Now()
	later := now.Add(1 * time.Hour)

	tests := []struct {
		name     string
		record   ParkingRecord
		expected bool
	}{
		{
			name: "Incomplete record",
			record: ParkingRecord{
				SpotID:     "1-1-1",
				ParkedAt:   now,
				UnparkedAt: nil,
			},
			expected: false,
		},
		{
			name: "Complete record",
			record: ParkingRecord{
				SpotID:     "1-1-1",
				ParkedAt:   now,
				UnparkedAt: &later,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.record.IsComplete() != tt.expected {
				t.Errorf("Expected IsComplete() = %v, got %v",
					tt.expected, tt.record.IsComplete())
			}
		})
	}
}

func TestParkingRecordDuration(t *testing.T) {
	now := time.Now()
	later := now.Add(1 * time.Hour)

	tests := []struct {
		name     string
		record   ParkingRecord
		expected time.Duration
	}{
		{
			name: "Complete record",
			record: ParkingRecord{
				SpotID:     "1-1-1",
				ParkedAt:   now,
				UnparkedAt: &later,
			},
			expected: 1 * time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			duration := tt.record.Duration()
			if duration != tt.expected {
				t.Errorf("Expected Duration() = %v, got %v", tt.expected, duration)
			}
		})
	}

	// Test ongoing parking (can't use exact duration comparison)
	t.Run("Ongoing parking", func(t *testing.T) {
		record := ParkingRecord{
			SpotID:     "1-1-1",
			ParkedAt:   time.Now().Add(-10 * time.Second),
			UnparkedAt: nil,
		}

		duration := record.Duration()
		if duration < 9*time.Second || duration > 11*time.Second {
			t.Errorf("Expected Duration() to be around 10 seconds, got %v", duration)
		}
	})
}

func TestVehicleHistory(t *testing.T) {
	vehicle, _ := NewVehicle(VehicleTypeAutomobile, "KA-01-HH-1234")
	history := NewVehicleHistory(vehicle)

	// Test initial state
	if history.Vehicle != vehicle {
		t.Errorf("Expected vehicle to be set correctly")
	}

	if len(history.Records) != 0 {
		t.Errorf("Expected empty records initially, got %d", len(history.Records))
	}

	if history.IsCurrentlyParked() {
		t.Errorf("New vehicle history should not be marked as parked")
	}

	if history.GetCurrentSpotID() != "" {
		t.Errorf("Expected empty spot ID for new history, got %s", history.GetCurrentSpotID())
	}

	if history.GetLastSpotID() != "" {
		t.Errorf("Expected empty last spot ID for new history, got %s", history.GetLastSpotID())
	}

	// Add parking record
	history.AddParkingRecord("1-2-3")

	if !history.IsCurrentlyParked() {
		t.Errorf("Vehicle should be marked as parked after adding record")
	}

	if history.GetCurrentSpotID() != "1-2-3" {
		t.Errorf("Expected current spot ID '1-2-3', got %s", history.GetCurrentSpotID())
	}

	if history.GetLastSpotID() != "1-2-3" {
		t.Errorf("Expected last spot ID '1-2-3', got %s", history.GetLastSpotID())
	}

	// Complete parking record
	if !history.CompleteLastParkingRecord() {
		t.Errorf("Expected to successfully complete last parking record")
	}

	if history.IsCurrentlyParked() {
		t.Errorf("Vehicle should not be marked as parked after completing record")
	}

	if history.GetCurrentSpotID() != "" {
		t.Errorf("Expected empty current spot ID, got %s", history.GetCurrentSpotID())
	}

	if history.GetLastSpotID() != "1-2-3" {
		t.Errorf("Expected last spot ID '1-2-3', got %s", history.GetLastSpotID())
	}

	// Try to complete again (should fail)
	if history.CompleteLastParkingRecord() {
		t.Errorf("Should not be able to complete an already completed record")
	}

	// Add another parking record
	history.AddParkingRecord("2-3-4")

	if !history.IsCurrentlyParked() {
		t.Errorf("Vehicle should be marked as parked after adding new record")
	}

	if history.GetCurrentSpotID() != "2-3-4" {
		t.Errorf("Expected current spot ID '2-3-4', got %s", history.GetCurrentSpotID())
	}

	if history.GetLastSpotID() != "2-3-4" {
		t.Errorf("Expected last spot ID '2-3-4', got %s", history.GetLastSpotID())
	}

	// Complete the second record
	if !history.CompleteLastParkingRecord() {
		t.Errorf("Expected to successfully complete last parking record")
	}

	// Verify record count
	if len(history.Records) != 2 {
		t.Errorf("Expected 2 records, got %d", len(history.Records))
	}

	// Verify all records are complete
	for i, record := range history.Records {
		if !record.IsComplete() {
			t.Errorf("Expected record %d to be complete", i)
		}
	}
}
