package controller

import (
	"github.com/gogf/gf/net/ghttp"
)

func index(r *ghttp.Request) {
	r.Response.Write("test")
}
