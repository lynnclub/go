package signal

import (
	"os"
	"os/signal"
	"syscall"
)

// Now 当前信号
var Now os.Signal

// ChannelOS 系统信号
var ChannelOS = make(chan os.Signal)

// Listen 监听
// SIGHUP 挂起（hangup），当终端关闭或者连接的会话结束时，由内核发送给进程
// SIGINT 中断（interrupt），通常由用户按下 Ctrl+C 产生，进程接收到信号后应立即停止当前的工作
// SIGQUIT 退出（quit），通常由用户按下 Ctrl+\ 产生，进程接收到信号后应立即退出，并清理自己占用的资源
// SIGTERM 终止（terminate），这是一个通用信号，通常用于要求进程正常终止
// SIGFPE 在发生致命的算术运算错误时发出，如除零操作、数据溢出等
// SIGKILL 立即结束程序的运行
// SIGALRM 时钟定时信号
// SIGBUS SIGSEGV 进程访问非法地址
func Listen(signals ...os.Signal) {
	if len(signals) == 0 {
		signals = []os.Signal{syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM}
	}

	go func(ch chan os.Signal) {
		Now = <-ch
		close(ch)
	}(ChannelOS)

	signal.Notify(ChannelOS, signals...)
}
