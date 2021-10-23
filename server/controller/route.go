package controller

import (
	"github.com/gogf/gf/net/ghttp"
)

func RegRoute(server *ghttp.Server) {
	server.BindHandler("/get_new", GetNewUrl)
	server.BindHandler("/", Index)
	server.BindHandler("/parse", Parse)

	//node
	node := server.Group("/node")
	node.POST("/ping", ping)

	//video
	video := server.Group("/video")
	video.GET("/list", listH)
	video.GET("/info", infoH)
}
