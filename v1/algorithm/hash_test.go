package algorithm

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"testing"
)

func MD5v2(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func MD5v3(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))

	return hex.EncodeToString(hash.Sum(nil))
}

func TestMD5(t *testing.T) {
	result := MD5("123")
	if result != "202cb962ac59075b964b07152d234b70" {
		panic("md5 incorrect")
	}

	resultV2 := MD5v2("123")
	if result != resultV2 {
		panic("md5v2 inconsistent")
	}

	resultV3 := MD5v3("123")
	if result != resultV3 {
		panic("md5v3 inconsistent")
	}
}

func BenchmarkMD5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MD5("123")
	}
}

func BenchmarkMD5V2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MD5v2("123")
	}
}

func BenchmarkMD5V3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MD5v3("123")
	}
}

func TestCRC32Hash(t *testing.T) {
	expectedHashValue := uint32(3964322768)
	hashValue := Crc32("Hello, World!")

	if hashValue != expectedHashValue {
		t.Errorf("Expected CRC32 hash: %d, got: %d", expectedHashValue, hashValue)
	}
}

// TestSHA1 测试SHA1哈希
func TestSHA1(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"", "da39a3ee5e6b4b0d3255bfef95601890afd80709"},
		{"123", "40bd001563085fc35165329ea1ff5c5ecbdbbeef"},
		{"hello", "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"},
		{"Hello, World!", "0a0a9f2a6772942557ab5355d76af442f8f65e01"},
		{"你好世界", "dabaa5fe7c47fb21be902480a13013f16a1ab6eb"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := SHA1(tc.input)
			if result != tc.expected {
				t.Errorf("SHA1(%q) = %s, 期望 %s", tc.input, result, tc.expected)
			}
		})
	}
}

// TestHmacSHA256 测试HmacSHA256
func TestHmacSHA256(t *testing.T) {
	testCases := []struct {
		name     string
		message  string
		secret   string
		expected string
	}{
		{
			"简单测试",
			"hello",
			"secret",
			"88aab3ede8d3adf94d26ab90d3bafd4a2083070c3bcce9c014ee04a443847c0b",
		},
		{
			"空消息",
			"",
			"secret",
			"f9e66e179b6747ae54108f82f8ade8b3c25d76fd30afde6c395822c530196169",
		},
		{
			"空密钥",
			"message",
			"",
			"eb08c1f56d5ddee07f7bdf80468083da06b64cf4fac64fe3a90883df5feacae4",
		},
		{
			"中文内容",
			"你好世界",
			"密钥",
			"bb51c7509de0e81fa5f963b0172d9276ddc33df1b5841c85cc7f868c8164fad0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := HmacSHA256(tc.message, tc.secret)
			if result != tc.expected {
				t.Errorf("HmacSHA256(%q, %q) = %s, 期望 %s", tc.message, tc.secret, result, tc.expected)
			}
		})
	}
}

// TestHmacSHA1 测试HmacSHA1
func TestHmacSHA1(t *testing.T) {
	testCases := []struct {
		name     string
		message  string
		secret   string
		expected string
	}{
		{
			"简单测试",
			"hello",
			"secret",
			"5112055c05f944f85755efc5cd8970e194e9f45b",
		},
		{
			"空消息",
			"",
			"secret",
			"25af6174a0fcecc4d346680a72b7ce644b9a88e8",
		},
		{
			"空密钥",
			"message",
			"",
			"d5d1ed05121417247616cfc8378f360a39da7cfa",
		},
		{
			"长消息",
			"The quick brown fox jumps over the lazy dog",
			"key",
			"de7c9b85b8b78aa6bc8a7a36f70a90701c9db4d9",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := HmacSHA1(tc.message, tc.secret)
			if result != tc.expected {
				t.Errorf("HmacSHA1(%q, %q) = %s, 期望 %s", tc.message, tc.secret, result, tc.expected)
			}
		})
	}
}

// TestMD5Extended 测试MD5扩展用例
func TestMD5Extended(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"", "d41d8cd98f00b204e9800998ecf8427e"},
		{"a", "0cc175b9c0f1b6a831c399e269772661"},
		{"abc", "900150983cd24fb0d6963f7d28e17f72"},
		{"message digest", "f96b697d7cb7938d525a2f31aaf161d0"},
		{"abcdefghijklmnopqrstuvwxyz", "c3fcd3d76192e4007dfb496cca67e13b"},
		{"你好", "7eca689f0d3389d9dea66ae112e5cfd7"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := MD5(tc.input)
			if result != tc.expected {
				t.Errorf("MD5(%q) = %s, 期望 %s", tc.input, result, tc.expected)
			}
		})
	}
}

// TestCrc32Extended 测试CRC32扩展用例
func TestCrc32Extended(t *testing.T) {
	testCases := []struct {
		input    string
		expected uint32
	}{
		{"", 0},
		{"a", 3904355907},
		{"abc", 891568578},
		{"123", 2286445522},
		{"hello world", 222957957},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := Crc32(tc.input)
			if result != tc.expected {
				t.Errorf("Crc32(%q) = %d, 期望 %d", tc.input, result, tc.expected)
			}
		})
	}
}

// TestHashConsistency 测试哈希一致性
func TestHashConsistency(t *testing.T) {
	input := "test string for consistency"

	// 多次计算同一输入，结果应该相同
	md5_1 := MD5(input)
	md5_2 := MD5(input)
	if md5_1 != md5_2 {
		t.Error("MD5哈希结果不一致")
	}

	sha1_1 := SHA1(input)
	sha1_2 := SHA1(input)
	if sha1_1 != sha1_2 {
		t.Error("SHA1哈希结果不一致")
	}

	hmac256_1 := HmacSHA256(input, "key")
	hmac256_2 := HmacSHA256(input, "key")
	if hmac256_1 != hmac256_2 {
		t.Error("HmacSHA256哈希结果不一致")
	}

	hmac1_1 := HmacSHA1(input, "key")
	hmac1_2 := HmacSHA1(input, "key")
	if hmac1_1 != hmac1_2 {
		t.Error("HmacSHA1哈希结果不一致")
	}
}

// TestHashUniqueness 测试哈希唯一性
func TestHashUniqueness(t *testing.T) {
	inputs := []string{"a", "b", "aa", "ab", "ba"}

	// MD5唯一性
	md5Results := make(map[string]bool)
	for _, input := range inputs {
		result := MD5(input)
		if md5Results[result] {
			t.Errorf("MD5冲突: 输入 %q 产生了重复的哈希", input)
		}
		md5Results[result] = true
	}

	// SHA1唯一性
	sha1Results := make(map[string]bool)
	for _, input := range inputs {
		result := SHA1(input)
		if sha1Results[result] {
			t.Errorf("SHA1冲突: 输入 %q 产生了重复的哈希", input)
		}
		sha1Results[result] = true
	}
}
