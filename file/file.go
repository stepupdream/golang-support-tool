package file

import (
	"log"
	"os"
	"path/filepath"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func FindFilePathParentDir(filename string) string {
	dirPath, _ := os.Getwd()

	for i := 0; i < 10; i++ {
		findPath := dirPath + "/" + filename
		if Exists(findPath) {
			return findPath
		}

		dirPath = filepath.Dir(dirPath)
	}

	log.Fatal("The specified file could not be found : ", filename)

	return ""
}

func GetNameWithoutExtension(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}
