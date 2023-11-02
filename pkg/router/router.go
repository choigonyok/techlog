package router

import (
	"net/http"

	"github.com/choigonyok/techlog/pkg/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine      *gin.Engine
	middlewares *middleware.Middleware
}

func NewRouter(m *middleware.Middleware) *Router {
	engine := gin.Default()
	engine.Use(m.Get()...)
	return &Router{
		engine:      engine,
		middlewares: m,
	}
}

func (r *Router) GetHTTPHandler() http.Handler {
	return r.engine.Handler()
}

func (r *Router) SetRoutes(routes []Route) {
	for _, v := range routes {
		switch v.Method {
		case "post":
			r.engine.POST(v.Path, v.Handler)
		case "get":
			r.engine.GET(v.Path, v.Handler)
		case "delete":
			r.engine.DELETE(v.Path, v.Handler)
		case "put":
			r.engine.PUT(v.Path, v.Handler)
		}
	}
}
