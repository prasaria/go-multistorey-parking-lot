# Multi-Storey Parking Lot CLI System

A command-line application for managing a multi-storey parking lot written in Go.

## Features

- Multiple floors with configurable rows and columns
- Support for different vehicle types (bicycles, motorcycles, automobiles)
- Park and unpark vehicles
- Find available parking spots
- Search for vehicles by license plate number
- Thread-safe operations for concurrent access
- Interactive CLI with colorized output
- JSON output option for programmatic use

## UML Diagram

```plaintext
+---------------------+       +---------------------+       +---------------------+
|      ParkingLot     |       |     ParkingFloor    |       |      ParkingSpot    |
+---------------------+       +---------------------+       +---------------------+
| - floors: []Floor   |1     *| - spots: [][]Spot   |1     *| - type: SpotType    |
| - parkedVehicles    |<----->| - floorNum: int     |<----->| - isActive: bool    |
| - vehicleHistory    |       | - numRows: int      |       | - isOccupied: bool  |
| - mutex: RWMutex    |       | - numCols: int      |       | - vehicleNumber: str|
+---------------------+       +---------------------+       | - row: int          |
| + Park()            |       | + GetSpot()         |       | - column: int       |
| + Unpark()          |       | + GetAvailableSpots |       +---------------------+
| + AvailableSpot()   |       | + FindVehicle()     |       | + GetSpotID()       |
| + SearchVehicle()   |       +---------------------+       | + Occupy()          |
+---------------------+                                     | + Vacate()          |
                                                            | + CanPark()         |
                                                            +---------------------+
                                                                      ^
                                                                      |
+---------------------+       +---------------------+                 |
|      Vehicle        |       |     SpotType        |<----------------+
+---------------------+       +---------------------+
| - type: VehicleType |       | BICYCLE (B-1)       |
| - number: string    |       | MOTORCYCLE (M-1)    |
+---------------------+       | AUTOMOBILE (A-1)    |
                              | INACTIVE (X-0)      |
+---------------------+       +---------------------+
|    VehicleType      |
+---------------------+
| BICYCLE             |
| MOTORCYCLE          |
| AUTOMOBILE          |
+---------------------+
```

## Installation

### Prerequisites

- Go 1.24.1 or higher

### Building from Source

1. Clone the repository:

   ```bash
   git clone https://github.com/prasaria/go-multistorey-parking-lot.git
   cd parking-lot
   ```

2. Build the application:

   ```bash
   make build
   ```

3. Run the application:

   ```bash
   make run
   ```

### Using Go Run

You can also run the application directly without building:

```bash
go run ./cmd/parking-lot
```

## Usage

### Interactive Mode

Start the application in interactive mode:

```bash
bin/parking-lot
```

This will present you with a prompt where you can enter commands:

```bash
Welcome to Parking Lot CLI
Type 'help' to see available commands or 'exit' to quit
Options:
  --json    Output results in JSON format
  --verbose Show detailed operation logs

> 
```

### Available Commands

#### Initialize Parking Lot

Create a new parking lot with specified dimensions:

```bash
> init <floors> <rows> <columns>
```

Example:

```bash
> init 3 5 10
```

#### Park Vehicle

Park a vehicle in the lot:

```bash
> park <vehicle_type> <vehicle_number>
```

Example:

```bash
> park automobile KA-01-HH-1234
```

#### Unpark Vehicle

Remove a vehicle from its parking spot:

```bash
> unpark <spot_id> <vehicle_number>
```

Example:

```bash
> unpark 1-2-3 KA-01-HH-1234
```

#### Find Available Spots

Display available spots for a vehicle type:

```bash
> available <vehicle_type>
```

Example:

```bash
> available motorcycle
```

#### Search Vehicle

Search for a vehicle by its number:

```bash
> search <vehicle_number>
```

Example:

```bash
> search KA-01-HH-1234
```

#### Check Status

Display the current status of the parking lot:

```bash
> status
```

#### Help

Display help information:

```bash
> help
```

### JSON Output

You can append `--json` to any command to get the output in JSON format:

```bash
> available bicycle --json
```

### Verbose Logging

Use the `--verbose` or `-v` flag to see detailed operation logs:

```bash
> park automobile KA-01-HH-1234 --verbose
```

## Constraints

- 1 <= floors <= 8
- 1 <= rows <= 1000
- 1 <= columns <= 1000
- Each floor will have the same number of rows
- Each rows will have the same number of columns
- Each parking spot is of the following type:
  - "B-1", active for bicycles
  - "M-1", active for motorcycles
  - "A-1", active for automobiles
  - "X-0", inactive

## Development

### Project Structure

```tree
parking-lot/
├── bin/                      # Compiled binaries
├── cmd/
│   └── parking-lot/          # Main application
│       ├── commands.go       # Command handling
│       ├── json_output.go    # JSON output formatting
│       ├── logger.go         # Logging utilities
│       ├── main.go           # Entry point
│       └── output.go         # Text output formatting
├── internal/
│   ├── model/                # Domain models
│   │   ├── parking_lot.go    # Parking lot implementation
│   │   ├── parking_floor.go  # Floor implementation
│   │   ├── parking_spot.go   # Spot implementation
│   │   ├── spot_type.go      # Spot type definitions
│   │   ├── vehicle.go        # Vehicle implementation
│   │   └── vehicle_type.go   # Vehicle type definitions
│   └── errors/               # Custom error types
├── test/                     # Integration tests
├── .gitignore
├── Makefile                  # Build commands
├── go.mod
└── README.md
```

### Testing

Run all tests:

```bash
make test
```

Run tests with coverage:

```bash
make coverage
```

Run long tests:

```bash
make test-long
```

Run performance tests:

```bash
make test-perf
```

### Building

Build for your current platform:

```bash
make build
```

Build for multiple platforms:

```bash
make build-all
```

## License

This project is licensed under the MIT License .

## Acknowledgments

- This project was created as a learning exercise for Go programming
- Inspired by real-world parking management systems
