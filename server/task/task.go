package task

import (
	"p00q.cn/video_cdn/comm/model"
	"p00q.cn/video_cdn/server/global"
	nodeServer "p00q.cn/video_cdn/server/service/node"
	"time"
)

func StartRunTimerTask() {
	pingT := time.NewTicker(time.Minute * 16)
	offlineTimeout := time.NewTicker(time.Minute)
	clearInvalidCacheTimeout := time.NewTicker(time.Hour)
	for {
		select {
		case <-pingT.C:
			delayTest()
		case <-offlineTimeout.C:
			offlineTimeoutF()
		case <-clearInvalidCacheTimeout.C:
			clearInvalidCache()
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

//超时心跳节点
func offlineTimeoutF() {
	global.MySQL.Exec(`UPDATE nodes SET on_line=0 WHERE  on_line=1 AND (updated_at<date_add(now(), interval -1 minute) OR updated_at=NULL)`)
}

//清除失效缓存
func clearInvalidCache() {
	caches := make([]model.Cache, 0)
	global.MySQL.Model(&model.Cache{}).Where("valid=0 and updated_at<date_add(now(), interval -7 DAY)").Find(&caches)
	keys := make([]string, 0)
	for _, cache := range caches {
		keys = append(keys, cache.VideoKey)
	}
	global.MySQL.Where("video_key IN ?", keys).Delete(&model.Data{})
	global.MySQL.Where("video_key IN ?", keys).Delete(&model.Cache{})
}
