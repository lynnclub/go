package safe

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func simulated() {
	panic("Simulated panic")
}

func TestGo(t *testing.T) {
	run := 0
	Retry(1, func() {
		run++
		simulated()
	})

	if run != 2 {
		t.Errorf("Retry:The number of executions is not 2 but %d", run)
	}

	var goRun int64

	Go(2, func() {
		atomic.AddInt64(&goRun, 1)
		fmt.Println(goRun)
		simulated()
	})

	time.Sleep(1 * time.Second)

	if goRun != 3 {
		t.Errorf("Go:The number of executions is not 3 but %d", goRun)
	}
}
