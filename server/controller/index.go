package controller

import (
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gcache"
	"net/url"
	"p00q.cn/video_cdn/server/global"
	m3u8Server "p00q.cn/video_cdn/server/service/m3u8"
	"strings"
)

var (
	users = gmap.New(true)       // 使用默认的并发安全Map
	names = gset.NewStrSet(true) // 使用并发安全的Set，用以用户昵称唯一性校验
	cache = gcache.New()         // 使用特定的缓存对象，不使用全局缓存对象
)

func Index(r *ghttp.Request) {
	r.Response.WriteTpl("index.html")
}
func Parse(r *ghttp.Request) {
	r.Response.WriteTpl("parse.html")
}

func GetNewUrl(r *ghttp.Request) {
	value := r.GetQueryString("url", "")
	if value != "" {
		parse, err2 := url.Parse(value)
		if err2 != nil {
			_ = r.Response.WriteJson(g.Map{"err": err2.Error(), "code": 1})
		}
		parse.Path = strings.ReplaceAll(parse.Path, "//", "/")
		global.Logs.Info(parse.String())
		transit, err := m3u8Server.CacheM3u8(parse.String())
		if err != nil {
			_ = r.Response.WriteJson(g.Map{"err": err.Error(), "code": 1})
			global.Logs.Error(err)
			return
		}
		_ = r.Response.WriteJson(g.Map{"err": "", "code": 0, "url": transit})
	}
}
