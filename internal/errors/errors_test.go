package errors

import (
	"errors"
	"strings"
	"testing"
)

func TestParkingError(t *testing.T) {
	// Test creating a simple error
	err := NewParkingError(CodeInvalidOperation, "Test error", nil)

	if err.Code != CodeInvalidOperation {
		t.Errorf("Expected code %s, got %s", CodeInvalidOperation, err.Code)
	}

	if err.Message != "Test error" {
		t.Errorf("Expected message 'Test error', got '%s'", err.Message)
	}

	// Test error message format
	expectedMsg := "INVALID_OPERATION: Test error"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}

	// Test wrapping an error
	underlyingErr := errors.New("underlying error")
	wrappedErr := WrapError(underlyingErr, CodeInvalidInput, "Wrapped error")

	// Test error message with wrapped error
	expectedWrappedMsg := "INVALID_INPUT: Wrapped error (underlying error)"
	if wrappedErr.Error() != expectedWrappedMsg {
		t.Errorf("Expected wrapped error message '%s', got '%s'",
			expectedWrappedMsg, wrappedErr.Error())
	}

	// Test unwrapping
	unwrappedErr := wrappedErr.Unwrap()
	if unwrappedErr != underlyingErr {
		t.Errorf("Unwrap didn't return the expected underlying error")
	}

	// Test errors.Is
	if !errors.Is(wrappedErr, wrappedErr) {
		t.Errorf("errors.Is failed for the same error")
	}

	sameCodeErr := NewParkingError(CodeInvalidInput, "Different message", nil)
	if !errors.Is(wrappedErr, sameCodeErr) {
		t.Errorf("errors.Is failed for errors with the same code")
	}

	differentCodeErr := NewParkingError(CodeInvalidOperation, "Different message", nil)
	if errors.Is(wrappedErr, differentCodeErr) {
		t.Errorf("errors.Is succeeded for errors with different codes")
	}
}

func TestSpecificErrors(t *testing.T) {
	// Test NoSpaceError
	noSpaceErr := NewNoSpaceError("AUTOMOBILE")
	if noSpaceErr.Code != CodeNoSpaceAvailable {
		t.Errorf("Expected code %s, got %s", CodeNoSpaceAvailable, noSpaceErr.Code)
	}

	if noSpaceErr.VehicleType != "AUTOMOBILE" {
		t.Errorf("Expected vehicle type AUTOMOBILE, got %s", noSpaceErr.VehicleType)
	}

	if !strings.Contains(noSpaceErr.Error(), "AUTOMOBILE") {
		t.Errorf("Error message doesn't contain vehicle type: %s", noSpaceErr.Error())
	}

	// Test InvalidOperationError
	invalidOpErr := NewInvalidOperationError("park", "vehicle already parked")
	if invalidOpErr.Code != CodeInvalidOperation {
		t.Errorf("Expected code %s, got %s", CodeInvalidOperation, invalidOpErr.Code)
	}

	if invalidOpErr.Operation != "park" {
		t.Errorf("Expected operation 'park', got '%s'", invalidOpErr.Operation)
	}

	if invalidOpErr.Reason != "vehicle already parked" {
		t.Errorf("Expected reason 'vehicle already parked', got '%s'", invalidOpErr.Reason)
	}

	// Test VehicleNotFoundError
	vehicleNotFoundErr := NewVehicleNotFoundError("ABC-123")
	if vehicleNotFoundErr.Code != CodeVehicleNotFound {
		t.Errorf("Expected code %s, got %s", CodeVehicleNotFound, vehicleNotFoundErr.Code)
	}

	if vehicleNotFoundErr.VehicleNumber != "ABC-123" {
		t.Errorf("Expected vehicle number ABC-123, got %s", vehicleNotFoundErr.VehicleNumber)
	}

	// Test SpotOccupancyError
	spotOccupiedErr := NewSpotAlreadyOccupiedError("1-2-3")
	if spotOccupiedErr.Code != CodeSpotAlreadyOccupied {
		t.Errorf("Expected code %s, got %s", CodeSpotAlreadyOccupied, spotOccupiedErr.Code)
	}

	if spotOccupiedErr.SpotID != "1-2-3" {
		t.Errorf("Expected spot ID 1-2-3, got %s", spotOccupiedErr.SpotID)
	}

	if !spotOccupiedErr.IsOccupying {
		t.Errorf("Expected IsOccupying to be true")
	}

	// Test ValidationError
	validationErr := NewValidationError("vehicleNumber", "ABC", "must be at least 4 characters")
	if validationErr.Code != CodeInvalidInput {
		t.Errorf("Expected code %s, got %s", CodeInvalidInput, validationErr.Code)
	}

	if validationErr.Field != "vehicleNumber" {
		t.Errorf("Expected field 'vehicleNumber', got '%s'", validationErr.Field)
	}

	if validationErr.Value != "ABC" {
		t.Errorf("Expected value 'ABC', got '%s'", validationErr.Value)
	}

	// Test standard errors.Is usage
	if !errors.Is(noSpaceErr, ErrNoSpaceAvailable) {
		t.Errorf("errors.Is failed for NoSpaceError and ErrNoSpaceAvailable")
	}

	if !errors.Is(invalidOpErr, ErrInvalidOperation) {
		t.Errorf("errors.Is failed for InvalidOperationError and ErrInvalidOperation")
	}

	if !errors.Is(vehicleNotFoundErr, ErrVehicleNotFound) {
		t.Errorf("errors.Is failed for VehicleNotFoundError and ErrVehicleNotFound")
	}
}
