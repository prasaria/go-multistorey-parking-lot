package model

import (
	"fmt"
	"sync"

	"github.com/prasaria/go-multistorey-parking-lot/internal/errors"
)

// ParkingSpot represents a single parking space in the parking lot
type ParkingSpot struct {
	// The type of spot (B-1, M-1, A-1, X-0)
	Type SpotType

	// Location information
	Floor  int
	Row    int
	Column int

	// State information
	isOccupied    bool
	vehicleNumber string

	// Mutex for thread-safety
	mu sync.RWMutex
}

// NewParkingSpot creates a new parking spot
func NewParkingSpot(spotType SpotType, floor, row, column int) (*ParkingSpot, error) {
	// Validate spot type
	_, err := ParseSpotType(string(spotType))
	if err != nil {
		return nil, errors.NewInvalidSpotTypeError(string(spotType))
	}

	// Validate location
	if floor < 0 {
		return nil, errors.NewValidationError("floor", fmt.Sprintf("%d", floor), "floor number cannot be negative")
	}

	if row < 0 {
		return nil, errors.NewValidationError("row", fmt.Sprintf("%d", row), "row number cannot be negative")
	}

	if column < 0 {
		return nil, errors.NewValidationError("column", fmt.Sprintf("%d", column), "column number cannot be negative")
	}

	return &ParkingSpot{
		Type:          spotType,
		Floor:         floor,
		Row:           row,
		Column:        column,
		isOccupied:    false,
		vehicleNumber: "",
	}, nil
}

// GetSpotID returns the ID of the spot in the format "floor-row-column"
func (s *ParkingSpot) GetSpotID() string {
	return fmt.Sprintf("%d-%d-%d", s.Floor, s.Row, s.Column)
}

// IsActive returns true if the spot is active
func (s *ParkingSpot) IsActive() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.Type.IsActive()
}

// IsOccupied returns true if the spot is currently occupied
func (s *ParkingSpot) IsOccupied() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.isOccupied
}

// GetVehicleNumber returns the number of the vehicle occupying this spot
// Returns empty string if spot is unoccupied
func (s *ParkingSpot) GetVehicleNumber() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.vehicleNumber
}

// CanPark checks if a vehicle of given type can park in this spot
func (s *ParkingSpot) CanPark(vehicleType VehicleType) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Must be active and unoccupied
	if !s.Type.IsActive() || s.isOccupied {
		return false
	}

	// Check if vehicle type can park in this spot type
	return s.Type.CanParkVehicleType(vehicleType)
}

// Occupy marks the spot as occupied by the given vehicle
func (s *ParkingSpot) Occupy(vehicleNumber string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if spot is active
	if !s.Type.IsActive() {
		return errors.NewSpotInactiveError(s.GetSpotID())
	}

	// Check if spot is already occupied
	if s.isOccupied {
		return errors.NewSpotAlreadyOccupiedError(s.GetSpotID())
	}

	// Validate vehicle number
	if err := ValidateVehicleNumber(vehicleNumber); err != nil {
		return errors.NewInvalidVehicleNumberError(vehicleNumber, err.Error())
	}

	// Mark as occupied
	s.isOccupied = true
	s.vehicleNumber = NormalizeVehicleNumber(vehicleNumber)

	return nil
}

// Vacate marks the spot as unoccupied
// Returns error if the spot is not occupied by the specified vehicle
func (s *ParkingSpot) Vacate(vehicleNumber string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if spot is active
	if !s.Type.IsActive() {
		return errors.NewSpotInactiveError(s.GetSpotID())
	}

	// Check if spot is occupied
	if !s.isOccupied {
		return errors.NewSpotNotOccupiedError(s.GetSpotID())
	}

	// Validate the vehicle number matches
	normalizedNumber := NormalizeVehicleNumber(vehicleNumber)
	if normalizedNumber != s.vehicleNumber {
		return errors.NewVehicleMismatchError(s.GetSpotID(), s.vehicleNumber, normalizedNumber)
	}

	// Mark as unoccupied
	s.isOccupied = false
	s.vehicleNumber = ""

	return nil
}

// String returns a string representation of the parking spot
func (s *ParkingSpot) String() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	spotType := GetSpotTypeDisplay(s.Type)
	location := s.GetSpotID()

	if s.isOccupied {
		return fmt.Sprintf("%s at %s (Occupied by %s)", spotType, location, s.vehicleNumber)
	}

	if s.Type.IsActive() {
		return fmt.Sprintf("%s at %s (Available)", spotType, location)
	}

	return fmt.Sprintf("%s at %s", spotType, location)
}

// ParseSpotID parses a spot ID in the format "floor-row-column"
// Returns floor, row, column and error
func ParseSpotID(spotID string) (int, int, int, error) {
	var floor, row, column int
	count, err := fmt.Sscanf(spotID, "%d-%d-%d", &floor, &row, &column)

	if err != nil || count != 3 {
		return 0, 0, 0, errors.NewInvalidSpotIDError(spotID, "invalid format, expected floor-row-column")
	}

	// Validate bounds
	if floor < 0 {
		return 0, 0, 0, errors.NewInvalidSpotIDError(spotID, "floor number cannot be negative")
	}

	if row < 0 {
		return 0, 0, 0, errors.NewInvalidSpotIDError(spotID, "row number cannot be negative")
	}

	if column < 0 {
		return 0, 0, 0, errors.NewInvalidSpotIDError(spotID, "column number cannot be negative")
	}

	return floor, row, column, nil
}
