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

func New(m *middleware.Middleware) *Router {
	engine := gin.Default()
	engine.Use(m.Get()...)

	return &Router{
		engine:      engine,
		middlewares: m,
	}
}

func (r *Router) GetHandler() http.Handler {
	return r.engine.Handler()
}
