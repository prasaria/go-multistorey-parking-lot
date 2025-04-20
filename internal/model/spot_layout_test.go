package model

import (
	"testing"
)

func TestNewSpotLayout(t *testing.T) {
	tests := []struct {
		name      string
		floors    int
		rows      int
		columns   int
		expectErr bool
	}{
		{
			name:      "Valid dimensions",
			floors:    3,
			rows:      10,
			columns:   20,
			expectErr: false,
		},
		{
			name:      "Minimum dimensions",
			floors:    1,
			rows:      1,
			columns:   1,
			expectErr: false,
		},
		{
			name:      "Maximum dimensions",
			floors:    8,
			rows:      1000,
			columns:   1000,
			expectErr: false,
		},
		{
			name:      "Too many floors",
			floors:    9,
			rows:      10,
			columns:   20,
			expectErr: true,
		},
		{
			name:      "Too many rows",
			floors:    3,
			rows:      1001,
			columns:   20,
			expectErr: true,
		},
		{
			name:      "Too many columns",
			floors:    3,
			rows:      10,
			columns:   1001,
			expectErr: true,
		},
		{
			name:      "Zero floors",
			floors:    0,
			rows:      10,
			columns:   20,
			expectErr: true,
		},
		{
			name:      "Zero rows",
			floors:    3,
			rows:      0,
			columns:   20,
			expectErr: true,
		},
		{
			name:      "Zero columns",
			floors:    3,
			rows:      10,
			columns:   0,
			expectErr: true,
		},
		{
			name:      "Negative floors",
			floors:    -1,
			rows:      10,
			columns:   20,
			expectErr: true,
		},
		{
			name:      "Negative rows",
			floors:    3,
			rows:      -1,
			columns:   20,
			expectErr: true,
		},
		{
			name:      "Negative columns",
			floors:    3,
			rows:      10,
			columns:   -1,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			layout, err := NewSpotLayout(tt.floors, tt.rows, tt.columns)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if layout == nil {
				t.Fatalf("Expected layout, got nil")
			}

			// Check dimensions
			if len(layout.SpotMap) != tt.floors {
				t.Errorf("Expected %d floors, got %d", tt.floors, len(layout.SpotMap))
			}

			for f := 0; f < tt.floors; f++ {
				if len(layout.SpotMap[f]) != tt.rows {
					t.Errorf("Floor %d: expected %d rows, got %d",
						f, tt.rows, len(layout.SpotMap[f]))
				}

				for r := 0; r < tt.rows; r++ {
					if len(layout.SpotMap[f][r]) != tt.columns {
						t.Errorf("Floor %d, Row %d: expected %d columns, got %d",
							f, r, tt.columns, len(layout.SpotMap[f][r]))
					}
				}
			}

			// Special case for minimum dimensions (1x1x1)
			// In this case, we know we'll only have one spot type
			if tt.name == "Minimum dimensions" {
				// Just verify that the spot type is valid and is one of the active types
				spotType := layout.SpotMap[0][0][0]
				if spotType != SpotTypeBicycle && spotType != SpotTypeMotorcycle && spotType != SpotTypeAutomobile {
					t.Errorf("The only spot in a 1x1x1 layout should be an active spot, got %s", spotType)
				}
				return
			}

			// For other layouts, check distribution of spot types
			counts := layout.CountSpotsByType()
			totalSpots := tt.floors * tt.rows * tt.columns

			// Ensure we have a reasonable distribution (at least some of each type)
			if counts[SpotTypeBicycle] == 0 {
				t.Errorf("No bicycle spots in layout")
			}

			if counts[SpotTypeMotorcycle] == 0 {
				t.Errorf("No motorcycle spots in layout")
			}

			if counts[SpotTypeAutomobile] == 0 {
				t.Errorf("No automobile spots in layout")
			}

			// Verify total count matches expected total
			totalCount := 0
			for _, count := range counts {
				totalCount += count
			}

			if totalCount != totalSpots {
				t.Errorf("Expected %d total spots, counted %d", totalSpots, totalCount)
			}
		})
	}
}

func TestSpotLayoutGetAndSetSpotType(t *testing.T) {
	layout, err := NewSpotLayout(3, 5, 10)
	if err != nil {
		t.Fatalf("Failed to create layout: %v", err)
	}

	// Test getting spot type
	_, err = layout.GetSpotType(1, 2, 3)
	if err != nil {
		t.Errorf("Failed to get spot type: %v", err)
	}

	// Test setting spot type
	err = layout.SetSpotType(1, 2, 3, SpotTypeInactive)
	if err != nil {
		t.Errorf("Failed to set spot type: %v", err)
	}

	// Verify the spot type was changed
	newSpotType, _ := layout.GetSpotType(1, 2, 3)
	if newSpotType != SpotTypeInactive {
		t.Errorf("Expected spot type %s, got %s", SpotTypeInactive, newSpotType)
	}

	// Test invalid location
	_, err = layout.GetSpotType(10, 0, 0)
	if err == nil {
		t.Errorf("Expected error for out of range floor")
	}

	_, err = layout.GetSpotType(0, 10, 0)
	if err == nil {
		t.Errorf("Expected error for out of range row")
	}

	_, err = layout.GetSpotType(0, 0, 20)
	if err == nil {
		t.Errorf("Expected error for out of range column")
	}

	// Test invalid spot type
	err = layout.SetSpotType(0, 0, 0, "INVALID")
	if err == nil {
		t.Errorf("Expected error for invalid spot type")
	}
}

func TestCreateParkingSpots(t *testing.T) {
	// Create a small layout for testing
	layout, err := NewSpotLayout(2, 3, 4)
	if err != nil {
		t.Fatalf("Failed to create layout: %v", err)
	}

	// Create spots
	spots, err := layout.CreateParkingSpots()
	if err != nil {
		t.Fatalf("Failed to create parking spots: %v", err)
	}

	// Verify number of floors
	if len(spots) != 2 {
		t.Fatalf("Expected 2 floors of spots, got %d", len(spots))
	}

	// Verify spot properties
	for f, floor := range spots {
		// Each floor should have 3 rows * 4 columns = 12 spots
		if len(floor) != 12 {
			t.Errorf("Floor %d: expected 12 spots, got %d", f, len(floor))
		}

		// Check a few random spots
		for _, spot := range floor {
			// Verify floor number
			if spot.Floor != f {
				t.Errorf("Spot has incorrect floor: expected %d, got %d", f, spot.Floor)
			}

			// Row and column should be within range
			if spot.Row < 0 || spot.Row >= 3 {
				t.Errorf("Spot has out of range row: %d", spot.Row)
			}

			if spot.Column < 0 || spot.Column >= 4 {
				t.Errorf("Spot has out of range column: %d", spot.Column)
			}

			// Spot type should match layout
			spotType, _ := layout.GetSpotType(spot.Floor, spot.Row, spot.Column)
			if spot.Type != spotType {
				t.Errorf("Spot type mismatch: expected %s, got %s",
					spotType, spot.Type)
			}

			// Verify initial state
			if spot.IsOccupied() {
				t.Errorf("New spot should not be occupied")
			}

			if spot.GetVehicleNumber() != "" {
				t.Errorf("New spot should not have a vehicle number")
			}
		}
	}
}

// Also add a test for the corner case of minimum layout with all active types
func TestMinimumLayoutWithAllActiveTypes(t *testing.T) {
	// Create a larger layout but with minimum dimensions that can fit all types
	layout, err := NewSpotLayout(3, 1, 1)
	if err != nil {
		t.Fatalf("Failed to create layout: %v", err)
	}

	// Count spot types
	counts := layout.CountSpotsByType()

	// Verify we have at least one of each active type
	hasAllTypes := counts[SpotTypeBicycle] > 0 &&
		counts[SpotTypeMotorcycle] > 0 &&
		counts[SpotTypeAutomobile] > 0

	if !hasAllTypes {
		t.Errorf("3x1x1 layout should have all active spot types, got: %v", counts)
	}

	// Verify total count is correct
	totalCount := 0
	for _, count := range counts {
		totalCount += count
	}

	if totalCount != 3 {
		t.Errorf("Expected 3 total spots, counted %d", totalCount)
	}
}
