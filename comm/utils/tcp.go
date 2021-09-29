package utils

import (
	"net"
	"time"
)

//CheckPort 检查端口是否开放
func CheckPort(ip string, port string) bool {
	conn, err := net.DialTimeout("tcp", ip+":"+port, time.Second)
	if err == nil {
		return true
	}
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	return false
}
