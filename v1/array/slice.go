package array

import "strings"

// In 是否存在
func In[T comparable](array []T, find T) bool {
	for _, item := range array {
		if item == find {
			return true
		}
	}
	return false
}

// NotIn 是否不存在
func NotIn[T comparable](array []T, find T) bool {
	return !In(array, find)
}

// Chunk 分组
func Chunk[T any](array []T, chunkSize int) [][]T {
	var chunks [][]T
	for {
		if len(array) == 0 {
			break
		}

		if len(array) < chunkSize {
			chunkSize = len(array)
		}

		chunk := make([]T, chunkSize)
		copy(chunk, array[:chunkSize])
		chunks = append(chunks, chunk)

		array = array[chunkSize:]
	}

	return chunks
}

// Keys 获取切片的key
func Keys[V comparable](array []V, find V) []int {
	keys := make([]int, 0)
	for key, val := range array {
		if val == find {
			keys = append(keys, key)
		}
	}

	return keys
}

// Diff 获取切片a中不存在于b的元素
func Diff[T comparable](a []T, b ...[]T) []T {
	excludeMap := make(map[T]bool)
	diff := []T{}

	for _, arr := range b {
		for _, val := range arr {
			excludeMap[val] = true
		}
	}

	for _, val := range a {
		if !excludeMap[val] {
			diff = append(diff, val)
		}
	}

	return diff
}

// ToLower 转小写
func ToLower(array []string) []string {
	newArray := make([]string, 0)
	for _, item := range array {
		newArray = append(newArray, strings.ToLower(item))
	}

	return newArray
}

// ToUpper 转大写
func ToUpper(array []string) []string {
	newArray := make([]string, 0)
	for _, item := range array {
		newArray = append(newArray, strings.ToUpper(item))
	}

	return newArray
}
