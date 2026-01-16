package json_struct

import (
	"encoding/json"
	"testing"
)

// TestRawSet 测试Raw结构的Set方法
func TestRawSet(t *testing.T) {
	r := &Raw{}
	data := map[string]interface{}{"key": "value"}

	// Raw结构只使用data参数，其他参数被忽略
	r.Set(200, "ignored", data, 1234567890)

	if r.Data == nil {
		t.Error("Data不应该为nil")
	}
}

// TestRawMarshalJSON 测试Raw结构的JSON序列化
func TestRawMarshalJSON(t *testing.T) {
	r := &Raw{}
	data := map[string]string{"name": "test", "value": "123"}
	r.Set(0, "", data, 0)

	jsonBytes, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	// Raw结构应该直接序列化为data的内容，不包含外层结构
	var result map[string]interface{}
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		t.Fatalf("JSON反序列化失败: %v", err)
	}

	// 应该直接是data的内容，不包含"data"字段
	if _, ok := result["data"]; ok {
		t.Error("Raw序列化不应该包含data字段")
	}
	if result["name"] != "test" {
		t.Errorf("期望name为test，实际为%v", result["name"])
	}
	if result["value"] != "123" {
		t.Errorf("期望value为123，实际为%v", result["value"])
	}
}

// TestRawWithArray 测试Raw结构序列化数组
func TestRawWithArray(t *testing.T) {
	r := &Raw{}
	data := []string{"item1", "item2", "item3"}
	r.Set(0, "", data, 0)

	jsonBytes, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	var result []interface{}
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		t.Fatalf("JSON反序列化失败: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("期望数组长度为3，实际为%d", len(result))
	}
	if result[0] != "item1" {
		t.Errorf("期望第一个元素为item1，实际为%v", result[0])
	}
}

// TestRawWithString 测试Raw结构序列化字符串
func TestRawWithString(t *testing.T) {
	r := &Raw{}
	data := "plain string"
	r.Set(0, "", data, 0)

	jsonBytes, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	var result string
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		t.Fatalf("JSON反序列化失败: %v", err)
	}

	if result != "plain string" {
		t.Errorf("期望字符串为'plain string'，实际为%s", result)
	}
}

// TestRawWithNumber 测试Raw结构序列化数字
func TestRawWithNumber(t *testing.T) {
	r := &Raw{}
	data := 12345
	r.Set(0, "", data, 0)

	jsonBytes, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	var result int
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		t.Fatalf("JSON反序列化失败: %v", err)
	}

	if result != 12345 {
		t.Errorf("期望数字为12345，实际为%d", result)
	}
}

// TestRawWithNil 测试Raw结构序列化nil
func TestRawWithNil(t *testing.T) {
	r := &Raw{}
	r.Set(0, "", nil, 0)

	jsonBytes, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	// nil应该序列化为"null"
	if string(jsonBytes) != "null" {
		t.Errorf("期望序列化为'null'，实际为%s", string(jsonBytes))
	}
}

// TestRawWithComplexStructure 测试Raw结构序列化复杂结构
func TestRawWithComplexStructure(t *testing.T) {
	r := &Raw{}
	data := map[string]interface{}{
		"users": []map[string]interface{}{
			{"id": 1, "name": "Alice"},
			{"id": 2, "name": "Bob"},
		},
		"total": 2,
		"metadata": map[string]interface{}{
			"page": 1,
			"size": 10,
		},
	}
	r.Set(0, "", data, 0)

	jsonBytes, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		t.Fatalf("JSON反序列化失败: %v", err)
	}

	// 验证复杂结构被正确序列化
	if result["total"] != float64(2) {
		t.Errorf("期望total为2，实际为%v", result["total"])
	}

	users, ok := result["users"].([]interface{})
	if !ok {
		t.Fatal("users应该是数组类型")
	}
	if len(users) != 2 {
		t.Errorf("期望users数组长度为2，实际为%d", len(users))
	}
}

// TestRawIgnoresOtherParameters 测试Raw结构忽略其他参数
func TestRawIgnoresOtherParameters(t *testing.T) {
	r := &Raw{}
	data := map[string]string{"key": "value"}

	// 设置各种参数，但只有data应该被使用
	r.Set(999, "should be ignored", data, 9999999999)

	jsonBytes, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		t.Fatalf("JSON反序列化失败: %v", err)
	}

	// 不应该包含status、msg、timestamp等字段
	if _, ok := result["status"]; ok {
		t.Error("Raw不应该包含status字段")
	}
	if _, ok := result["msg"]; ok {
		t.Error("Raw不应该包含msg字段")
	}
	if _, ok := result["timestamp"]; ok {
		t.Error("Raw不应该包含timestamp字段")
	}

	// 应该只包含data的内容
	if result["key"] != "value" {
		t.Errorf("期望key为value，实际为%v", result["key"])
	}
}
