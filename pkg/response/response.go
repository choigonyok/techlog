package response

import (
	"encoding/json"
	"fmt"
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

func Response500(c *gin.Context, err error) {
	fmt.Println(err.Error())
	c.Writer.WriteHeader(http.StatusInternalServerError)
}

func Response401(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusUnauthorized)
}

func Response400(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusBadRequest)
}
