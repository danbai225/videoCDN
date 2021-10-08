package controller

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	nodeServer "p00q.cn/video_cdn/server/service/node"
)

func ping(r *ghttp.Request) {

	host := r.GetFormString("host")
	if host == "" {
		_ = r.Response.WriteJson(g.Map{"err": "没有host参数", "code": 1})
		return
	}
	nodeServer.DelayTest(host)
	_ = r.Response.WriteJson(g.Map{"err": "", "code": 0})
}
