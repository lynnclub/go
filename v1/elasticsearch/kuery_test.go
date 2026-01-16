package elasticsearch

import (
	"strings"
	"testing"
)

// TestGetKuery 测试生成Kuery查询字符串
func TestGetKuery(t *testing.T) {
	result := GetKuery("status", "active")
	expected := `status:"active"`

	if result != expected {
		t.Errorf("期望%s，实际为%s", expected, result)
	}
}

// TestGetKueryWithQuotes 测试包含引号的值
func TestGetKueryWithQuotes(t *testing.T) {
	result := GetKuery("message", `test "quoted" value`)
	expected := `message:"test quoted value"`

	if result != expected {
		t.Errorf("期望%s，实际为%s", expected, result)
	}
}

// TestGetKueryLongValue 测试长值被截断
func TestGetKueryLongValue(t *testing.T) {
	longValue := strings.Repeat("a", 100)
	result := GetKuery("field", longValue)

	// 应该被截断到90个字符以内
	if len(result) > 100 { // "field:" + 引号 + 最多90个字符
		t.Errorf("结果太长: %d个字符", len(result))
	}

	if !strings.HasPrefix(result, `field:"`) {
		t.Errorf("结果应该以field:\"开头")
	}

	if !strings.HasSuffix(result, `"`) {
		t.Errorf("结果应该以\"结尾")
	}
}

// TestGetKueryWithSpace 测试包含空格的值被截断
func TestGetKueryWithSpace(t *testing.T) {
	// 创建一个包含空格的长字符串
	longValue := strings.Repeat("word ", 20) + "end"
	result := GetKuery("field", longValue)

	// 应该在空格处截断
	if len(result) > 100 {
		t.Errorf("结果太长: %d个字符", len(result))
	}
}

// TestGetKueryWithBackslash 测试包含反斜杠的值被截断
func TestGetKueryWithBackslash(t *testing.T) {
	// 创建一个包含反斜杠的长字符串
	longValue := strings.Repeat("a", 85) + `\test`
	result := GetKuery("field", longValue)

	// 应该在反斜杠处截断
	if len(result) > 100 {
		t.Errorf("结果太长: %d个字符", len(result))
	}
}

// TestGetKueryShortValue 测试短值不被截断
func TestGetKueryShortValue(t *testing.T) {
	shortValue := "short"
	result := GetKuery("field", shortValue)
	expected := `field:"short"`

	if result != expected {
		t.Errorf("期望%s，实际为%s", expected, result)
	}
}

// TestGetKueryEmptyValue 测试空值
func TestGetKueryEmptyValue(t *testing.T) {
	result := GetKuery("field", "")
	expected := `field:""`

	if result != expected {
		t.Errorf("期望%s，实际为%s", expected, result)
	}
}

// TestGetKuerySpecialChars 测试特殊字符
func TestGetKuerySpecialChars(t *testing.T) {
	result := GetKuery("field", "value@#$%")
	expected := `field:"value@#$%"`

	if result != expected {
		t.Errorf("期望%s，实际为%s", expected, result)
	}
}

// TestGetKibanaUrl 测试生成Kibana URL
func TestGetKibanaUrl(t *testing.T) {
	url := "http://kibana.example.com"
	index := "logs-*"
	querys := []string{
		`status:"error"`,
		`level:"critical"`,
	}

	result := GetKibanaUrl(url, index, querys)

	// 验证URL包含基础部分
	if !strings.HasPrefix(result, url+"#/?_a=") {
		t.Errorf("URL应该以%s#/?_a=开头", url)
	}

	// 验证包含index
	if !strings.Contains(result, "index") {
		t.Error("URL应该包含index参数")
	}

	// 验证包含query
	if !strings.Contains(result, "query") {
		t.Error("URL应该包含query参数")
	}

	// 验证包含_g参数
	if !strings.Contains(result, "&_g=") {
		t.Error("URL应该包含_g参数")
	}
}

// TestGetKibanaUrlSingleQuery 测试单个查询
func TestGetKibanaUrlSingleQuery(t *testing.T) {
	url := "http://kibana.example.com"
	index := "logs-*"
	querys := []string{`status:"error"`}

	result := GetKibanaUrl(url, index, querys)

	if !strings.HasPrefix(result, url+"#/?_a=") {
		t.Errorf("URL应该以%s#/?_a=开头", url)
	}
}

// TestGetKibanaUrlEmptyQuery 测试空查询
func TestGetKibanaUrlEmptyQuery(t *testing.T) {
	url := "http://kibana.example.com"
	index := "logs-*"
	querys := []string{}

	result := GetKibanaUrl(url, index, querys)

	if !strings.HasPrefix(result, url+"#/?_a=") {
		t.Errorf("URL应该以%s#/?_a=开头", url)
	}
}

// TestCheckIncompleteness 测试检查不完整性
func TestCheckIncompleteness(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		expected int
	}{
		{"包含空格", "hello world test", 11}, // 最后一个空格的位置
		{"包含反斜杠", "path\\to\\file", 7},   // 最后一个反斜杠的位置
		{"包含双冒号", "namespace::class", 9}, // 双冒号的位置
		{"无特殊字符", "helloworld", -1},      // 没有特殊字符
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := checkIncompleteness(tc.text)
			if result != tc.expected {
				t.Errorf("期望%d，实际为%d", tc.expected, result)
			}
		})
	}
}

// TestToUrlRison 测试转换为URL Rison格式
func TestToUrlRison(t *testing.T) {
	m := map[string]interface{}{
		"key1": "value1",
		"key2": map[string]interface{}{
			"nested": "value2",
		},
	}

	result := toUrlRison(m)

	// 验证包含key1
	if !strings.Contains(result, "key1") {
		t.Error("结果应该包含key1")
	}

	// 验证包含key2
	if !strings.Contains(result, "key2") {
		t.Error("结果应该包含key2")
	}

	// 验证包含nested
	if !strings.Contains(result, "nested") {
		t.Error("结果应该包含nested")
	}
}

// TestToUrlRisonWithQuotes 测试包含单引号的值
func TestToUrlRisonWithQuotes(t *testing.T) {
	m := map[string]interface{}{
		"key": "value's test",
	}

	result := toUrlRison(m)

	// 单引号被替换为!' 然后被URL编码，所以应该包含%21（!的编码）
	if !strings.Contains(result, "%21") && !strings.Contains(result, "!") {
		t.Logf("结果: %s", result)
		t.Error("结果应该包含转义的单引号（! 或 %21）")
	}
}

// TestToUrlRisonEmpty 测试空map
func TestToUrlRisonEmpty(t *testing.T) {
	m := map[string]interface{}{}

	result := toUrlRison(m)

	if result != "" {
		t.Errorf("空map应该返回空字符串，实际为%s", result)
	}
}

// TestToUrlRisonNested 测试嵌套map
func TestToUrlRisonNested(t *testing.T) {
	m := map[string]interface{}{
		"level1": map[string]interface{}{
			"level2": map[string]interface{}{
				"level3": "value",
			},
		},
	}

	result := toUrlRison(m)

	// 验证嵌套结构
	if !strings.Contains(result, "level1") {
		t.Error("结果应该包含level1")
	}
	if !strings.Contains(result, "level2") {
		t.Error("结果应该包含level2")
	}
	if !strings.Contains(result, "level3") {
		t.Error("结果应该包含level3")
	}
}
