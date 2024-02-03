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
