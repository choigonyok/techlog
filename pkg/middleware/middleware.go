package middleware

import (
	"errors"

	"github.com/choigonyok/techlog/pkg/middleware/cors"
	"github.com/gin-gonic/gin"
)

type MiddlewareInterface interface {
	AllowCORS(middlewares []string) error
	Get() gin.HandlerFunc
}

type Middleware struct {
	cors cors.CorsInterface
}

var (
	allowOrigin      = []string{"http://localhost", "http://localhost:3000", "http://frontend", "http://fronend:3000", "http://www.choigonyok.com", "https://www.choigonyok.com", "http://choigonyok.com", "https://choigonyok.com"}
	allowMethods     = []string{"GET", "POST", "DELETE", "PUT"}
	allowHeaders     = []string{"Content-type"}
	allowCredentials = true
	allowWildCard    = true
)

func New() *Middleware {
	middleware := &Middleware{
		cors: cors.NewCORS(),
	}
	return middleware
}

func (m *Middleware) AllowCORS(middlewares []string) error {
	for _, v := range middlewares {
		switch v {
		case "origin":
			m.cors.AllowOrigins(allowOrigin)
		case "method":
			m.cors.AllowMethods(allowMethods)
		case "header":
			m.cors.AllowHeaders(allowHeaders)
		case "credential":
			m.cors.AllowCredentials(allowCredentials)
		case "wildcard":
			m.cors.AllowWildcard(allowWildCard)
		default:
			return errors.New("invalid string input")
		}
	}
	return nil
}

func (m *Middleware) Get() gin.HandlerFunc {
	cors := m.cors.Get()
	return cors
}
