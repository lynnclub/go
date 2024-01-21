package logger

import (
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/lynnclub/go/v1/datetime"
	"github.com/lynnclub/go/v1/encoding/json"
	"github.com/lynnclub/go/v1/notice"
)

// 日志，实时告警
type logger struct {
	Raw        *log.Logger       // 原生log
	level      int               // 起始级别
	env        string            // 环境
	timezone   string            // 时区
	timeFormat string            // 时间格式
	feishu     map[string]string // 飞书通知配置
}

var (
	callerDepth = 3
	levelFlags  = []string{"DEBUG", "INFO", "NOTICE", "WARN", "ERROR", "FATAL", "PANIC"}
	Logger      = New(log.New(os.Stderr, "", log.Lmsgprefix), DEBUG, "local", "asia/shanghai", datetime.LayoutDateTimeZoneT, nil)
	feiShuGroup *notice.FeiShuGroup
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

// New 初始化
func New(raw *log.Logger, level int, env string, timezone, timeFormat string, feishu map[string]string) *logger {
	if feishu != nil {
		// 飞书群实例化
		feiShuGroup = notice.NewFeiShuGroup(
			feishu["webhook"],
			feishu["sign_key"],
			"",
		)
	}

	return &logger{
		Raw:        raw,
		level:      level,
		env:        env,
		timezone:   timezone,
		timeFormat: timeFormat,
		feishu:     feishu,
	}
}

// preprocessing 预处理
func (l *logger) preprocessing(message string, level int, v ...interface{}) string {
	message = l.Raw.Prefix() + message
	full := map[string]interface{}{
		"message":    message,
		"context":    v,
		"level":      level,
		"level_name": levelFlags[level],
		"env":        l.env,
		"channel":    "",
		"datetime":   datetime.Any(l.timezone, l.timeFormat),
		"trace":      "",
		"command":    strings.Join(os.Args, " "),
		"method":     "",
		"url":        "",
		"ua":         "",
		"referer":    "",
		"ip":         "",
	}

	if level > 2 {
		pcs := make([]uintptr, 10)
		deeps := runtime.Callers(callerDepth, pcs)

		trace := make([]string, 0)
		for deep := 0; deep < deeps; deep++ {
			function := runtime.FuncForPC(pcs[deep])
			file, line := function.FileLine(pcs[deep])
			trace = append(trace, "["+strconv.Itoa(deep)+"] "+function.Name()+"()")
			trace = append(trace, file+":"+strconv.Itoa(line))
		}

		full["extra"] = trace
	}

	// 自动告警
	if level >= 2 && l.feishu != nil {
		// Send 飞书群发送
		content := map[string]interface{}{
			"tag":  "text",
			"text": json.Encode(full),
		}
		feiShuGroup.Send(message, content, l.feishu["user_id"])
	}

	return json.Encode(full)
}

// SetLevel 等级
func (l *logger) SetLevel(level int) {
	l.level = level
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

// SetLevel 等级
func SetLevel(level int) {
	Logger.SetLevel(level)
}

// Debug 调试
func Debug(message string, v ...interface{}) {
	Logger.Debug(message, v...)
}

// Info 信息
func Info(message string, v ...interface{}) {
	Logger.Info(message, v...)
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
