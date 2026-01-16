package redis

import (
	"bytes"
	"context"
	"log"
	"os"
	"strings"
	"testing"
)

// TestNewStdoutLogger 测试创建stdout日志记录器
func TestNewStdoutLogger(t *testing.T) {
	logger := newStdoutLogger()

	if logger == nil {
		t.Error("创建stdout日志记录器失败")
	}

	if logger.log == nil {
		t.Error("日志记录器的log字段为nil")
	}
}

// TestStdoutLoggerPrintf 测试日志输出
func TestStdoutLoggerPrintf(t *testing.T) {
	// 创建一个缓冲区来捕获输出
	var buf bytes.Buffer

	// 创建自定义logger
	logger := &stdoutLogger{
		log: log.New(&buf, "redis: ", log.LstdFlags|log.Lshortfile),
	}

	// 测试输出
	ctx := context.Background()
	testMessage := "test message"
	logger.Printf(ctx, testMessage)

	// 检查输出
	output := buf.String()
	if !strings.Contains(output, testMessage) {
		t.Errorf("日志输出不包含期望的消息。期望包含: %s, 实际输出: %s", testMessage, output)
	}

	if !strings.Contains(output, "redis:") {
		t.Errorf("日志输出不包含前缀 'redis:'。实际输出: %s", output)
	}
}

// TestStdoutLoggerWithFormat 测试格式化日志输出
func TestStdoutLoggerWithFormat(t *testing.T) {
	var buf bytes.Buffer

	logger := &stdoutLogger{
		log: log.New(&buf, "redis: ", log.LstdFlags|log.Lshortfile),
	}

	ctx := context.Background()
	logger.Printf(ctx, "connection %s status: %d", "localhost:6379", 200)

	output := buf.String()
	if !strings.Contains(output, "localhost:6379") {
		t.Errorf("日志输出不包含期望的地址。实际输出: %s", output)
	}

	if !strings.Contains(output, "200") {
		t.Errorf("日志输出不包含期望的状态码。实际输出: %s", output)
	}
}

// TestStdoutLoggerOutput 测试日志输出目标
func TestStdoutLoggerOutput(t *testing.T) {
	// 保存原始的stdout
	oldStdout := os.Stdout

	// 创建一个管道来捕获stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// 使用实际的newStdoutLogger函数
	logger := newStdoutLogger()
	logger.Printf(context.Background(), "output test")

	// 关闭写入端并恢复stdout
	w.Close()
	os.Stdout = oldStdout

	// 读取输出
	var buf bytes.Buffer
	buf.ReadFrom(r)

	output := buf.String()
	if !strings.Contains(output, "output test") {
		t.Errorf("stdout输出不包含期望的消息。实际输出: %s", output)
	}
}
