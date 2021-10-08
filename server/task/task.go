package task

import (
	"p00q.cn/video_cdn/server/global"
	nodeServer "p00q.cn/video_cdn/server/service/node"
	"time"
)

func StartRunTimerTask() {
	pingT := time.NewTicker(time.Minute * 16)
	offlineTimeout := time.NewTicker(time.Minute)
	for {
		select {
		case <-pingT.C:
			delayTest()
		case <-offlineTimeout.C:
			offlineTimeoutF()
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
func offlineTimeoutF() {
	global.MySQL.Exec(`UPDATE nodes SET on_line=0 WHERE  on_line=1 AND (updated_at<date_add(now(), interval -1 minute) OR updated_at=NULL)`)
}
