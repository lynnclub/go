package rand

import (
	"testing"
)

// TestRange 区间随机
func TestRange(t *testing.T) {
	num := Range(6, 8)
	if num < 6 || num > 8 {
		panic("rand range error")
	}
}
