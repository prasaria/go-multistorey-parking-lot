package model

import (
	"fmt"
	"sync"

	"github.com/prasaria/go-multistorey-parking-lot/internal/errors"
)

// ParkingLot represents a multi-storey parking lot
type ParkingLot struct {
	// Name of the parking lot
	Name string

	// Floors in the parking lot
	floors []*ParkingFloor

	// Map to track parked vehicles by vehicle number
	// Key: vehicle number, Value: spot ID
	parkedVehicles sync.Map

	// Map to track vehicle history
	// Key: vehicle number, Value: *VehicleHistory
	vehicleHistory sync.Map

	// Read-write mutex for thread-safety
	mu sync.RWMutex
}

// NewParkingLot creates a new parking lot with the given floors
func NewParkingLot(name string, floors []*ParkingFloor) (*ParkingLot, error) {
	if len(floors) == 0 {
		return nil, errors.NewValidationError("floors", "[]", "parking lot must have at least one floor")
	}

	if len(floors) > 8 {
		return nil, errors.NewValidationError("floors", fmt.Sprintf("%d", len(floors)),
			"parking lot cannot have more than 8 floors")
	}

	// Validate floor numbers
	floorNumbers := make(map[int]bool)
	for _, floor := range floors {
		if floor == nil {
			return nil, errors.NewValidationError("floors", "nil", "floor cannot be nil")
		}

		if floorNumbers[floor.FloorNumber] {
			return nil, errors.NewValidationError("floors",
				fmt.Sprintf("%d", floor.FloorNumber),
				"duplicate floor number")
		}

		floorNumbers[floor.FloorNumber] = true
	}

	return &ParkingLot{
		Name:   name,
		floors: floors,
	}, nil
}

// CreateParkingLot creates a new parking lot with the specified dimensions
func CreateParkingLot(name string, numFloors, rows, columns int) (*ParkingLot, error) {
	if numFloors < 1 || numFloors > 8 {
		return nil, errors.NewValidationError("numFloors",
			fmt.Sprintf("%d", numFloors),
			"number of floors must be between 1 and 8")
	}

	if rows < 1 || rows > 1000 {
		return nil, errors.NewValidationError("rows",
			fmt.Sprintf("%d", rows),
			"number of rows must be between 1 and 1000")
	}

	if columns < 1 || columns > 1000 {
		return nil, errors.NewValidationError("columns",
			fmt.Sprintf("%d", columns),
			"number of columns must be between 1 and 1000")
	}

	// Create spot layout
	layout, err := NewSpotLayout(numFloors, rows, columns)
	if err != nil {
		return nil, err
	}

	// Create floors
	floors := make([]*ParkingFloor, numFloors)
	for i := 0; i < numFloors; i++ {
		floor, err := CreateParkingFloor(i, rows, columns, layout.SpotMap[i])
		if err != nil {
			return nil, errors.WrapError(err, "CREATION_ERROR",
				fmt.Sprintf("failed to create floor %d", i))
		}
		floors[i] = floor
	}

	return NewParkingLot(name, floors)
}

// GetFloor returns the floor with the given number
func (p *ParkingLot) GetFloor(floorNum int) (*ParkingFloor, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for _, floor := range p.floors {
		if floor.FloorNumber == floorNum {
			return floor, nil
		}
	}

	return nil, errors.NewValidationError("floorNum",
		fmt.Sprintf("%d", floorNum),
		"floor not found")
}

// GetSpot returns the parking spot at the given floor, row, and column
func (p *ParkingLot) GetSpot(floorNum, row, column int) (*ParkingSpot, error) {
	floor, err := p.GetFloor(floorNum)
	if err != nil {
		return nil, err
	}

	return floor.GetSpot(row, column)
}

// GetSpotByID returns the parking spot with the given ID
func (p *ParkingLot) GetSpotByID(spotID string) (*ParkingSpot, error) {
	floor, row, column, err := ParseSpotID(spotID)
	if err != nil {
		return nil, err
	}

	return p.GetSpot(floor, row, column)
}

// GetNumFloors returns the number of floors in the parking lot
func (p *ParkingLot) GetNumFloors() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return len(p.floors)
}

// GetFloors returns all floors in the parking lot
func (p *ParkingLot) GetFloors() []*ParkingFloor {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Return a copy to prevent modification of internal state
	result := make([]*ParkingFloor, len(p.floors))
	copy(result, p.floors)

	return result
}

// GetTotalSpotCount returns the total number of parking spots in the lot
func (p *ParkingLot) GetTotalSpotCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	total := 0
	for _, floor := range p.floors {
		total += floor.GetSpotCount()
	}

	return total
}

// GetActiveSpotCount returns the number of active parking spots in the lot
func (p *ParkingLot) GetActiveSpotCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	total := 0
	for _, floor := range p.floors {
		total += floor.GetActiveSpotCount()
	}

	return total
}

// GetOccupiedSpotCount returns the number of occupied parking spots in the lot
func (p *ParkingLot) GetOccupiedSpotCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	total := 0
	for _, floor := range p.floors {
		total += floor.GetOccupiedSpotCount()
	}

	return total
}

// GetAvailableSpotCount returns the number of available parking spots in the lot
func (p *ParkingLot) GetAvailableSpotCount() int {
	return p.GetActiveSpotCount() - p.GetOccupiedSpotCount()
}

// GetSpotCountByType returns the number of spots of each type in the lot
func (p *ParkingLot) GetSpotCountByType() map[SpotType]int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	totalCounts := make(map[SpotType]int)

	// Initialize counts for all spot types
	totalCounts[SpotTypeBicycle] = 0
	totalCounts[SpotTypeMotorcycle] = 0
	totalCounts[SpotTypeAutomobile] = 0
	totalCounts[SpotTypeInactive] = 0

	// Sum counts from all floors
	for _, floor := range p.floors {
		floorCounts := floor.GetSpotCountByType()
		for spotType, count := range floorCounts {
			totalCounts[spotType] += count
		}
	}

	return totalCounts
}

// GetVehicleHistory returns the parking history for a vehicle
func (p *ParkingLot) GetVehicleHistory(vehicleNumber string) (*VehicleHistory, bool) {
	normalizedNumber := NormalizeVehicleNumber(vehicleNumber)
	historyObj, found := p.vehicleHistory.Load(normalizedNumber)
	if !found {
		return nil, false
	}

	history, ok := historyObj.(*VehicleHistory)
	return history, ok
}

// FindVehicle searches for a vehicle by number in the parking lot
// Returns the spot if found, nil otherwise
func (p *ParkingLot) FindVehicle(vehicleNumber string) (*ParkingSpot, error) {
	normalizedNumber := NormalizeVehicleNumber(vehicleNumber)

	// Check if vehicle is currently parked
	spotIDObj, found := p.parkedVehicles.Load(normalizedNumber)
	if !found {
		// If not currently parked, check if it has parking history
		history, historyFound := p.GetVehicleHistory(normalizedNumber)
		if !historyFound || history.GetLastSpotID() == "" {
			return nil, errors.NewVehicleNotFoundError(vehicleNumber)
		}

		// Return information about the last spot (vehicle was previously parked)
		return nil, errors.NewInvalidOperationError("findVehicle",
			fmt.Sprintf("vehicle %s is not currently parked, but was last seen at spot %s",
				vehicleNumber, history.GetLastSpotID()))
	}

	spotID, ok := spotIDObj.(string)
	if !ok {
		return nil, errors.NewParkingError("INTERNAL_ERROR",
			"stored spot ID has invalid type", nil)
	}

	// Get the actual spot
	return p.GetSpotByID(spotID)
}

// IsVehicleParked checks if a vehicle is currently parked
func (p *ParkingLot) IsVehicleParked(vehicleNumber string) bool {
	normalizedNumber := NormalizeVehicleNumber(vehicleNumber)
	_, found := p.parkedVehicles.Load(normalizedNumber)
	return found
}

// String returns a string representation of the parking lot
func (p *ParkingLot) String() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return fmt.Sprintf("%s: %d floors, %d total spots, %d active, %d occupied, %d available",
		p.Name, len(p.floors), p.GetTotalSpotCount(), p.GetActiveSpotCount(),
		p.GetOccupiedSpotCount(), p.GetAvailableSpotCount())
}
