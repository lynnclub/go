package array

import (
	"sync"
	"sync/atomic"
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

	keysInt := Keys(testInt64, 342)
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

	keys := KeysMap(testMap)
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

func TestSliceWithMutex_Add(t *testing.T) {
	slice := SliceWithMutex[int]{}
	slice.Add(42)

	if len(slice.slice) != 1 || slice.slice[0] != 42 {
		t.Errorf("Add method failed. Expected: [42], Got: %v", slice.slice)
	}
}

func TestSliceWithMutex_Get(t *testing.T) {
	slice := SliceWithMutex[int]{}
	slice.Add(42)

	result := slice.Get(0)

	if result != 42 {
		t.Errorf("Get method failed. Expected: 42, Got: %v", result)
	}
}

func TestSliceWithMutex_All(t *testing.T) {
	slice := SliceWithMutex[int]{}
	slice.Add(42)
	slice.Add(24)

	result := slice.All()

	expected := []int{42, 24}
	if len(result) != len(expected) {
		t.Errorf("All method failed. Length mismatch. Expected: %v, Got: %v", expected, result)
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("All method failed. Expected: %v, Got: %v", expected, result)
			break
		}
	}
}

func TestSliceWithMutex_Pop(t *testing.T) {
	slice := SliceWithMutex[int]{}
	slice.Add(42)
	slice.Add(24)

	result := slice.Pop()

	expected := []int{42, 24}
	if len(result) != len(expected) {
		t.Errorf("Pop method failed. Length mismatch. Expected: %v, Got: %v", expected, result)
	}

	if len(slice.slice) != 0 {
		t.Errorf("Pop method failed. The underlying slice is not empty after Pop.")
	}
}

func TestSliceWithMutex_ConcurrentAdd(t *testing.T) {
	slice := SliceWithMutex[int]{}
	var wg sync.WaitGroup
	numGoroutines := 100

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(value int) {
			defer wg.Done()
			slice.Add(value)
		}(i)
	}

	wg.Wait()

	if len(slice.slice) != numGoroutines {
		t.Errorf("Concurrent Add method failed. Expected length: %d, Got: %d", numGoroutines, len(slice.slice))
	}
}

func TestSliceWithMutex_ConcurrentGet(t *testing.T) {
	slice := SliceWithMutex[int]{}
	numElements := 100

	for i := 0; i < numElements; i++ {
		slice.Add(i)
	}

	var wg sync.WaitGroup
	wg.Add(numElements)
	for i := 0; i < numElements; i++ {
		go func(index int) {
			defer wg.Done()
			result := slice.Get(index)
			if result != index {
				t.Errorf("Concurrent Get method failed. Expected: %d, Got: %v", index, result)
			}
		}(i)
	}

	wg.Wait()
}

func TestSliceWithMutex_ConcurrentAll(t *testing.T) {
	slice := SliceWithMutex[int]{}
	numElements := 100

	// Populate the slice with elements
	for i := 0; i < numElements; i++ {
		slice.Add(i)
	}

	var wg sync.WaitGroup
	wg.Add(numElements)
	for i := 0; i < numElements; i++ {
		go func() {
			defer wg.Done()
			result := slice.All()
			// Validate the length of the result slice
			if len(result) != numElements {
				t.Errorf("Concurrent All method failed. Expected length: %d, Got: %d", numElements, len(result))
			}
		}()
	}

	wg.Wait()
}

func TestSliceWithMutex_ConcurrentPop(t *testing.T) {
	slice := SliceWithMutex[int]{}
	numElements := 100

	// Populate the slice with elements
	for i := 0; i < numElements; i++ {
		slice.Add(i)
	}

	var counter int64

	var wg sync.WaitGroup
	wg.Add(numElements)
	for i := 0; i < numElements; i++ {
		go func() {
			defer wg.Done()
			result := slice.Pop()
			// Validate the length of the result slice
			if len(result) > 0 {
				atomic.AddInt64(&counter, 1)
				if len(result) != numElements {
					t.Errorf("Concurrent Pop method failed. Expected length: %d, Got: %d", numElements, len(result))
				}
			}
		}()
	}

	wg.Wait()

	if counter != 1 {
		t.Errorf("Concurrent Pop method failed. Slice Pop not once. counter: %d", counter)
	}

	// After Pop, the underlying slice should be empty
	if len(slice.slice) != 0 {
		t.Errorf("Concurrent Pop method failed. The underlying slice is not empty after Pop.")
	}
}
