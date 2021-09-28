package controller

import (
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gcache"
	"p00q.cn/video_cdn/server/global"
	nodeserver "p00q.cn/video_cdn/server/service/node"
)

var Index = &indexApi{}

type indexApi struct{}

var (
	users = gmap.New(true)       // 使用默认的并发安全Map
	names = gset.NewStrSet(true) // 使用并发安全的Set，用以用户昵称唯一性校验
	cache = gcache.New()         // 使用特定的缓存对象，不使用全局缓存对象
)

func (a *indexApi) Index(r *ghttp.Request) {
	r.Response.WriteTpl("index.html")
}
func GetNewUrl(r *ghttp.Request) {
	value := r.GetQueryString("url", "")
	if value != "" {
		transit, err := nodeserver.NewCache(value, "127.0.0.1")
		if err != nil {
			global.Logs.Error(err)
			return
		}
		r.Response.Write(transit)
	}
}
