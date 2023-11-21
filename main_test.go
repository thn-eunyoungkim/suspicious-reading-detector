package main

import (
	"encoding/csv"
	"os"
	"sort"
	"strconv"
	"testing"
)

func TestCSVReadingAndAverage(t *testing.T) {
	// Assuming your CSV file is named "yourfile.csv" and is in the "inputs" folder
	filePath := "inputs/2016-readings.csv"

	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("Error opening CSV file: %v", err)
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all records from the CSV file
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Error reading CSV records: %v", err)
	}

	// Check if there are at least two rows (header + data)
	if len(records) < 2 {
		t.Fatalf("Expected at least two rows in the CSV file, but got %d", len(records))
	}

	// Sort the data array based on the first column (assuming the first column is numeric)
	sort.SliceStable(records[1:], func(i, j int) bool {
		valI, _ := strconv.Atoi(records[1:][i][0])
		valJ, _ := strconv.Atoi(records[1:][j][0])
		return valI < valJ
	})

	// Calculate the average of readings at indices 5 and 6
	if len(records) > 6 {
		reading5, _ := strconv.Atoi(records[5][2])
		reading6, _ := strconv.Atoi(records[6][2])
		average := (reading5 + reading6) / 2

		// Add any additional assertions or checks based on your requirements
		if average < 0 {
			t.Fatalf("Expected average to be non-negative, but got %d", average)
		}
	} else {
		t.Fatal("Insufficient data for calculating the average at indices 5 and 6")
	}
}
