package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/prasaria/go-multistorey-parking-lot/internal/model"
)

// JSON output structures

// JSONResult is the root structure for JSON output
type JSONResult struct {
	Success bool        `json:"success"`
	Command string      `json:"command"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Time    string      `json:"time"`
}

// InitResult contains data for init command output
type InitResult struct {
	Floors  int            `json:"floors"`
	Rows    int            `json:"rows"`
	Columns int            `json:"columns"`
	Total   int            `json:"totalSpots"`
	Counts  map[string]int `json:"spotCounts"`
}

// ParkResult contains data for park command output
type ParkResult struct {
	VehicleType   string `json:"vehicleType"`
	VehicleNumber string `json:"vehicleNumber"`
	SpotID        string `json:"spotId"`
}

// UnparkResult contains data for unpark command output
type UnparkResult struct {
	VehicleNumber string `json:"vehicleNumber"`
	SpotID        string `json:"spotId"`
}

// AvailableResult contains data for available command output
type AvailableResult struct {
	VehicleType string   `json:"vehicleType"`
	SpotIDs     []string `json:"spotIds"`
	Count       int      `json:"count"`
}

// SearchResult contains data for search command output
type SearchResult struct {
	VehicleNumber string `json:"vehicleNumber"`
	SpotID        string `json:"spotId"`
	IsParked      bool   `json:"isParked"`
}

// StatusResult contains data for status command output
type StatusResult struct {
	Name            string            `json:"name"`
	Floors          int               `json:"floors"`
	TotalSpots      int               `json:"totalSpots"`
	ActiveSpots     int               `json:"activeSpots"`
	OccupiedSpots   int               `json:"occupiedSpots"`
	AvailableSpots  int               `json:"availableSpots"`
	SpotCounts      map[string]int    `json:"spotCounts"`
	AvailableCounts map[string]int    `json:"availableCounts"`
	ParkedVehicles  map[string]string `json:"parkedVehicles"`
}

// Helper functions

// PrintJSON outputs a result as JSON
func PrintJSON(command string, data interface{}, err error) {
	result := JSONResult{
		Success: err == nil,
		Command: command,
		Time:    time.Now().Format(time.RFC3339),
	}

	if err != nil {
		result.Error = err.Error()
	} else {
		result.Data = data
	}

	// Marshal to JSON
	jsonBytes, jsonErr := json.MarshalIndent(result, "", "  ")
	if jsonErr != nil {
		fmt.Printf("Error marshaling JSON: %v\n", jsonErr)
		return
	}

	fmt.Println(string(jsonBytes))
}

// Convert SpotType map to string map for JSON
func convertSpotTypeMap(m map[model.SpotType]int) map[string]int {
	result := make(map[string]int)
	for k, v := range m {
		result[string(k)] = v
	}
	return result
}

// Convert VehicleType map to string map for JSON
func convertVehicleTypeMap(m map[model.VehicleType]int) map[string]int {
	result := make(map[string]int)
	for k, v := range m {
		result[string(k)] = v
	}
	return result
}
