package controller

import "github.com/gogf/gf/net/ghttp"

func RegRoute(server *ghttp.Server) {
	server.BindHandler("/", index)
	server.BindHandler("/video/*", video)
	server.BindHandler("/get_new", GetNewUrl)
}
