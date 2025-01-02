package safe

import (
	"github.com/lynnclub/go/v1/logger"
)

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
		logger.Error("异常重试", err, retry)

		if retry > 0 {
			Retry(retry-1, fn)
		}
	})
}

// Go 安全运行协程
func Go(retry int, fn func()) {
	go Catch(fn, func(err any) {
		logger.Error("协程异常重试", err, retry)

		if retry > 0 {
			Go(retry-1, fn)
		}
	})
}
