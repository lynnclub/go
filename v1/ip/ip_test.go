package ip

import (
	"fmt"
	"testing"
)

// TestLocal 获取本地IP
func TestLocal(t *testing.T) {
	ips := Local(true)
	fmt.Println(ips)

	if len(ips) <= 0 {
		panic("ip Local error")
	}
}
