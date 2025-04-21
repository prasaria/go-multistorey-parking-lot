package model

import (
	"fmt"
	"testing"
)

// BenchmarkParkingOperations benchmarks basic parking operations
func BenchmarkParkingOperations(b *testing.B) {
	lot, _ := CreateParkingLot("Benchmark", 3, 10, 10)

	// Create test data
	vehicleNumbers := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		vehicleNumbers[i] = fmt.Sprintf("BENCH-%04d", i)
	}

	b.Run("ParkUnpark", func(b *testing.B) {
		// Reset the timer to exclude setup
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			idx := i % 1000
			vehicleType := VehicleTypeBicycle

			// Rotate vehicle types
			if idx%3 == 1 {
				vehicleType = VehicleTypeMotorcycle
			} else if idx%3 == 2 {
				vehicleType = VehicleTypeAutomobile
			}

			vehicleNumber := vehicleNumbers[idx]
			spotID, err := lot.Park(vehicleType, vehicleNumber)
			if err != nil {
				b.Fatalf("Failed to park vehicle: %v", err)
			}

			err = lot.Unpark(spotID, vehicleNumber)
			if err != nil {
				b.Fatalf("Failed to unpark vehicle: %v", err)
			}
		}
	})

	b.Run("Search", func(b *testing.B) {
		// First park some vehicles
		parkedVehicles := make(map[string]string) // map vehicleNumber to spotID

		for i := 0; i < 100; i++ {
			vehicleType := VehicleTypeBicycle

			// Rotate vehicle types
			if i%3 == 1 {
				vehicleType = VehicleTypeMotorcycle
			} else if i%3 == 2 {
				vehicleType = VehicleTypeAutomobile
			}

			vehicleNumber := vehicleNumbers[i]
			spotID, err := lot.Park(vehicleType, vehicleNumber)
			if err != nil {
				b.Fatalf("Failed to park vehicle: %v", err)
			}

			parkedVehicles[vehicleNumber] = spotID
		}

		// Reset the timer to exclude setup
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			idx := i % 100
			vehicleNumber := vehicleNumbers[idx]

			_, isParked, err := lot.SearchVehicle(vehicleNumber)
			if err != nil {
				b.Fatalf("Failed to search vehicle: %v", err)
			}

			if !isParked {
				b.Fatalf("Expected vehicle %s to be parked", vehicleNumber)
			}
		}

		// Cleanup
		b.StopTimer()
		for vehicleNumber, spotID := range parkedVehicles {
			_ = lot.Unpark(spotID, vehicleNumber)
		}
	})

	b.Run("AvailableSpots", func(b *testing.B) {
		// Reset the timer
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			vehicleType := VehicleTypeBicycle

			// Rotate vehicle types
			if i%3 == 1 {
				vehicleType = VehicleTypeMotorcycle
			} else if i%3 == 2 {
				vehicleType = VehicleTypeAutomobile
			}

			_, err := lot.AvailableSpot(vehicleType)
			if err != nil {
				b.Fatalf("Failed to get available spots: %v", err)
			}
		}
	})

	b.Run("ConcurrentOperations", func(b *testing.B) {
		// This benchmark tests concurrent access
		lot, _ := CreateParkingLot("ConcurrentBench", 3, 20, 20)

		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				i++
				vehicleNumber := fmt.Sprintf("CONCURRENT-%04d", i)
				vehicleType := VehicleTypeBicycle

				// Rotate vehicle types
				if i%3 == 1 {
					vehicleType = VehicleTypeMotorcycle
				} else if i%3 == 2 {
					vehicleType = VehicleTypeAutomobile
				}

				// Alternate between operations
				switch i % 4 {
				case 0: // Park
					spotID, err := lot.Park(vehicleType, vehicleNumber)
					if err == nil {
						// Remember to unpark later
						defer func(id, vn string) {
							lot.Unpark(id, vn)
						}(spotID, vehicleNumber)
					}
				case 1: // Search
					lot.SearchVehicle(vehicleNumber)
				case 2: // AvailableSpot
					lot.AvailableSpot(vehicleType)
				case 3: // GetSpotCountByType
					lot.GetSpotCountByType()
				}
			}
		})
	})
}

// BenchmarkCreation benchmarks parking lot creation
func BenchmarkCreation(b *testing.B) {
	b.Run("SmallLot", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := CreateParkingLot("Small", 2, 5, 5)
			if err != nil {
				b.Fatalf("Failed to create small lot: %v", err)
			}
		}
	})

	b.Run("MediumLot", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := CreateParkingLot("Medium", 3, 10, 20)
			if err != nil {
				b.Fatalf("Failed to create medium lot: %v", err)
			}
		}
	})

	b.Run("LargeLot", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := CreateParkingLot("Large", 5, 30, 40)
			if err != nil {
				b.Fatalf("Failed to create large lot: %v", err)
			}
		}
	})
}
