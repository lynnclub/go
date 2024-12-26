package safe

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
)

// Trace 执行链路
func Trace(skip int, deep int) []string {
	trace := make([]string, 0)

	pcs := make([]uintptr, deep)
	deeps := runtime.Callers(skip, pcs)
	for current := 0; current < deeps; current++ {
		function := runtime.FuncForPC(pcs[current])
		file, line := function.FileLine(pcs[current])
		trace = append(trace, "["+strconv.Itoa(current)+"] "+function.Name()+"()")
		trace = append(trace, file+":"+strconv.Itoa(line))
	}

	return trace
}

// Recover 错误处理
func Recover(onErr func(err any)) {
	if err := recover(); err != nil {
		onErr(err)
	}
}

// Catch 捕获错误
func Catch(fn func(), onErr func(err any)) {
	defer Recover(onErr)

	// 执行
	fn()
}

// Retry 重试
func Retry(retry int, fn func()) {
	Catch(fn, func(err any) {
		fmt.Fprintln(os.Stderr, "异常重试", retry, err, Trace(10, 3))

		if retry > 0 {
			Retry(retry-1, fn)
		}
	})
}

// Go 安全运行协程
func Go(retry int, fn func()) {
	go Catch(fn, func(err any) {
		fmt.Fprintln(os.Stderr, "协程异常重试", retry, err, Trace(10, 3))

		if retry > 0 {
			Go(retry-1, fn)
		}
	})
}
