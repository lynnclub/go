package ip

import (
	"net/http/httptest"
	"testing"
)

// TestLocal 获取本地IP (IPv4)
func TestLocal(t *testing.T) {
	ips := Local(false)

	if len(ips) <= 0 {
		t.Error("Local应该返回至少一个IPv4地址")
	}

	// 验证返回的是IP地址格式
	for _, ip := range ips {
		if ip == "" {
			t.Error("IP地址不应该为空")
		}
	}
}

// TestLocalIPv6 获取本地IP (包含IPv6)
func TestLocalIPv6(t *testing.T) {
	ips := Local(true)

	if len(ips) <= 0 {
		t.Error("Local应该返回至少一个IP地址")
	}

	// 验证返回的是IP地址格式
	for _, ip := range ips {
		if ip == "" {
			t.Error("IP地址不应该为空")
		}
	}
}

// TestGetClientsWithXForwardedFor 测试从X-Forwarded-For获取客户端IP
func TestGetClientsWithXForwardedFor(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("X-Forwarded-For", "192.168.1.1, 10.0.0.1")
	req.RemoteAddr = "127.0.0.1:12345"

	ips := GetClients(req)

	if len(ips) < 2 {
		t.Errorf("期望至少2个IP地址，实际得到%d个", len(ips))
	}

	if ips[0] != "192.168.1.1" {
		t.Errorf("期望第一个IP为192.168.1.1，实际为%s", ips[0])
	}

	if ips[1] != "10.0.0.1" {
		t.Errorf("期望第二个IP为10.0.0.1，实际为%s", ips[1])
	}
}

// TestGetClientsWithXRealIP 测试从X-Real-IP获取客户端IP
func TestGetClientsWithXRealIP(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("X-Real-IP", "203.0.113.1")
	req.RemoteAddr = "127.0.0.1:12345"

	ips := GetClients(req)

	if len(ips) < 1 {
		t.Errorf("期望至少1个IP地址，实际得到%d个", len(ips))
	}

	if ips[0] != "203.0.113.1" {
		t.Errorf("期望第一个IP为203.0.113.1，实际为%s", ips[0])
	}
}

// TestGetClientsWithCustomHeaders 测试使用自定义请求头
func TestGetClientsWithCustomHeaders(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("CF-Connecting-IP", "198.51.100.1")
	req.Header.Set("X-Forwarded-For", "192.168.1.1")
	req.RemoteAddr = "127.0.0.1:12345"

	// 只检查自定义头
	ips := GetClients(req, "CF-Connecting-IP")

	if len(ips) < 1 {
		t.Errorf("期望至少1个IP地址，实际得到%d个", len(ips))
	}

	if ips[0] != "198.51.100.1" {
		t.Errorf("期望第一个IP为198.51.100.1，实际为%s", ips[0])
	}
}

// TestGetClientsWithMultipleHeaders 测试多个请求头
func TestGetClientsWithMultipleHeaders(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("X-Forwarded-For", "192.168.1.1, 10.0.0.1")
	req.Header.Set("X-Real-IP", "203.0.113.1")
	req.RemoteAddr = "127.0.0.1:12345"

	ips := GetClients(req)

	// 应该包含所有IP
	if len(ips) < 4 {
		t.Errorf("期望至少4个IP地址，实际得到%d个", len(ips))
	}
}

// TestGetClientsWithRemoteAddrOnly 测试只有RemoteAddr的情况
func TestGetClientsWithRemoteAddrOnly(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.RemoteAddr = "192.0.2.1:54321"

	ips := GetClients(req)

	if len(ips) != 1 {
		t.Errorf("期望1个IP地址，实际得到%d个", len(ips))
	}

	if ips[0] != "192.0.2.1" {
		t.Errorf("期望IP为192.0.2.1，实际为%s", ips[0])
	}
}

// TestGetClientsWithRemoteAddrNoPort 测试RemoteAddr没有端口的情况
func TestGetClientsWithRemoteAddrNoPort(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.RemoteAddr = "192.0.2.1"

	ips := GetClients(req)

	if len(ips) != 1 {
		t.Errorf("期望1个IP地址，实际得到%d个", len(ips))
	}

	if ips[0] != "192.0.2.1" {
		t.Errorf("期望IP为192.0.2.1，实际为%s", ips[0])
	}
}

// TestGetClientsWithEmptyHeaders 测试空请求头
func TestGetClientsWithEmptyHeaders(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.RemoteAddr = "127.0.0.1:12345"

	ips := GetClients(req)

	if len(ips) != 1 {
		t.Errorf("期望1个IP地址，实际得到%d个", len(ips))
	}

	if ips[0] != "127.0.0.1" {
		t.Errorf("期望IP为127.0.0.1，实际为%s", ips[0])
	}
}

// TestGetClientsWithSpacesInHeader 测试请求头中有空格的情况
func TestGetClientsWithSpacesInHeader(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("X-Forwarded-For", " 192.168.1.1 , 10.0.0.1 , 172.16.0.1 ")
	req.RemoteAddr = "127.0.0.1:12345"

	ips := GetClients(req)

	if len(ips) < 3 {
		t.Errorf("期望至少3个IP地址，实际得到%d个", len(ips))
	}

	// 验证空格被正确去除
	if ips[0] != "192.168.1.1" {
		t.Errorf("期望第一个IP为192.168.1.1，实际为%s", ips[0])
	}

	if ips[1] != "10.0.0.1" {
		t.Errorf("期望第二个IP为10.0.0.1，实际为%s", ips[1])
	}

	if ips[2] != "172.16.0.1" {
		t.Errorf("期望第三个IP为172.16.0.1，实际为%s", ips[2])
	}
}

// TestGetClientsWithIPv6 测试IPv6地址
func TestGetClientsWithIPv6(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("X-Forwarded-For", "2001:db8::1, 2001:db8::2")
	req.RemoteAddr = "[2001:db8::3]:12345"

	ips := GetClients(req)

	if len(ips) < 3 {
		t.Errorf("期望至少3个IP地址，实际得到%d个", len(ips))
	}

	if ips[0] != "2001:db8::1" {
		t.Errorf("期望第一个IP为2001:db8::1，实际为%s", ips[0])
	}

	if ips[2] != "2001:db8::3" {
		t.Errorf("期望第三个IP为2001:db8::3，实际为%s", ips[2])
	}
}

// TestGetClientsEmptyHeaderValue 测试空的请求头值
func TestGetClientsEmptyHeaderValue(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("X-Forwarded-For", "")
	req.Header.Set("X-Real-IP", "")
	req.RemoteAddr = "127.0.0.1:12345"

	ips := GetClients(req)

	// 应该只返回RemoteAddr
	if len(ips) != 1 {
		t.Errorf("期望1个IP地址，实际得到%d个", len(ips))
	}

	if ips[0] != "127.0.0.1" {
		t.Errorf("期望IP为127.0.0.1，实际为%s", ips[0])
	}
}

// TestGetClientsWithEmptyParts 测试包含空部分的逗号分隔值
func TestGetClientsWithEmptyParts(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("X-Forwarded-For", "192.168.1.1,,10.0.0.1,  ,172.16.0.1")
	req.RemoteAddr = "127.0.0.1:12345"

	ips := GetClients(req)

	// 应该过滤掉空部分
	if len(ips) < 3 {
		t.Errorf("期望至少3个有效IP地址，实际得到%d个", len(ips))
	}

	// 验证只包含有效IP
	validIPs := []string{"192.168.1.1", "10.0.0.1", "172.16.0.1"}
	for i, expectedIP := range validIPs {
		if ips[i] != expectedIP {
			t.Errorf("期望第%d个IP为%s，实际为%s", i+1, expectedIP, ips[i])
		}
	}
}
