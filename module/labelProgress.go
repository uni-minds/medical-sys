package module

var progressData map[int]string

func ProgressQueryString(key int) string {
	value, ok := progressData[key]
	if ok {
		return value
	} else {
		return ""
	}
}

func ProgressQueryValue(str string) int {
	for k, v := range progressData {
		if v == str {
			return k
		}
	}
	return 0
}
