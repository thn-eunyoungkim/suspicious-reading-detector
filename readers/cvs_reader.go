package readers

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"main.go/models"
)

type CSVReader struct {
}

func NewCSVReader() *CSVReader {
	return &CSVReader{}
}

func (c *CSVReader) GetReadings(filePath string) (*models.ReadingCollection, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening CSV file: %w", err)
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)
	rawData := make([][]string, 0)

	// Read all CSV records
	for {
		records, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("error reading CSV records: %w", err) // fail fast
		}
		rawData = append(rawData, records)
		// fmt.Println(records)
	}

	readingCollection := models.NewReadingCollection()
	// Process CSV records
	for _, row := range rawData {
		customerId := row[0]
		yearMonth := row[1]
		consumption, _ := strconv.Atoi(row[2])
		reading := models.Reading{CustomerId: string(customerId), YearMonth: yearMonth, Consumption: consumption}

		readingCollection.AddReading(reading)
	}
	return readingCollection, nil
}
