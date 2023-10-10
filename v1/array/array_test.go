package array

import (
	"testing"
)

// TestIn 是否存在
func TestIn(t *testing.T) {
	testArray := []string{
		"adc",
		"mon",
		"测试",
	}
	if In(testArray, "测试abc") {
		panic("must not in")
	}
	if NotIn(testArray, "测试") {
		panic("must in")
	}
}
