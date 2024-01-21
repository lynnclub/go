package ip

import (
	"testing"
)

// TestLocal 获取本地IP
func TestLocal(t *testing.T) {
	ips := Local(true)
	if len(ips) <= 0 {
		panic("ip Local error")
	}
}
