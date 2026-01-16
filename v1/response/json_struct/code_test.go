package json_struct

import (
	"encoding/json"
	"testing"
)

// TestCodeSet 测试Code结构的Set方法
func TestCodeSet(t *testing.T) {
	c := &Code{}
	status := 200
	msg := "success"
	data := map[string]interface{}{"key": "value"}
	timestamp := int64(1234567890)

	c.Set(status, msg, data, timestamp)

	if c.Code != status {
		t.Errorf("期望Code为%d，实际为%d", status, c.Code)
	}
	if c.Msg != msg {
		t.Errorf("期望Msg为%s，实际为%s", msg, c.Msg)
	}
	if c.Data == nil {
		t.Error("Data不应该为nil")
	}
	if c.Timestamp != timestamp {
		t.Errorf("期望Timestamp为%d，实际为%d", timestamp, c.Timestamp)
	}
}

// TestCodeJSON 测试Code结构的JSON序列化
func TestCodeJSON(t *testing.T) {
	c := &Code{}
	c.Set(200, "success", map[string]string{"name": "test"}, 1234567890)

	jsonBytes, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		t.Fatalf("JSON反序列化失败: %v", err)
	}

	// 验证字段名为"code"而不是"status"
	if _, ok := result["code"]; !ok {
		t.Error("JSON应该包含code字段")
	}
	if int(result["code"].(float64)) != 200 {
		t.Errorf("期望code为200，实际为%v", result["code"])
	}
	if result["msg"] != "success" {
		t.Errorf("期望msg为success，实际为%v", result["msg"])
	}
}

// TestCodeWithErrorStatus 测试错误状态码
func TestCodeWithErrorStatus(t *testing.T) {
	c := &Code{}
	c.Set(500, "internal server error", nil, 1234567890)

	if c.Code != 500 {
		t.Errorf("期望Code为500，实际为%d", c.Code)
	}
	if c.Msg != "internal server error" {
		t.Errorf("期望Msg为'internal server error'，实际为%s", c.Msg)
	}
}

// TestCodeWithDifferentStatusCodes 测试不同的状态码
func TestCodeWithDifferentStatusCodes(t *testing.T) {
	testCases := []struct {
		code int
		msg  string
	}{
		{200, "OK"},
		{201, "Created"},
		{400, "Bad Request"},
		{401, "Unauthorized"},
		{403, "Forbidden"},
		{404, "Not Found"},
		{500, "Internal Server Error"},
	}

	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			c := &Code{}
			c.Set(tc.code, tc.msg, nil, 1234567890)

			if c.Code != tc.code {
				t.Errorf("期望Code为%d，实际为%d", tc.code, c.Code)
			}
			if c.Msg != tc.msg {
				t.Errorf("期望Msg为%s，实际为%s", tc.msg, c.Msg)
			}
		})
	}
}
