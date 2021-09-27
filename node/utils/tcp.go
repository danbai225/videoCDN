package utils

import (
	logs "github.com/danbai225/go-logs"
	"net"
	"time"
)

//CheckPort 检查端口是否开放
func CheckPort(ip string, port string) bool {
	conn, err := net.DialTimeout("tcp", ip+":"+port, time.Second)
	if err == nil {
		return true
	}
	logs.Info("主机：", ip, " 端口:", port, "无法链接，详情：", err.Error())
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	return false
}
