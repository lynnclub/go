package gosafe

import "fmt"

// Run 安全运行协程
func Run(isRecover bool, fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("协程异常退出", err)

				if isRecover {
					fn()
				}
			}
		}()

		// 执行
		fn()
	}()
}
