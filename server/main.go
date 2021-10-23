package main

import (
	"github.com/gogf/gf/frame/g"
	"log"
	"os"
	"os/signal"
	"p00q.cn/video_cdn/server/controller"
	"p00q.cn/video_cdn/server/global"
	"p00q.cn/video_cdn/server/middleware"
	"p00q.cn/video_cdn/server/service/node"
	"p00q.cn/video_cdn/server/task"
	"syscall"
)

func main() {
	c := make(chan os.Signal)
	// 监听信号
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL:
				global.Logs.Info("退出:", s)
				ExitFunc()
			}
		}
	}()
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	global.InitDB()
	go node.Run()
	go task.StartRunTimerTask()
	server := g.Server()
	controller.RegRoute(server)
	middleware.RegMiddleware(server)
	server.EnableHTTPS("./https/server.crt", "./https/server.key")
	server.SetHTTPSPort(443)
	server.SetPort(80)
	server.SetPort(8082)
	server.Run()
}
func ExitFunc() {
	global.Logs.Info("开始退出...")
	global.Logs.Info("执行清理...")
	db, _ := global.MySQL.DB()
	db.Close()
	global.Logs.Info("结束退出...")
	os.Exit(0)
}
