package task

import (
	"p00q.cn/video_cdn/node/service"
	"time"
)

func StartRunTimerTask() {
	pingT := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-pingT.C:
			service.Ping()
		}
	}
}
