package router

import (
	"net/http"

	"github.com/choigonyok/techlog/pkg/handler"
	"github.com/choigonyok/techlog/pkg/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine      *gin.Engine
	middlewares *middleware.Middleware
}

func New(m *middleware.Middleware) *Router {
	engine := gin.Default()
	engine.Use(m.Get()...)

	return &Router{
		engine:      engine,
		middlewares: m,
	}
}

func (r *Router) GetHTTPHandlers() http.Handler {
	return r.engine.Handler()
}

func (r *Router) SetHandlers(h []handler.Handler) {
	for _, v := range h {
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
