package controller

import "github.com/gogf/gf/net/ghttp"

func RegRoute(server *ghttp.Server) {
	server.BindHandler("/video/*", video)
}
