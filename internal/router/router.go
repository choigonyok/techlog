package router

import (
	"net/http"

	"github.com/choigonyok/techlog/internal/router/route"
	"github.com/choigonyok/techlog/pkg/middleware"
	"github.com/gin-gonic/gin"
)

type RouterInterface interface {
	setRoutes()
	setMiddlewares() error
	GetHTTPHandler() http.Handler
}

type Router struct {
	engine     *gin.Engine
	middleware middleware.MiddlewareInterface
	route      *route.Routes
}

var (
	endpointPrefix   = "/api/"
	allowMiddlewares = []string{"origin", "method", "header", "credential", "wildcard"}
)

// NewRouter returns middlewares applied new gin engine
func New() (*Router, error) {
	router := &Router{
		engine:     gin.Default(),
		middleware: middleware.New(),
		route:      route.New(endpointPrefix),
	}

	if err := router.setMiddlewares(); err != nil {
		return nil, err
	}
	router.setRoutes()

	return router, nil
}

// GetHTTPHandler converts gin.Engine to http.Handler type
func (r *Router) GetHTTPHandler() http.Handler {
	return r.engine.Handler()
}

func (r *Router) setMiddlewares() error {
	err := r.middleware.AllowCORS(allowMiddlewares)
	middlewares := r.middleware.Get()
	r.engine.Use(middlewares)
	return err
}

// SetRoutes set gin handler with specific methods and paths
func (r *Router) setRoutes() {
	for _, v := range *r.route {
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
