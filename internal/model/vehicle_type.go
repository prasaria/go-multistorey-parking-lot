package model

import (
	"strings"

	"github.com/prasaria/go-multistorey-parking-lot/internal/errors"
)

// VehicleType represents the type of a vehicle
type VehicleType string

const (
	// VehicleTypeBicycle represents a bicycle
	VehicleTypeBicycle VehicleType = "BICYCLE"

	// VehicleTypeMotorcycle represents a motorcycle
	VehicleTypeMotorcycle VehicleType = "MOTORCYCLE"

	// VehicleTypeAutomobile represents an automobile
	VehicleTypeAutomobile VehicleType = "AUTOMOBILE"
)

// GetPreferredSpotType returns the preferred spot type for this vehicle type
func (v VehicleType) GetPreferredSpotType() SpotType {
	switch v {
	case VehicleTypeBicycle:
		return SpotTypeBicycle
	case VehicleTypeMotorcycle:
		return SpotTypeMotorcycle
	case VehicleTypeAutomobile:
		return SpotTypeAutomobile
	default:
		// This should never happen if VehicleType is properly validated
		return SpotTypeInactive
	}
}

// GetCompatibleSpotTypes returns all spot types that can accommodate this vehicle type
func (v VehicleType) GetCompatibleSpotTypes() []SpotType {
	switch v {
	case VehicleTypeBicycle:
		return []SpotType{SpotTypeBicycle}
	case VehicleTypeMotorcycle:
		return []SpotType{SpotTypeMotorcycle}
	case VehicleTypeAutomobile:
		return []SpotType{SpotTypeAutomobile}
	default:
		return []SpotType{}
	}
}

// ParseVehicleType converts a string to VehicleType
func ParseVehicleType(s string) (VehicleType, error) {
	switch strings.ToUpper(s) {
	case "BICYCLE", "B", "BIKE":
		return VehicleTypeBicycle, nil
	case "MOTORCYCLE", "M", "MOTORBIKE":
		return VehicleTypeMotorcycle, nil
	case "AUTOMOBILE", "A", "CAR", "AUTO":
		return VehicleTypeAutomobile, nil
	default:
		return "", errors.NewInvalidVehicleTypeError(s)
	}
}

// String returns the string representation of VehicleType
func (v VehicleType) String() string {
	return string(v)
}
