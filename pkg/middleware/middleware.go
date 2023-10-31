package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	middlewares []gin.HandlerFunc
}

func (m *Middleware) AddTestMiddleware() {
	fmt.Println("HELLO")
}

func (m *Middleware) Get() []gin.HandlerFunc {
	return m.middlewares
}
