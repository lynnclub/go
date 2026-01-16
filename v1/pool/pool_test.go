package pool

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

type TestData struct {
	Value int
}

func TestPool(t *testing.T) {
	myPool := Pool[TestData]{
		Create: func(key any) (*TestData, error) { return &TestData{Value: 42}, nil },
		Close: func(key, value any) bool {
			value.(*TestData).Value = 0
			return true
		},
	}

	key := "testKey"
	instance, _ := myPool.Get(key)
	if instance.Value != 42 {
		t.Errorf("Expected Value: 42, Got: %d", instance.Value)
	}

	myPool.CloseAll()
	instance, _ = myPool.Get(key)
	if instance.Value != 0 {
		t.Errorf("Expected Value: 42, Got: %d", instance.Value)
	}

	myPool.Pool.Delete(key)
	_, ok := myPool.Pool.Load(key)
	if ok {
		t.Error("Remove method failed. Key still exists in the pool.")
	}
}

// TestPool_CreateError 测试 Create 函数返回错误
func TestPool_CreateError(t *testing.T) {
	expectedError := errors.New("create failed")
	myPool := Pool[TestData]{
		Create: func(key any) (*TestData, error) {
			return nil, expectedError
		},
	}

	instance, err := myPool.Get("testKey")
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err != expectedError {
		t.Errorf("Expected error: %v, Got: %v", expectedError, err)
	}
	if instance != nil {
		t.Errorf("Expected nil instance, Got: %v", instance)
	}
}

// TestPool_MultipleKeys 测试多个不同 key
func TestPool_MultipleKeys(t *testing.T) {
	myPool := Pool[TestData]{
		Create: func(key any) (*TestData, error) {
			keyInt := key.(int)
			return &TestData{Value: keyInt * 10}, nil
		},
		Close: func(key, value any) bool {
			return true
		},
	}

	// 创建多个不同的实例
	keys := []int{1, 2, 3, 4, 5}
	for _, key := range keys {
		instance, err := myPool.Get(key)
		if err != nil {
			t.Errorf("Unexpected error for key %d: %v", key, err)
		}
		expectedValue := key * 10
		if instance.Value != expectedValue {
			t.Errorf("For key %d, Expected Value: %d, Got: %d", key, expectedValue, instance.Value)
		}
	}

	// 验证再次获取相同 key 返回相同实例
	for _, key := range keys {
		instance, _ := myPool.Get(key)
		expectedValue := key * 10
		if instance.Value != expectedValue {
			t.Errorf("For key %d on second get, Expected Value: %d, Got: %d", key, expectedValue, instance.Value)
		}
	}
}

// TestPool_ReuseSameKey 测试重复获取同一个 key
func TestPool_ReuseSameKey(t *testing.T) {
	createCount := 0
	myPool := Pool[TestData]{
		Create: func(key any) (*TestData, error) {
			createCount++
			return &TestData{Value: 100}, nil
		},
	}

	key := "sameKey"

	// 第一次获取
	instance1, err := myPool.Get(key)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if createCount != 1 {
		t.Errorf("Expected Create to be called once, called %d times", createCount)
	}

	// 第二次获取应该返回相同实例，不应该再次调用 Create
	instance2, err := myPool.Get(key)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if createCount != 1 {
		t.Errorf("Expected Create to be called once, called %d times", createCount)
	}

	// 验证返回的是同一个实例
	if instance1 != instance2 {
		t.Error("Expected same instance for same key")
	}
}

// TestPool_LockForCreate_True 测试 LockForCreate 为 true
func TestPool_LockForCreate_True(t *testing.T) {
	var createCount int32
	myPool := Pool[TestData]{
		LockForCreate: true,
		Create: func(key any) (*TestData, error) {
			atomic.AddInt32(&createCount, 1)
			return &TestData{Value: int(key.(int))}, nil
		},
	}

	key := 42
	numGoroutines := 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// 多个 goroutine 同时获取同一个 key
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			instance, err := myPool.Get(key)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if instance.Value != key {
				t.Errorf("Expected Value: %d, Got: %d", key, instance.Value)
			}
		}()
	}

	wg.Wait()

	// 由于 LockForCreate 为 true，Create 应该只被调用一次
	if createCount != 1 {
		t.Errorf("With LockForCreate=true, Expected Create to be called once, called %d times", createCount)
	}
}

// TestPool_LockForCreate_False 测试 LockForCreate 为 false
func TestPool_LockForCreate_False(t *testing.T) {
	var createCount int32
	myPool := Pool[TestData]{
		LockForCreate: false,
		Create: func(key any) (*TestData, error) {
			atomic.AddInt32(&createCount, 1)
			return &TestData{Value: int(key.(int))}, nil
		},
	}

	key := 42
	numGoroutines := 10

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// 多个 goroutine 同时获取同一个 key
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			instance, err := myPool.Get(key)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if instance.Value != key {
				t.Errorf("Expected Value: %d, Got: %d", key, instance.Value)
			}
		}()
	}

	wg.Wait()

	// 由于 LockForCreate 为 false，在并发情况下 Create 可能被调用多次
	// 但至少应该被调用一次
	if createCount < 1 {
		t.Errorf("Expected Create to be called at least once, called %d times", createCount)
	}
}

// TestPool_ConcurrentMultipleKeys 测试并发访问多个 key
func TestPool_ConcurrentMultipleKeys(t *testing.T) {
	myPool := Pool[TestData]{
		LockForCreate: true,
		Create: func(key any) (*TestData, error) {
			return &TestData{Value: key.(int) * 100}, nil
		},
	}

	numKeys := 10
	numGoroutinesPerKey := 10

	var wg sync.WaitGroup
	wg.Add(numKeys * numGoroutinesPerKey)

	for key := 0; key < numKeys; key++ {
		for i := 0; i < numGoroutinesPerKey; i++ {
			go func(k int) {
				defer wg.Done()
				instance, err := myPool.Get(k)
				if err != nil {
					t.Errorf("Unexpected error for key %d: %v", k, err)
				}
				expectedValue := k * 100
				if instance.Value != expectedValue {
					t.Errorf("For key %d, Expected Value: %d, Got: %d", k, expectedValue, instance.Value)
				}
			}(key)
		}
	}

	wg.Wait()

	// 验证池中有正确数量的实例
	count := 0
	myPool.Pool.Range(func(key, value any) bool {
		count++
		return true
	})

	if count != numKeys {
		t.Errorf("Expected %d instances in pool, Got: %d", numKeys, count)
	}
}

// TestPool_CloseAll 测试 CloseAll 方法
func TestPool_CloseAll(t *testing.T) {
	closedKeys := make(map[int]bool)
	var mutex sync.Mutex

	myPool := Pool[TestData]{
		Create: func(key any) (*TestData, error) {
			return &TestData{Value: key.(int)}, nil
		},
		Close: func(key, value any) bool {
			mutex.Lock()
			defer mutex.Unlock()
			closedKeys[key.(int)] = true
			return true
		},
	}

	// 创建多个实例
	keys := []int{1, 2, 3, 4, 5}
	for _, key := range keys {
		_, _ = myPool.Get(key)
	}

	// 关闭所有实例
	myPool.CloseAll()

	// 验证所有 key 都被 Close 调用
	mutex.Lock()
	defer mutex.Unlock()
	for _, key := range keys {
		if !closedKeys[key] {
			t.Errorf("Key %d was not closed", key)
		}
	}
}

// TestPool_CloseAllWithError 测试 CloseAll 时 Close 返回 false
func TestPool_CloseAllWithError(t *testing.T) {
	myPool := Pool[TestData]{
		Create: func(key any) (*TestData, error) {
			return &TestData{Value: key.(int)}, nil
		},
		Close: func(key, value any) bool {
			// Close 返回 false 表示停止迭代
			if key.(int) == 3 {
				return false
			}
			return true
		},
	}

	// 创建多个实例
	for i := 1; i <= 5; i++ {
		_, _ = myPool.Get(i)
	}

	// CloseAll 应该在某个点停止
	myPool.CloseAll()
	// 这个测试主要确保 CloseAll 不会 panic
}

// TestPool_EmptyPool 测试空池的操作
func TestPool_EmptyPool(t *testing.T) {
	myPool := Pool[TestData]{
		Create: func(key any) (*TestData, error) {
			return &TestData{Value: 1}, nil
		},
		Close: func(key, value any) bool {
			return true
		},
	}

	// 在空池上调用 CloseAll 不应该 panic
	myPool.CloseAll()

	// 验证池仍然可用
	instance, err := myPool.Get("key")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if instance.Value != 1 {
		t.Errorf("Expected Value: 1, Got: %d", instance.Value)
	}
}

// TestPool_DifferentTypes 测试不同类型的 key
func TestPool_DifferentTypes(t *testing.T) {
	myPool := Pool[TestData]{
		Create: func(key any) (*TestData, error) {
			return &TestData{Value: 999}, nil
		},
	}

	// 测试 string key
	instance1, _ := myPool.Get("stringKey")
	if instance1.Value != 999 {
		t.Errorf("Expected Value: 999, Got: %d", instance1.Value)
	}

	// 测试 int key
	instance2, _ := myPool.Get(123)
	if instance2.Value != 999 {
		t.Errorf("Expected Value: 999, Got: %d", instance2.Value)
	}

	// 测试 struct key
	type customKey struct {
		id int
	}
	instance3, _ := myPool.Get(customKey{id: 1})
	if instance3.Value != 999 {
		t.Errorf("Expected Value: 999, Got: %d", instance3.Value)
	}

	// 验证不同类型的 key 产生不同的实例
	if instance1 == instance2 || instance2 == instance3 || instance1 == instance3 {
		t.Error("Different key types should produce different instances")
	}
}

// TestPool_CreateWithKeyInfo 测试 Create 函数使用 key 信息
func TestPool_CreateWithKeyInfo(t *testing.T) {
	myPool := Pool[TestData]{
		Create: func(key any) (*TestData, error) {
			keyStr, ok := key.(string)
			if !ok {
				return nil, fmt.Errorf("key must be string")
			}
			// 根据 key 创建不同的值
			var value int
			switch keyStr {
			case "small":
				value = 10
			case "medium":
				value = 50
			case "large":
				value = 100
			default:
				value = 0
			}
			return &TestData{Value: value}, nil
		},
	}

	testCases := []struct {
		key           string
		expectedValue int
	}{
		{"small", 10},
		{"medium", 50},
		{"large", 100},
		{"unknown", 0},
	}

	for _, tc := range testCases {
		instance, err := myPool.Get(tc.key)
		if err != nil {
			t.Errorf("Unexpected error for key %s: %v", tc.key, err)
		}
		if instance.Value != tc.expectedValue {
			t.Errorf("For key %s, Expected Value: %d, Got: %d", tc.key, tc.expectedValue, instance.Value)
		}
	}
}

// TestPool_NilCreate 测试 Create 返回 nil 但没有错误的情况
func TestPool_NilCreate(t *testing.T) {
	myPool := Pool[TestData]{
		Create: func(key any) (*TestData, error) {
			// 返回 nil 但没有错误
			return nil, nil
		},
	}

	instance, err := myPool.Get("key")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	// 虽然不推荐，但这种情况应该能处理
	if instance != nil {
		t.Error("Expected nil instance")
	}
}
