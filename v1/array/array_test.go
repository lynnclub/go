package array

import (
	"sort"
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

	arrayIntersect := []string{
		"测试",
		"mon",
	}
	intersect := Intersect(arrayString, arrayDiff, arrayIntersect)
	if NotIn(intersect, "测试") {
		panic("must in")
	}
	if len(intersect) != 1 {
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

// TestSliceWithMutex_Len 测试 Len 方法
func TestSliceWithMutex_Len(t *testing.T) {
	slice := SliceWithMutex[int]{}
	if slice.Len() != 0 {
		t.Errorf("Expected length: 0, Got: %d", slice.Len())
	}

	slice.Add(1)
	slice.Add(2)
	slice.Add(3)

	if slice.Len() != 3 {
		t.Errorf("Expected length: 3, Got: %d", slice.Len())
	}
}

// TestCombineNum_Parse 测试 Parse 方法
func TestCombineNum_Parse(t *testing.T) {
	c := CombineNum{Length: 3, Sep: "."}

	// 测试空字符串
	result := c.Parse("")
	if len(result) != 3 {
		t.Errorf("Expected length: 3, Got: %d", len(result))
	}

	// 测试正常解析
	result = c.Parse("1.2.3")
	expected := []string{"1", "2", "3"}
	if len(result) != len(expected) {
		t.Errorf("Expected length: %d, Got: %d", len(expected), len(result))
	}
	for i, v := range expected {
		if result[i] != v {
			t.Errorf("At index %d, Expected: %s, Got: %s", i, v, result[i])
		}
	}

	// 测试不同分隔符
	c2 := CombineNum{Length: 2, Sep: "-"}
	result = c2.Parse("10-20")
	expected = []string{"10", "20"}
	for i, v := range expected {
		if result[i] != v {
			t.Errorf("At index %d, Expected: %s, Got: %s", i, v, result[i])
		}
	}
}

// TestCombineNum_Get 测试 Get 方法
func TestCombineNum_Get(t *testing.T) {
	c := CombineNum{Length: 3, Sep: "."}

	// 测试正常获取
	num, err := c.Get("1.2.3", 1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if num != 2 {
		t.Errorf("Expected: 2, Got: %d", num)
	}

	// 测试第一个元素
	num, err = c.Get("100.200.300", 0)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if num != 100 {
		t.Errorf("Expected: 100, Got: %d", num)
	}

	// 测试最后一个元素
	num, err = c.Get("100.200.300", 2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if num != 300 {
		t.Errorf("Expected: 300, Got: %d", num)
	}

	// 测试下标不存在
	_, err = c.Get("1.2.3", 3)
	if err == nil {
		t.Error("Expected error for out of bounds index")
	}

	// 测试下标不存在（负数）
	_, err = c.Get("1.2.3", 10)
	if err == nil {
		t.Error("Expected error for out of bounds index")
	}
}

// TestCombineNum_Set 测试 Set 方法
func TestCombineNum_Set(t *testing.T) {
	c := CombineNum{Length: 3, Sep: "."}

	// 测试正常设置
	result, err := c.Set("1.2.3", 1, 999)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != "1.999.3" {
		t.Errorf("Expected: 1.999.3, Got: %s", result)
	}

	// 测试设置第一个元素
	result, err = c.Set("10.20.30", 0, 100)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != "100.20.30" {
		t.Errorf("Expected: 100.20.30, Got: %s", result)
	}

	// 测试设置最后一个元素
	result, err = c.Set("10.20.30", 2, 300)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != "10.20.300" {
		t.Errorf("Expected: 10.20.300, Got: %s", result)
	}

	// 测试下标不存在
	_, err = c.Set("1.2.3", 5, 999)
	if err == nil {
		t.Error("Expected error for out of bounds index")
	}

	// 测试空字符串
	result, err = c.Set("", 0, 100)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != "100.." {
		t.Errorf("Expected: 100.., Got: %s", result)
	}
}

// TestToLower 测试 ToLower 方法
func TestToLower(t *testing.T) {
	// 测试正常转换
	input := []string{"ABC", "Def", "GHI"}
	expected := []string{"abc", "def", "ghi"}
	result := ToLower(input)

	if len(result) != len(expected) {
		t.Errorf("Expected length: %d, Got: %d", len(expected), len(result))
	}
	for i, v := range expected {
		if result[i] != v {
			t.Errorf("At index %d, Expected: %s, Got: %s", i, v, result[i])
		}
	}

	// 测试空切片
	emptyResult := ToLower([]string{})
	if len(emptyResult) != 0 {
		t.Errorf("Expected empty slice, Got length: %d", len(emptyResult))
	}

	// 测试混合字符
	mixed := []string{"Hello123", "WORLD!", "测试TeSt"}
	result = ToLower(mixed)
	if result[0] != "hello123" || result[1] != "world!" || result[2] != "测试test" {
		t.Errorf("Mixed case conversion failed, Got: %v", result)
	}
}

// TestToUpper 测试 ToUpper 方法
func TestToUpper(t *testing.T) {
	// 测试正常转换
	input := []string{"abc", "Def", "ghi"}
	expected := []string{"ABC", "DEF", "GHI"}
	result := ToUpper(input)

	if len(result) != len(expected) {
		t.Errorf("Expected length: %d, Got: %d", len(expected), len(result))
	}
	for i, v := range expected {
		if result[i] != v {
			t.Errorf("At index %d, Expected: %s, Got: %s", i, v, result[i])
		}
	}

	// 测试空切片
	emptyResult := ToUpper([]string{})
	if len(emptyResult) != 0 {
		t.Errorf("Expected empty slice, Got length: %d", len(emptyResult))
	}

	// 测试混合字符
	mixed := []string{"hello123", "world!", "测试TeSt"}
	result = ToUpper(mixed)
	if result[0] != "HELLO123" || result[1] != "WORLD!" || result[2] != "测试TEST" {
		t.Errorf("Mixed case conversion failed, Got: %v", result)
	}
}

// TestUnique 测试 Unique 方法
func TestUnique(t *testing.T) {
	// 测试整数切片去重
	intSlice := []int{1, 2, 3, 2, 4, 3, 5, 1}
	uniqueInt := Unique(intSlice)
	if len(uniqueInt) != 5 {
		t.Errorf("Expected length: 5, Got: %d", len(uniqueInt))
	}
	// 验证元素是否正确
	expectedMap := map[int]bool{1: true, 2: true, 3: true, 4: true, 5: true}
	for _, v := range uniqueInt {
		if !expectedMap[v] {
			t.Errorf("Unexpected value in unique slice: %d", v)
		}
	}

	// 测试字符串切片去重
	strSlice := []string{"a", "b", "a", "c", "b", "d"}
	uniqueStr := Unique(strSlice)
	if len(uniqueStr) != 4 {
		t.Errorf("Expected length: 4, Got: %d", len(uniqueStr))
	}

	// 测试空切片
	emptySlice := []int{}
	uniqueEmpty := Unique(emptySlice)
	if len(uniqueEmpty) != 0 {
		t.Errorf("Expected empty slice, Got length: %d", len(uniqueEmpty))
	}

	// 测试已经是唯一的切片
	uniqueInput := []int{1, 2, 3, 4, 5}
	result := Unique(uniqueInput)
	if len(result) != len(uniqueInput) {
		t.Errorf("Expected length: %d, Got: %d", len(uniqueInput), len(result))
	}

	// 测试所有元素相同
	sameSlice := []string{"a", "a", "a", "a"}
	uniqueSame := Unique(sameSlice)
	if len(uniqueSame) != 1 || uniqueSame[0] != "a" {
		t.Errorf("Expected: [a], Got: %v", uniqueSame)
	}
}

// TestChunkEdgeCases 测试 Chunk 边界情况
func TestChunkEdgeCases(t *testing.T) {
	// 测试空切片
	emptySlice := []int{}
	chunks := Chunk(emptySlice, 2)
	if len(chunks) != 0 {
		t.Errorf("Expected 0 chunks, Got: %d", len(chunks))
	}

	// 测试 chunkSize 大于切片长度
	smallSlice := []int{1, 2, 3}
	chunks = Chunk(smallSlice, 10)
	if len(chunks) != 1 || len(chunks[0]) != 3 {
		t.Errorf("Expected 1 chunk with 3 elements, Got: %d chunks", len(chunks))
	}

	// 测试 chunkSize 为 1
	slice := []int{1, 2, 3}
	chunks = Chunk(slice, 1)
	if len(chunks) != 3 {
		t.Errorf("Expected 3 chunks, Got: %d", len(chunks))
	}
	for i, chunk := range chunks {
		if len(chunk) != 1 || chunk[0] != i+1 {
			t.Errorf("Chunk %d failed", i)
		}
	}

	// 测试切片长度恰好是 chunkSize 的倍数
	slice = []int{1, 2, 3, 4, 5, 6}
	chunks = Chunk(slice, 3)
	if len(chunks) != 2 {
		t.Errorf("Expected 2 chunks, Got: %d", len(chunks))
	}
}

// TestChunkMapEdgeCases 测试 ChunkMap 边界情况
func TestChunkMapEdgeCases(t *testing.T) {
	// 测试空 map
	emptyMap := map[string]int{}
	chunks := ChunkMap(emptyMap, 2)
	if len(chunks) != 0 {
		t.Errorf("Expected 0 chunks, Got: %d", len(chunks))
	}

	// 测试 chunkSize 大于 map 大小
	smallMap := map[string]int{"a": 1, "b": 2}
	chunks = ChunkMap(smallMap, 10)
	if len(chunks) != 1 || len(chunks[0]) != 2 {
		t.Errorf("Expected 1 chunk with 2 elements, Got: %d chunks", len(chunks))
	}

	// 测试 chunkSize 为 1
	testMap := map[string]int{"a": 1, "b": 2, "c": 3}
	chunks = ChunkMap(testMap, 1)
	if len(chunks) != 3 {
		t.Errorf("Expected 3 chunks, Got: %d", len(chunks))
	}
	for _, chunk := range chunks {
		if len(chunk) != 1 {
			t.Errorf("Expected each chunk to have 1 element")
		}
	}
}

// TestIntersectEdgeCases 测试 Intersect 边界情况
func TestIntersectEdgeCases(t *testing.T) {
	// 测试与空切片的交集
	slice1 := []int{1, 2, 3}
	slice2 := []int{}
	result := Intersect(slice1, slice2)
	if len(result) != 0 {
		t.Errorf("Expected empty result, Got: %v", result)
	}

	// 测试没有交集
	slice3 := []int{1, 2, 3}
	slice4 := []int{4, 5, 6}
	result = Intersect(slice3, slice4)
	if len(result) != 0 {
		t.Errorf("Expected empty result, Got: %v", result)
	}

	// 测试完全相同的切片
	slice5 := []int{1, 2, 3}
	slice6 := []int{1, 2, 3}
	result = Intersect(slice5, slice6)
	if len(result) != 3 {
		t.Errorf("Expected 3 elements, Got: %d", len(result))
	}

	// 测试三个切片的交集
	slice7 := []int{1, 2, 3, 4}
	slice8 := []int{2, 3, 4, 5}
	slice9 := []int{3, 4, 5, 6}
	result = Intersect(slice7, slice8, slice9)
	expectedCount := 0
	for _, v := range result {
		if v == 3 || v == 4 {
			expectedCount++
		}
	}
	if expectedCount != len(result) || len(result) != 2 {
		t.Errorf("Expected [3, 4] in result, Got: %v", result)
	}
}

// TestDiffEdgeCases 测试 Diff 边界情况
func TestDiffEdgeCases(t *testing.T) {
	// 测试与空切片的差集
	slice1 := []int{1, 2, 3}
	slice2 := []int{}
	result := Diff(slice1, slice2)
	if len(result) != 3 {
		t.Errorf("Expected 3 elements, Got: %d", len(result))
	}

	// 测试空切片
	slice3 := []int{}
	slice4 := []int{1, 2, 3}
	result = Diff(slice3, slice4)
	if len(result) != 0 {
		t.Errorf("Expected empty result, Got: %v", result)
	}

	// 测试完全排除
	slice5 := []int{1, 2, 3}
	slice6 := []int{1, 2, 3, 4, 5}
	result = Diff(slice5, slice6)
	if len(result) != 0 {
		t.Errorf("Expected empty result, Got: %v", result)
	}
}

// TestSort 测试 Sort 排序功能
func TestSort(t *testing.T) {
	// 创建测试数据
	data := Sort{
		struct{ Id int }{Id: 3},
		struct{ Id int }{Id: 1},
		struct{ Id int }{Id: 4},
		struct{ Id int }{Id: 2},
	}

	// 排序
	sort.Sort(data)

	// 验证排序结果
	expected := []int{1, 2, 3, 4}
	for i, item := range data {
		s := item.(struct{ Id int })
		if s.Id != expected[i] {
			t.Errorf("At index %d, Expected Id: %d, Got: %d", i, expected[i], s.Id)
		}
	}
}

// TestSort_Len 测试 Len 方法
func TestSort_Len(t *testing.T) {
	data := Sort{
		struct{ Id int }{Id: 1},
		struct{ Id int }{Id: 2},
		struct{ Id int }{Id: 3},
	}

	if data.Len() != 3 {
		t.Errorf("Expected length: 3, Got: %d", data.Len())
	}

	emptyData := Sort{}
	if emptyData.Len() != 0 {
		t.Errorf("Expected length: 0, Got: %d", emptyData.Len())
	}
}

// TestSort_Swap 测试 Swap 方法
func TestSort_Swap(t *testing.T) {
	data := Sort{
		struct{ Id int }{Id: 1},
		struct{ Id int }{Id: 2},
		struct{ Id int }{Id: 3},
	}

	// 交换第 0 和第 2 个元素
	data.Swap(0, 2)

	// 验证交换结果
	if data[0].(struct{ Id int }).Id != 3 {
		t.Errorf("Expected Id: 3 at index 0, Got: %d", data[0].(struct{ Id int }).Id)
	}
	if data[2].(struct{ Id int }).Id != 1 {
		t.Errorf("Expected Id: 1 at index 2, Got: %d", data[2].(struct{ Id int }).Id)
	}
}

// TestSort_Less 测试 Less 方法
func TestSort_Less(t *testing.T) {
	data := Sort{
		struct{ Id int }{Id: 5},
		struct{ Id int }{Id: 3},
		struct{ Id int }{Id: 8},
	}

	// 测试比较
	if !data.Less(1, 2) {
		t.Error("Expected data[1] < data[2]")
	}
	if data.Less(2, 1) {
		t.Error("Expected data[2] >= data[1]")
	}
	if data.Less(0, 1) {
		t.Error("Expected data[0] >= data[1]")
	}
}

// TestSort_EmptySlice 测试空切片排序
func TestSort_EmptySlice(t *testing.T) {
	data := Sort{}

	// 空切片排序不应该 panic
	sort.Sort(data)

	if len(data) != 0 {
		t.Errorf("Expected empty slice, Got length: %d", len(data))
	}
}

// TestSort_SingleElement 测试单元素排序
func TestSort_SingleElement(t *testing.T) {
	data := Sort{
		struct{ Id int }{Id: 42},
	}

	sort.Sort(data)

	if data[0].(struct{ Id int }).Id != 42 {
		t.Errorf("Expected Id: 42, Got: %d", data[0].(struct{ Id int }).Id)
	}
}

// TestSort_DuplicateIds 测试重复 Id 的排序
func TestSort_DuplicateIds(t *testing.T) {
	data := Sort{
		struct{ Id int }{Id: 3},
		struct{ Id int }{Id: 1},
		struct{ Id int }{Id: 3},
		struct{ Id int }{Id: 2},
		struct{ Id int }{Id: 1},
	}

	sort.Sort(data)

	// 验证排序后的顺序（应该是非递减的）
	for i := 0; i < len(data)-1; i++ {
		current := data[i].(struct{ Id int }).Id
		next := data[i+1].(struct{ Id int }).Id
		if current > next {
			t.Errorf("Sort failed: data[%d]=%d > data[%d]=%d", i, current, i+1, next)
		}
	}
}

// TestKeysEdgeCases 测试 Keys 函数边界情况
func TestKeysEdgeCases(t *testing.T) {
	// 测试空切片
	emptySlice := []int{}
	keys := Keys(emptySlice, 1)
	if len(keys) != 0 {
		t.Errorf("Expected empty result, Got: %v", keys)
	}

	// 测试没有匹配的值
	slice := []string{"a", "b", "c"}
	keys2 := Keys(slice, "d")
	if len(keys2) != 0 {
		t.Errorf("Expected empty result, Got: %v", keys2)
	}

	// 测试所有值都匹配
	slice2 := []int{5, 5, 5, 5}
	keys3 := Keys(slice2, 5)
	if len(keys3) != 4 {
		t.Errorf("Expected 4 keys, Got: %d", len(keys3))
	}
	for i := 0; i < 4; i++ {
		if !In(keys3, i) {
			t.Errorf("Expected key %d to be in result", i)
		}
	}
}

// TestColumnEdgeCases 测试 Column 函数边界情况
func TestColumnEdgeCases(t *testing.T) {
	// 测试空 map
	emptyMap := map[int]map[string]int{}
	columns := Column(emptyMap, "key")
	if len(columns) != 0 {
		t.Errorf("Expected empty result, Got: %v", columns)
	}

	// 测试不存在的 column key
	testMap := map[int]map[string]int{
		1: {"a": 10, "b": 20},
		2: {"a": 30, "b": 40},
	}
	columns = Column(testMap, "c")
	if len(columns) != 0 {
		t.Errorf("Expected empty result, Got: %v", columns)
	}

	// 测试部分存在的 column key
	testMap2 := map[int]map[string]int{
		1: {"a": 10, "b": 20},
		2: {"a": 30},
		3: {"b": 40},
	}
	columns = Column(testMap2, "a")
	if len(columns) != 2 {
		t.Errorf("Expected 2 columns, Got: %d", len(columns))
	}
}

// TestKeysFindEdgeCases 测试 KeysFind 函数边界情况
func TestKeysFindEdgeCases(t *testing.T) {
	// 测试空 map
	emptyMap := map[string]int{}
	keys := KeysFind(emptyMap, 1)
	if len(keys) != 0 {
		t.Errorf("Expected empty result, Got: %v", keys)
	}

	// 测试没有匹配的值
	testMap := map[string]int{"a": 1, "b": 2, "c": 3}
	keys = KeysFind(testMap, 4)
	if len(keys) != 0 {
		t.Errorf("Expected empty result, Got: %v", keys)
	}

	// 测试多个匹配的值
	testMap2 := map[string]int{"a": 1, "b": 1, "c": 2, "d": 1}
	keys = KeysFind(testMap2, 1)
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, Got: %d", len(keys))
	}
	for _, key := range keys {
		if testMap2[key] != 1 {
			t.Errorf("Key %s has wrong value: %d", key, testMap2[key])
		}
	}
}
