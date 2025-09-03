package signal

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	ChannelOS  = make(chan os.Signal, 1)       // 信号通道
	Now        os.Signal                       // 当前信号
	handlerMap = make(map[os.Signal][]Handler) // 信号处理器映射
	handlerMu  sync.RWMutex                    // 读写锁
)

// Handler 信号处理器
type Handler func(os.Signal)

// Listen 监听信号
// SIGHUP 挂起（hangup），当终端关闭或者连接的会话结束时，由内核发送给进程
// SIGINT 中断（interrupt），通常由用户按下 Ctrl+C 产生，进程接收到信号后应立即停止当前的工作
// SIGQUIT 退出（quit），通常由用户按下 Ctrl+\ 产生，进程接收到信号后应立即退出，并清理自己占用的资源
// SIGTERM 终止（terminate），这是一个通用信号，通常用于要求进程正常终止
// SIGFPE 在发生致命的算术运算错误时发出，如除零操作、数据溢出等
// SIGKILL 立即结束程序，无法被捕获
// SIGALRM 时钟定时信号
// SIGBUS 总线错误，通常是内存对齐问题或硬件故障
// SIGSEGV 段错误，访问未映射的内存区域
func Listen(signals ...os.Signal) {
	if len(signals) == 0 {
		signals = []os.Signal{syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP}
	}

	go func(ch chan os.Signal) {
		for sig := range ch {
			Now = sig

			handlerMu.RLock()
			handlers := handlerMap[sig]
			handlerMu.RUnlock()

			for _, handler := range handlers {
				if handler != nil {
					handler(sig)
				}
			}
		}
	}(ChannelOS)

	signal.Notify(ChannelOS, signals...)
}

// SetHandler 设置信号处理器
func SetHandler(handler Handler, signals ...os.Signal) {
	handlerMu.Lock()
	defer handlerMu.Unlock()

	for _, sig := range signals {
		handlerMap[sig] = append(handlerMap[sig], handler)
	}
}
