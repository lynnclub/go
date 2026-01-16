package json

import (
	"reflect"
	"testing"
)

type test struct {
	ABC string `json:"adc"`
	No  int
}

type nestedTest struct {
	Name  string `json:"name"`
	Inner test   `json:"inner"`
}

type sliceTest struct {
	Items []string `json:"items"`
	Nums  []int    `json:"nums"`
}

var jsonStr = "{\"adc\":\"123\",\"No\":1234}"
var jsonStrReverse = "{\"No\":1234,\"adc\":\"123\"}"

// TestEncode Json编码
func TestEncode(t *testing.T) {
	result := Encode(map[string]interface{}{"adc": "123", "No": 1234})
	if result != jsonStrReverse {
		panic("encoding json.Encode map error")
	}

	result = Encode(test{ABC: "123", No: 1234})
	if result != jsonStr {
		panic("encoding json.Encode struct error")
	}
}

// TestEncode Json编码
func TestEncodeToByte(t *testing.T) {
	result := EncodeToByte(map[string]interface{}{"adc": "123", "No": 1234})
	if string(result) != jsonStrReverse {
		panic("encoding json.Encode map error")
	}

	result = EncodeToByte(test{ABC: "123", No: 1234})
	if string(result) != jsonStr {
		panic("encoding json.Encode struct error")
	}
}

// TestDecode Json解码
func TestDecode(t *testing.T) {
	var data interface{}
	err := Decode(jsonStr, &data)
	if err != nil {
		panic("encoding json.Decode map " + err.Error())
	}

	var testOne test
	err = Decode(jsonStr, &testOne)
	if err != nil {
		panic("encoding json.Decode struct " + err.Error())
	}
}

// TestDecodeFromByte 测试从字节数组解码
func TestDecodeFromByte(t *testing.T) {
	jsonBytes := []byte(jsonStr)

	var data interface{}
	err := DecodeFromByte(jsonBytes, &data)
	if err != nil {
		t.Fatalf("DecodeFromByte failed for interface{}: %v", err)
	}

	var testOne test
	err = DecodeFromByte(jsonBytes, &testOne)
	if err != nil {
		t.Fatalf("DecodeFromByte failed for struct: %v", err)
	}

	if testOne.ABC != "123" {
		t.Errorf("Expected ABC to be '123', got '%s'", testOne.ABC)
	}
	if testOne.No != 1234 {
		t.Errorf("Expected No to be 1234, got %d", testOne.No)
	}
}

// TestEncodeNil 测试编码nil值
func TestEncodeNil(t *testing.T) {
	result := Encode(nil)
	expected := "null"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

// TestEncodeSlice 测试编码切片
func TestEncodeSlice(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	result := Encode(slice)
	expected := "[1,2,3,4,5]"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}

	stringSlice := []string{"hello", "world"}
	result = Encode(stringSlice)
	expected = "[\"hello\",\"world\"]"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

// TestEncodeArray 测试编码数组
func TestEncodeArray(t *testing.T) {
	arr := [3]int{1, 2, 3}
	result := Encode(arr)
	expected := "[1,2,3]"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

// TestEncodeNestedStruct 测试编码嵌套结构体
func TestEncodeNestedStruct(t *testing.T) {
	nested := nestedTest{
		Name: "outer",
		Inner: test{
			ABC: "inner_value",
			No:  999,
		},
	}
	result := Encode(nested)
	expected := "{\"name\":\"outer\",\"inner\":{\"adc\":\"inner_value\",\"No\":999}}"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

// TestEncodeEmptyStruct 测试编码空结构体
func TestEncodeEmptyStruct(t *testing.T) {
	empty := struct{}{}
	result := Encode(empty)
	expected := "{}"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

// TestEncodePointer 测试编码指针
func TestEncodePointer(t *testing.T) {
	value := test{ABC: "pointer_test", No: 777}
	ptr := &value
	result := Encode(ptr)
	expected := "{\"adc\":\"pointer_test\",\"No\":777}"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

// TestEncodeBoolean 测试编码布尔值
func TestEncodeBoolean(t *testing.T) {
	result := Encode(true)
	if result != "true" {
		t.Errorf("Expected 'true', got %s", result)
	}

	result = Encode(false)
	if result != "false" {
		t.Errorf("Expected 'false', got %s", result)
	}
}

// TestEncodeNumber 测试编码数字类型
func TestEncodeNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"int", 123, "123"},
		{"int64", int64(9223372036854775807), "9223372036854775807"},
		{"float32", float32(3.14), "3.14"},
		{"float64", float64(2.718281828), "2.718281828"},
		{"negative", -456, "-456"},
		{"zero", 0, "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Encode(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// TestEncodeString 测试编码字符串
func TestEncodeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple", "hello", "\"hello\""},
		{"empty", "", "\"\""},
		{"with spaces", "hello world", "\"hello world\""},
		{"with quotes", "say \"hello\"", "\"say \\\"hello\\\"\""},
		{"with newline", "line1\nline2", "\"line1\\nline2\""},
		{"with tab", "col1\tcol2", "\"col1\\tcol2\""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Encode(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// TestDecodeInvalidJSON 测试解码无效的JSON
func TestDecodeInvalidJSON(t *testing.T) {
	tests := []struct {
		name       string
		jsonString string
	}{
		{"invalid format", "{invalid json}"},
		{"unclosed brace", "{\"key\": \"value\""},
		{"trailing comma", "{\"key\": \"value\",}"},
		{"single quotes", "{'key': 'value'}"},
		{"unquoted key", "{key: \"value\"}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result interface{}
			err := Decode(tt.jsonString, &result)
			if err == nil {
				t.Errorf("Expected error for invalid JSON: %s", tt.jsonString)
			}
		})
	}
}

// TestDecodeFromByteInvalidJSON 测试从字节数组解码无效JSON
func TestDecodeFromByteInvalidJSON(t *testing.T) {
	invalidJSON := []byte("{invalid json}")
	var result interface{}
	err := DecodeFromByte(invalidJSON, &result)
	if err == nil {
		t.Error("Expected error for invalid JSON bytes")
	}
}

// TestDecodeToMap 测试解码到map
func TestDecodeToMap(t *testing.T) {
	jsonString := "{\"name\":\"John\",\"age\":30,\"city\":\"New York\"}"
	var result map[string]interface{}
	err := Decode(jsonString, &result)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if result["name"] != "John" {
		t.Errorf("Expected name to be 'John', got %v", result["name"])
	}
	// JSON numbers are decoded as float64
	if result["age"] != float64(30) {
		t.Errorf("Expected age to be 30, got %v", result["age"])
	}
}

// TestDecodeArray 测试解码数组
func TestDecodeArray(t *testing.T) {
	jsonString := "[1,2,3,4,5]"
	var result []int
	err := Decode(jsonString, &result)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	expected := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestDecodeNestedStruct 测试解码嵌套结构体
func TestDecodeNestedStruct(t *testing.T) {
	jsonString := "{\"name\":\"test\",\"inner\":{\"adc\":\"value\",\"No\":123}}"
	var result nestedTest
	err := Decode(jsonString, &result)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if result.Name != "test" {
		t.Errorf("Expected name to be 'test', got '%s'", result.Name)
	}
	if result.Inner.ABC != "value" {
		t.Errorf("Expected inner ABC to be 'value', got '%s'", result.Inner.ABC)
	}
	if result.Inner.No != 123 {
		t.Errorf("Expected inner No to be 123, got %d", result.Inner.No)
	}
}

// TestDecodeSliceStruct 测试解码包含切片的结构体
func TestDecodeSliceStruct(t *testing.T) {
	jsonString := "{\"items\":[\"a\",\"b\",\"c\"],\"nums\":[1,2,3]}"
	var result sliceTest
	err := Decode(jsonString, &result)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	expectedItems := []string{"a", "b", "c"}
	if !reflect.DeepEqual(result.Items, expectedItems) {
		t.Errorf("Expected items %v, got %v", expectedItems, result.Items)
	}

	expectedNums := []int{1, 2, 3}
	if !reflect.DeepEqual(result.Nums, expectedNums) {
		t.Errorf("Expected nums %v, got %v", expectedNums, result.Nums)
	}
}

// TestEncodeDecodeRoundTrip 测试编码解码往返
func TestEncodeDecodeRoundTrip(t *testing.T) {
	original := test{ABC: "test_value", No: 12345}

	encoded := Encode(original)
	var decoded test
	err := Decode(encoded, &decoded)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if !reflect.DeepEqual(original, decoded) {
		t.Errorf("Round trip failed. Original: %+v, Decoded: %+v", original, decoded)
	}
}

// TestEncodeToByteDecodeFromByteRoundTrip 测试字节编码解码往返
func TestEncodeToByteDecodeFromByteRoundTrip(t *testing.T) {
	original := nestedTest{
		Name: "round_trip",
		Inner: test{
			ABC: "inner_test",
			No:  999,
		},
	}

	encoded := EncodeToByte(original)
	var decoded nestedTest
	err := DecodeFromByte(encoded, &decoded)
	if err != nil {
		t.Fatalf("DecodeFromByte failed: %v", err)
	}

	if !reflect.DeepEqual(original, decoded) {
		t.Errorf("Round trip failed. Original: %+v, Decoded: %+v", original, decoded)
	}
}

// TestEncodeEmptySlice 测试编码空切片
func TestEncodeEmptySlice(t *testing.T) {
	emptySlice := []string{}
	result := Encode(emptySlice)
	expected := "[]"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

// TestEncodeEmptyMap 测试编码空map
func TestEncodeEmptyMap(t *testing.T) {
	emptyMap := map[string]interface{}{}
	result := Encode(emptyMap)
	expected := "{}"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

// TestDecodeNull 测试解码null值
func TestDecodeNull(t *testing.T) {
	var result interface{}
	err := Decode("null", &result)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}
}

// TestDecodeBoolean 测试解码布尔值
func TestDecodeBoolean(t *testing.T) {
	var result bool
	err := Decode("true", &result)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if result != true {
		t.Errorf("Expected true, got %v", result)
	}

	err = Decode("false", &result)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if result != false {
		t.Errorf("Expected false, got %v", result)
	}
}

// TestEncodeSpecialCharacters 测试编码特殊字符
func TestEncodeSpecialCharacters(t *testing.T) {
	special := struct {
		Unicode string `json:"unicode"`
		Emoji   string `json:"emoji"`
	}{
		Unicode: "你好世界",
		Emoji:   "😀🎉",
	}

	result := Encode(special)
	var decoded struct {
		Unicode string `json:"unicode"`
		Emoji   string `json:"emoji"`
	}
	err := Decode(result, &decoded)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if decoded.Unicode != special.Unicode {
		t.Errorf("Unicode mismatch: expected %s, got %s", special.Unicode, decoded.Unicode)
	}
	if decoded.Emoji != special.Emoji {
		t.Errorf("Emoji mismatch: expected %s, got %s", special.Emoji, decoded.Emoji)
	}
}
