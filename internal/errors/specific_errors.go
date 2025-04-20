package errors

import "fmt"

// NoSpaceError is returned when there's no available space for a vehicle
type NoSpaceError struct {
	ParkingError
	VehicleType string
}

// NewNoSpaceError creates a new NoSpaceError
func NewNoSpaceError(vehicleType string) *NoSpaceError {
	return &NoSpaceError{
		ParkingError: ParkingError{
			Code:    CodeNoSpaceAvailable,
			Message: "No available parking spot for vehicle type: " + vehicleType,
			Err:     ErrNoSpaceAvailable,
		},
		VehicleType: vehicleType,
	}
}

// InvalidOperationError is returned when an operation is invalid in the current state
type InvalidOperationError struct {
	ParkingError
	Operation string
	Reason    string
}

// NewInvalidOperationError creates a new InvalidOperationError
func NewInvalidOperationError(operation, reason string) *InvalidOperationError {
	return &InvalidOperationError{
		ParkingError: ParkingError{
			Code:    CodeInvalidOperation,
			Message: "Invalid operation '" + operation + "': " + reason,
			Err:     ErrInvalidOperation,
		},
		Operation: operation,
		Reason:    reason,
	}
}

// VehicleNotFoundError is returned when a vehicle is not found
type VehicleNotFoundError struct {
	ParkingError
	VehicleNumber string
}

// NewVehicleNotFoundError creates a new VehicleNotFoundError
func NewVehicleNotFoundError(vehicleNumber string) *VehicleNotFoundError {
	return &VehicleNotFoundError{
		ParkingError: ParkingError{
			Code:    CodeVehicleNotFound,
			Message: "Vehicle not found: " + vehicleNumber,
			Err:     ErrVehicleNotFound,
		},
		VehicleNumber: vehicleNumber,
	}
}

// SpotOccupancyError is returned for errors related to spot occupancy
type SpotOccupancyError struct {
	ParkingError
	SpotID        string
	VehicleNumber string
	IsOccupying   bool // true for occupying, false for vacating
}

// NewSpotAlreadyOccupiedError creates a new SpotOccupancyError for an already occupied spot
func NewSpotAlreadyOccupiedError(spotID string) *SpotOccupancyError {
	return &SpotOccupancyError{
		ParkingError: ParkingError{
			Code:    CodeSpotAlreadyOccupied,
			Message: "Spot already occupied: " + spotID,
			Err:     ErrSpotAlreadyOccupied,
		},
		SpotID:      spotID,
		IsOccupying: true,
	}
}

// NewSpotNotOccupiedError creates a new SpotOccupancyError for an unoccupied spot
func NewSpotNotOccupiedError(spotID string) *SpotOccupancyError {
	return &SpotOccupancyError{
		ParkingError: ParkingError{
			Code:    CodeSpotNotOccupied,
			Message: "Spot not occupied: " + spotID,
			Err:     ErrSpotNotOccupied,
		},
		SpotID:      spotID,
		IsOccupying: false,
	}
}

// NewVehicleMismatchError creates a new SpotOccupancyError for a vehicle mismatch
func NewVehicleMismatchError(spotID, expected, actual string) *SpotOccupancyError {
	return &SpotOccupancyError{
		ParkingError: ParkingError{
			Code: CodeVehicleMismatch,
			Message: fmt.Sprintf("Vehicle mismatch at spot %s: expected %s, got %s",
				spotID, expected, actual),
			Err: ErrVehicleMismatch,
		},
		SpotID:        spotID,
		VehicleNumber: actual,
		IsOccupying:   false,
	}
}

// ValidationError is returned for validation errors
type ValidationError struct {
	ParkingError
	Field string
	Value string
}

// NewValidationError creates a new ValidationError
func NewValidationError(field, value, message string) *ValidationError {
	return &ValidationError{
		ParkingError: ParkingError{
			Code:    CodeInvalidInput,
			Message: "Invalid " + field + ": " + message,
			Err:     ErrInvalidInput,
		},
		Field: field,
		Value: value,
	}
}

// NewInvalidSpotIDError creates a ValidationError for an invalid spot ID
func NewInvalidSpotIDError(spotID, reason string) *ValidationError {
	return &ValidationError{
		ParkingError: ParkingError{
			Code:    CodeInvalidSpotID,
			Message: "Invalid spot ID '" + spotID + "': " + reason,
			Err:     ErrInvalidSpotID,
		},
		Field: "spotID",
		Value: spotID,
	}
}

// NewInvalidVehicleTypeError creates a ValidationError for an invalid vehicle type
func NewInvalidVehicleTypeError(vehicleType string) *ValidationError {
	return &ValidationError{
		ParkingError: ParkingError{
			Code:    CodeInvalidVehicleType,
			Message: "Invalid vehicle type: " + vehicleType,
			Err:     ErrInvalidVehicleType,
		},
		Field: "vehicleType",
		Value: vehicleType,
	}
}

// NewInvalidVehicleNumberError creates a ValidationError for an invalid vehicle number
func NewInvalidVehicleNumberError(vehicleNumber, reason string) *ValidationError {
	return &ValidationError{
		ParkingError: ParkingError{
			Code:    CodeInvalidVehicleNumber,
			Message: "Invalid vehicle number '" + vehicleNumber + "': " + reason,
			Err:     ErrInvalidVehicleNumber,
		},
		Field: "vehicleNumber",
		Value: vehicleNumber,
	}
}

// VehicleAlreadyParkedError is returned when trying to park a vehicle that's already parked
type VehicleAlreadyParkedError struct {
	ParkingError
	VehicleNumber string
	CurrentSpotID string
}

// NewVehicleAlreadyParkedError creates a new VehicleAlreadyParkedError
func NewVehicleAlreadyParkedError(vehicleNumber, spotID string) *VehicleAlreadyParkedError {
	return &VehicleAlreadyParkedError{
		ParkingError: ParkingError{
			Code:    CodeVehicleAlreadyParked,
			Message: "Vehicle " + vehicleNumber + " is already parked at spot " + spotID,
			Err:     ErrVehicleAlreadyParked,
		},
		VehicleNumber: vehicleNumber,
		CurrentSpotID: spotID,
	}
}

// SpotTypeError is returned for errors related to spot types
type SpotTypeError struct {
	ParkingError
	SpotType    string
	VehicleType string
}

// NewInvalidSpotTypeError creates a SpotTypeError for an invalid spot type
func NewInvalidSpotTypeError(spotType string) *SpotTypeError {
	return &SpotTypeError{
		ParkingError: ParkingError{
			Code:    CodeInvalidSpotType,
			Message: "Invalid spot type: " + spotType,
			Err:     ErrInvalidSpotType,
		},
		SpotType: spotType,
	}
}

// NewSpotInactiveError creates a SpotTypeError for an inactive spot
func NewSpotInactiveError(spotID string) *SpotTypeError {
	return &SpotTypeError{
		ParkingError: ParkingError{
			Code:    CodeSpotInactive,
			Message: "Spot is inactive: " + spotID,
			Err:     ErrSpotInactive,
		},
		SpotType: "X-0",
	}
}

// NewVehicleSpotTypeMismatchError creates a SpotTypeError for a type mismatch
func NewVehicleSpotTypeMismatchError(vehicleType, spotType string) *SpotTypeError {
	return &SpotTypeError{
		ParkingError: ParkingError{
			Code:    CodeInvalidOperation,
			Message: "Vehicle type " + vehicleType + " cannot park in spot type " + spotType,
			Err:     ErrInvalidOperation,
		},
		VehicleType: vehicleType,
		SpotType:    spotType,
	}
}
