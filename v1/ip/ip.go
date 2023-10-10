package ip

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
)

// Local 获取本地IP
func Local(ipv4 bool) []string {
	var ips []string

	address, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("net.InterfaceAddrs error:", err.Error())
		return ips
	}

	for _, addr := range address {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipv4 {
				// 获取IPv4
				if ipNet.IP.To4() != nil {
					ips = append(ips, ipNet.IP.String())
				}
			} else {
				// 获取IPv6
				if ipNet.IP.To16() != nil {
					ips = append(ips, ipNet.IP.String())
				}
			}
		}
	}

	return ips
}

// GetClientIP 获取Header client-ip 的内容
func GetClientIP(c *gin.Context) (string, error) {
	ip := c.Request.Header.Get("client-ip")
	netIp := net.ParseIP(ip)
	if netIp != nil {
		return ip, nil
	}

	return c.ClientIP(), nil
}
