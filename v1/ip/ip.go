package ip

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

// Local 获取本地IP
func Local(ipv6 bool) []string {
	var ips []string

	address, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("net.InterfaceAddrs error:", err.Error())
		return ips
	}

	for _, addr := range address {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			// 获取IPv4
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}

	if ipv6 {
		for _, addr := range address {
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				// 获取IPv6
				if ipNet.IP.To16() != nil {
					ips = append(ips, ipNet.IP.String())
				}
			}
		}
	}

	return ips
}

// GetClientIPs 获取用户IP
func GetClientIPs(r *http.Request, trustedHeaders ...string) []string {
	ips := make([]string, 0)

	if len(trustedHeaders) == 0 {
		trustedHeaders = []string{"X-Forwarded-For", "X-Real-IP"}
	}

	for _, header := range trustedHeaders {
		if value := r.Header.Get(header); value != "" {
			for _, part := range strings.Split(value, ",") {
				if ip := strings.TrimSpace(part); ip != "" {
					ips = append(ips, ip)
				}
			}
		}
	}

	if remoteAddr := r.RemoteAddr; remoteAddr != "" {
		if ip, _, err := net.SplitHostPort(remoteAddr); err == nil {
			ips = append(ips, ip)
		} else {
			ips = append(ips, remoteAddr)
		}
	}

	return ips
}
