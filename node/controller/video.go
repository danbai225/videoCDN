package controller

import (
	logs "github.com/danbai225/go-logs"
	"github.com/gogf/gf/net/ghttp"
	"p00q.cn/video_cdn/node/service"
	"strings"
)

func video(r *ghttp.Request) {
	path := strings.ReplaceAll(r.RequestURI, "//", "/")
	r.Response.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	resources, err := service.GetResources(path)
	if err == nil {
		r.Response.Write(resources)
	} else {
		logs.Info("err get", err)
		r.Response.Write("err")
	}
}
