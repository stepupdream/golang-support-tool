package array

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
