package model

import (
	"errors"
	"strings"
)

// SpotType represents the type of a parkign spot
type SpotType string

const (
	// SpotTypeBicycle represents a bicycle parking spot
	SpotTypeBicycle SpotType = "B-1"

	// SpotTypeMotorcycle represents a motorcycle parking spot
	SpotTypeMotorcycle SpotType = "M-1"

	// SpotTypeAutomobile represents an automobile parking spot
	SpotTypeAutomobile SpotType = "A-1"

	// SpotTypeInactive represents an inactive parking spot
	SpotTypeInactive SpotType = "X-0"
)

// Error definitions
var (
	ErrInvalidSpotType = errors.New("invalid spot type")
)

// IsActive returns true if the spot type is active
func (s SpotType) IsActive() bool {
	return s != SpotTypeInactive
}

// CanParkVehicleType checks if a vehicle type can park in this spot type
func (s SpotType) CanParkVehicleType(vt VehicleType) bool {
	if !s.IsActive() {
		return false
	}

	switch s {
	case SpotTypeBicycle:
		return vt == VehicleTypeBicycle
	case SpotTypeMotorcycle:
		return vt == VehicleTypeBicycle || vt == VehicleTypeMotorcycle
	case SpotTypeAutomobile:
		return vt == VehicleTypeAutomobile
	default:
		return false
	}
}

// ParseSpotType converts a string to SpotType
func ParseSpotType(s string) (SpotType, error) {
	switch strings.ToUpper(s) {
	case string(SpotTypeBicycle):
		return SpotTypeBicycle, nil
	case string(SpotTypeMotorcycle):
		return SpotTypeMotorcycle, nil
	case string(SpotTypeAutomobile):
		return SpotTypeAutomobile, nil
	case string(SpotTypeInactive):
		return SpotTypeInactive, nil
	default:
		return "", ErrInvalidSpotType
	}
}
