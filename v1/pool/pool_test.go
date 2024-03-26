package pool

import (
	"testing"
)

type TestData struct {
	Value int
}

func TestPool(t *testing.T) {
	myPool := Pool[TestData]{
		Create: func(key any) *TestData { return &TestData{Value: 42} },
		Close: func(key, value any) bool {
			value.(*TestData).Value = 0
			return true
		},
	}

	key := "testKey"
	instance := myPool.Get(key)
	if instance.Value != 42 {
		t.Errorf("Expected Value: 42, Got: %d", instance.Value)
	}

	myPool.CloseAll()
	instance = myPool.Get(key)
	if instance.Value != 0 {
		t.Errorf("Expected Value: 42, Got: %d", instance.Value)
	}

	myPool.Pool.Delete(key)
	_, ok := myPool.Pool.Load(key)
	if ok {
		t.Error("Remove method failed. Key still exists in the pool.")
	}
}
