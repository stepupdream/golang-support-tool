package csv

import (
	"bufio"
	"encoding/csv"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"
	"github.com/stepupdream/golang-support-tool/array"
)

// Key Make keys into structures to achieve multidimensional arrays.
type Key struct {
	Id  int
	Key string
}

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

// ConvertMap
// Replacing CSV data (two-dimensional array of height and width) into a multidimensional associative array in a format
// that facilitates direct value specification by key.
func ConvertMap(rows [][]string, columnNumbers []int, filepath string) map[Key]string {
	result := make(map[Key]string)
	keyName := map[int]string{}
	findIdColumn := false
	idColumnNumber := 0

	for rowNumber, row := range rows {
		for columnNumber, value := range row {
			// The first line is the key.
			if rowNumber == 0 {
				if value == "id" {
					findIdColumn = true
					idColumnNumber = columnNumber
				}
				keyName[columnNumber] = value
				continue
			}

			if array.IntContains(columnNumbers, columnNumber) {
				id, _ := strconv.Atoi(row[idColumnNumber])

				if _, flg := result[Key{id, keyName[columnNumber]}]; flg {
					log.Fatal("ID is not unique : ", filepath)
				}
				result[Key{id, keyName[columnNumber]}] = value
			}
		}
	}

	if !findIdColumn {
		log.Fatal("CSV without ID column cannot be read : ", filepath)
	}

	return result
}

func PluckId(csv map[Key]string) []int {
	var ids []int

	for mapKey, _ := range csv {
		if mapKey.Key == "id" {
			ids = append(ids, mapKey.Id)
		}
	}
	return ids
}

func GetCSVFilePaths(path string) ([]string, error) {
	var paths []string

	// Recursively retrieve directories and files. (use WalkDir since Walk is now deprecated)
	err := filepath.WalkDir(path, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return errors.Wrap(err, "failed filepath.WalkDir")
		}

		if info.IsDir() {
			return nil
		}

		extension := filepath.Ext(path)
		if extension != ".csv" {
			return nil
		}

		paths = append(paths, path)

		return nil
	})

	return paths, err
}

func NewFile(path string, rows [][]string) {
	// create allows you to create a new file and overwrite a new file.
	csvFile, err := os.Create(path)
	if err != nil {
		log.Fatalln(err)
	}

	// Make it with BOM to avoid garbled characters.
	_, err = csvFile.Write([]byte{0xEF, 0xBB, 0xBF})
	if err != nil {
		log.Fatal(err)
	}

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	if err := writer.WriteAll(rows); err != nil {
		log.Fatal(err)
	}
}
