package utils

func IndexOf(key string, arr []string) int {
	for k, v := range arr {
		if key == v {
			return k
		}
	}
	return -1
}
