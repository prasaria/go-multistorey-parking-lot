package config

// ParkingLotConfig represents the configuration for parking lot
type ParkingLotConfig struct {
	Floors  int
	Rows    int
	Columns int
}

// Validate checks if the parking lot configuration is valid
func (c *ParkingLotConfig) Validate() error {
	if c.Floors < 1 || c.Floors > 8 {
		return ErrInvalidFloorCount
	}

	if c.Rows < 1 || c.Rows > 1000 {
		return ErrInvalidRowCount
	}

	if c.Columns < 1 || c.Columns > 1000 {
		return ErrInvalidColumnCount
	}

	return nil
}

// DefaultConfig returns a default configuration
func DefaultConfig() ParkingLotConfig {
	return ParkingLotConfig{
		Floors:  3,
		Rows:    5,
		Columns: 10,
	}
}
