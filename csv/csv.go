package csv

import (
	"bufio"
	"encoding/csv"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/pkg/errors"
	"github.com/stepupdream/golang-support-tool/array"
	"github.com/stepupdream/golang-support-tool/directory"
	supportFile "github.com/stepupdream/golang-support-tool/file"
)

// Key Make keys into structures to achieve multidimensional arrays.
type Key struct {
	Id  int
	Key string
}

func LoadCsvMap(filePath string, filterNames []string, isColumnExclusion bool) map[Key]string {
	var rows [][]string
	if !supportFile.Exists(filePath) {
		return make(map[Key]string)
	}

	var filterColumnNumbers []int
	rows = LoadCsv(filePath, true, isColumnExclusion)
	if len(filterNames) != 0 {
		filterColumnNumbers = FilterColumnNumbers(filePath, filterNames)
	}

	return ConvertMap(rows, filterColumnNumbers, filePath)
}

// LoadCsv Reading CSV files
func LoadCsv(filepath string, isRowExclusion bool, isColumnExclusion bool) [][]string {
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
	if isRowExclusion {
		csvReader.Comment = '#'
	}
	rows, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("CSVReadAllError: ", err)
	}

	if isColumnExclusion {
		return exclusionColumn(rows, isColumnExclusion)
	}

	return rows
}

func exclusionColumn(rows [][]string, isExclusion bool) [][]string {
	var disableColumnIndexes []int
	for index, value := range rows[0] {
		if isExclusion && value == "#" {
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
func ConvertMap(rows [][]string, filterColumnNumbers []int, filepath string) map[Key]string {
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

			if len(filterColumnNumbers) != 0 && !array.IntContains(filterColumnNumbers, columnNumber) {
				continue
			}

			id, _ := strconv.Atoi(row[idColumnNumber])
			if _, flg := result[Key{id, keyName[columnNumber]}]; flg {
				log.Fatal("ID is not unique : ", filepath)
			}
			if value == "" {
				log.Fatal("Blank space is prohibited because it is impossible to determine if you forgot to enter the information. : ", filepath, " rowNumber : ", rowNumber)
			}
			result[Key{id, keyName[columnNumber]}] = value
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

	sort.Ints(ids)

	return ids
}

func PluckKey(csv map[Key]string, key string) []string {
	var values []string

	for mapKey, value := range csv {
		if mapKey.Key == key {
			values = append(values, value)
		}
	}

	return values
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

func DeleteCSV(baseCSV map[Key]string, editCSV map[Key]string, csvFilePath string) map[Key]string {
	baseIds := PluckId(baseCSV)

	for key, _ := range editCSV {
		if key.Key == "id" {
			if !array.IntContains(baseIds, key.Id) {
				log.Fatal("Attempted to delete a non-existent ID : id ", key.Id, " ", csvFilePath)
			}
		}
		delete(baseCSV, Key{Id: key.Id, Key: key.Key})
	}

	return baseCSV
}

func InsertCSV(baseCSV map[Key]string, editCSV map[Key]string, csvFilePath string) map[Key]string {
	baseIds := PluckId(baseCSV)
	editIds := PluckId(editCSV)

	for _, id := range editIds {
		if array.IntContains(baseIds, id) {
			log.Fatal("Tried to do an insert on an existing ID : id ", id, " ", csvFilePath)
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

func UpdateCSV(baseCSV map[Key]string, editCSV map[Key]string, csvFilePath string) map[Key]string {
	baseIds := PluckId(baseCSV)
	editIds := PluckId(editCSV)
	for _, id := range editIds {
		if !array.IntContains(baseIds, id) {
			log.Fatal("Tried to update a non-existent ID : id ", id, " ", csvFilePath)
		}
	}

	baseCSV = DeleteCSV(baseCSV, editCSV, csvFilePath)
	baseCSV = InsertCSV(baseCSV, editCSV, csvFilePath)

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

func FilterColumnNumbers(filepath string, filterColumnNames []string) []int {
	rows := LoadCsv(filepath, true, false)

	// Get the column number of the column to filter
	var columnNumbers []int
	for columnNumber, columnName := range rows[0] {
		if array.StrContains(filterColumnNames, columnName) {
			columnNumbers = append(columnNumbers, columnNumber)
		}
	}

	return columnNumbers
}

func LoadNewCsvByDirectoryPath(directoryPath string, fileName string, baseCsvMap map[Key]string, filterNames []string) map[Key]string {
	// Avoid immediately UPDATING an INSET record within the same version (since it is an unintended update).
	loadTypes := []string{"delete", "update", "insert"}
	if !directory.Exist(directoryPath+"/"+loadTypes[0]+"/") &&
		!directory.Exist(directoryPath+"/"+loadTypes[1]+"/") &&
		!directory.Exist(directoryPath+"/"+loadTypes[2]+"/") {
		log.Fatal("Neither insert/update/delete directories were found : ", directoryPath)
	}

	var editIdsAll []int

	for _, loadType := range loadTypes {
		loadTypePath := directoryPath + "/" + loadType + "/"
		if !directory.Exist(loadTypePath) {
			continue
		}

		csvFilePaths, err := GetFilePathRecursive(loadTypePath)
		if err != nil {
			log.Fatal("GetFilePathRecursiveError: ", err)
		}

		for _, csvFilePath := range csvFilePaths {
			if fileName != filepath.Base(csvFilePath) {
				continue
			}

			var editCsvMap map[Key]string
			if len(filterNames) != 0 {
				editCsvMap = LoadCsvMap(csvFilePath, filterNames, false)
			} else {
				editCsvMap = LoadCsvMap(csvFilePath, filterNames, true)
			}

			editIds := PluckId(editCsvMap)
			editIdsAll = append(editIdsAll, editIds...)

			switch loadType {
			case "insert":
				baseCsvMap = InsertCSV(baseCsvMap, editCsvMap, csvFilePath)
			case "update":
				baseCsvMap = UpdateCSV(baseCsvMap, editCsvMap, csvFilePath)
			case "delete":
				baseCsvMap = DeleteCSV(baseCsvMap, editCsvMap, csvFilePath)
			}
		}
	}

	if !array.IsArrayUnique(editIdsAll) {
		log.Fatal("ID is not unique : ", directoryPath, " ", fileName)
	}

	return baseCsvMap
}
