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
func Chunk[T any](slice []T, chunkSize int) [][]T {
	var chunks [][]T
	for {
		if len(slice) == 0 {
			break
		}

		if len(slice) < chunkSize {
			chunkSize = len(slice)
		}

		chunk := make([]T, chunkSize)
		copy(chunk, slice[:chunkSize])
		chunks = append(chunks, chunk)

		slice = slice[chunkSize:]
	}

	return chunks
}

// ChunkMap Map分组
func ChunkMap[K comparable, V any](elements map[K]V, chunkSize int) []map[K]V {
	var chunks []map[K]V
	var keys []K

	for k := range elements {
		keys = append(keys, k)
	}

	for {
		if len(keys) == 0 {
			break
		}

		if len(keys) < chunkSize {
			chunkSize = len(keys)
		}

		chunk := make(map[K]V)
		for _, key := range keys[:chunkSize] {
			chunk[key] = elements[key]
		}

		chunks = append(chunks, chunk)

		keys = keys[chunkSize:]
	}

	return chunks
}

// Keys 获取Map的key
func Keys[K comparable, V any](elements map[K]V) []K {
	keys := make([]K, len(elements))
	i := 0
	for key := range elements {
		keys[i] = key
		i++
	}

	return keys
}

// KeysFind 获取Map指定值的key
func KeysFind[K comparable, V comparable](elements map[K]V, findValue V) []K {
	keys := make([]K, 0)
	for key, val := range elements {
		if val == findValue {
			keys = append(keys, key)
		}
	}

	return keys
}

// KeysArray 获取数组的key
func KeysArray[V comparable](elements []V, findValue V) []int {
	keys := make([]int, 0)
	for key, val := range elements {
		if val == findValue {
			keys = append(keys, key)
		}
	}

	return keys
}

// Values 获取Map的value
func Values[K comparable, V any](elements map[K]V) []V {
	vals := make([]V, len(elements))
	i := 0
	for _, val := range elements {
		vals[i] = val
		i++
	}

	return vals
}

// Column 获取Map指定column
func Column[T any, N comparable, K comparable](input map[N]map[K]T, columnKey K) []T {
	columns := make([]T, 0, len(input))
	for _, val := range input {
		if v, ok := val[columnKey]; ok {
			columns = append(columns, v)
		}
	}

	return columns
}

// Diff 获取数组a中不存在于b的元素
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
