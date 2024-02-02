package writer

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"main.go/models"
)

type ReadingWriter struct {
}

func NewReadingWriter() *ReadingWriter {
	return &ReadingWriter{}
}

func (rw *ReadingWriter) WriteReadingsToFile(filename string, readings *models.ReadingCollection) error {
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

	median, err := readings.Median()

	// Write data
	for _, sr := range readings.Readings {
		record := []string{
			sr.CustomerId,
			sr.YearMonth,
			strconv.Itoa(sr.Consumption),
			fmt.Sprintf("%.2f", float64(median)),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing record: %w", err)
		}
	}

	return nil
}

// func writeSuspiciousReadingsToFile(filename string, readings []SuspiciousReading) error {
// 	file, err := os.Create(filename)
// 	if err != nil {
// 		return fmt.Errorf("error creating file: %w", err)
// 	}
// 	defer file.Close()

// 	writer := csv.NewWriter(file)
// 	defer writer.Flush()

// 	// Write header
// 	if err := writer.Write([]string{"Client", "Month", "Suspicious", "Median"}); err != nil {
// 		return fmt.Errorf("error writing header: %w", err)
// 	}

// 	// Write data
// 	for _, sr := range readings {
// 		record := []string{
// 			sr.Reading.CustomerId,
// 			sr.Reading.YearMonth,
// 			strconv.Itoa(sr.Reading.Consumption),
// 			fmt.Sprintf("%.2f", float64(sr.Median)),
// 		}
// 		if err := writer.Write(record); err != nil {
// 			return fmt.Errorf("error writing record: %w", err)
// 		}
// 	}

// 	return nil
// }
