package array

import (
	"log"
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

func IsArrayUnique(args []int) bool {
	encountered := map[int]bool{}
	count := len(args)
	for i := 0; i < count; i++ {
		if !encountered[args[i]] {
			encountered[args[i]] = true
		} else {
			return false
		}
	}

	return true
}

func NextArrayValue(allValues []string, nowValue string) string {
	if !StrContains(allValues, nowValue) {
		log.Fatal("Incorrect value specified. The specified value does not exist in the array : ", nowValue)
	}

	var nowKey int
	for key, value := range allValues {
		if value == nowValue {
			nowKey = key
		}
	}

	if len(allValues) < nowKey+2 {
		return ""
	}

	return allValues[nowKey+1]
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

func SliceString(all []string, start string, end string) []string {
	var tmp []string
	if start == "" {
		start = all[0]
	}

	isStart := false
	for _, value := range all {
		if value == start {
			isStart = true
		}

		if isStart {
			tmp = append(tmp, value)
		}
	}

	var result []string
	isEnd := false
	for _, value := range tmp {
		switch end {
		case "next":
			return []string{value}
		case "max":
			result = append(result, value)
		default:
			if !StrContains(all, end) {
				log.Fatal("The specified value could not be found : ", end)
			}
			if !isEnd {
				result = append(result, value)
			}

			if value == end {
				isEnd = true
			}
		}
	}

	return result
}
