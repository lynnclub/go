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

// TestResolve 解析IP
func TestResolve(t *testing.T) {
	response, _ := Resolve("14.1.44.228")
	if response.Code != 0 || response.Data.CountryCode != "NZ" {
		panic("ip Resolve error")
	}
}
