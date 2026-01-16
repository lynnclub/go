package algorithm

import (
	"testing"
)

// TestBase64Encode 测试Base64编码
func TestBase64Encode(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"空字符串", "", ""},
		{"简单文本", "hello", "aGVsbG8="},
		{"包含空格", "hello world", "aGVsbG8gd29ybGQ="},
		{"中文", "你好世界", "5L2g5aW95LiW55WM"},
		{"数字", "123456", "MTIzNDU2"},
		{"特殊字符", "!@#$%^&*()", "IUAjJCVeJiooKQ=="},
		{"混合内容", "Hello世界123!@#", "SGVsbG/kuJbnlYwxMjMhQCM="},
		{"长文本", "This is a long text for testing base64 encoding functionality", "VGhpcyBpcyBhIGxvbmcgdGV4dCBmb3IgdGVzdGluZyBiYXNlNjQgZW5jb2RpbmcgZnVuY3Rpb25hbGl0eQ=="},
		{"JSON格式", `{"name":"test","value":123}`, "eyJuYW1lIjoidGVzdCIsInZhbHVlIjoxMjN9"},
		{"URL", "https://example.com/path?key=value", "aHR0cHM6Ly9leGFtcGxlLmNvbS9wYXRoP2tleT12YWx1ZQ=="},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Base64Encode(tc.input)
			if result != tc.expected {
				t.Errorf("Base64Encode(%q) = %q, 期望 %q", tc.input, result, tc.expected)
			}
		})
	}
}

// TestBase64Decode 测试Base64解码
func TestBase64Decode(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{"空字符串", "", "", false},
		{"简单文本", "aGVsbG8=", "hello", false},
		{"包含空格", "aGVsbG8gd29ybGQ=", "hello world", false},
		{"中文", "5L2g5aW95LiW55WM", "你好世界", false},
		{"数字", "MTIzNDU2", "123456", false},
		{"特殊字符", "IUAjJCVeJiooKQ==", "!@#$%^&*()", false},
		{"混合内容", "SGVsbG/kuJbnlYwxMjMhQCM=", "Hello世界123!@#", false},
		{"长文本", "VGhpcyBpcyBhIGxvbmcgdGV4dCBmb3IgdGVzdGluZyBiYXNlNjQgZW5jb2RpbmcgZnVuY3Rpb25hbGl0eQ==", "This is a long text for testing base64 encoding functionality", false},
		{"JSON格式", "eyJuYW1lIjoidGVzdCIsInZhbHVlIjoxMjN9", `{"name":"test","value":123}`, false},
		{"URL", "aHR0cHM6Ly9leGFtcGxlLmNvbS9wYXRoP2tleT12YWx1ZQ==", "https://example.com/path?key=value", false},
		{"无效Base64", "invalid base64!!!", "", true},
		{"不完整Base64", "aGVsbG8", "", true}, // 标准Base64需要正确的padding
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Base64Decode(tc.input)

			if tc.hasError {
				if err == nil {
					t.Errorf("Base64Decode(%q) 应该返回错误，但没有", tc.input)
				}
				return
			}

			if err != nil {
				t.Errorf("Base64Decode(%q) 返回错误: %v", tc.input, err)
				return
			}

			if result != tc.expected {
				t.Errorf("Base64Decode(%q) = %q, 期望 %q", tc.input, result, tc.expected)
			}
		})
	}
}

// TestBase64EncodeDecode 测试编码解码往返
func TestBase64EncodeDecode(t *testing.T) {
	testCases := []string{
		"",
		"a",
		"hello world",
		"你好世界",
		"1234567890",
		"!@#$%^&*()",
		"The quick brown fox jumps over the lazy dog",
		"混合content123!@#中文English数字Special字符",
		`{"key":"value","nested":{"array":[1,2,3]}}`,
	}

	for _, original := range testCases {
		t.Run(original, func(t *testing.T) {
			encoded := Base64Encode(original)
			decoded, err := Base64Decode(encoded)

			if err != nil {
				t.Errorf("解码失败: %v", err)
				return
			}

			if decoded != original {
				t.Errorf("往返测试失败: 原始=%q, 编码=%q, 解码=%q", original, encoded, decoded)
			}
		})
	}
}

// TestBase64EncodeBytes 测试编码二进制数据
func TestBase64EncodeBytes(t *testing.T) {
	// 测试各种字节值
	testBytes := []byte{0, 1, 127, 128, 255}
	encoded := Base64Encode(string(testBytes))
	decoded, err := Base64Decode(encoded)

	if err != nil {
		t.Errorf("解码二进制数据失败: %v", err)
		return
	}

	decodedBytes := []byte(decoded)
	if len(decodedBytes) != len(testBytes) {
		t.Errorf("字节长度不匹配: 期望 %d, 实际 %d", len(testBytes), len(decodedBytes))
		return
	}

	for i, b := range testBytes {
		if decodedBytes[i] != b {
			t.Errorf("字节[%d]不匹配: 期望 %d, 实际 %d", i, b, decodedBytes[i])
		}
	}
}

// TestBase64LargeData 测试大数据编码解码
func TestBase64LargeData(t *testing.T) {
	// 创建一个较大的字符串（1MB）
	largeData := make([]byte, 1024*1024)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	original := string(largeData)
	encoded := Base64Encode(original)
	decoded, err := Base64Decode(encoded)

	if err != nil {
		t.Errorf("解码大数据失败: %v", err)
		return
	}

	if decoded != original {
		t.Error("大数据往返测试失败")
	}
}

// TestBase64DecodeInvalidInput 测试各种无效输入
func TestBase64DecodeInvalidInput(t *testing.T) {
	invalidInputs := []string{
		"!!!",
		"====",
		"abc@def",
		"123#456",
		"中文不是Base64",
	}

	for _, input := range invalidInputs {
		t.Run(input, func(t *testing.T) {
			_, err := Base64Decode(input)
			if err == nil {
				t.Errorf("Base64Decode(%q) 应该返回错误，但没有", input)
			}
		})
	}
}
