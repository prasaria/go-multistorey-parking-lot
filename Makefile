# Makefile for Multi-Storey Parking Lot CLI System

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOVET=$(GOCMD) vet
GOFMT=gofmt
BINARY_NAME=parking-lot
BINARY_UNIX=$(BINARY_NAME)_unix
MAIN_DIR=./cmd/go-multistorey-parking-lot

# Build flags
LDFLAGS=-ldflags "-X main.Version=1.0.0"

# Make targets
.PHONY: all build clean test coverage fmt vet run help tidy vendor

all: test build

build:
	$(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME) $(MAIN_DIR)

clean:
	$(GOCLEAN)
	rm -rf bin/

test:
	$(GOTEST) -v ./...

# Run test with coverage
coverage:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

# Format code
fmt:
	$(GOFMT) -w .

# Run Go vet
vet:
	$(GOVET) ./...

# Run the application
run:
	$(GORUN) $(MAIN_DIR)

# Initialize the go module
init-mod:
	$(GOMOD) init github.com/yourusername/parking-lot

# Update dependencies
tidy:
	$(GOMOD) tidy

# Vendor dependencies
vendor:
	$(GOMOD) vendor

# Build for multiple platforms
build-all: build-linux build-windows build-mac

build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME)_linux_amd64 $(MAIN_DIR)

build-windows:
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME)_windows_amd64.exe $(MAIN_DIR)

build-mac:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME)_darwin_amd64 $(MAIN_DIR)

# Run long tests
test-long:
	$(GOTEST) -v -run=. -long ./...

# Run performance tests
test-perf:
	$(GOTEST) -v -run=. -perf ./...

# Run race detector
test-race:
	$(GOTEST) -race -v ./...

# Display help
help:
	@echo "Multi-Storey Parking Lot CLI System"
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all          Run tests and build binary"
	@echo "  build        Build the binary"
	@echo "  clean        Remove binary and coverage files"
	@echo "  test         Run all tests"
	@echo "  coverage     Generate test coverage report"
	@echo "  fmt          Format code"
	@echo "  vet          Run Go vet"
	@echo "  run          Run the application"
	@echo "  tidy         Update dependencies"
	@echo "  vendor       Vendor dependencies"
	@echo "  build-all    Build for multiple platforms"
	@echo "  test-long    Run long tests"
	@echo "  test-perf    Run performance tests"
	@echo "  test-race    Run tests with race detector"
	@echo "  help         Display this help"