package directory

import (
	"log"
	"os"

	"github.com/stepupdream/golang-support-tool/array"
)

func Exist(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func GetNames(path string, exclusionTexts []string) []string {
	dir, err := os.Open(path)
	if err != nil {
		log.Fatal("Not found : ", err)
	}
	defer func(dir *os.File) {
		err := dir.Close()
		if err != nil {

		}
	}(dir)

	names, err := dir.Readdirnames(-1)
	if err != nil {
		log.Fatal("ReadDirError: ", err)
	}

	var result []string
	for _, name := range names {
		if !array.StrContains(exclusionTexts, name) {
			result = append(result, name)
		}
	}

	return result
}

func ExistMulti(parentPaths []string) bool {
	isExist := false

	for _, path := range parentPaths {
		if Exist(path) {
			isExist = true
		}
	}

	return isExist
}

func MaxFileName(directoryPath string) string {
	maxName := ""
	dirEntries, _ := os.ReadDir(directoryPath)
	for _, dirEntry := range dirEntries {
		if maxName == "" {
			maxName = dirEntry.Name()
			continue
		}

		if maxName < dirEntry.Name() {
			maxName = dirEntry.Name()
		}
	}

	return maxName
}
