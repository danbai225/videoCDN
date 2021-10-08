package main

import (
	"github.com/gogf/gf/frame/g"
	"log"
	"p00q.cn/video_cdn/server/controller"
	"p00q.cn/video_cdn/server/global"
	"p00q.cn/video_cdn/server/middleware"
	"p00q.cn/video_cdn/server/service/node"
	"p00q.cn/video_cdn/server/task"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	global.InitDB()
	test()
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
func test() {
	//nodes := make([]model.Node, 0)
	//global.MySQL.Debug().Model(&model.Node{}).Find(&nodes)
	//global.DB().Model(&model.Node{}).Insert(&model.Node{
	//	IP:        "12.2.3.4",
	//	Area:      "test",
	//	Bandwidth: 10,
	//})
	//node := model.Node{}
}
