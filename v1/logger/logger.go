package logger

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/lynnclub/go/v1/datetime"
	"github.com/lynnclub/go/v1/encoding/json"
	"github.com/lynnclub/go/v1/ip"
	"github.com/lynnclub/go/v1/safe"
)

var (
	levelFlags = []string{"DEBUG", "INFO", "NOTICE", "WARN", "ERROR", "FATAL", "PANIC"}
	Logger     = New(log.New(os.Stderr, "", log.Lmsgprefix), DEBUG, "local", "asia/shanghai", datetime.LayoutDateTimeZoneT, nil)
)

const (
	DEBUG int = iota
	INFO
	NOTICE
	WARN
	ERROR
	FATAL
	PANIC
)

type logger struct {
	Raw        *log.Logger                      // 原生log
	env        string                           // 环境
	level      int                              // 起始级别
	trace      string                           // 追踪标识，traceId/userId/orderId等
	timezone   string                           // 时区
	timeFormat string                           // 时间格式
	request    *http.Request                    // 请求
	callback   func(log map[string]interface{}) // 回调
}

func New(raw *log.Logger, env string, level int, timezone, timeFormat string, callback func(log map[string]interface{})) *logger {
	return &logger{
		Raw:        raw,
		env:        env,
		level:      level,
		timezone:   timezone,
		timeFormat: timeFormat,
		callback:   callback,
	}
}

// SetLevel 起始等级
func (l *logger) SetLevel(level int) {
	l.level = level
}

// SetTrace 追踪
func (l *logger) SetTrace(trace string) {
	l.trace = trace
}

// SetRequest 请求
func (l *logger) SetRequest(request *http.Request) {
	l.request = request
}

// Debug 调试
func (l *logger) Debug(message string, v ...interface{}) {
	if l.level > DEBUG {
		return
	}
	l.Raw.Println(l.preprocessing(message, DEBUG, v...))
}

// Info 信息
func (l *logger) Info(message string, v ...interface{}) {
	if l.level > INFO {
		return
	}
	l.Raw.Println(l.preprocessing(message, INFO, v...))
}

// Notice 通知
func (l *logger) Notice(message string, v ...interface{}) {
	if l.level > NOTICE {
		return
	}
	l.Raw.Println(l.preprocessing(message, NOTICE, v...))
}

// Warn 警告
func (l *logger) Warn(message string, v ...interface{}) {
	if l.level > WARN {
		return
	}
	l.Raw.Println(l.preprocessing(message, WARN, v...))
}

// Error 错误
func (l *logger) Error(message string, v ...interface{}) {
	if l.level > ERROR {
		return
	}
	l.Raw.Println(l.preprocessing(message, ERROR, v...))
}

// Fatal 致命错误
func (l *logger) Fatal(message string, v ...interface{}) {
	if l.level > FATAL {
		return
	}
	l.Raw.Fatalln(l.preprocessing(message, FATAL, v...))
}

// Panic 恐慌
func (l *logger) Panic(message string, v ...interface{}) {
	if l.level > PANIC {
		return
	}
	l.Raw.Panicln(l.preprocessing(message, PANIC, v...))
}

// SetLevel 起始等级
func SetLevel(level int) {
	Logger.SetLevel(level)
}

// SetTrace 追踪
func SetTrace(trace string) {
	Logger.SetTrace(trace)
}

// SetRequest 请求
func SetRequest(request *http.Request) {
	Logger.SetRequest(request)
}

// Debug 调试
func Debug(message string, v ...interface{}) {
	Logger.Debug(message, v...)
}

// Info 信息
func Info(message string, v ...interface{}) {
	Logger.Info(message, v...)
}

// Notice 通知
func Notice(message string, v ...interface{}) {
	Logger.Notice(message, v...)
}

// Warn 警告
func Warn(message string, v ...interface{}) {
	Logger.Warn(message, v...)
}

// Error 错误
func Error(message string, v ...interface{}) {
	Logger.Error(message, v...)
}

// Fatal 致命错误
func Fatal(message string, v ...interface{}) {
	Logger.Fatal(message, v...)
}

// Panic 恐慌
func Panic(message string, v ...interface{}) {
	Logger.Panic(message, v...)
}

func (l *logger) preprocessing(message string, level int, v ...interface{}) string {
	full := map[string]interface{}{
		"datetime":   datetime.Any(l.timezone, l.timeFormat),
		"env":        l.env,
		"channel":    "",
		"level":      level,
		"level_name": levelFlags[level],
		"trace":      l.trace,
		"ip":         "",
		"command":    strings.Join(os.Args, " "),
		"message":    l.Raw.Prefix() + message,
		"context":    json.Encode(v),
		"method":     "",
		"url":        "",
		"ua":         "",
		"referer":    "",
	}

	ips := ip.Local(true)
	if len(ips) > 0 {
		full["ip"] = ips[0]
	}

	if level > 2 {
		full["extra"] = safe.Trace(10)
	}

	if l.request != nil {
		full["method"] = l.request.Method
		full["url"] = l.request.URL.String()
		full["ua"] = l.request.UserAgent()
		full["referer"] = l.request.Referer()
		// todo:: client_ip
	}

	if l.callback != nil {
		l.callback(full)
	}

	return json.Encode(full)
}
