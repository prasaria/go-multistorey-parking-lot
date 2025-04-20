package model

import (
	"errors"
	"fmt"
	"sync"
)

// Error definitions of parking spots
var (
	ErrSpotAlreadyOccupied   = errors.New("parking spot is already occupied")
	ErrSpotNotOccupied       = errors.New("parking spot is not occupied")
	ErrSpotInactive          = errors.New("parking spot is inactive")
	ErrInvalidVehicleForSpot = errors.New("vehicle type not allowed in this spot")
	ErrInvalidSpotID         = errors.New("invalid spot ID format")
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
		return nil, err
	}

	// Validate location
	if floor < 0 {
		return nil, errors.New("floor number cannot be negative")
	}

	if row < 0 {
		return nil, errors.New("row number cannot be negative")
	}

	if column < 0 {
		return nil, errors.New("column number cannot be negative")
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
		return ErrSpotInactive
	}

	// Check if spot is already occupied
	if s.isOccupied {
		return ErrSpotAlreadyOccupied
	}

	// Validate vehicle number
	if err := ValidateVehicleNumber(vehicleNumber); err != nil {
		return err
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
		return ErrSpotInactive
	}

	// Check if spot is occupied
	if !s.isOccupied {
		return ErrSpotNotOccupied
	}

	// Validate the vehicle number matches
	if NormalizeVehicleNumber(vehicleNumber) != s.vehicleNumber {
		return fmt.Errorf("spot is occupied by %s, not %s",
			s.vehicleNumber, NormalizeVehicleNumber(vehicleNumber))
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
	_, err := fmt.Sscanf(spotID, "%d-%d-%d", &floor, &row, &column)
	if err != nil {
		return 0, 0, 0, ErrInvalidSpotID
	}

	// Validate bounds
	if floor < 0 {
		return 0, 0, 0, errors.New("floor number cannot be negative")
	}

	if row < 0 {
		return 0, 0, 0, errors.New("row number cannot be negative")
	}

	if column < 0 {
		return 0, 0, 0, errors.New("column number cannot be negative")
	}

	return floor, row, column, nil
}
