package logger

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/lynnclub/go/v1/datetime"
	"gopkg.in/natefinch/lumberjack.v2"
)

// TestLogWithLumberjack 测试使用 lumberjack 记录日志到文件
func TestLogWithLumberjack(t *testing.T) {
	// 使用临时文件名
	tempFile := "test_logger_temp.log"
	defer os.Remove(tempFile) // 测试结束后清理

	var buf bytes.Buffer
	// 先测试基本的 logger
	testLogger := New(
		log.New(&buf, "", log.Lmsgprefix),
		"local",
		INFO,
		"asia/shanghai",
		datetime.LayoutDateTimeZoneT,
		nil,
	)

	testLogger.Info("test info message")
	if !strings.Contains(buf.String(), "test info message") {
		t.Error("Expected log to contain info message")
	}

	// Debug 不应该输出（级别是 INFO）
	buf.Reset()
	testLogger.Debug("debug message")
	if buf.Len() > 0 {
		t.Error("Debug message should not appear when level is INFO")
	}

	// 测试使用 lumberjack 写入文件
	lumberjackLogger := &lumberjack.Logger{
		Filename:   tempFile,
		MaxSize:    10, // megabytes
		MaxBackups: 1,
		MaxAge:     1,     // days
		Compress:   false, // 不压缩，方便测试读取
	}
	defer lumberjackLogger.Close()

	fileLogger := New(
		log.New(lumberjackLogger, "", log.Lmsgprefix),
		"local",
		INFO,
		"asia/shanghai",
		datetime.LayoutDateTimeZoneT,
		nil,
	)

	fileLogger.Info("file log message")
	fileLogger.Error("error message")

	// 确保内容已写入
	lumberjackLogger.Close()

	// 验证文件存在且有内容
	if _, err := os.Stat(tempFile); os.IsNotExist(err) {
		t.Error("Log file was not created")
	}

	// 读取文件内容验证
	content, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "file log message") {
		t.Error("Log file does not contain expected info message")
	}
	if !strings.Contains(contentStr, "error message") {
		t.Error("Log file does not contain expected error message")
	}
}

// TestLogLevels 测试所有日志级别
func TestLogLevels(t *testing.T) {
	var buf bytes.Buffer
	testLogger := New(
		log.New(&buf, "", log.Lmsgprefix),
		"test",
		DEBUG,
		"asia/shanghai",
		datetime.LayoutDateTimeZoneT,
		nil,
	)

	tests := []struct {
		name     string
		logFunc  func(string, ...interface{})
		message  string
		expected string
	}{
		{"Debug", testLogger.Debug, "debug message", "DEBUG"},
		{"Info", testLogger.Info, "info message", "INFO"},
		{"Notice", testLogger.Notice, "notice message", "NOTICE"},
		{"Warn", testLogger.Warn, "warn message", "WARN"},
		{"Error", testLogger.Error, "error message", "ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc(tt.message)
			output := buf.String()
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected log to contain %s, got: %s", tt.expected, output)
			}
			if !strings.Contains(output, tt.message) {
				t.Errorf("Expected log to contain message %s, got: %s", tt.message, output)
			}
		})
	}
}

// TestLogFormattedMethods 测试格式化日志方法
func TestLogFormattedMethods(t *testing.T) {
	var buf bytes.Buffer
	testLogger := New(
		log.New(&buf, "", log.Lmsgprefix),
		"test",
		DEBUG,
		"asia/shanghai",
		datetime.LayoutDateTimeZoneT,
		nil,
	)

	tests := []struct {
		name     string
		logFunc  func(string, ...interface{})
		format   string
		args     []interface{}
		expected string
	}{
		{"Debugf", testLogger.Debugf, "debug: %s %d", []interface{}{"test", 123}, "debug: test 123"},
		{"Infof", testLogger.Infof, "info: %s", []interface{}{"test"}, "info: test"},
		{"Noticef", testLogger.Noticef, "notice: %v", []interface{}{true}, "notice: true"},
		{"Warnf", testLogger.Warnf, "warn: %d", []interface{}{456}, "warn: 456"},
		{"Errorf", testLogger.Errorf, "error: %s", []interface{}{"fail"}, "error: fail"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc(tt.format, tt.args...)
			output := buf.String()
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected log to contain %s, got: %s", tt.expected, output)
			}
		})
	}
}

// TestSetLevel 测试设置日志级别
func TestSetLevel(t *testing.T) {
	var buf bytes.Buffer
	testLogger := New(
		log.New(&buf, "", log.Lmsgprefix),
		"test",
		INFO,
		"asia/shanghai",
		datetime.LayoutDateTimeZoneT,
		nil,
	)

	// 默认级别是 INFO，DEBUG 不应该输出
	testLogger.Debug("should not appear")
	if buf.Len() > 0 {
		t.Error("Debug log should not appear when level is INFO")
	}

	// 设置为 DEBUG 级别
	testLogger.SetLevel(DEBUG)
	testLogger.Debug("should appear")
	if buf.Len() == 0 {
		t.Error("Debug log should appear when level is DEBUG")
	}

	// 设置为 ERROR 级别
	buf.Reset()
	testLogger.SetLevel(ERROR)
	testLogger.Info("should not appear")
	if buf.Len() > 0 {
		t.Error("Info log should not appear when level is ERROR")
	}

	testLogger.Error("should appear")
	if buf.Len() == 0 {
		t.Error("Error log should appear when level is ERROR")
	}
}

// TestSetTrace 测试设置追踪标识
func TestSetTrace(t *testing.T) {
	var buf bytes.Buffer
	testLogger := New(
		log.New(&buf, "", log.Lmsgprefix),
		"test",
		INFO,
		"asia/shanghai",
		datetime.LayoutDateTimeZoneT,
		nil,
	)

	traceId := "trace-123-456"
	testLogger.SetTrace(traceId)
	testLogger.Info("test message")

	output := buf.String()
	if !strings.Contains(output, traceId) {
		t.Errorf("Expected log to contain trace id %s, got: %s", traceId, output)
	}
}

// TestSetRequest 测试设置请求
func TestSetRequest(t *testing.T) {
	var buf bytes.Buffer
	testLogger := New(
		log.New(&buf, "", log.Lmsgprefix),
		"test",
		INFO,
		"asia/shanghai",
		datetime.LayoutDateTimeZoneT,
		nil,
	)

	req, err := http.NewRequest("GET", "http://example.com/test?param=value", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("User-Agent", "TestAgent/1.0")
	req.Header.Set("Referer", "http://example.com")

	testLogger.SetRequest(req)
	testLogger.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "GET") {
		t.Errorf("Expected log to contain method GET, got: %s", output)
	}
	if !strings.Contains(output, "/test") {
		t.Errorf("Expected log to contain URL path, got: %s", output)
	}
	if !strings.Contains(output, "api") {
		t.Errorf("Expected log to contain channel 'api', got: %s", output)
	}
}

// TestGlobalFunctions 测试全局函数
func TestGlobalFunctions(t *testing.T) {
	var buf bytes.Buffer
	Logger = New(
		log.New(&buf, "", log.Lmsgprefix),
		"test",
		DEBUG,
		"asia/shanghai",
		datetime.LayoutDateTimeZoneT,
		nil,
	)

	SetLevel(INFO)
	SetTrace("global-trace")

	Info("global info")
	output := buf.String()
	if !strings.Contains(output, "global info") {
		t.Error("Global Info function failed")
	}
	if !strings.Contains(output, "global-trace") {
		t.Error("Global SetTrace function failed")
	}

	buf.Reset()
	Infof("formatted: %d", 789)
	if !strings.Contains(buf.String(), "formatted: 789") {
		t.Error("Global Infof function failed")
	}
}

// TestLogWithContext 测试带上下文的日志
func TestLogWithContext(t *testing.T) {
	var buf bytes.Buffer
	testLogger := New(
		log.New(&buf, "", log.Lmsgprefix),
		"test",
		INFO,
		"asia/shanghai",
		datetime.LayoutDateTimeZoneT,
		nil,
	)

	contextData := map[string]interface{}{
		"user_id": 123,
		"action":  "login",
	}
	testLogger.Info("user action", contextData)

	output := buf.String()
	if !strings.Contains(output, "user action") {
		t.Error("Expected log to contain message")
	}
	// 上下文数据会被序列化到 context 字段
	if !strings.Contains(output, "context") {
		t.Error("Expected log to contain context field")
	}
}

// TestTrace 测试追踪函数
func TestTrace(t *testing.T) {
	trace := Trace(1, 5)

	if len(trace) == 0 {
		t.Error("Expected trace to have entries")
	}

	// 检查追踪信息格式
	hasFileInfo := false
	for _, entry := range trace {
		if strings.Contains(entry, "logger_test.go") {
			hasFileInfo = true
			break
		}
	}
	if !hasFileInfo {
		t.Error("Expected trace to contain file information")
	}
}

// TestCallback 测试回调函数
func TestCallback(t *testing.T) {
	var buf bytes.Buffer
	callbackCalled := false
	var capturedLog LogEntry

	testLogger := New(
		log.New(&buf, "", log.Lmsgprefix),
		"test",
		INFO,
		"asia/shanghai",
		datetime.LayoutDateTimeZoneT,
		func(log LogEntry) {
			callbackCalled = true
			capturedLog = log
		},
	)

	testLogger.Info("callback test")

	if !callbackCalled {
		t.Error("Expected callback to be called")
	}

	if capturedLog.Message != "callback test" {
		t.Errorf("Expected callback to receive message 'callback test', got: %s", capturedLog.Message)
	}

	if capturedLog.Level != INFO {
		t.Errorf("Expected callback to receive INFO level, got: %d", capturedLog.Level)
	}
}

// TestLogEntry 测试日志条目结构
func TestLogEntry(t *testing.T) {
	entry := LogEntry{
		Datetime:  time.Now().Format(time.RFC3339),
		Env:       "test",
		Channel:   "api",
		Level:     ERROR,
		LevelName: "ERROR",
		Trace:     "trace-123",
		IP:        "127.0.0.1",
		Command:   "test command",
		Message:   "test message",
		Context:   "{}",
		Memory:    1024,
		Method:    "POST",
		URL:       "/api/test",
		UserAgent: "TestAgent",
		Referer:   "http://example.com",
	}

	if entry.Level != ERROR {
		t.Errorf("Expected level ERROR, got: %d", entry.Level)
	}
	if entry.LevelName != "ERROR" {
		t.Errorf("Expected level name ERROR, got: %s", entry.LevelName)
	}
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

// TestFeishuAlertAdd 测试添加飞书告警配置
func TestFeishuAlertAdd(t *testing.T) {
	alert := &FeishuAlert{}
	option := Option{
		Levels:    []string{"ERROR", "WARN"},
		Webhook:   "http://test.webhook",
		SignKey:   "test_key",
		UserId:    "user123",
		KibanaUrl: "http://kibana.test",
		EsIndex:   "test-index",
	}

	alert.Add("test_app", option)

	if alert.options == nil {
		t.Error("Expected options to be initialized")
	}

	if _, ok := alert.options["test_app"]; !ok {
		t.Error("Expected test_app to be added to options")
	}
}

// TestFeishuAlertAddMap 测试通过 map 添加配置
func TestFeishuAlertAddMap(t *testing.T) {
	alert := &FeishuAlert{}
	setting := map[string]interface{}{
		"levels":     []string{"error", "warn"},
		"webhook":    "http://test.webhook",
		"sign_key":   "test_key",
		"user_id":    "user123",
		"kibana_url": "http://kibana.test",
		"es_index":   "test-index",
	}

	alert.AddMap("test_app", setting)

	if _, ok := alert.options["test_app"]; !ok {
		t.Error("Expected test_app to be added to options")
	}

	opt := alert.options["test_app"]
	if len(opt.Levels) != 2 {
		t.Errorf("Expected 2 levels, got: %d", len(opt.Levels))
	}
	// 级别应该被转换为大写
	if opt.Levels[0] != "ERROR" {
		t.Errorf("Expected level to be uppercase ERROR, got: %s", opt.Levels[0])
	}
}

// TestFeishuAlertAddMapBatch 测试批量添加配置
func TestFeishuAlertAddMapBatch(t *testing.T) {
	alert := &FeishuAlert{}
	batch := map[string]interface{}{
		"app1": map[string]interface{}{
			"levels":     []string{"error"},
			"webhook":    "http://webhook1",
			"sign_key":   "key1",
			"user_id":    "user1",
			"kibana_url": "http://kibana1",
			"es_index":   "index1",
		},
		"app2": map[string]interface{}{
			"levels":     []string{"warn"},
			"webhook":    "http://webhook2",
			"sign_key":   "key2",
			"user_id":    "user2",
			"kibana_url": "http://kibana2",
			"es_index":   "index2",
		},
	}

	alert.AddMapBatch(batch)

	if len(alert.options) != 2 {
		t.Errorf("Expected 2 options, got: %d", len(alert.options))
	}

	if _, ok := alert.options["app1"]; !ok {
		t.Error("Expected app1 to be added")
	}
	if _, ok := alert.options["app2"]; !ok {
		t.Error("Expected app2 to be added")
	}
}

// TestFeishuAlertFindOption 测试查找配置
func TestFeishuAlertFindOption(t *testing.T) {
	alert := &FeishuAlert{}
	alert.Add("user", Option{
		Levels:  []string{"ERROR"},
		Webhook: "http://webhook1",
		SignKey: "key1",
	})
	alert.Add("order", Option{
		Levels:  []string{"WARN"},
		Webhook: "http://webhook2",
		SignKey: "key2",
	})
	alert.Add("default_api", Option{
		Levels:  []string{},
		Webhook: "http://webhook_default",
		SignKey: "key_default",
	})

	tests := []struct {
		name         string
		levelName    string
		entry        string
		defaultName  string
		expectedName string
	}{
		{
			name:         "Match user with ERROR",
			levelName:    "ERROR",
			entry:        "/api/user/login",
			defaultName:  "default_api",
			expectedName: "user",
		},
		{
			name:         "Match order with WARN",
			levelName:    "WARN",
			entry:        "/api/order/create",
			defaultName:  "default_api",
			expectedName: "order",
		},
		{
			name:         "No match, return default",
			levelName:    "INFO",
			entry:        "/api/unknown",
			defaultName:  "default_api",
			expectedName: "default_api",
		},
		{
			name:         "Level not in allowed list",
			levelName:    "INFO",
			entry:        "/api/user/login",
			defaultName:  "default_api",
			expectedName: "default_api",
		},
		{
			name:         "Match with empty levels (any level)",
			levelName:    "INFO",
			entry:        "/api/default_api/test",
			defaultName:  "fallback",
			expectedName: "default_api",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := alert.FindOption(tt.levelName, tt.entry, tt.defaultName)
			if result != tt.expectedName {
				t.Errorf("Expected %s, got: %s", tt.expectedName, result)
			}
		})
	}
}

// TestFeishuAlertFormat 测试格式化消息
func TestFeishuAlertFormat(t *testing.T) {
	alert := &FeishuAlert{}
	log := LogEntry{
		Level:     ERROR,
		LevelName: "ERROR",
		Datetime:  "2026-01-16T10:00:00Z",
		Env:       "production",
		IP:        "192.168.1.100",
		Trace:     "trace-123",
		Command:   "/app/server",
		Message:   "Database connection failed",
		Extra:     []string{"trace info"},
	}

	content := alert.Format(log, "http://kibana.example.com", "app-logs-*")

	// 检查格式化内容包含关键信息
	expectedFields := []string{
		"production",
		"ERROR",
		"2026-01-16T10:00:00Z",
		"192.168.1.100",
		"trace-123",
		"Database connection failed",
	}

	for _, field := range expectedFields {
		if !strings.Contains(content, field) {
			t.Errorf("Expected formatted content to contain %s, got: %s", field, content)
		}
	}
}

// TestFeishuAlertSendLowLevel 测试低级别日志不发送
func TestFeishuAlertSendLowLevel(t *testing.T) {
	alert := &FeishuAlert{}
	option := Option{
		Levels:  []string{"ERROR"},
		Webhook: "http://test.webhook",
		SignKey: "test_key",
	}
	alert.Add("default_api", option)

	// INFO 级别的日志不应该触发发送
	log := LogEntry{
		Level:     INFO,
		LevelName: "INFO",
		Message:   "info message",
		Command:   "test",
		Datetime:  time.Now().Format(time.RFC3339),
	}

	alert.Send(log)

	// 验证没有添加到历史记录（因为级别太低被跳过）
	if len(alert.lastHashs) != 0 {
		t.Error("Expected no logs in history for low level logs")
	}
}

// TestNewFeishu 测试创建飞书告警实例
func TestNewFeishu(t *testing.T) {
	options := map[string]interface{}{
		"app1": map[string]interface{}{
			"levels":     []string{"error"},
			"webhook":    "http://webhook1",
			"sign_key":   "key1",
			"user_id":    "user1",
			"kibana_url": "http://kibana1",
			"es_index":   "index1",
		},
	}

	instance := NewFeishu(options)

	if instance == nil {
		t.Error("Expected non-nil instance")
	}

	if len(instance.options) != 1 {
		t.Errorf("Expected 1 option, got: %d", len(instance.options))
	}
}

// TestMoreGlobalFunctions 测试更多的全局函数
func TestMoreGlobalFunctions(t *testing.T) {
	var buf bytes.Buffer
	Logger = New(
		log.New(&buf, "", log.Lmsgprefix),
		"test",
		DEBUG,
		"asia/shanghai",
		datetime.LayoutDateTimeZoneT,
		nil,
	)

	// 测试 Notice
	buf.Reset()
	Notice("notice message")
	if !strings.Contains(buf.String(), "notice message") {
		t.Error("Global Notice function failed")
	}

	// 测试 Warn
	buf.Reset()
	Warn("warn message")
	if !strings.Contains(buf.String(), "warn message") {
		t.Error("Global Warn function failed")
	}

	// 测试 Debugf
	buf.Reset()
	Debugf("debug: %s", "test")
	if !strings.Contains(buf.String(), "debug: test") {
		t.Error("Global Debugf function failed")
	}

	// 测试 Noticef
	buf.Reset()
	Noticef("notice: %d", 123)
	if !strings.Contains(buf.String(), "notice: 123") {
		t.Error("Global Noticef function failed")
	}

	// 测试 Warnf
	buf.Reset()
	Warnf("warn: %v", true)
	if !strings.Contains(buf.String(), "warn: true") {
		t.Error("Global Warnf function failed")
	}

	// 测试 Errorf
	buf.Reset()
	Errorf("error: %s", "fail")
	if !strings.Contains(buf.String(), "error: fail") {
		t.Error("Global Errorf function failed")
	}

	// 测试 SetRequest
	req, err := http.NewRequest("POST", "http://test.com/api", nil)
	if err != nil {
		t.Fatal(err)
	}
	SetRequest(req)
	buf.Reset()
	Info("request test")
	if !strings.Contains(buf.String(), "POST") {
		t.Error("Global SetRequest function failed")
	}
}

// TestPanicMethod 测试 Panic 方法
func TestPanicMethod(t *testing.T) {
	var buf bytes.Buffer
	testLogger := New(
		log.New(&buf, "", log.Lmsgprefix),
		"test",
		DEBUG,
		"asia/shanghai",
		datetime.LayoutDateTimeZoneT,
		nil,
	)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic to be called")
		} else {
			if !strings.Contains(buf.String(), "panic message") {
				t.Error("Expected panic log to contain message")
			}
		}
	}()

	testLogger.Panic("panic message")
}

// TestPanicfMethod 测试 Panicf 方法
func TestPanicfMethod(t *testing.T) {
	var buf bytes.Buffer
	testLogger := New(
		log.New(&buf, "", log.Lmsgprefix),
		"test",
		DEBUG,
		"asia/shanghai",
		datetime.LayoutDateTimeZoneT,
		nil,
	)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panicf to be called")
		} else {
			if !strings.Contains(buf.String(), "panic: test") {
				t.Error("Expected panicf log to contain formatted message")
			}
		}
	}()

	testLogger.Panicf("panic: %s", "test")
}

// TestGlobalPanic 测试全局 Panic 函数
func TestGlobalPanic(t *testing.T) {
	var buf bytes.Buffer
	Logger = New(
		log.New(&buf, "", log.Lmsgprefix),
		"test",
		DEBUG,
		"asia/shanghai",
		datetime.LayoutDateTimeZoneT,
		nil,
	)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected global panic to be called")
		}
	}()

	Panic("global panic")
}

// TestGlobalPanicf 测试全局 Panicf 函数
func TestGlobalPanicf(t *testing.T) {
	var buf bytes.Buffer
	Logger = New(
		log.New(&buf, "", log.Lmsgprefix),
		"test",
		DEBUG,
		"asia/shanghai",
		datetime.LayoutDateTimeZoneT,
		nil,
	)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected global panicf to be called")
		}
	}()

	Panicf("global panic: %s", "test")
}

// TestLogLevelFiltering 测试日志级别过滤
func TestLogLevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	testLogger := New(
		log.New(&buf, "", log.Lmsgprefix),
		"test",
		WARN,
		"asia/shanghai",
		datetime.LayoutDateTimeZoneT,
		nil,
	)

	// 低于 WARN 的日志不应该输出
	testLogger.Debug("debug")
	testLogger.Info("info")
	testLogger.Notice("notice")

	if buf.Len() > 0 {
		t.Error("Expected no output for logs below WARN level")
	}

	// WARN 及以上的日志应该输出
	testLogger.Warn("warn")
	if buf.Len() == 0 {
		t.Error("Expected output for WARN level")
	}

	buf.Reset()
	testLogger.Error("error")
	if buf.Len() == 0 {
		t.Error("Expected output for ERROR level")
	}
}

// TestFeishuAlertAddPanic 测试添加空 webhook 应该 panic
func TestFeishuAlertAddPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when webhook is empty")
		}
	}()

	alert := &FeishuAlert{}
	option := Option{
		Webhook: "", // 空 webhook 应该触发 panic
		SignKey: "test",
	}
	alert.Add("test", option)
}

// TestFeishuAlertSendWithTrace 测试带追踪信息的日志发送
func TestFeishuAlertSendWithTrace(t *testing.T) {
	alert := &FeishuAlert{}
	option := Option{
		Levels:  []string{"ERROR"},
		Webhook: "http://test.webhook",
		SignKey: "test_key",
	}
	alert.Add("default_command", option)

	// 带 Trace 的日志
	log := LogEntry{
		Level:     ERROR,
		LevelName: "ERROR",
		Trace:     "trace-123",
		Message:   "test with trace",
		Command:   "test_command",
		Datetime:  time.Now().Format(time.RFC3339),
		Env:       "test",
		IP:        "127.0.0.1",
	}

	alert.Send(log)

	if len(alert.lastHashs) != 1 {
		t.Errorf("Expected 1 log in history, got %d", len(alert.lastHashs))
	}
}

// TestFeishuAlertSendWithURL 测试带 URL 的日志发送
func TestFeishuAlertSendWithURL(t *testing.T) {
	alert := &FeishuAlert{}
	option := Option{
		Levels:  []string{"ERROR"},
		Webhook: "http://test.webhook",
		SignKey: "test_key",
	}
	alert.Add("api", option)

	// 带 URL 的日志
	log := LogEntry{
		Level:     ERROR,
		LevelName: "ERROR",
		URL:       "/api/test",
		Message:   "test with url",
		Command:   "server",
		Datetime:  time.Now().Format(time.RFC3339),
		Env:       "test",
		IP:        "127.0.0.1",
	}

	alert.Send(log)

	if len(alert.lastHashs) != 1 {
		t.Errorf("Expected 1 log in history, got %d", len(alert.lastHashs))
	}
}

// TestFeishuAlertSendWithExtra 测试带 Extra 字段的日志
func TestFeishuAlertSendWithExtra(t *testing.T) {
	alert := &FeishuAlert{}
	option := Option{
		Levels:  []string{"ERROR"},
		Webhook: "http://test.webhook",
		SignKey: "test_key",
	}
	alert.Add("default_command", option)

	// 带 Extra 的日志
	log := LogEntry{
		Level:     ERROR,
		LevelName: "ERROR",
		Message:   "test with extra",
		Command:   "test_command",
		Datetime:  time.Now().Format(time.RFC3339),
		Env:       "test",
		IP:        "127.0.0.1",
		Extra:     []string{"trace line 1", "trace line 2"},
	}

	alert.Send(log)

	if len(alert.lastHashs) != 1 {
		t.Errorf("Expected 1 log in history, got %d", len(alert.lastHashs))
	}
}
