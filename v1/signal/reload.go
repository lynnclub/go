package signal

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

var (
	mu         sync.Mutex
	shutdownCh = make(chan struct{})
)

// Reload 热重载
func Reload() error {
	mu.Lock()
	defer mu.Unlock()

	// 启动新进程
	if err := NewProcess(os.Args[0], os.Args[1:]); err != nil {
		return fmt.Errorf("启动新进程失败: %v", err)
	}

	// 通知旧进程关闭
	shutdown()

	return nil
}

// NewProcess 启动新进程
func NewProcess(executable string, args []string) error {
	cmd := exec.Command(executable, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动新进程失败: %v", err)
	}

	// 等待新进程稳定运行
	for i := 0; i < 3; i++ {
		time.Sleep(1 * time.Second)
		if err := cmd.Process.Signal(syscall.Signal(0)); err != nil {
			return fmt.Errorf("新进程 PID %d 不存在: %v", cmd.Process.Pid, err)
		}
	}

	return nil
}

// Wait 等待关闭信号
func Wait() <-chan struct{} {
	return shutdownCh
}

// shutdown 处理正常关闭
func shutdown() {
	select {
	case <-shutdownCh:
		// 已关闭
		return
	default:
		close(shutdownCh)
	}
}
