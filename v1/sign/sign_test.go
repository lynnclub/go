package sign

import (
	"sort"
	"strings"
	"testing"

	"github.com/lynnclub/go/v1/algorithm"
)

var paramTest = map[string]interface{}{"test": "123"}

// MD5v2 常规md5 get拼接（类型转换要求原始数据是字符类型，否则会报错）
func MD5v2(params map[string]interface{}, secret string) string {
	var keys []string
	for key := range params {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	var list []string
	for _, key := range keys {
		list = append(list, key+"="+params[key].(string))
	}

	return algorithm.MD5(strings.Join(list, "&") + secret)
}

// TestMD5 常规md5
func TestMD5(t *testing.T) {
	result := MD5(paramTest, "123")
	if result != "8e0b84aa9445962fb44aec55d189ffce" {
		panic("sign md5 error")
	}

	resultV2 := MD5v2(paramTest, "123")
	if resultV2 != result {
		panic("sign md5v2 error")
	}
}

// TestFeiShu 飞书
func TestFeiShu(t *testing.T) {
	result, err := FeiShu("123", 1667820457)
	if err != nil {
		panic("sign md5 error:" + err.Error())
	}
	if result != "A0jG/9oVYCU86IAS2umF0ZlrZKzG+J16TG+WiNh2eaw=" {
		panic("sign md5 error")
	}
}

func BenchmarkMD5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MD5(paramTest, "123")
	}
}

func BenchmarkMD5v2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MD5v2(paramTest, "123")
	}
}

// TestSHA1 测试SHA1签名
func TestSHA1(t *testing.T) {
	testCases := []struct {
		name     string
		params   map[string]interface{}
		secret   string
		expected string
	}{
		{
			"简单参数",
			map[string]interface{}{"test": "123"},
			"secret",
			"06f0ac9e893f2fe0082908065f41a7f5dafb417c",
		},
		{
			"多个参数",
			map[string]interface{}{"a": "1", "b": "2", "c": "3"},
			"key",
			"fd9b1d5c82367782d7fddfc30f4cccc035a46a44",
		},
		{
			"空secret",
			map[string]interface{}{"key": "value"},
			"",
			"dc30e0faa89a70dec1033016a2309bbeb9efcbd0",
		},
		{
			"包含数字",
			map[string]interface{}{"id": 123, "name": "test"},
			"abc",
			"28626d2247e0fad609b7657e8c1a34c1c135f2a1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := SHA1(tc.params, tc.secret)
			if result != tc.expected {
				t.Errorf("SHA1(%v, %q) = %s, 期望 %s", tc.params, tc.secret, result, tc.expected)
			}
		})
	}
}

// TestMD5Extended 测试MD5扩展用例
func TestMD5Extended(t *testing.T) {
	testCases := []struct {
		name     string
		params   map[string]interface{}
		secret   string
		expected string
	}{
		{
			"空参数",
			map[string]interface{}{},
			"secret",
			"5ebe2294ecd0e0f08eab7690d2a6ee69", // MD5("secret")
		},
		{
			"多个参数按字母排序",
			map[string]interface{}{"z": "last", "a": "first", "m": "middle"},
			"key",
			"e691162a5a32163a18d9af424f65c5a1",
		},
		{
			"包含特殊字符",
			map[string]interface{}{"key": "value!@#"},
			"secret",
			"4628fa1b665be3de865ba982927013d0",
		},
		{
			"包含中文",
			map[string]interface{}{"name": "测试"},
			"密钥",
			"649908b401fe4bdc14a64f883e65075a",
		},
		{
			"包含bool值",
			map[string]interface{}{"active": true, "enabled": false},
			"key",
			"60cf273bf8143c5be4f37a702d8ae263",
		},
		{
			"包含float64",
			map[string]interface{}{"price": 99.99, "discount": 0.1},
			"key",
			"8e648233a98a3c6b40f62586f02361fa",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := MD5(tc.params, tc.secret)
			if result != tc.expected {
				t.Errorf("MD5(%v, %q) = %s, 期望 %s", tc.params, tc.secret, result, tc.expected)
			}
		})
	}
}

// TestFeiShuExtended 测试飞书签名扩展用例
func TestFeiShuExtended(t *testing.T) {
	testCases := []struct {
		name      string
		secret    string
		timestamp int64
		expected  string
	}{
		{
			"基本测试",
			"test_secret",
			1234567890,
			"3H7JNC7ltBAwibQHFO1KFVN9HTkLtm2virjdsmGcAzw=",
		},
		{
			"空secret",
			"",
			1234567890,
			"E86wjcxHwES7WIjIIyomJSOtYniDW/JoZQHgVaaSoUw=",
		},
		{
			"不同时间戳",
			"secret",
			1000000000,
			"mFumEWnYYs8k3vr9JnQHybV+wRA4JjUGkpZGXGvimEc=",
		},
		{
			"当前时间戳",
			"my_secret",
			1667820457,
			"TAh/qbzJlS/g+EfOUyPyFgVD/HjgCuSc5o3BjiRd81U=",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := FeiShu(tc.secret, tc.timestamp)
			if err != nil {
				t.Errorf("FeiShu返回错误: %v", err)
				return
			}
			if result != tc.expected {
				t.Errorf("FeiShu(%q, %d) = %s, 期望 %s", tc.secret, tc.timestamp, result, tc.expected)
			}
		})
	}
}

// TestMD5Consistency 测试MD5签名一致性
func TestMD5Consistency(t *testing.T) {
	params := map[string]interface{}{
		"user_id": 12345,
		"action":  "login",
		"time":    1667820457,
	}
	secret := "test_secret"

	// 多次计算应该返回相同结果
	result1 := MD5(params, secret)
	result2 := MD5(params, secret)
	result3 := MD5(params, secret)

	if result1 != result2 || result2 != result3 {
		t.Error("MD5签名结果不一致")
	}
}

// TestSHA1Consistency 测试SHA1签名一致性
func TestSHA1Consistency(t *testing.T) {
	params := map[string]interface{}{
		"user_id": 12345,
		"action":  "login",
		"time":    1667820457,
	}
	secret := "test_secret"

	// 多次计算应该返回相同结果
	result1 := SHA1(params, secret)
	result2 := SHA1(params, secret)
	result3 := SHA1(params, secret)

	if result1 != result2 || result2 != result3 {
		t.Error("SHA1签名结果不一致")
	}
}

// TestFeiShuConsistency 测试飞书签名一致性
func TestFeiShuConsistency(t *testing.T) {
	secret := "test_secret"
	timestamp := int64(1667820457)

	// 多次计算应该返回相同结果
	result1, _ := FeiShu(secret, timestamp)
	result2, _ := FeiShu(secret, timestamp)
	result3, _ := FeiShu(secret, timestamp)

	if result1 != result2 || result2 != result3 {
		t.Error("飞书签名结果不一致")
	}
}

// TestMD5WithComplexTypes 测试MD5处理复杂类型
func TestMD5WithComplexTypes(t *testing.T) {
	params := map[string]interface{}{
		"string": "hello",
		"int":    123,
		"float":  45.67,
		"bool":   true,
		"nil":    nil,
	}

	// 应该能够处理各种类型而不panic
	result := MD5(params, "secret")
	if result == "" {
		t.Error("MD5应该返回非空结果")
	}
}

// TestSHA1WithComplexTypes 测试SHA1处理复杂类型
func TestSHA1WithComplexTypes(t *testing.T) {
	params := map[string]interface{}{
		"string": "hello",
		"int":    123,
		"float":  45.67,
		"bool":   true,
		"nil":    nil,
	}

	// 应该能够处理各种类型而不panic
	result := SHA1(params, "secret")
	if result == "" {
		t.Error("SHA1应该返回非空结果")
	}
}

// TestSignatureWithEmptyParams 测试空参数
func TestSignatureWithEmptyParams(t *testing.T) {
	emptyParams := map[string]interface{}{}

	md5Result := MD5(emptyParams, "secret")
	if md5Result == "" {
		t.Error("MD5空参数应该返回非空结果")
	}

	sha1Result := SHA1(emptyParams, "secret")
	if sha1Result == "" {
		t.Error("SHA1空参数应该返回非空结果")
	}
}

// TestFeiShuWithDifferentTimestamps 测试不同时间戳产生不同签名
func TestFeiShuWithDifferentTimestamps(t *testing.T) {
	secret := "test_secret"

	result1, _ := FeiShu(secret, 1000000000)
	result2, _ := FeiShu(secret, 2000000000)

	if result1 == result2 {
		t.Error("不同时间戳应该产生不同的签名")
	}
}

// TestSignatureKeyOrdering 测试参数键排序
func TestSignatureKeyOrdering(t *testing.T) {
	// 无论参数顺序如何，结果应该相同
	params1 := map[string]interface{}{"a": "1", "b": "2", "c": "3"}
	params2 := map[string]interface{}{"c": "3", "a": "1", "b": "2"}
	params3 := map[string]interface{}{"b": "2", "c": "3", "a": "1"}

	md5_1 := MD5(params1, "key")
	md5_2 := MD5(params2, "key")
	md5_3 := MD5(params3, "key")

	if md5_1 != md5_2 || md5_2 != md5_3 {
		t.Error("参数顺序不同不应该影响MD5签名结果")
	}

	sha1_1 := SHA1(params1, "key")
	sha1_2 := SHA1(params2, "key")
	sha1_3 := SHA1(params3, "key")

	if sha1_1 != sha1_2 || sha1_2 != sha1_3 {
		t.Error("参数顺序不同不应该影响SHA1签名结果")
	}
}
