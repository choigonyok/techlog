package controller

import (
	"fmt"
	"net/http"
	"os"

	"github.com/choigonyok/blog-project-backend/internal/model"
	"github.com/gin-gonic/gin"
)

func ConnectDB(driverName, dbData string) {
	err := model.OpenDB(driverName, dbData)
	if err != nil {
		fmt.Println("ERROR #73 : ", err.Error())
	}
}

func UnConnectDB() {
	err := model.CloseDB()
	if err != nil {
		fmt.Println("ERROR #74 : ", err.Error())
	}
}

func CheckIDAndPW(c *gin.Context){
	data := model.LoginData{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		fmt.Println("ERROR #1 : ", err.Error())
	}
	if data.Id == os.Getenv("BLOG_ID") && data.Password == os.Getenv("BLOG_PW"){
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.SetCookie("admin", "authorized",60*60*12,"/","choigonyok.com",false,true)
			c.String(http.StatusOK, "COOKIE SENDED")
	}
}