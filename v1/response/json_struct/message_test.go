package json_struct

import (
	"encoding/json"
	"testing"
)

// TestMessageSet 测试Message结构的Set方法
func TestMessageSet(t *testing.T) {
	m := &Message{}
	status := 200
	msg := "success"
	data := map[string]interface{}{"key": "value"}
	timestamp := int64(1234567890)

	m.Set(status, msg, data, timestamp)

	if m.Status != status {
		t.Errorf("期望Status为%d，实际为%d", status, m.Status)
	}
	if m.Msg != msg {
		t.Errorf("期望Msg为%s，实际为%s", msg, m.Msg)
	}
	if m.Data == nil {
		t.Error("Data不应该为nil")
	}
	if m.Timestamp != timestamp {
		t.Errorf("期望Timestamp为%d，实际为%d", timestamp, m.Timestamp)
	}
}

// TestMessageJSON 测试Message结构的JSON序列化
func TestMessageJSON(t *testing.T) {
	m := &Message{}
	m.Set(200, "success", map[string]string{"name": "test"}, 1234567890)

	jsonBytes, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		t.Fatalf("JSON反序列化失败: %v", err)
	}

	// 验证字段名为"message"而不是"msg"
	if _, ok := result["message"]; !ok {
		t.Error("JSON应该包含message字段")
	}
	if result["message"] != "success" {
		t.Errorf("期望message为success，实际为%v", result["message"])
	}
	if int(result["status"].(float64)) != 200 {
		t.Errorf("期望status为200，实际为%v", result["status"])
	}
}

// TestMessageWithLongMessage 测试长消息
func TestMessageWithLongMessage(t *testing.T) {
	m := &Message{}
	longMsg := "这是一个很长的消息，用于测试Message结构是否能正确处理长文本内容。" +
		"包含中文、English、数字123、特殊字符!@#$%^&*()等各种内容。"
	m.Set(200, longMsg, nil, 1234567890)

	if m.Msg != longMsg {
		t.Errorf("长消息处理不正确")
	}

	jsonBytes, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		t.Fatalf("JSON反序列化失败: %v", err)
	}

	if result["message"] != longMsg {
		t.Error("长消息序列化/反序列化不正确")
	}
}

// TestMessageWithEmptyMessage 测试空消息
func TestMessageWithEmptyMessage(t *testing.T) {
	m := &Message{}
	m.Set(200, "", map[string]string{"key": "value"}, 1234567890)

	if m.Msg != "" {
		t.Errorf("期望Msg为空字符串，实际为%s", m.Msg)
	}

	jsonBytes, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		t.Fatalf("JSON反序列化失败: %v", err)
	}

	if result["message"] != "" {
		t.Errorf("期望message为空字符串，实际为%v", result["message"])
	}
}
