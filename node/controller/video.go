package controller

import (
	logs "github.com/danbai225/go-logs"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"p00q.cn/video_cdn/node/service"
	"strings"
	"sync/atomic"
)

var connIPCountMap = gmap.New(true)

func video(r *ghttp.Request) {
	split := strings.Split(r.RemoteAddr, ":")
	ip := split[0]
	df := int64(1)
	count := connIPCountMap.GetOrSet(ip, &df).(*int64)
	//logs.Info(atomic.LoadInt64(count))
	if atomic.LoadInt64(count) > 5 {
		r.Response.Status = 500
		r.Response.WriteJsonExit(g.Map{"err": "conn max"})
		return
	} else {
		atomic.AddInt64(count, 1)
		defer func() {
			atomic.AddInt64(count, -1)
		}()
	}
	path := strings.ReplaceAll(r.RequestURI, "//", "/")
	resources, err := service.GetResources(path)
	if err == nil {
		if strings.HasSuffix(r.RequestURI, "ts") {
			r.Response.Header().Set("Content-Type", "video/mp2t")
		} else {
			r.Response.Header().Set("Content-Type", "application/x-www-form-urlencoded")
		}
		r.Response.Write(resources)
	} else {
		if err.Error() == "redirect" {
			r.Response.RedirectTo(string(resources))
			return
		}
		r.Response.Status = 500
		logs.Info("err get", err)
		r.Response.Write("err")
	}
}
