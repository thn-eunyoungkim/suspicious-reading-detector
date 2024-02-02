package detector

import (
	"main.go/models"
)

type ReadingDetector struct {
}

func NewReadingDetector() *ReadingDetector {
	return &ReadingDetector{}
}

func (rd *ReadingDetector) DetectSuspiciousReadings(readingCollection *models.ReadingCollection) (*models.ReadingCollection, error) {
	// Create a new ReadingCollection to store suspicious readings
	suspiciousReadingCollection := models.NewReadingCollection()

	// Get the median value from the original readingCollection
	median, err := readingCollection.Median()
	if err != nil {
		return nil, err
	}

	// Iterate through the readings in the original collection
	for _, reading := range readingCollection.Readings {
		// Check if the reading is suspicious based on the median
		if models.IsSuspicious(reading.Consumption, median) {
			// Add suspicious reading to the new collection
			suspiciousReadingCollection.AddReading(reading)
		}
	}

	return suspiciousReadingCollection, nil
}
