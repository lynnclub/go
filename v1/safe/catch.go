package safe

import (
	"fmt"
	"os"
	"runtime"
)

// Catch 捕获错误
func Catch(fn func(), onErr func(err any)) {
	defer func() {
		if err := recover(); err != nil {
			onErr(err)
		}
	}()

	// 执行
	fn()
}

// Retry 重试
func Retry(retry int, fn func()) {
	Catch(fn, func(err any) {
		fmt.Fprintln(os.Stderr, "异常重试", retry, err, runtime.Callers(3, make([]uintptr, 10)))

		if retry > 0 {
			Retry(retry-1, fn)
		}
	})
}

// Go 安全运行协程
func Go(retry int, fn func()) {
	go Catch(fn, func(err any) {
		fmt.Fprintln(os.Stderr, "协程异常重试", retry, err, runtime.Callers(3, make([]uintptr, 10)))

		if retry > 0 {
			Go(retry-1, fn)
		}
	})
}
