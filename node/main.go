package main

import (
	"flag"
	"fmt"
	logs "github.com/danbai225/go-logs"
	"github.com/gogf/gf/frame/g"
	"log"
	"p00q.cn/video_cdn/node/service"
	"p00q.cn/video_cdn/node/task"
	"time"

	"p00q.cn/video_cdn/node/config"
	"p00q.cn/video_cdn/node/controller"
	"p00q.cn/video_cdn/node/middleware"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	//启动与服务端通信
	Path := flag.String("c", "", "指定配置文件路径，默认为./config.json")
	flag.Parse()
	if *Path != "" {
		config.Path = *Path
	}
	err := config.LoadConfig()
	if err != nil {
		logs.Err(err)
		return
	}
	go task.StartRunTimerTask()
	go func() {
		for {
			service.Run()
			time.Sleep(time.Second * 5)
			logs.Info("连接断开,尝试重新连接...")
		}
	}()

	//启动http服务
	logs.Info(fmt.Sprintf("启动服务 :%d", config.GlobalConfig.Port))
	server := g.Server()
	server.SetPort(config.GlobalConfig.Port)
	controller.RegRoute(server)
	middleware.RegMiddleware(server)
	if config.GlobalConfig.CertFile != "" {
		server.EnableHTTPS(config.GlobalConfig.CertFile, config.GlobalConfig.KeyFile)
	}
	server.SetReadTimeout(time.Second * 30)
	server.SetWriteTimeout(time.Second * 30)
	server.Run()
}
