package readers

import (
	"fmt"
	"path/filepath"
	"strings"

	"main.go/models"
)

type FileReader struct {
	csvReader *CSVReader
	xmlReader *XMLReader
}

func NewFileReader() *FileReader {
	return &FileReader{
		csvReader: NewCSVReader(),
		xmlReader: NewXMLReader(),
	}
}

func (fr *FileReader) ReadReadings(filePath string) (*models.ReadingCollection, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".csv":
		return fr.csvReader.GetReadings(filePath)
	case ".xml":
		return fr.xmlReader.GetReadings(filePath)
	default:
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}
}

// TODO: new function with for loop the multiple files
// TODO: combine the results
// ? try to find a way not to use pointers
// func obtainReadings(filePath string) (map[models.CustomerId]models.ReadingCollection, error) {
// 	// cvs or xml
// 	fileType := filepath.Ext(filePath)

// 	switch fileType {
// 	case ".csv":
// 		readings, err := GetReadingsByCustomerIdFromCsv(filePath)
// 		if err != nil {
// 			return nil, fmt.Errorf("error obtaining readings from CSV: %w", err)
// 		}
// 		return readings, nil
// 	case ".xml":
// 		readings, err := GetReadingsByCustonerIdFromXml(filePath)
// 		if err != nil {
// 			return nil, fmt.Errorf("error obtaining readings from XML: %w", err)
// 		}
// 		return readings, nil
// 	default:
// 		return nil, fmt.Errorf("unsupported file type: %s", fileType)
// 	}
// }
