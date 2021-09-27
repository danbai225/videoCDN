package middleware

import "github.com/gogf/gf/net/ghttp"

func RegMiddleware(server *ghttp.Server) {
	server.Use(middlewareCORS)
}

func middlewareCORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}
