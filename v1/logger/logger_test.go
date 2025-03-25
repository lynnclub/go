package logger

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/lynnclub/go/v1/datetime"
	"gopkg.in/natefinch/lumberjack.v2"
)

func TestLog(t *testing.T) {
	Logger = New(
		log.New(os.Stderr, "", log.Lmsgprefix),
		"local",
		INFO,
		"asia/shanghai",
		datetime.LayoutDateTimeZoneT,
		nil,
	)
	Info("123")
	Debug("123")

	lumberjack := &lumberjack.Logger{
		Filename:   "foo.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     14,   //days
		Compress:   true, // disabled by default
	}
	Logger = New(
		log.New(lumberjack, "", log.Lmsgprefix),
		"local",
		INFO,
		"asia/shanghai",
		datetime.LayoutDateTimeZoneT,
		nil,
	)
	Info("321")
	Debug("321")
	Error("4321")
}

func TestSendDuplicateLog(t *testing.T) {
	alert := &FeishuAlert{}
	option := Option{
		Levels:  []string{"ERROR"},
		Webhook: "http://test.webhook",
		SignKey: "test_sign_key",
		UserId:  "test_user_id",
	}

	alert.Add("test", option)

	log := LogEntry{
		Level:     400,
		LevelName: "ERROR",
		Command:   "test_command",
		Message:   "test_message",
		Datetime:  time.Now().Format(time.RFC3339),
		Trace:     "",
		URL:       "",
		Env:       "production",
		IP:        "127.0.0.1",
	}

	alert.Send(log)

	// 立即再次发送相同的日志，应该被过滤掉
	alert.Send(log)
	if len(alert.lastHashs) != 1 {
		t.Errorf("Expected 1 log in history, got %d", len(alert.lastHashs))
	}
}
