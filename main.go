package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/jedib0t/go-pretty/v6/table"
)

func main() {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return
	}

	// Construct the full path to the CSV file in the "inputs" folder
	filePath := filepath.Join(currentDir, "inputs", "2016-readings.csv")

	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all records from the CSV file
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV records:", err)
		return
	}

	// Create a new table
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// Convert the header (records[0]) to table.Row
	headerRow := make(table.Row, len(records[0]))
	for i, header := range records[0] {
		headerRow[i] = header
	}

	// Set column names based on the first row of the CSV file
	t.AppendHeader(headerRow)

	// Append data rows to the table
	for _, record := range records[1:] {
		// Convert each data row to table.Row
		dataRow := make(table.Row, len(record))
		for i, value := range record {
			dataRow[i] = value
		}

		t.AppendRow(dataRow)
	}

	// Render the table
	t.Render()

	if len(records) > 6 {
		reading5, _ := strconv.Atoi(records[5][2])
		reading6, _ := strconv.Atoi(records[6][2])
		average := (reading5 + reading6) / 2

		headerRow := make(table.Row, len(records[0]))
		for i, header := range records[0] {
			headerRow[i] = header
			fmt.Println(header)
		}

		// Render the table
		t.Render()
		fmt.Printf("Average of readings at indices 5 and 6: %d\n", average)
	} else {
		fmt.Println("Insufficient data for calculating the average at indices 5 and 6")
	}
}
