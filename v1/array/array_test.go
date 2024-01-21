package array

import (
	"testing"
)

// TestIn 是否存在
func TestIn(t *testing.T) {
	arrayString := []string{
		"adc",
		"mon",
		"测试",
	}
	if In(arrayString, "测试abc") {
		panic("must not in")
	}
	if NotIn(arrayString, "测试") {
		panic("must in")
	}

	arrayDiff := []string{
		"测试",
	}
	diff := Diff(arrayString, arrayDiff)
	if In(diff, "测试") {
		panic("must not in")
	}
	if len(diff) != 2 {
		panic("keys num error")
	}
	arrayDiff2 := []string{
		"adc",
		"mon",
	}
	diff = Diff(arrayString, arrayDiff, arrayDiff2)
	if In(diff, "adc") {
		panic("must not in")
	}
	if len(diff) != 0 {
		panic("keys num error")
	}

	testInt64 := []int64{
		123,
		342,
		6,
		342,
	}
	if In(testInt64, 1342) {
		panic("must not in")
	}
	if NotIn(testInt64, 342) {
		panic("must in")
	}

	keysInt := KeysArray(testInt64, 342)
	if NotIn(keysInt, 3) {
		panic("key not in")
	}
	if len(keysInt) != 2 {
		panic("keys num error")
	}

	chunkString := Chunk(arrayString, 2)
	if &chunkString[0][0] == &arrayString[0] {
		panic("slices may overflow memory")
	}
	count := 0
	for _, chunk := range chunkString {
		count += len(chunk)
	}
	if count != 3 {
		panic("chunk num error")
	}

	testMap := map[string]int64{
		"123": 123,
		"342": 342,
		"6":   6,
	}
	chunkMap := ChunkMap(testMap, 2)
	count = 0
	for _, chunk := range chunkMap {
		count += len(chunk)
	}
	if count != 3 {
		panic("chunk map num error")
	}

	keys := Keys(testMap)
	if NotIn(keys, "6") {
		panic("key not in")
	}
	if len(keys) != 3 {
		panic("keys num error")
	}

	keys = KeysFind(testMap, 342)
	if NotIn(keys, "342") {
		panic("key not in")
	}
	if len(keys) != 1 {
		panic("keys num error")
	}

	values := Values(testMap)
	if NotIn(values, 6) {
		panic("value not in")
	}
	if len(values) != 3 {
		panic("values num error")
	}

	testMap2 := map[int64]map[string]int64{
		123:   {"23": 123},
		1234:  {"234": 1234},
		12345: {"234": 12345},
	}
	values = Column(testMap2, "234")
	if NotIn(values, 12345) {
		panic("value not in")
	}
	if len(values) != 2 {
		panic("values num error")
	}
}
