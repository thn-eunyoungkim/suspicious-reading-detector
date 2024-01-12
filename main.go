package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

type Reading struct {
	CustomerId  string
	YearMonth   string
	Consumption int
}

type ReadingCollection struct {
	Readings []Reading
}

type SuspiciousReading struct {
	Reading Reading
	Median  int
}

type CustomerId string

func GetReadingsByCustomerIdFromCsv(filePath string) (map[CustomerId]ReadingCollection, error) {
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
			log.Fatal(err)
		}
		rawData = append(rawData, records)
		// fmt.Println(records)
	}

	// Readings map initialization
	byCustomerId := make(map[CustomerId]ReadingCollection)

	// Process CSV records
	for _, row := range rawData {
		customerId := CustomerId(row[0])
		yearMonth := row[1]
		consumption, _ := strconv.Atoi(row[2])
		reading := Reading{CustomerId: string(customerId), YearMonth: yearMonth, Consumption: consumption}

		// Check if CustomerId exists in the map
		if collection, ok := byCustomerId[customerId]; !ok {
			byCustomerId[customerId] = ReadingCollection{Readings: []Reading{reading}}
		} else {
			collection.Readings = append(collection.Readings, reading)
			byCustomerId[customerId] = collection
		}
	}
	return byCustomerId, nil
}

func GetReadingsByCustonerIdFromXml(filePath string) (map[CustomerId]ReadingCollection, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening XML file: %w", err)
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)

	// Readings map initialization
	byCustomerId := make(map[CustomerId]ReadingCollection)

	var currentReading Reading
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("error reading XML tokens: %w", err)
		}

		switch token := token.(type) {
		case xml.StartElement:
			switch token.Name.Local {
			case "reading":
				// Start of a new Reading element, reset currentReading
				currentReading = Reading{}

				// Extract attributes (clientID, period)
				for _, attr := range token.Attr {
					switch attr.Name.Local {
					case "clientID":
						currentReading.CustomerId = attr.Value
					case "period":
						currentReading.YearMonth = attr.Value
					}
				}

				// Read the content of the reading element
				if err := decoder.DecodeElement(&currentReading.Consumption, &token); err != nil {
					return nil, fmt.Errorf("error decoding Consumption: %w", err)
				}

				// Add the currentReading to the map
				customerID := CustomerId(currentReading.CustomerId)
				if collection, ok := byCustomerId[customerID]; !ok {
					byCustomerId[customerID] = ReadingCollection{Readings: []Reading{currentReading}}
				} else {
					collection.Readings = append(collection.Readings, currentReading)
					byCustomerId[customerID] = collection
				}
			}
		}
	}

	return byCustomerId, nil
}

func (c ReadingCollection) Median() (int, bool) {
	n := len(c.Readings)
	if n < 6 {
		return 0, false
	}

	// Sort the readings by consumption
	sortedReadings := make([]int, n)
	for i, reading := range c.Readings {
		sortedReadings[i] = reading.Consumption
	}
	sort.Ints(sortedReadings)

	// Calculate the median value as the average of the two middle values
	medianIndex := (n - 1) / 2
	median := sortedReadings[medianIndex]
	return median, true
}

func obtainReadings(filePath string) (map[CustomerId]ReadingCollection, error) {
	// cvs or xml
	fileType := filepath.Ext(filePath)

	switch fileType {
	case ".csv":
		readings, err := GetReadingsByCustomerIdFromCsv(filePath)
		if err != nil {
			return nil, fmt.Errorf("error obtaining readings from CSV: %w", err)
		}
		return readings, nil
	case ".xml":
		readings, err := GetReadingsByCustonerIdFromXml(filePath)
		if err != nil {
			return nil, fmt.Errorf("error obtaining readings from XML: %w", err)
		}
		return readings, nil
	default:
		return nil, fmt.Errorf("unsupported file type: %s", fileType)
	}
}

// TODO separate this logic
// TODO we need an object inclues Median
// TODO indtroduce a new object a list of suspicious readings
// Reading, Median -> SuspiciousReading (Reading, Median)
// function - readings - filter  return weird readings
// create a function that creates a file with a suspicious readings
func GetSuspiciousReadings(readings map[CustomerId]ReadingCollection) []SuspiciousReading {
	suspiciousReadings := make([]SuspiciousReading, 0)

	for _, collection := range readings {
		median, ok := collection.Median()
		if !ok {
			continue
		}
		for _, reading := range collection.Readings {
			if isSuspicious(reading.Consumption, median) {
				suspiciousReadings = append(suspiciousReadings, SuspiciousReading{Reading: reading, Median: median})
			}
		}
	}

	return suspiciousReadings
}

func isSuspicious(consumption, median int) bool {
	const lowerThreshold = 0.5
	const upperThreshold = 1.5

	return float64(consumption) > upperThreshold*float64(median) || float64(consumption) < lowerThreshold*float64(median)
}

func printSuspiciousReadings(readings []Reading, median int) {

	fmt.Printf("| %-15s | %-10s | %-10s | %-10s |\n", "Client", "Month", "Suspicious", "Median")
	fmt.Println(" -------------------------------------------------------------------------------")

	for _, reading := range readings {

		// Check if the reading is either higher or lower than the annual median ± 50%
		fmt.Printf("| %-15s | %-10s | %-10d | %-10.2f |\n", reading.CustomerId, reading.YearMonth, reading.Consumption, float64(median))

	}
}

func writeSuspiciousReadingsToFile(filename string, readings []SuspiciousReading) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"Client", "Month", "Suspicious", "Median"}); err != nil {
		return fmt.Errorf("error writing header: %w", err)
	}

	// Write data
	for _, sr := range readings {
		record := []string{
			sr.Reading.CustomerId,
			sr.Reading.YearMonth,
			strconv.Itoa(sr.Reading.Consumption),
			fmt.Sprintf("%.2f", float64(sr.Median)),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing record: %w", err)
		}
	}

	return nil
}

// filtering the readings - we need to have filtered readings

func main() {
	filePath := "inputs/2016-readings.xml"

	readings, err := obtainReadings(filePath)
	if err != nil {
		log.Fatal(err)
	}

	filteredReadings := GetSuspiciousReadings(readings)

	if err := writeSuspiciousReadingsToFile("suspicious_readings.csv", filteredReadings); err != nil {
		log.Fatal(err)
	}
}

// OOP
// refactoring
// concepts : events & messages
// event sourcing - state : whole story of what happened instead of keeping the last state
