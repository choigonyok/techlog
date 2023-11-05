package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	middlewares []gin.HandlerFunc
}

// Get returns middleware pool
func (m *Middleware) Get() []gin.HandlerFunc {
	return m.middlewares
}

// AllowConfig returns settings server allows as gin.Handlefunc() type
func (m *Middleware) AllowConfig(allowAddress, allowMethods, allowHeaders []string, allowCredentials bool) {
	cfg := cors.DefaultConfig()
	cfg.AllowOrigins = allowAddress
	cfg.AllowMethods = allowMethods
	cfg.AllowHeaders = allowHeaders
	cfg.AllowCredentials = allowCredentials
	m.middlewares = append(m.middlewares, cors.New(cfg))
}
