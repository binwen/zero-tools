package stringx

// 判断字符串是否在列表中
func Contains(list []string, str string) bool {
	for _, each := range list {
		if each == str {
			return true
		}
	}
	return false
}
