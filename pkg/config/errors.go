package config

import "errors"

// Configuration-related errors
var (
	ErrInvalidFloorCount  = errors.New("invalid floor count: must be between 1 and 8")
	ErrInvalidRowCount    = errors.New("invalid row count: must be between 1 and 1000")
	ErrInvalidColumnCount = errors.New("invalid column count: must be between 1 and 1000")
)
