# Multi-Storey Parking Lot Management System

A robust command-line interface (CLI) application designed to manage a multi-storey parking facility with support for various vehicle types, written in Go 1.24.1.

## Project Overview

This system simulates a real-world parking lot with multiple floors, rows, and columns of parking spots. It's designed to handle concurrent access through multiple entry/exit gates while maintaining data consistency.

## Key Features

- Configurable parking layout with support for up to 8 floors
- Dedicated spots for different vehicle types (bicycles, motorcycles, automobiles)
- Thread-safe operations for concurrent access
- Customizable spot allocation strategy
- Comprehensive vehicle tracking and reporting

## System Requirements

- Go 1.24.1 or higher
- No external dependencies for core functionality

## Installation

```bash
# Get the code
git clone https://github.com/prasaria/go-multistorey-parking-lot.git

# Move into the project directory
cd go-multistorey-parking-lot

# Build the application
go build -o bin/parking-lot cmd/go-multistorey-parking-lot/main.go

# Run the application
./bin/parking-lot
```

## Usage Examples

### Initializing a Parking Lot

Start a new parking session by configuring the lot dimensions:

```bash
# Create a default parking lot (3 floors, 5 rows, 10 columns)
./bin/parking-lot init

# Create a custom-sized parking lot
./bin/parking-lot init -floors 4 -rows 8 -columns 15
```

### Parking Operations

```bash
# Park a vehicle
./bin/parking-lot park -type automobile -number "KA01AB1234"

# Remove a vehicle from a spot
./bin/parking-lot unpark -spot "2-3-5" -number "KA01AB1234"

# Find available parking spots for a specific vehicle type
./bin/parking-lot available -type motorcycle

# Locate a vehicle in the parking lot
./bin/parking-lot search -number "KA01AB1234"
```

## Design Approach

The parking lot system follows a modular design with several key components:

- **CLI Interface**: Handles user commands and displays results
- **Parking Manager**: Core component managing all parking operations
- **Spot Allocator**: Assigns optimal spots based on vehicle type
- **Data Store**: Maintains the state of all parking spots
- **Concurrency Control**: Ensures thread-safety for multiple gates

## Thread Safety

The system implements a comprehensive concurrency control strategy using read-write mutexes to ensure that multiple parking gates can operate simultaneously without data corruption. Read operations can proceed in parallel, while write operations obtain exclusive locks.

## Project Structure

```tree
parking-lot/
├── cmd/                  # Application entrypoints
├── internal/             # Private application code
│   ├── model/            # Domain models
│   ├── service/          # Business logic
│   └── errors/           # Custom error types
├── pkg/                  # Public library code
│   ├── config/           # Configuration handling
│   └── utils/            # Utility functions
└── test/                 # Integration tests
```

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with race detector
go test -race ./...

# Run specific test package
go test ./pkg/config
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
