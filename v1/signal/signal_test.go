package signal

import (
	"fmt"
	"sync"
	"syscall"
	"testing"
	"time"
)

// TestListen 监听
func TestListen(t *testing.T) {
	Listen()

	// 模拟业务协程
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for loop := 0; loop < 100; loop++ {
			// 收到停机信号，主动退出业务
			if Now != nil {
				wg.Done()
				fmt.Println("business stop signal:", Now)
				break
			}

			// do something...

			time.Sleep(100 * time.Millisecond)
		}

		if Now == nil {
			panic("no signal")
		}
	}()

	// 发送信号
	err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	if err != nil {
		panic("signal send failed")
	}

	// 主程形式一：循环
	for {
		// 主程需要等待协程停止
		if Now != nil {
			wg.Wait()
			fmt.Println("main stop signal:", Now)
			break
		}

		// do something...

		time.Sleep(100 * time.Millisecond)
	}

	// 主程形式二：阻塞
	// select {
	// case <-ChannelOS:
	// 	// 主程需要等待协程停止
	// 	wg.Wait()
	// 	fmt.Println("main stop signal:", Now)
	// 	break
	// }
}
