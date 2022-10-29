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
	"github.com/stepupdream/golang-support-tool/directory"
)

// Key Make keys into structures to achieve multidimensional arrays.
type Key struct {
	Id  int
	Key string
}

// LoadCsv Reading CSV files
func LoadCsv(filepath string, isRowFilter bool, isColumnFilter bool) [][]string {
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
	if isRowFilter {
		csvReader.Comment = '#'
	}
	rows, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("CSVReadAllError: ", err)
	}

	if isColumnFilter {
		return filterColumn(rows, isColumnFilter)
	}

	return rows
}

func filterColumn(rows [][]string, isFilter bool) [][]string {
	var disableColumnIndexes []int
	for index, value := range rows[0] {
		if isFilter && value == "#" {
			disableColumnIndexes = append(disableColumnIndexes, index)
		}
	}

	var newRows [][]string
	for _, row := range rows {
		var newRow []string
		for index, value := range row {
			if !array.IntContains(disableColumnIndexes, index) {
				newRow = append(newRow, value)
			}
		}
		newRows = append(newRows, newRow)
	}

	return newRows
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

			if len(columnNumbers) == 0 || array.IntContains(columnNumbers, columnNumber) {
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

func GetFilePathRecursive(path string) ([]string, error) {
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

func DeleteCSV(baseCSV map[Key]string, editCSV map[Key]string) map[Key]string {
	baseIds := PluckId(baseCSV)

	for key, _ := range editCSV {
		if key.Key == "id" {
			if !array.IntContains(baseIds, key.Id) {
				log.Fatal("Attempted to delete a non-existent ID : id ", key.Id)
			}
		}
		delete(baseCSV, Key{Id: key.Id, Key: key.Key})
	}

	return baseCSV
}

func InsertCSV(baseCSV map[Key]string, editCSV map[Key]string) map[Key]string {
	baseIds := PluckId(baseCSV)
	editIds := PluckId(editCSV)

	for _, id := range editIds {
		if array.IntContains(baseIds, id) {
			log.Fatal("Tried to do an insert on an existing ID : id ", id)
		}
	}

	result := make(map[Key]string)

	for mapKey, value := range baseCSV {
		result[Key{Id: mapKey.Id, Key: mapKey.Key}] = value
	}
	for mapKey, value := range editCSV {
		result[Key{Id: mapKey.Id, Key: mapKey.Key}] = value
	}

	return result
}

func UpdateCSV(baseCSV map[Key]string, editCSV map[Key]string) map[Key]string {
	baseIds := PluckId(baseCSV)
	editIds := PluckId(editCSV)
	for _, id := range editIds {
		if !array.IntContains(baseIds, id) {
			log.Fatal("Tried to update a non-existent ID : id ", id)
		}
	}

	baseCSV = DeleteCSV(baseCSV, editCSV)
	baseCSV = InsertCSV(baseCSV, editCSV)

	return baseCSV
}

func LoadFileFirstContent(directoryPath string, fileName string) string {
	if !directory.Exist(directoryPath) {
		log.Fatal("The directory could not be found : ", directoryPath)
	}
	baseCsvFilePaths, err := GetFilePathRecursive(directoryPath)
	if err != nil {
		log.Fatal("LoadFileFirstContentError: ", err)
	}

	if len(baseCsvFilePaths) == 0 {
		return ""
	}

	var result string
	for _, path := range baseCsvFilePaths {
		if filepath.Base(path) == fileName {
			rows := LoadCsv(path, true, false)
			row := rows[0]
			result = row[0]
			break
		}
	}

	if result == "" {
		log.Fatal("The content could not be found : ", fileName)
	}

	return result
}
