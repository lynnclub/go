package logger

import (
	"log"
	"os"
	"testing"

	"github.com/lynnclub/go/v1/datetime"
	"gopkg.in/natefinch/lumberjack.v2"
)

// TestLog
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
