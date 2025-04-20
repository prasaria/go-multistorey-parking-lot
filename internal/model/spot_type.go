package model

import (
	"errors"
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
