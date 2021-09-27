package controller

import (
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gcache"
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
