package main

import (
	logs "github.com/danbai225/go-logs"
	"github.com/gogf/gf/frame/g"
	"videoCDN/controller"
	"videoCDN/middleware"
)

func main() {
	logs.Info("启动服务:http://127.0.0.1:8080")
	server := g.Server()
	server.SetPort(8080)
	controller.RegRoute(server)
	middleware.RegMiddleware(server)
	server.Run()
}
