package notice

import (
	"testing"
	"time"
)

func TestNewFeishu(t *testing.T) {
	options := map[string]interface{}{
		"alert1": map[string]interface{}{
			"levels":     []string{"ERROR", "WARNING"},
			"webhook":    "http://test.webhook",
			"sign_key":   "test_sign_key",
			"user_id":    "test_user_id",
			"kibana_url": "http://kibana.test",
			"es_index":   "test_index",
		},
	}

	alert := NewFeishu(options)
	if alert == nil {
		t.Errorf("Expected FeishuAlert instance, got nil")
	}

	if len(alert.options) != 1 {
		t.Errorf("Expected 1 option, got %d", len(alert.options))
	}
}

func TestAddAndFindOption(t *testing.T) {
	alert := &FeishuAlert{}
	option := Option{
		Levels:    []string{"ERROR"},
		Webhook:   "http://test.webhook",
		SignKey:   "test_sign_key",
		UserId:    "test_user_id",
		KibanaUrl: "http://kibana.test",
		EsIndex:   "test_index",
	}

	alert.Add("test", option)

	name := alert.FindOption("ERROR", "test_command", "default")
	if name != "test" {
		t.Errorf("Expected 'test', got '%s'", name)
	}

	name = alert.FindOption("INFO", "test_command", "default")
	if name != "default" {
		t.Errorf("Expected 'default', got '%s'", name)
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

	log := map[string]interface{}{
		"level":      3,
		"level_name": "ERROR",
		"command":    "test_command",
		"message":    "test_message",
		"datetime":   time.Now().Format(time.RFC3339),
		"trace":      "",
		"url":        "",
		"env":        "production",
		"ip":         "127.0.0.1",
	}

	alert.Send(log)

	// 立即再次发送相同的日志，应该被过滤掉
	alert.Send(log)
	if len(alert.lastHashs) != 1 {
		t.Errorf("Expected 1 log in history, got %d", len(alert.lastHashs))
	}
}
