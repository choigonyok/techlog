package response

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Response200(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
}

func ResponseDataWith200(c *gin.Context, data interface{}) error {
	result, err := json.Marshal(data)
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(result)
	return err
}

func Response500(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusInternalServerError)
}

func Response401(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusUnauthorized)
}
