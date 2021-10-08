package controller

import (
	"github.com/gogf/gf/net/ghttp"
)

func RegRoute(server *ghttp.Server) {
	server.BindHandler("/get_new", GetNewUrl)
	node := server.Group("/node")
	node.POST("/ping", ping)
}
