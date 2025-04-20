package model

import (
	"errors"
	"fmt"
	"sync"
)

// ParkingFloor represents a single floor in the parking lot
type ParkingFloor struct {
	// Floor number (0-indexed)
	FloorNumber int

	// Grid of parking spots
	spots [][]*ParkingSpot

	// Dimensions
	numRows    int
	numColumns int

	// Read-write mutex for thread safety
	mu sync.RWMutex
}

// NewParkingFloor creates a new parking floor with the given spots
func NewParkingFloor(floorNumber int, spots [][]*ParkingSpot) (*ParkingFloor, error) {
	if floorNumber < 0 {
		return nil, errors.New("floor number cannot be negative")
	}

	if len(spots) == 0 {
		return nil, errors.New("floor must have at least one row")
	}

	numRows := len(spots)
	numColumns := 0

	// Verify all rows have the same number of columns
	for i, row := range spots {
		if i == 0 {
			numColumns = len(row)
			if numColumns == 0 {
				return nil, errors.New("each row must have at least one spot")
			}
		} else if len(row) != numColumns {
			return nil, errors.New("all rows must have the same number of columns")
		}
	}

	return &ParkingFloor{
		FloorNumber: floorNumber,
		spots:       spots,
		numRows:     numRows,
		numColumns:  numColumns,
	}, nil
}

// CreateParkingFloor creates a new parking floor with the given dimensions and spot types
func CreateParkingFloor(floorNumber int, rows, columns int, spotTypes [][]SpotType) (*ParkingFloor, error) {
	// Validate dimensions
	if rows <= 0 || rows > 1000 {
		return nil, fmt.Errorf("number of rows must be between 1 and 1000, got %d", rows)
	}

	if columns <= 0 || columns > 1000 {
		return nil, fmt.Errorf("number of columns must be between 1 and 1000, got %d", columns)
	}

	// Validate spot types if provided
	if spotTypes != nil {
		if len(spotTypes) != rows {
			return nil, fmt.Errorf("spot types has %d rows, expected %d", len(spotTypes), rows)
		}

		for i, row := range spotTypes {
			if len(row) != columns {
				return nil, fmt.Errorf("spot types row %d has %d columns, expected %d", i, len(row), columns)
			}
		}
	}

	// Create spots
	spots := make([][]*ParkingSpot, rows)
	for r := 0; r < rows; r++ {
		spots[r] = make([]*ParkingSpot, columns)
		for c := 0; c < columns; c++ {
			var spotType SpotType
			if spotTypes != nil {
				spotType = spotTypes[r][c]
			} else {
				// Default distribution if spot types not provided
				switch {
				case c < columns/4:
					spotType = SpotTypeBicycle
				case c < columns/2:
					spotType = SpotTypeMotorcycle
				default:
					spotType = SpotTypeAutomobile
				}
				// Make some spots inactive (e.g., for pillars)
				if (r%7 == 0 && c%7 == 0) || (r%7 == 0 && c%7 == 1) {
					spotType = SpotTypeInactive
				}
			}

			spot, err := NewParkingSpot(spotType, floorNumber, r, c)
			if err != nil {
				return nil, fmt.Errorf("failed to create spot at floor %d, row %d, column %d: %w",
					floorNumber, r, c, err)
			}
			spots[r][c] = spot
		}
	}

	return &ParkingFloor{
		FloorNumber: floorNumber,
		spots:       spots,
		numRows:     rows,
		numColumns:  columns,
	}, nil
}

// GetSpot returns the parking spot at the given row and column
func (f *ParkingFloor) GetSpot(row, column int) (*ParkingSpot, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if row < 0 || row >= f.numRows {
		return nil, fmt.Errorf("row %d out of range [0-%d]", row, f.numRows-1)
	}

	if column < 0 || column >= f.numColumns {
		return nil, fmt.Errorf("column %d out of range [0-%d]", column, f.numColumns-1)
	}

	return f.spots[row][column], nil
}

// GetAvailableSpots returns all available spots for the given vehicle type
func (f *ParkingFloor) GetAvailableSpots(vehicleType VehicleType) []*ParkingSpot {
	f.mu.RLock()
	defer f.mu.RUnlock()

	var availableSpots []*ParkingSpot

	for r := 0; r < f.numRows; r++ {
		for c := 0; c < f.numColumns; c++ {
			spot := f.spots[r][c]
			if spot.CanPark(vehicleType) {
				availableSpots = append(availableSpots, spot)
			}
		}
	}

	return availableSpots
}

// GetSpotCount returns the total number of spots on this floor
func (f *ParkingFloor) GetSpotCount() int {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return f.numRows * f.numColumns
}

// GetActiveSpotCount returns the number of active spots on this floor
func (f *ParkingFloor) GetActiveSpotCount() int {
	f.mu.RLock()
	defer f.mu.RUnlock()

	count := 0
	for r := 0; r < f.numRows; r++ {
		for c := 0; c < f.numColumns; c++ {
			if f.spots[r][c].IsActive() {
				count++
			}
		}
	}

	return count
}

// GetSpotCountByType returns the number of spots of each type on this floor
func (f *ParkingFloor) GetSpotCountByType() map[SpotType]int {
	f.mu.RLock()
	defer f.mu.RUnlock()

	counts := make(map[SpotType]int)

	// Initialize counts for all spot types
	counts[SpotTypeBicycle] = 0
	counts[SpotTypeMotorcycle] = 0
	counts[SpotTypeAutomobile] = 0
	counts[SpotTypeInactive] = 0

	// Count spots
	for r := 0; r < f.numRows; r++ {
		for c := 0; c < f.numColumns; c++ {
			spotType := f.spots[r][c].Type
			counts[spotType]++
		}
	}

	return counts
}

// GetOccupiedSpotCount returns the number of occupied spots on this floor
func (f *ParkingFloor) GetOccupiedSpotCount() int {
	f.mu.RLock()
	defer f.mu.RUnlock()

	count := 0
	for r := 0; r < f.numRows; r++ {
		for c := 0; c < f.numColumns; c++ {
			if f.spots[r][c].IsOccupied() {
				count++
			}
		}
	}

	return count
}

// FindVehicle searches for a vehicle by number on this floor
// Returns the spot if found, nil otherwise
func (f *ParkingFloor) FindVehicle(vehicleNumber string) *ParkingSpot {
	f.mu.RLock()
	defer f.mu.RUnlock()

	normalizedNumber := NormalizeVehicleNumber(vehicleNumber)

	for r := 0; r < f.numRows; r++ {
		for c := 0; c < f.numColumns; c++ {
			spot := f.spots[r][c]
			if spot.IsOccupied() && spot.GetVehicleNumber() == normalizedNumber {
				return spot
			}
		}
	}

	return nil
}

// GetNumRows returns the number of rows on this floor
func (f *ParkingFloor) GetNumRows() int {
	return f.numRows
}

// GetNumColumns returns the number of columns on this floor
func (f *ParkingFloor) GetNumColumns() int {
	return f.numColumns
}

// GetLayout returns a copy of the spot type layout for this floor
func (f *ParkingFloor) GetLayout() [][]SpotType {
	f.mu.RLock()
	defer f.mu.RUnlock()

	layout := make([][]SpotType, f.numRows)
	for r := 0; r < f.numRows; r++ {
		layout[r] = make([]SpotType, f.numColumns)
		for c := 0; c < f.numColumns; c++ {
			layout[r][c] = f.spots[r][c].Type
		}
	}

	return layout
}

// GetDisplayState returns a string grid representing the current state of the floor
// Each cell contains:
// - 'B', 'M', 'A' for available spots of each type
// - 'b', 'm', 'a' for occupied spots of each type
// - 'X' for inactive spots
func (f *ParkingFloor) GetDisplayState() [][]string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	display := make([][]string, f.numRows)
	for r := 0; r < f.numRows; r++ {
		display[r] = make([]string, f.numColumns)
		for c := 0; c < f.numColumns; c++ {
			spot := f.spots[r][c]
			switch spot.Type {
			case SpotTypeBicycle:
				if spot.IsOccupied() {
					display[r][c] = "b"
				} else {
					display[r][c] = "B"
				}
			case SpotTypeMotorcycle:
				if spot.IsOccupied() {
					display[r][c] = "m"
				} else {
					display[r][c] = "M"
				}
			case SpotTypeAutomobile:
				if spot.IsOccupied() {
					display[r][c] = "a"
				} else {
					display[r][c] = "A"
				}
			case SpotTypeInactive:
				display[r][c] = "X"
			}
		}
	}

	return display
}

// GetDimensions returns the number of rows and columns on this floor
func (f *ParkingFloor) GetDimensions() (int, int) {
	return f.numRows, f.numColumns
}

// String returns a string representation of the parking floor
func (f *ParkingFloor) String() string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return fmt.Sprintf("Floor %d (%dx%d): %d total spots, %d active, %d occupied",
		f.FloorNumber, f.numRows, f.numColumns,
		f.GetSpotCount(), f.GetActiveSpotCount(), f.GetOccupiedSpotCount())
}
