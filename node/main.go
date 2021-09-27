package main

import (
	"flag"
	logs "github.com/danbai225/go-logs"
	"p00q.cn/video_cdn/node/config"
	"p00q.cn/video_cdn/node/service"
	"p00q.cn/video_cdn/node/task"
	"time"
)

func main() {
	configPath := flag.String("c", "", "指定配置文件路径，默认为./config.json")
	flag.Parse()
	err := config.LoadConfig(*configPath)
	if err != nil {
		logs.Err(err)
		return
	}
	go task.StartRunTimerTask()
	for {
		service.Run()
		time.Sleep(time.Second * 5)
		logs.Info("连接断开,尝试重新连接...")
	}
}
