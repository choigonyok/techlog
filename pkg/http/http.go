package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Response200(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusUnauthorized)
}

func Response500(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusInternalServerError)
}

func Response401(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusUnauthorized)
}
