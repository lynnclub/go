package logger

import (
	"fmt"
	"github.com/lynnclub/go/v1/datetime"
	"github.com/lynnclub/go/v1/encoding/json"
	"github.com/lynnclub/go/v1/notice"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"
)

// 日志，自动告警
type logger struct {
	raw      *log.Logger
	level    int
	prefix   string
	timezone string
	feishu   map[string]string // 飞书通知配置
}

var (
	callerDepth = 3
	levelFlags  = []string{"DEBUG", "INFO", "WARN", "ERROR", "PANIC", "FATAL"}
	Logger      = New("", DEBUG, "asia/shanghai", nil)
)

const (
	DEBUG int = iota
	INFO
	WARN
	ERROR
	PANIC
	FATAL
)

// New 初始化
func New(prefix string, level int, timezone string, feishu map[string]string) *logger {
	raw := log.New(os.Stderr, "", log.LstdFlags)
	raw.SetFlags(0)

	return &logger{
		raw:      raw,
		level:    level,
		prefix:   prefix,
		timezone: timezone,
		feishu:   feishu,
	}
}

// SetLevel 等级
func (l *logger) SetLevel(level int) {
	l.level = level
}

// SetPrefix 前缀
func (l *logger) SetPrefix(prefix string) {
	l.prefix = prefix
}

// preprocessing 预处理
func (l *logger) preprocessing(level int, v ...interface{}) string {
	flag := levelFlags[level]

	var source interface{}
	if level > 2 {
		pcs := make([]uintptr, 10)
		deeps := runtime.Callers(callerDepth, pcs)

		trace := make([]map[string]string, deeps)
		for deep := 0; deep < deeps; deep++ {
			function := runtime.FuncForPC(pcs[deep])
			file, line := function.FileLine(pcs[deep])
			trace[deep] = map[string]string{
				"deep": strconv.Itoa(deep),
				"file": file + ":" + strconv.Itoa(line),
				"func": function.Name(),
			}
		}

		source = trace
	} else {
		_, file, line, _ := runtime.Caller(callerDepth)
		source = fmt.Sprintf("%s:%d", file, line)
	}

	title := l.prefix + v[0].(string)
	full := map[string]interface{}{
		"TIMESTAMP": time.Now().Unix(),
		"TIME":      datetime.DateTime(l.timezone),
		"LEVEL":     flag,
		"SOURCE":    source,
		"MESSAGE":   v[1:],
	}

	// 自动告警
	if level >= 2 && l.feishu != nil {
		// NewFeiShuGroup 飞书群实例化
		group := notice.NewFeiShuGroup(
			l.feishu["webhook"],
			l.feishu["sign_key"],
			l.feishu["env"],
		)

		// Send 飞书群发送
		content := map[string]interface{}{
			"tag":  "text",
			"text": json.Encode(full),
		}
		group.Send(title, content, l.feishu["user_id"])
	}

	full["TITLE"] = title
	return json.Encode(full)
}

func (l *logger) Printf(msg string, v ...interface{}) {
	tmp := map[string]interface{}{
		"use": v[1],
		"row": v[2],
		"sql": v[3],
	}
	l.raw.Println(l.preprocessing(INFO, "[SQL]", tmp))
}

// Debug 调试
func (l *logger) Debug(v ...interface{}) {
	if l.level > DEBUG {
		return
	}

	l.raw.Println(l.preprocessing(DEBUG, v...))
}

// Info 信息
func (l *logger) Info(v ...interface{}) {
	if l.level > INFO {
		return
	}

	l.raw.Println(l.preprocessing(INFO, v...))
}

// Warn 警告
func (l *logger) Warn(v ...interface{}) {
	if l.level > WARN {
		return
	}

	l.raw.Println(l.preprocessing(WARN, v...))
}

// Error 错误
func (l *logger) Error(v ...interface{}) {
	if l.level > ERROR {
		return
	}

	l.raw.Println(l.preprocessing(ERROR, v...))
}

// Panic 恐慌
func (l *logger) Panic(v ...interface{}) {
	if l.level > PANIC {
		return
	}

	l.raw.Panicln(l.preprocessing(PANIC, v...))
}

// Fatal 致命错误
func (l *logger) Fatal(v ...interface{}) {
	if l.level > FATAL {
		return
	}

	l.raw.Fatalln(l.preprocessing(FATAL, v...))
}

// SetLevel 等级
func SetLevel(level int) {
	Logger.SetLevel(level)
}

// SetPrefix 前缀
func SetPrefix(prefix string) {
	Logger.SetPrefix(prefix)
}

// Debug 调试
func Debug(v ...interface{}) {
	Logger.Debug(v...)
}

// Info 信息
func Info(v ...interface{}) {
	Logger.Info(v...)
}

// Warn 警告
func Warn(v ...interface{}) {
	Logger.Warn(v...)
}

// Error 错误
func Error(v ...interface{}) {
	Logger.Error(v...)
}

// Panic 恐慌
func Panic(v ...interface{}) {
	Logger.Panic(v...)
}

// Fatal 致命错误
func Fatal(v ...interface{}) {
	Logger.Fatal(v...)
}
