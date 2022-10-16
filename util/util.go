package util

import (
	"bufio"
	"fmt"
	"os"
)

func StrContains(slice []string, target string) bool {
	for _, value := range slice {
		if value == target {
			return true
		}
	}

	return false
}

func IntContains(slice []int, target int) bool {
	for _, value := range slice {
		if value == target {
			return true
		}
	}

	return false
}

func ExistDirectory(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func MergeMap(m1, m2 map[string]interface{}) map[string]interface{} {
	ans := make(map[string]interface{})

	for k, v := range m1 {
		ans[k] = v
	}
	for k, v := range m2 {
		ans[k] = v
	}

	return ans
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func KeyWait(message string) {
	fmt.Println(message)
	bufio.NewScanner(os.Stdin).Scan()
}
