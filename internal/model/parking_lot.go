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

// Park parks a vehicle of the given type and number in an available spot
// Returns the assigned spot ID or an error if no spot is available
func (p *ParkingLot) Park(vehicleType VehicleType, vehicleNumber string) (string, error) {
	// Validate inputs
	if vehicleType != VehicleTypeBicycle &&
		vehicleType != VehicleTypeMotorcycle &&
		vehicleType != VehicleTypeAutomobile {
		return "", errors.NewInvalidVehicleTypeError(string(vehicleType))
	}

	if err := ValidateVehicleNumber(vehicleNumber); err != nil {
		return "", err
	}

	normalizedNumber := NormalizeVehicleNumber(vehicleNumber)

	// Check if the vehicle is already parked
	if spotIDObj, found := p.parkedVehicles.Load(normalizedNumber); found {
		spotID := spotIDObj.(string)
		return "", errors.NewVehicleAlreadyParkedError(vehicleNumber, spotID)
	}

	// Find an available parking spot for this vehicle type
	p.mu.RLock()
	var availableSpot *ParkingSpot

	// Try to find a spot on each floor
	for _, floor := range p.floors {
		spots := floor.GetAvailableSpots(vehicleType)
		if len(spots) > 0 {
			availableSpot = spots[0]
			break
		}
	}
	p.mu.RUnlock()

	if availableSpot == nil {
		return "", errors.NewNoSpaceError(string(vehicleType))
	}

	// Occupy the spot
	if err := availableSpot.Occupy(normalizedNumber); err != nil {
		return "", errors.WrapError(err, "OCCUPATION_ERROR",
			fmt.Sprintf("failed to occupy spot %s", availableSpot.GetSpotID()))
	}

	// Record the parking in the maps
	spotID := availableSpot.GetSpotID()
	p.parkedVehicles.Store(normalizedNumber, spotID)

	// Update vehicle history
	historyObj, found := p.vehicleHistory.Load(normalizedNumber)
	var history *VehicleHistory

	if !found {
		// Create new vehicle
		vehicle, _ := NewVehicle(vehicleType, normalizedNumber)
		history = NewVehicleHistory(vehicle)
	} else {
		history = historyObj.(*VehicleHistory)
	}

	history.AddParkingRecord(spotID)
	p.vehicleHistory.Store(normalizedNumber, history)

	return spotID, nil
}

// Unpark removes a vehicle from its parking spot
// Returns an error if the vehicle is not parked or if the spot ID doesn't match
func (p *ParkingLot) Unpark(spotID, vehicleNumber string) error {
	// Validate inputs
	if err := ValidateVehicleNumber(vehicleNumber); err != nil {
		return err
	}

	normalizedNumber := NormalizeVehicleNumber(vehicleNumber)

	// Check if the vehicle is parked
	spotIDObj, found := p.parkedVehicles.Load(normalizedNumber)
	if !found {
		return errors.NewVehicleNotFoundError(vehicleNumber)
	}

	currentSpotID := spotIDObj.(string)
	if currentSpotID != spotID {
		return errors.NewInvalidOperationError("unpark",
			fmt.Sprintf("vehicle %s is parked at spot %s, not %s",
				vehicleNumber, currentSpotID, spotID))
	}

	// Get the spot
	spot, err := p.GetSpotByID(spotID)
	if err != nil {
		return errors.WrapError(err, "RETRIEVAL_ERROR",
			fmt.Sprintf("failed to get spot %s", spotID))
	}

	// Vacate the spot
	if err := spot.Vacate(normalizedNumber); err != nil {
		return errors.WrapError(err, "VACATION_ERROR",
			fmt.Sprintf("failed to vacate spot %s", spotID))
	}

	// Remove from parked vehicles map
	p.parkedVehicles.Delete(normalizedNumber)

	// Update vehicle history
	historyObj, found := p.vehicleHistory.Load(normalizedNumber)
	if found {
		history := historyObj.(*VehicleHistory)
		if err := history.CompleteLastParkingRecord(); err != nil {
			// Log this error but don't fail the operation
			fmt.Printf("Warning: failed to complete parking record: %v\n", err)
		}
		p.vehicleHistory.Store(normalizedNumber, history)
	}

	return nil
}

// AvailableSpot returns the list of available spot IDs for the given vehicle type
func (p *ParkingLot) AvailableSpot(vehicleType VehicleType) ([]string, error) {
	// Validate input
	if vehicleType != VehicleTypeBicycle &&
		vehicleType != VehicleTypeMotorcycle &&
		vehicleType != VehicleTypeAutomobile {
		return nil, errors.NewInvalidVehicleTypeError(string(vehicleType))
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	var availableSpots []string

	// Gather available spots from each floor
	for _, floor := range p.floors {
		spots := floor.GetAvailableSpots(vehicleType)
		for _, spot := range spots {
			availableSpots = append(availableSpots, spot.GetSpotID())
		}
	}

	return availableSpots, nil
}

// GetAvailableSpotCount returns the number of available spots for each vehicle type
func (p *ParkingLot) GetAvailableSpotCountByType() map[VehicleType]int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	counts := make(map[VehicleType]int)
	counts[VehicleTypeBicycle] = 0
	counts[VehicleTypeMotorcycle] = 0
	counts[VehicleTypeAutomobile] = 0

	// Count available spots of each type
	for _, floor := range p.floors {
		for _, vehicleType := range []VehicleType{
			VehicleTypeBicycle,
			VehicleTypeMotorcycle,
			VehicleTypeAutomobile,
		} {
			spots := floor.GetAvailableSpots(vehicleType)
			counts[vehicleType] += len(spots)
		}
	}

	return counts
}

// SearchVehicle finds a vehicle in the parking lot by its number
// Returns the spot ID where the vehicle is parked, or the last spot ID if unparked
func (p *ParkingLot) SearchVehicle(vehicleNumber string) (string, bool, error) {
	if err := ValidateVehicleNumber(vehicleNumber); err != nil {
		return "", false, err
	}

	normalizedNumber := NormalizeVehicleNumber(vehicleNumber)

	// Check if the vehicle is currently parked
	if spotIDObj, found := p.parkedVehicles.Load(normalizedNumber); found {
		return spotIDObj.(string), true, nil
	}

	// Check if the vehicle has parking history
	historyObj, found := p.vehicleHistory.Load(normalizedNumber)
	if !found {
		return "", false, errors.NewVehicleNotFoundError(vehicleNumber)
	}

	history := historyObj.(*VehicleHistory)
	lastSpotID := history.GetLastSpotID()
	if lastSpotID == "" {
		return "", false, errors.NewVehicleNotFoundError(vehicleNumber)
	}

	return lastSpotID, false, nil
}

// GetParkedVehicleCount returns the total number of vehicles currently parked
func (p *ParkingLot) GetParkedVehicleCount() int {
	count := 0
	p.parkedVehicles.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
}

// GetAllParkedVehicles returns all currently parked vehicles
func (p *ParkingLot) GetAllParkedVehicles() map[string]string {
	vehicles := make(map[string]string)
	p.parkedVehicles.Range(func(k, v interface{}) bool {
		vehicleNumber := k.(string)
		spotID := v.(string)
		vehicles[vehicleNumber] = spotID
		return true
	})
	return vehicles
}

// Reset removes all vehicles from the parking lot and clears history
func (p *ParkingLot) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Reset parking spots
	for _, floor := range p.floors {
		rows, cols := floor.GetDimensions()
		for r := 0; r < rows; r++ {
			for c := 0; c < cols; c++ {
				spot, err := floor.GetSpot(r, c)
				if err != nil {
					continue
				}

				if spot.IsOccupied() {
					// Ignore errors as we're forcefully resetting
					_ = spot.Vacate(spot.GetVehicleNumber())
				}
			}
		}
	}

	// Clear maps
	p.parkedVehicles = sync.Map{}
	p.vehicleHistory = sync.Map{}
}
