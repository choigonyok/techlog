package response

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response200 writes http.StatusOK header
func Response200(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
}

// ResponseDataWith200 writes data to client with http.StatusOK header
func ResponseDataWith200(c *gin.Context, data interface{}) error {
	result, err := json.Marshal(data)
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(result)
	return err
}

// Response500 writes http.StatusInternalServerError header and print error
func Response500(c *gin.Context, err error) {
	fmt.Println(err.Error())
	c.Writer.WriteHeader(http.StatusInternalServerError)
}

// Response401 writes http.StatusUnauthorizaed header
func Response401(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusUnauthorized)
}

// Response400 writes http.StatusBadRequest header
func Response400(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusBadRequest)
}
