package model

import (
	"strings"

	"github.com/prasaria/go-multistorey-parking-lot/internal/errors"
)

// IsValidSpotCode checks if a string is a valid spot code in the format "TYPE-STATUS"
// where TYPE is one of B, M, A, X and STATUS is 0 or 1
func IsValidSpotCode(code string) bool {
	parts := strings.Split(code, "-")
	if len(parts) != 2 {
		return false
	}

	typeChar := parts[0]
	status := parts[1]

	// Check type
	if typeChar != "B" && typeChar != "M" && typeChar != "A" && typeChar != "X" {
		return false
	}

	// Check status
	if status != "0" && status != "1" {
		return false
	}

	return true
}

// SpotCodeToSpotType converts a spot code (e.g., "B-1") to a SpotType
func SpotCodeToSpotType(code string) (SpotType, error) {
	if !IsValidSpotCode(code) {
		return "", errors.NewInvalidSpotTypeError(code)
	}

	return SpotType(code), nil
}

// GetVehicleTypeDisplay returns a user-friendly display name for a vehicle type
func GetVehicleTypeDisplay(vt VehicleType) string {
	switch vt {
	case VehicleTypeBicycle:
		return "Bicycle"
	case VehicleTypeMotorcycle:
		return "Motorcycle"
	case VehicleTypeAutomobile:
		return "Automobile"
	default:
		return "Unknown"
	}
}

// GetSpotTypeDisplay returns a user-friendly display name for a spot type
func GetSpotTypeDisplay(st SpotType) string {
	switch st {
	case SpotTypeBicycle:
		return "Bicycle Spot"
	case SpotTypeMotorcycle:
		return "Motorcycle Spot"
	case SpotTypeAutomobile:
		return "Automobile Spot"
	case SpotTypeInactive:
		return "Inactive Spot"
	default:
		return "Unknown Spot"
	}
}
