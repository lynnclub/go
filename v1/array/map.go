package array

// ChunkMap Map分组
func ChunkMap[K comparable, V any](array map[K]V, chunkSize int) []map[K]V {
	var chunks []map[K]V
	var keys []K

	for k := range array {
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
			chunk[key] = array[key]
		}

		chunks = append(chunks, chunk)

		keys = keys[chunkSize:]
	}

	return chunks
}

// KeysMap 获取Map的key
func KeysMap[K comparable, V any](array map[K]V) []K {
	keys := make([]K, len(array))
	i := 0
	for key := range array {
		keys[i] = key
		i++
	}

	return keys
}

// KeysFind 获取Map指定值的key
func KeysFind[K comparable, V comparable](array map[K]V, find V) []K {
	keys := make([]K, 0)
	for key, val := range array {
		if val == find {
			keys = append(keys, key)
		}
	}

	return keys
}

// Values 获取Map的value
func Values[K comparable, V any](array map[K]V) []V {
	vals := make([]V, len(array))
	i := 0
	for _, val := range array {
		vals[i] = val
		i++
	}

	return vals
}

// Column 获取Map指定column
func Column[T any, N comparable, K comparable](array map[N]map[K]T, columnKey K) []T {
	columns := make([]T, 0)
	for _, val := range array {
		if v, ok := val[columnKey]; ok {
			columns = append(columns, v)
		}
	}

	return columns
}
