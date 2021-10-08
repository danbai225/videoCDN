package task

import (
	"p00q.cn/video_cdn/server/global"
	nodeServer "p00q.cn/video_cdn/server/service/node"
	"time"
)

func StartRunTimerTask() {
	pingT := time.NewTicker(time.Minute * 16)
	for {
		select {
		case <-pingT.C:
			delayTest()
		}
	}
}

//定时任务延迟测试
func delayTest() {
	hosts := make([]string, 0)
	global.MySQL.Raw("SELECT DISTINCT `host` FROM delays").Scan(&hosts)
	for _, host := range hosts {
		nodeServer.DelayTest(host)
	}
}
