package model

import (
	"regexp"
	"strings"

	"github.com/prasaria/go-multistorey-parking-lot/internal/errors"
)

// Vehicle represents a vehicle that can be parked in the parking lot
type Vehicle struct {
	// Type of the vehicle (bicycle, motorcycle, automobile)
	Type VehicleType

	// Registration number of the vehicle (must be unique)
	Number string
}

// Vehicle number validation pattern
// This is a simple pattern that accepts alphanumeric characters with possible hyphens and spaces
var vehicleNumberPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9\s-]*$`)

// NewVehicle creates a new vehicle with the given type and number
func NewVehicle(vehicleType VehicleType, number string) (*Vehicle, error) {
	// Validate vehicle type
	if vehicleType != VehicleTypeBicycle &&
		vehicleType != VehicleTypeMotorcycle &&
		vehicleType != VehicleTypeAutomobile {
		return nil, errors.NewInvalidVehicleTypeError(string(vehicleType))
	}

	// Validate vehicle number
	if err := ValidateVehicleNumber(number); err != nil {
		return nil, err
	}

	return &Vehicle{
		Type:   vehicleType,
		Number: NormalizeVehicleNumber(number),
	}, nil
}

// ValidateVehicleNumber checks if the vehicle number is valid
func ValidateVehicleNumber(number string) error {
	// Check for empty string
	if strings.TrimSpace(number) == "" {
		return errors.NewInvalidVehicleNumberError(number, "vehicle number cannot be empty")
	}

	// Check pattern
	if !vehicleNumberPattern.MatchString(number) {
		return errors.NewInvalidVehicleNumberError(number, "invalid format")
	}

	return nil
}

// NormalizeVehicleNumber standardizes a vehicle number by trimming spaces
// and converting to uppercase
func NormalizeVehicleNumber(number string) string {
	// Convert to uppercase and trim spaces
	normalized := strings.ToUpper(strings.TrimSpace(number))

	// Standardize spacing by replacing multiple spaces with a single space
	spaceRegex := regexp.MustCompile(`\s+`)
	normalized = spaceRegex.ReplaceAllString(normalized, " ")

	return normalized
}

// String returns a string representation of the vehicle
func (v *Vehicle) String() string {
	return GetVehicleTypeDisplay(v.Type) + " [" + v.Number + "]"
}

// Equal checks if two vehicles are equal (same type and number)
func (v *Vehicle) Equal(other *Vehicle) bool {
	if v == nil || other == nil {
		return v == other
	}

	return v.Type == other.Type && v.Number == other.Number
}
