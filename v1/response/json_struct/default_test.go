package json_struct

import (
	"encoding/json"
	"testing"
)

// TestDefaultSet 测试Default结构的Set方法
func TestDefaultSet(t *testing.T) {
	d := &Default{}
	status := 200
	msg := "success"
	data := map[string]interface{}{"key": "value"}
	timestamp := int64(1234567890)

	d.Set(status, msg, data, timestamp)

	if d.Status != status {
		t.Errorf("期望Status为%d，实际为%d", status, d.Status)
	}
	if d.Msg != msg {
		t.Errorf("期望Msg为%s，实际为%s", msg, d.Msg)
	}
	if d.Data == nil {
		t.Error("Data不应该为nil")
	}
	if d.Timestamp != timestamp {
		t.Errorf("期望Timestamp为%d，实际为%d", timestamp, d.Timestamp)
	}
}

// TestDefaultJSON 测试Default结构的JSON序列化
func TestDefaultJSON(t *testing.T) {
	d := &Default{}
	d.Set(200, "success", map[string]string{"name": "test"}, 1234567890)

	jsonBytes, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		t.Fatalf("JSON反序列化失败: %v", err)
	}

	if int(result["status"].(float64)) != 200 {
		t.Errorf("期望status为200，实际为%v", result["status"])
	}
	if result["msg"] != "success" {
		t.Errorf("期望msg为success，实际为%v", result["msg"])
	}
	if result["data"] == nil {
		t.Error("data不应该为nil")
	}
	if int64(result["timestamp"].(float64)) != 1234567890 {
		t.Errorf("期望timestamp为1234567890，实际为%v", result["timestamp"])
	}
}

// TestDefaultWithNilData 测试nil数据
func TestDefaultWithNilData(t *testing.T) {
	d := &Default{}
	d.Set(404, "not found", nil, 1234567890)

	if d.Status != 404 {
		t.Errorf("期望Status为404，实际为%d", d.Status)
	}
	if d.Data != nil {
		t.Error("Data应该为nil")
	}
}

// TestDefaultWithArrayData 测试数组数据
func TestDefaultWithArrayData(t *testing.T) {
	d := &Default{}
	data := []string{"item1", "item2", "item3"}
	d.Set(200, "success", data, 1234567890)

	jsonBytes, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		t.Fatalf("JSON反序列化失败: %v", err)
	}

	dataArray, ok := result["data"].([]interface{})
	if !ok {
		t.Fatal("data应该是数组类型")
	}

	if len(dataArray) != 3 {
		t.Errorf("期望data数组长度为3，实际为%d", len(dataArray))
	}
}

// TestDefaultWithComplexData 测试复杂数据结构
func TestDefaultWithComplexData(t *testing.T) {
	d := &Default{}
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"id":   123,
			"name": "张三",
			"tags": []string{"admin", "user"},
		},
		"count": 100,
	}
	d.Set(200, "success", data, 1234567890)

	jsonBytes, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	// 验证可以成功序列化复杂结构
	if len(jsonBytes) == 0 {
		t.Error("序列化结果不应该为空")
	}
}
