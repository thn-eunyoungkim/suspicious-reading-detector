package readers

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"main.go/models"
)

type XMLReader struct {
}

func NewXMLReader() *XMLReader {
	return &XMLReader{}
}

func (x *XMLReader) GetReadings(filePath string) (*models.ReadingCollection, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening XML file: %w", err)
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)

	readingCollection := models.NewReadingCollection()

	var currentReading models.Reading
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
				currentReading = models.Reading{}

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

				readingCollection.AddReading(currentReading)
			}
		}
	}

	return readingCollection, nil
}
