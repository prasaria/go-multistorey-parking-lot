package model

import (
	"errors"
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

// Error definitions
var (
	ErrInvalidVehicleType = errors.New("invalid vehicle type")
)
