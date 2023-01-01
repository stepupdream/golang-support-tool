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

type SeparatedValue struct {
	separatedType string
	extension string
}

// Key Make keys into structures to achieve multidimensional arrays.
type Key struct {
	Id  int
	Key string
}

func (separatedValue *SeparatedValue) loadSeparatedValueMap(filePath string, filterNames []string, isColumnExclusion bool) map[Key]string {
	var rows [][]string
	if !supportFile.Exists(filePath) {
		return make(map[Key]string)
	}

	var filterColumnNumbers []int
	rows = separatedValue.LoadSeparatedValue(filePath, true, isColumnExclusion)
	if len(filterNames) != 0 {
		filterColumnNumbers = separatedValue.filterColumnNumbers(filePath, filterNames)
	}

	return separatedValue.convertMap(rows, filterColumnNumbers, filePath)
}

// LoadSeparatedValue Reading separated value files
func (separatedValue *SeparatedValue) LoadSeparatedValue(filepath string, isRowExclusion bool, isColumnExclusion bool) [][]string {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("LoadSeparatedValueOpenError: ", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("LoadSeparatedValueCloseError: ", err)
		}
	}(file)

	// If BOM is included, delete the BOM
	// https://pinzolo.github.io/2017/03/29/utf8-csv-with-bom-on-golang.html
	reader := bufio.NewReader(file)
	bytes, err := reader.Peek(3)
	if err != nil {
		log.Fatal("LoadSeparatedValueNewReaderError: ", err)
	} else if bytes[0] == 0xEF && bytes[1] == 0xBB && bytes[2] == 0xBF {
		_, err := reader.Discard(3)
		if err != nil {
			log.Fatal("SeparatedValueDiscardError: ", err)
		}
	}

	separatedValueReader := csv.NewReader(reader)
	if isRowExclusion {
		separatedValueReader.Comment = '#'
	}
	if separatedValue.separatedType == "tsv" {
		separatedValueReader.Comma = '\t'
		separatedValueReader.LazyQuotes = true
	}
	rows, err := separatedValueReader.ReadAll()
	if err != nil {
		log.Fatal("SeparatedValueReadAllError: ", err)
	}

	if isColumnExclusion {
		return separatedValue.exclusionColumn(rows, isColumnExclusion)
	}

	return rows
}

func (separatedValue *SeparatedValue) exclusionColumn(rows [][]string, isExclusion bool) [][]string {
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

// convertMap
// Replacing separated value data (two-dimensional array of height and width) into a multidimensional associative array in a format
// that facilitates direct value specification by key.
func (separatedValue *SeparatedValue) convertMap(rows [][]string, filterColumnNumbers []int, filepath string) map[Key]string {
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
		log.Fatal("Separated value without ID column cannot be read : ", filepath)
	}

	return result
}

func (separatedValue *SeparatedValue) PluckId(separatedValueMap map[Key]string) []int {
	var ids []int

	for mapKey, _ := range separatedValueMap {
		if mapKey.Key == "id" {
			ids = append(ids, mapKey.Id)
		}
	}

	sort.Ints(ids)

	return ids
}

func (separatedValue *SeparatedValue) PluckKey(separatedValueMap map[Key]string, key string) []string {
	var values []string

	for mapKey, value := range separatedValueMap {
		if mapKey.Key == key {
			values = append(values, value)
		}
	}

	return values
}

func (separatedValue *SeparatedValue) GetFilePathRecursive(path string) ([]string, error) {
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
		if extension != separatedValue.extension {
			return nil
		}

		paths = append(paths, path)

		return nil
	})

	return paths, err
}

func (separatedValue *SeparatedValue) NewFile(path string, rows [][]string) {
	// create allows you to create a new file and overwrite a new file.
	separatedFile, err := os.Create(path)
	if err != nil {
		log.Fatalln(err)
	}

	// Make it with BOM to avoid garbled characters.
	_, err = separatedFile.Write([]byte{0xEF, 0xBB, 0xBF})
	if err != nil {
		log.Fatal(err)
	}

	writer := csv.NewWriter(separatedFile)
	defer writer.Flush()

	if err := writer.WriteAll(rows); err != nil {
		log.Fatal(err)
	}
}

func (separatedValue *SeparatedValue) deleteSeparatedValue(baseSeparatedValue map[Key]string, editSeparatedValue map[Key]string, filePath string) map[Key]string {
	baseIds := separatedValue.PluckId(baseSeparatedValue)

	for key, _ := range editSeparatedValue {
		if key.Key == "id" {
			if !array.IntContains(baseIds, key.Id) {
				log.Fatal("Attempted to delete a non-existent ID : id ", key.Id, " ", filePath)
			}
		}
		delete(baseSeparatedValue, Key{Id: key.Id, Key: key.Key})
	}

	return baseSeparatedValue
}

func (separatedValue *SeparatedValue) insertSeparatedValue(baseSeparatedValue map[Key]string, editSeparatedValue map[Key]string, separatedValueFilePath string) map[Key]string {
	baseIds := separatedValue.PluckId(baseSeparatedValue)
	editIds := separatedValue.PluckId(editSeparatedValue)

	for _, id := range editIds {
		if array.IntContains(baseIds, id) {
			log.Fatal("Tried to do an insert on an existing ID : id ", id, " ", separatedValueFilePath)
		}
	}

	result := make(map[Key]string)

	for mapKey, value := range baseSeparatedValue {
		result[Key{Id: mapKey.Id, Key: mapKey.Key}] = value
	}
	for mapKey, value := range editSeparatedValue {
		result[Key{Id: mapKey.Id, Key: mapKey.Key}] = value
	}

	return result
}

func (separatedValue *SeparatedValue) updateSeparatedValue(baseSeparatedValue map[Key]string, editSeparatedValue map[Key]string, separatedValueFilePath string) map[Key]string {
	baseIds := separatedValue.PluckId(baseSeparatedValue)
	editIds := separatedValue.PluckId(editSeparatedValue)
	for _, id := range editIds {
		if !array.IntContains(baseIds, id) {
			log.Fatal("Tried to update a non-existent ID : id ", id, " ", separatedValueFilePath)
		}
	}

	baseSeparatedValue = separatedValue.deleteSeparatedValue(baseSeparatedValue, editSeparatedValue, separatedValueFilePath)
	baseSeparatedValue = separatedValue.insertSeparatedValue(baseSeparatedValue, editSeparatedValue, separatedValueFilePath)

	return baseSeparatedValue
}

func (separatedValue *SeparatedValue) LoadFileFirstContent(directoryPath string, fileName string) string {
	if !directory.Exist(directoryPath) {
		log.Fatal("The directory could not be found : ", directoryPath)
	}
	baseSeparatedValueFilePaths, err := separatedValue.GetFilePathRecursive(directoryPath)
	if err != nil {
		log.Fatal("LoadFileFirstContentError: ", err)
	}

	if len(baseSeparatedValueFilePaths) == 0 {
		return ""
	}

	var result string
	for _, path := range baseSeparatedValueFilePaths {
		if filepath.Base(path) == fileName {
			rows := separatedValue.LoadSeparatedValue(path, true, false)
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

func (separatedValue *SeparatedValue) filterColumnNumbers(filepath string, filterColumnNames []string) []int {
	rows := separatedValue.LoadSeparatedValue(filepath, true, false)

	// Get the column number of the column to filter
	var columnNumbers []int
	for columnNumber, columnName := range rows[0] {
		if array.StrContains(filterColumnNames, columnName) {
			columnNumbers = append(columnNumbers, columnNumber)
		}
	}

	return columnNumbers
}

func (separatedValue *SeparatedValue) LoadNewSeparatedValueByDirectoryPath(directoryPath string, fileName string, baseSeparatedValueMap map[Key]string, filterNames []string) map[Key]string {
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

		separatedValueFilePaths, err := separatedValue.GetFilePathRecursive(loadTypePath)
		if err != nil {
			log.Fatal("GetFilePathRecursiveError: ", err)
		}

		for _, filePath := range separatedValueFilePaths {
			if fileName != filepath.Base(filePath) {
				continue
			}

			var editSeparatedValueMap map[Key]string
			if len(filterNames) != 0 {
				editSeparatedValueMap = separatedValue.loadSeparatedValueMap(filePath, filterNames, false)
			} else {
				editSeparatedValueMap = separatedValue.loadSeparatedValueMap(filePath, filterNames, true)
			}

			editIds := separatedValue.PluckId(editSeparatedValueMap)
			editIdsAll = append(editIdsAll, editIds...)

			switch loadType {
			case "insert":
				baseSeparatedValueMap = separatedValue.insertSeparatedValue(baseSeparatedValueMap, editSeparatedValueMap, filePath)
			case "update":
				baseSeparatedValueMap = separatedValue.updateSeparatedValue(baseSeparatedValueMap, editSeparatedValueMap, filePath)
			case "delete":
				baseSeparatedValueMap = separatedValue.deleteSeparatedValue(baseSeparatedValueMap, editSeparatedValueMap, filePath)
			}
		}
	}

	if !array.IsArrayUnique(editIdsAll) {
		log.Fatal("ID is not unique : ", directoryPath, " ", fileName)
	}

	return baseSeparatedValueMap
}
