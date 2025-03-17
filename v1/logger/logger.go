package logger

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/lynnclub/go/v1/datetime"
	"github.com/lynnclub/go/v1/encoding/json"
	"github.com/lynnclub/go/v1/ip"
)

const (
	DEBUG  = 100
	INFO   = 200
	NOTICE = 250
	WARN   = 300
	ERROR  = 400
	PANIC  = 500
	FATAL  = 600
)

var (
	levelFlags = map[int]string{DEBUG: "DEBUG", INFO: "INFO", NOTICE: "NOTICE", WARN: "WARN", ERROR: "ERROR", PANIC: "PANIC", FATAL: "FATAL"}
	Logger     = New(log.New(os.Stderr, "", log.Lmsgprefix), "local", DEBUG, "asia/shanghai", datetime.LayoutDateTimeZoneT, nil)
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

// Panic 恐慌
func (l *logger) Panic(message string, v ...interface{}) {
	if l.level > PANIC {
		return
	}
	l.Raw.Panicln(l.preprocessing(message, PANIC, v...))
}

// Fatal 致命错误
func (l *logger) Fatal(message string, v ...interface{}) {
	if l.level > FATAL {
		return
	}
	l.Raw.Fatalln(l.preprocessing(message, FATAL, v...))
}

// Debugf 调试
func (l *logger) Debugf(format string, v ...interface{}) {
	l.Debug(fmt.Sprintf(format, v...))
}

// Infof 信息
func (l *logger) Infof(format string, v ...interface{}) {
	l.Info(fmt.Sprintf(format, v...))
}

// Noticef 通知
func (l *logger) Noticef(format string, v ...interface{}) {
	l.Notice(fmt.Sprintf(format, v...))
}

// Warnf 警告
func (l *logger) Warnf(format string, v ...interface{}) {
	l.Warn(fmt.Sprintf(format, v...))
}

// Errorf 错误
func (l *logger) Errorf(format string, v ...interface{}) {
	l.Error(fmt.Sprintf(format, v...))
}

// Panicf 恐慌
func (l *logger) Panicf(format string, v ...interface{}) {
	l.Panic(fmt.Sprintf(format, v...))
}

// Fatalf 致命错误
func (l *logger) Fatalf(format string, v ...interface{}) {
	l.Fatal(fmt.Sprintf(format, v...))
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

// Panic 恐慌
func Panic(message string, v ...interface{}) {
	Logger.Panic(message, v...)
}

// Fatal 致命错误
func Fatal(message string, v ...interface{}) {
	Logger.Fatal(message, v...)
}

// Debugf 调试
func Debugf(format string, v ...interface{}) {
	Logger.Debug(fmt.Sprintf(format, v...))
}

// Infof 信息
func Infof(format string, v ...interface{}) {
	Logger.Info(fmt.Sprintf(format, v...))
}

// Noticef 通知
func Noticef(format string, v ...interface{}) {
	Logger.Notice(fmt.Sprintf(format, v...))
}

// Warnf 警告
func Warnf(format string, v ...interface{}) {
	Logger.Warn(fmt.Sprintf(format, v...))
}

// Errorf 错误
func Errorf(format string, v ...interface{}) {
	Logger.Error(fmt.Sprintf(format, v...))
}

// Panicf 恐慌
func Panicf(format string, v ...interface{}) {
	Logger.Panic(fmt.Sprintf(format, v...))
}

// Fatalf 致命错误
func Fatalf(format string, v ...interface{}) {
	Logger.Fatal(fmt.Sprintf(format, v...))
}

func (l *logger) preprocessing(message string, level int, v ...interface{}) string {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

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
		"memory":     memStats.Alloc,
		"method":     "",
		"url":        "",
		"ua":         "",
		"referer":    "",
	}

	ips := ip.Local(true)
	if len(ips) > 0 {
		full["ip"] = ips[0]
	}

	if level > 250 {
		full["extra"] = Trace(4, 10)
	}

	if l.request == nil {
		full["channel"] = "script"
	} else {
		full["channel"] = "api"
		full["method"] = l.request.Method
		full["url"] = l.request.URL.String()
		full["ua"] = l.request.UserAgent()
		full["referer"] = l.request.Referer()

		ips = ip.GetClients(l.request)
		if len(ips) > 0 {
			full["ip"] = ips[0]
		}
	}

	if l.callback != nil {
		l.callback(full)
	}

	return json.Encode(full)
}

// Trace 执行链路
func Trace(skip, deep int) []string {
	trace := make([]string, 0)

	pcs := make([]uintptr, deep)
	deeps := runtime.Callers(skip, pcs)
	for current := range deeps {
		function := runtime.FuncForPC(pcs[current])
		file, line := function.FileLine(pcs[current])
		trace = append(trace, "["+strconv.Itoa(current)+"] "+function.Name()+"()")
		trace = append(trace, file+":"+strconv.Itoa(line))
	}

	return trace
}
