package main

import (
	"github.com/gogf/gf/frame/g"
	"p00q.cn/video_cdn/comm/model"
	"p00q.cn/video_cdn/server/controller"
	"p00q.cn/video_cdn/server/global"
	"p00q.cn/video_cdn/server/middleware"
	"p00q.cn/video_cdn/server/service/node"
)

func main() {
	global.InitDB()
	test()
	go node.Run()
	server := g.Server()
	controller.RegRoute(server)
	middleware.RegMiddleware(server)
	server.Run()
}
func test() {
	nodes := make([]model.Node, 0)
	global.MySQL.Debug().Model(&model.Node{}).Find(&nodes)
	//global.DB().Model(&model.Node{}).Insert(&model.Node{
	//	IP:        "12.2.3.4",
	//	Area:      "test",
	//	Bandwidth: 10,
	//})
	//node := model.Node{}
}
