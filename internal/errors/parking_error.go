package errors

import (
	"errors"
	"fmt"
)

// ParkingError is the base error type for all parking-related errors
type ParkingError struct {
	// Error code for categorizing errors
	Code string

	// Human-readable error message
	Message string

	// Optional underlying error
	Err error
}

// Error returns the string representation of the error
func (e *ParkingError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Err.Error())
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *ParkingError) Unwrap() error {
	return e.Err
}

// Is reports whether target matches this error
func (e *ParkingError) Is(target error) bool {
	t, ok := target.(*ParkingError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// NewParkingError creates a new ParkingError with the given code and message
func NewParkingError(code, message string, err error) *ParkingError {
	return &ParkingError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// WrapError wraps an existing error with a parking error
func WrapError(err error, code, message string) *ParkingError {
	return &ParkingError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Predefined error codes
const (
	CodeInvalidOperation     = "INVALID_OPERATION"
	CodeInvalidInput         = "INVALID_INPUT"
	CodeNoSpaceAvailable     = "NO_SPACE_AVAILABLE"
	CodeVehicleNotFound      = "VEHICLE_NOT_FOUND"
	CodeVehicleAlreadyParked = "VEHICLE_ALREADY_PARKED"
	CodeSpotAlreadyOccupied  = "SPOT_ALREADY_OCCUPIED"
	CodeSpotNotOccupied      = "SPOT_NOT_OCCUPIED"
	CodeInvalidSpotID        = "INVALID_SPOT_ID"
	CodeInvalidVehicleType   = "INVALID_VEHICLE_TYPE"
	CodeInvalidVehicleNumber = "INVALID_VEHICLE_NUMBER"
	CodeVehicleMismatch      = "VEHICLE_MISMATCH"
	CodeSpotInactive         = "SPOT_INACTIVE"
	CodeInvalidSpotType      = "INVALID_SPOT_TYPE"
	CodeInvalidFloor         = "INVALID_FLOOR"
	CodeInternalError        = "INTERNAL_ERROR"
)

// Errors using the standard errors package for use with errors.Is
var (
	ErrInvalidOperation     = errors.New("invalid operation")
	ErrInvalidInput         = errors.New("invalid input")
	ErrNoSpaceAvailable     = errors.New("no space available")
	ErrVehicleNotFound      = errors.New("vehicle not found")
	ErrVehicleAlreadyParked = errors.New("vehicle already parked")
	ErrSpotAlreadyOccupied  = errors.New("spot already occupied")
	ErrSpotNotOccupied      = errors.New("spot not occupied")
	ErrInvalidSpotID        = errors.New("invalid spot ID")
	ErrInvalidVehicleType   = errors.New("invalid vehicle type")
	ErrInvalidVehicleNumber = errors.New("invalid vehicle number")
	ErrVehicleMismatch      = errors.New("vehicle mismatch")
	ErrSpotInactive         = errors.New("spot is inactive")
	ErrInvalidSpotType      = errors.New("invalid spot type")
	ErrInvalidFloor         = errors.New("invalid floor")
	ErrInternalError        = errors.New("internal error")
)
