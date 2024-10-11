package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/lynnclub/go/v1/algorithm"
	"github.com/lynnclub/go/v1/datetime"
	"github.com/lynnclub/go/v1/elasticsearch"
	"github.com/lynnclub/go/v1/encoding/json"
	"github.com/lynnclub/go/v1/ip"
	"github.com/lynnclub/go/v1/notice"
	"github.com/lynnclub/go/v1/safe"
)

// 日志，实时告警
type logger struct {
	Raw        *log.Logger       // 原生log
	level      int               // 起始级别
	env        string            // 环境
	timezone   string            // 时区
	timeFormat string            // 时间格式
	feishu     map[string]string // 飞书通知配置
	mu         sync.Mutex
	lastHashs  []AlertHash
}

type AlertHash struct {
	hash string
	Time time.Time
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
		"context":    json.Encode(v),
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

	ips := ip.Local(true)
	if len(ips) > 0 {
		full["ip"] = ips[0]
	}

	if level > 2 {
		full["extra"] = safe.Trace(10)
	}

	// 自动告警
	if level >= 2 && l.feishu != nil {
		l.alert(full)
	}

	return json.Encode(full)
}

func (l *logger) alert(full map[string]interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	hash := AlertHash{
		hash: algorithm.MD5(full["message"].(string)),
		Time: time.Now(),
	}

	// 阻止重复报警
	for _, item := range l.lastHashs {
		if item.hash == hash.hash {
			if item.Time.Add(10 * time.Minute).After(hash.Time) {
				return
			}
		}
	}

	// 如果超过最大容量，移除最老的
	l.lastHashs = append(l.lastHashs, hash)
	if len(l.lastHashs) > 10 {
		l.lastHashs = l.lastHashs[1:]
	}

	safe.Catch(func() {
		// Send 飞书群发送
		content := map[string]interface{}{
			"tag":  "text",
			"text": l.formatText(full),
		}
		feiShuGroup.Send("", content, l.feishu["user_id"])
	}, func(err any) {
		println(err)
	})
}

func (l *logger) formatText(full map[string]interface{}) string {
	querys := []string{}
	if trace, exists := full["trace"].(string); exists && trace != "" {
		querys = append(querys, elasticsearch.GetKuery("trace", trace))
	}
	querys = append(querys, elasticsearch.GetKuery("message", full["message"].(string)))

	return fmt.Sprintf(`环境：%s
级别：%s
时间：%s
IP：%s
追踪：%s
入口：%s
	
%s

%s

如有问题请尽快处理 []~(￣▽￣)~*`,
		full["env"],
		full["level_name"],
		full["datetime"],
		full["ip"],
		full["trace"],
		full["command"],
		full["message"],
		elasticsearch.GetKibanaUrl(l.feishu["kibana_url"], l.feishu["es_index"], querys),
	)
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
