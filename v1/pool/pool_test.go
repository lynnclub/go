package pool

import (
	"sync"
	"testing"
)

type TestData struct {
	Value int
}

func TestPool(t *testing.T) {
	myPool := Pool[TestData]{
		Pool:   &sync.Map{},
		Create: func(key any) TestData { return TestData{Value: 42} },
	}

	key := "testKey"
	instance := myPool.Get(key)
	if instance.Value != 42 {
		t.Errorf("Expected Value: 42, Got: %d", instance.Value)
	}

	myPool.Pool.Delete(key)
	_, ok := myPool.Pool.Load(key)
	if ok {
		t.Error("Remove method failed. Key still exists in the pool.")
	}
}
