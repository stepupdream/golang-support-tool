package excel

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"reflect"

	"github.com/pkg/errors"
	"github.com/stepupdream/golang-support-tool/csv"
	"github.com/stepupdream/golang-support-tool/file"
	"github.com/xuri/excelize/v2"
)

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
		if extension != ".xlsx" && extension != ".xlsm" {
			return nil
		}

		paths = append(paths, path)

		return nil
	})

	return paths, err
}

func ToCsv(excelPath string, csvPath string, sheetIndex int) (string, error) {
	excel, err := excelize.OpenFile(excelPath)
	if err != nil {
		return csvPath, nil
	}
	defer func() {
		if err := excel.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	rowsExcel, err := excel.GetRows(excel.GetSheetName(sheetIndex))
	if err != nil {
		log.Fatal("Failed to open Excel file : ", err)
	}

	if file.Exists(csvPath) {
		rowsCsv := csv.LoadCsv(csvPath, false)
		if reflect.DeepEqual(rowsExcel, rowsCsv) {
			fmt.Println("[  SKIP  ] ", excelPath)
			return "", nil
		}
	}

	csv.NewFile(csvPath, rowsExcel)
	fmt.Println("[COMPLETE] ", excelPath)

	return csvPath, nil
}
