package model

import (
	"fmt"

	"github.com/prasaria/go-multistorey-parking-lot/internal/errors"
)

// SpotLayout represents a layout of parking spots
type SpotLayout struct {
	SpotMap [][][]SpotType
}

// NewSpotLayout creates a new spot layout with the given dimensions
func NewSpotLayout(floors, rows, columns int) (*SpotLayout, error) {
	// Validate dimensions
	if floors < 1 || floors > 8 {
		return nil, errors.NewValidationError("floors",
			fmt.Sprintf("%d", floors), "number of floors must be between 1 and 8")
	}

	if rows < 1 || rows > 1000 {
		return nil, errors.NewValidationError("rows",
			fmt.Sprintf("%d", rows), "number of rows must be between 1 and 1000")
	}

	if columns < 1 || columns > 1000 {
		return nil, errors.NewValidationError("columns",
			fmt.Sprintf("%d", columns), "number of columns must be between 1 and 1000")
	}

	// Initialize layout with default distribution
	layout := &SpotLayout{
		SpotMap: make([][][]SpotType, floors),
	}

	for f := 0; f < floors; f++ {
		layout.SpotMap[f] = make([][]SpotType, rows)

		for r := 0; r < rows; r++ {
			layout.SpotMap[f][r] = make([]SpotType, columns)

			for c := 0; c < columns; c++ {
				// Special case for very small layouts to ensure at least one of each spot type
				if rows == 1 && columns < 4 {
					// For minimal layouts, ensure we have at least one of each active type
					if columns == 1 {
						// With just one column, pick a type based on the floor
						if f%3 == 0 {
							layout.SpotMap[f][r][c] = SpotTypeBicycle
						} else if f%3 == 1 {
							layout.SpotMap[f][r][c] = SpotTypeMotorcycle
						} else {
							layout.SpotMap[f][r][c] = SpotTypeAutomobile
						}
					} else if columns == 2 {
						// With two columns, alternate between types
						if c == 0 {
							if f%2 == 0 {
								layout.SpotMap[f][r][c] = SpotTypeBicycle
							} else {
								layout.SpotMap[f][r][c] = SpotTypeMotorcycle
							}
						} else {
							if f%2 == 0 {
								layout.SpotMap[f][r][c] = SpotTypeAutomobile
							} else {
								layout.SpotMap[f][r][c] = SpotTypeBicycle
							}
						}
					} else { // columns == 3
						// With three columns, one of each type
						if c == 0 {
							layout.SpotMap[f][r][c] = SpotTypeBicycle
						} else if c == 1 {
							layout.SpotMap[f][r][c] = SpotTypeMotorcycle
						} else {
							layout.SpotMap[f][r][c] = SpotTypeAutomobile
						}
					}
				} else {
					// Regular distribution for larger layouts
					// Make some spots inactive (e.g., pillars, structural elements)
					if (r%7 == 0 && c%7 == 0) || (r%7 == 0 && c%7 == 1) {
						layout.SpotMap[f][r][c] = SpotTypeInactive
						continue
					}

					// Distribute different vehicle types
					switch {
					case c < columns/4:
						// First 25% of columns for bicycles
						layout.SpotMap[f][r][c] = SpotTypeBicycle
					case c < columns/2:
						// Next 25% for motorcycles
						layout.SpotMap[f][r][c] = SpotTypeMotorcycle
					default:
						// Remaining 50% for automobiles
						layout.SpotMap[f][r][c] = SpotTypeAutomobile
					}
				}
			}
		}
	}

	// Ensure we have at least one of each spot type in the entire layout
	counts := layout.CountSpotsByType()
	if counts[SpotTypeBicycle] == 0 || counts[SpotTypeMotorcycle] == 0 || counts[SpotTypeAutomobile] == 0 {
		// Force at least one of each type
		if floors > 0 && rows > 0 && columns >= 3 {
			layout.SpotMap[0][0][0] = SpotTypeBicycle
			layout.SpotMap[0][0][1] = SpotTypeMotorcycle
			layout.SpotMap[0][0][2] = SpotTypeAutomobile
		} else if floors > 0 && rows > 0 && columns == 2 {
			layout.SpotMap[0][0][0] = SpotTypeBicycle
			layout.SpotMap[0][0][1] = SpotTypeMotorcycle
			// If we have another row or floor, add the automobile there
			if rows > 1 {
				layout.SpotMap[0][1][0] = SpotTypeAutomobile
			} else if floors > 1 {
				layout.SpotMap[1][0][0] = SpotTypeAutomobile
			}
		} else if floors > 0 && rows > 0 && columns == 1 {
			// With only one column, we need multiple rows or floors
			layout.SpotMap[0][0][0] = SpotTypeBicycle
			if rows > 1 {
				layout.SpotMap[0][1][0] = SpotTypeMotorcycle
				if rows > 2 {
					layout.SpotMap[0][2][0] = SpotTypeAutomobile
				} else if floors > 1 {
					layout.SpotMap[1][0][0] = SpotTypeAutomobile
				}
			} else if floors > 1 {
				layout.SpotMap[1][0][0] = SpotTypeMotorcycle
				if floors > 2 {
					layout.SpotMap[2][0][0] = SpotTypeAutomobile
				}
			}
		}
	}

	return layout, nil
}

// GetSpotType returns the spot type at the given location
func (l *SpotLayout) GetSpotType(floor, row, column int) (SpotType, error) {
	// Validate location
	if floor < 0 || floor >= len(l.SpotMap) {
		return "", errors.NewValidationError("floor",
			fmt.Sprintf("%d", floor),
			fmt.Sprintf("floor out of range [0-%d]", len(l.SpotMap)-1))
	}

	if row < 0 || row >= len(l.SpotMap[floor]) {
		return "", errors.NewValidationError("row",
			fmt.Sprintf("%d", row),
			fmt.Sprintf("row out of range [0-%d]", len(l.SpotMap[floor])-1))
	}

	if column < 0 || column >= len(l.SpotMap[floor][row]) {
		return "", errors.NewValidationError("column",
			fmt.Sprintf("%d", column),
			fmt.Sprintf("column out of range [0-%d]", len(l.SpotMap[floor][row])-1))
	}

	return l.SpotMap[floor][row][column], nil
}

// SetSpotType sets the spot type at the given location
func (l *SpotLayout) SetSpotType(floor, row, column int, spotType SpotType) error {
	// Validate spot type
	_, err := ParseSpotType(string(spotType))
	if err != nil {
		return err
	}

	// Validate location
	if floor < 0 || floor >= len(l.SpotMap) {
		return errors.NewValidationError("floor",
			fmt.Sprintf("%d", floor),
			fmt.Sprintf("floor out of range [0-%d]", len(l.SpotMap)-1))
	}

	if row < 0 || row >= len(l.SpotMap[floor]) {
		return errors.NewValidationError("row",
			fmt.Sprintf("%d", row),
			fmt.Sprintf("row out of range [0-%d]", len(l.SpotMap[floor])-1))
	}

	if column < 0 || column >= len(l.SpotMap[floor][row]) {
		return errors.NewValidationError("column",
			fmt.Sprintf("%d", column),
			fmt.Sprintf("column out of range [0-%d]", len(l.SpotMap[floor][row])-1))
	}

	l.SpotMap[floor][row][column] = spotType
	return nil
}

// CountSpotsByType counts the number of spots of each type in the layout
func (l *SpotLayout) CountSpotsByType() map[SpotType]int {
	counts := make(map[SpotType]int)

	// Initialize counts for all spot types
	counts[SpotTypeBicycle] = 0
	counts[SpotTypeMotorcycle] = 0
	counts[SpotTypeAutomobile] = 0
	counts[SpotTypeInactive] = 0

	// Count spots
	for f := range l.SpotMap {
		for r := range l.SpotMap[f] {
			for c := range l.SpotMap[f][r] {
				spotType := l.SpotMap[f][r][c]
				counts[spotType]++
			}
		}
	}

	return counts
}

// CreateParkingSpots creates actual parking spot objects from the layout
func (l *SpotLayout) CreateParkingSpots() ([][]*ParkingSpot, error) {
	floors := len(l.SpotMap)
	if floors == 0 {
		return nil, errors.NewValidationError("floorNumber",
			fmt.Sprintf("%d", floors), "floor number cannot be less than 1")
	}

	// Create spots
	spots := make([][]*ParkingSpot, floors)

	for f := range l.SpotMap {
		rows := len(l.SpotMap[f])
		spots[f] = make([]*ParkingSpot, 0, rows*len(l.SpotMap[f][0]))

		for r := range l.SpotMap[f] {
			for c := range l.SpotMap[f][r] {
				spotType := l.SpotMap[f][r][c]

				spot, err := NewParkingSpot(spotType, f, r, c)
				if err != nil {
					return nil, fmt.Errorf("error creating spot at %d-%d-%d: %w", f, r, c, err)
				}

				spots[f] = append(spots[f], spot)
			}
		}
	}

	return spots, nil
}
