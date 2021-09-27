package middleware

import "github.com/gogf/gf/net/ghttp"

// Middleware 中间件管理服务
var Middleware = middlewareService{}

type middlewareService struct{}

// CORS 允许接口跨域请求
func (s *middlewareService) CORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}
