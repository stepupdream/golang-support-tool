package csv

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
)

// LoadCsv Reading CSV files
func LoadCsv(filepath string, isFilter bool) [][]string {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("CSVFileOpenError: ", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("CSVFileCloseError: ", err)
		}
	}(file)

	// If BOM is included, delete the BOM
	// https://pinzolo.github.io/2017/03/29/utf8-csv-with-bom-on-golang.html
	reader := bufio.NewReader(file)
	bytes, err := reader.Peek(3)
	if err != nil {
		log.Fatal("CSVFileNewReaderError: ", err)
	} else if bytes[0] == 0xEF && bytes[1] == 0xBB && bytes[2] == 0xBF {
		_, err := reader.Discard(3)
		if err != nil {
			log.Fatal("CSVFileDiscardError: ", err)
		}
	}

	csvReader := csv.NewReader(reader)
	if isFilter {
		csvReader.Comment = '#'
	}
	rows, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("CSVReadAllError: ", err)
	}

	return rows
}
