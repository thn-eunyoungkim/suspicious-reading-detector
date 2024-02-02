package models

import (
	"errors"
	"sort"
)

type Reading struct {
	CustomerId  string
	YearMonth   string
	Consumption int
}

type SuspiciousReading struct {
	Reading Reading
	Median  int
}

type CustomerId string

type ReadingsReader interface {
	GetReadings(filePath string) (*ReadingCollection, error)
}

type ReadingCollection struct {
	Readings []Reading
}

func NewReadingCollection() *ReadingCollection {
	return &ReadingCollection{}
}

func (c *ReadingCollection) AddReading(reading Reading) {
	c.Readings = append(c.Readings, reading)
}

func (c *ReadingCollection) Median() (int, error) {
	n := len(c.Readings)
	if n < 6 {
		return 0, errors.New("not enough readings")
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
	return median, nil
}

func (c *ReadingCollection) GetSuspiciousReadings(customerIds ...CustomerId) []SuspiciousReading {
	suspiciousReadings := make([]SuspiciousReading, 0)

	median, err := c.Median()
	if err != nil {
		panic(err)
	}

	if len(customerIds) == 0 {
		for _, reading := range c.Readings {
			if IsSuspicious(reading.Consumption, median) {
				suspiciousReadings = append(suspiciousReadings, SuspiciousReading{Reading: reading, Median: median})
			}
		}
	} else {
		for _, customerId := range customerIds {
			for _, reading := range c.Readings {
				if reading.CustomerId == string(customerId) && IsSuspicious(reading.Consumption, median) {
					suspiciousReadings = append(suspiciousReadings, SuspiciousReading{Reading: reading, Median: median})
				}
			}
		}
	}

	return suspiciousReadings
}

func IsSuspicious(consumption, median int) bool {
	const lowerThreshold = 0.5
	const upperThreshold = 1.5

	return float64(consumption) > upperThreshold*float64(median) || float64(consumption) < lowerThreshold*float64(median)
}
