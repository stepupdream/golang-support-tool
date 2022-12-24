package file

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/stepupdream/golang-support-tool/directory"
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

func BaseNamesByArray(paths []string, withExtension bool) []string {
	var result []string
	for _, name := range paths {
		if withExtension {
			result = append(result, filepath.Base(name))
		} else {
			result = append(result, GetNameWithoutExtension(name))
		}
	}

	return result
}

func fileCopy(basedPath string, targetPath string) {
	if !directory.Exist(filepath.Dir(targetPath)) {
		err := os.MkdirAll(filepath.Dir(targetPath), 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	newFile, err := os.Create(targetPath)
	if err != nil {
		log.Fatal(err)
	}

	oldFile, err := os.Open(basedPath)
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(newFile, oldFile)
	if err != nil {
		log.Fatal(err)
	}
}
