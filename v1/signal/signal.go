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
func Listen(signals ...os.Signal) {
	if len(signals) == 0 {
		signals = []os.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM}
	}

	go func(ch chan os.Signal) {
		Now = <-ch
		close(ch)
	}(ChannelOS)

	signal.Notify(ChannelOS, signals...)
}
