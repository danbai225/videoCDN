package controller

import (
	logs "github.com/danbai225/go-logs"
	"github.com/gogf/gf/net/ghttp"
	m3u8Server "p00q.cn/video_cdn/http/service/m3u8"
	"strings"
)

func video(r *ghttp.Request) {
	path := strings.ReplaceAll(r.RequestURI, "//", "/")
	r.Response.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	resources, err := m3u8Server.GetResources(path)
	if err == nil {
		r.Response.Write(resources)
	} else {
		r.Response.Write("err")
	}
}
func GetNewUrl(r *ghttp.Request) {
	value := r.PostFormValue("url")
	if value != "" {
		transit, err := m3u8Server.NewTransit(value)
		if err != nil {
			logs.Err(err)
			return
		}
		r.Response.Write(transit)
	}
}
