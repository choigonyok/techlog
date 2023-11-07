package cors

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type CorsInterface interface {
	AllowOrigins(allowAddress []string)
	AllowMethods(allowMethods []string)
	AllowHeaders(allowHeaders []string)
	AllowCredentials(allowCredentials bool)
	AllowWildcard(allowWildCard bool)
	Get() gin.HandlerFunc
}

func (c *GinCORS) Get() gin.HandlerFunc {
	return cors.New(c.cors)
}

type GinCORS struct {
	cors cors.Config
}

func NewCORS() CorsInterface {
	return &GinCORS{
		cors: cors.DefaultConfig(),
	}
}

// AllowConfig returns settings server allows as gin.Handlefunc() type
func (c *GinCORS) AllowOrigins(allowAddress []string) {
	c.cors.AllowOrigins = allowAddress
}

func (c *GinCORS) AllowMethods(allowMethods []string) {
	c.cors.AllowMethods = allowMethods
}

func (c *GinCORS) AllowHeaders(allowHeaders []string) {
	c.cors.AllowHeaders = allowHeaders
}

func (c *GinCORS) AllowCredentials(allowCredentials bool) {
	c.cors.AllowCredentials = allowCredentials
}

func (c *GinCORS) AllowWildcard(allowWildCard bool) {
	c.cors.AllowWildcard = allowWildCard
}
