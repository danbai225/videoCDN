package controller

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"p00q.cn/video_cdn/server/service/middleware"
)

func init() {
	s := g.Server()
	// 分组路由注册方式
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(
			middleware.Middleware.CORS,
		)
		//group.ALL("/chat", Index)
		group.Group("/", func(group *ghttp.RouterGroup) {
			//group.Middleware(service.Middleware.Auth)
			//group.ALL("/", Index)
		})
	})
}
