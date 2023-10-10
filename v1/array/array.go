package array

// In 是否存在
// todo 考虑升1.18，使用泛型
func In(array interface{}, find interface{}) bool {
	switch key := find.(type) {
	case string:
		for _, item := range array.([]string) {
			if key == item {
				return true
			}
		}
	case int:
		for _, item := range array.([]int) {
			if key == item {
				return true
			}
		}
	case int64:
		for _, item := range array.([]int64) {
			if key == item {
				return true
			}
		}
	default:
		return false
	}

	return false
}

// NotIn 是否不存在
func NotIn(array interface{}, find interface{}) bool {
	return !In(array, find)
}
