package model

import (
	"time"

	"github.com/prasaria/go-multistorey-parking-lot/internal/errors"
)

// ParkingRecord represent a single parking or unparking event
type ParkingRecord struct {
	// The spot ID where the vehicle was parked
	SpotID string

	// Timestamps for parking and unparking
	ParkedAt   time.Time
	UnparkedAt *time.Time // nil if still parked
}

// IsComplete returns true if the parking record has both parking and unparking time
func (r *ParkingRecord) IsComplete() bool {
	return r.UnparkedAt != nil
}

// Duration returns the duration for which the vehicle was parked
// If the vehicle is still parked, it returns the duration until now
func (r *ParkingRecord) Duration() time.Duration {
	if r.IsComplete() {
		return r.UnparkedAt.Sub(r.ParkedAt)
	}
	return time.Since(r.ParkedAt)
}

// VehicleHistory tracks the parking history of a vehicle
type VehicleHistory struct {
	// The vehicle that being tracked
	Vehicle *Vehicle

	// History of parking records of this vehicle
	Records []ParkingRecord
}

// NewVehicleHistory creates a new vehicle history for a vehicle
func NewVehicleHistory(vehicle *Vehicle) *VehicleHistory {
	return &VehicleHistory{
		Vehicle: vehicle,
		Records: make([]ParkingRecord, 0),
	}
}

// AddParkingRecord adds a new parking record to the history
func (h *VehicleHistory) AddParkingRecord(spotID string) {
	record := ParkingRecord{
		SpotID:     spotID,
		ParkedAt:   time.Now(),
		UnparkedAt: nil,
	}

	h.Records = append(h.Records, record)
}

// CompleteLastParkingRecord marks the last parking record as complete
func (h *VehicleHistory) CompleteLastParkingRecord() error {
	if len(h.Records) == 0 {
		return errors.NewInvalidOperationError("completeRecord",
			"no parking records exist for this vehicle")
	}

	lastIndex := len(h.Records) - 1
	if h.Records[lastIndex].IsComplete() {
		return errors.NewInvalidOperationError("completeRecord",
			"last record is already complete")
	}

	now := time.Now()
	h.Records[lastIndex].UnparkedAt = &now
	return nil
}

// GetLastParkingRecord returns the last parking record for the vehicle
// Returns nil if there is no history
func (h *VehicleHistory) GetLastParkingRecord() *ParkingRecord {
	if len(h.Records) == 0 {
		return nil
	}

	return &h.Records[len(h.Records)-1]
}

// IsCurrentlyParked returns true if the vehicle is currently parked
func (h *VehicleHistory) IsCurrentlyParked() bool {
	record := h.GetLastParkingRecord()
	return record != nil && !record.IsComplete()
}

// GetCurrentSpotID returns the current spot ID where the vehicle is parked
// Returns empty string if the vehicle is not currently parked
func (h *VehicleHistory) GetCurrentSpotID() string {
	if !h.IsCurrentlyParked() {
		return ""
	}

	return h.GetLastParkingRecord().SpotID
}

// GetLastSpotID returns the last spot ID where the vehicle was parked
// regardless of whether it's still parked or not
// Returns empty string if there's no parking history
func (h *VehicleHistory) GetLastSpotID() string {
	record := h.GetLastParkingRecord()
	if record == nil {
		return ""
	}

	return record.SpotID
}
