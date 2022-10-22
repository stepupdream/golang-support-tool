package directory

import (
	"log"
	"os"
)

func ExistDirectory(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func GetDirectoryNames(path string) []string {
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

	return names
}
