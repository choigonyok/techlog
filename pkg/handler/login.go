package handler

import (
	"os"

	"github.com/choigonyok/techlog/pkg/database"
	"github.com/choigonyok/techlog/pkg/model"
	resp "github.com/choigonyok/techlog/pkg/response"
	"github.com/choigonyok/techlog/pkg/service"
	"github.com/gin-gonic/gin"
)

// VerifyAdminIDAndPW chekcs the input id/password of client is correct
func VerifyAdminIDAndPW(c *gin.Context) {
	user := model.User{}
	cookie := &AdminCookie{}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		resp.Response500(c, err)
		return
	}

	if user.ID != os.Getenv("BLOG_ID") || user.Password != os.Getenv("BLOG_PW") {
		resp.Response500(c, err)
		return
	}
	cookie.setCookie(c, false, true)
	resp.Response200(c)
}

// VerifyAdminUser checks the client has already logged in
func VerifyAdminUser(c *gin.Context) {
	pvrSlave := database.NewMysqlProvider(database.GetReadConnector())
	svcSlave := service.NewService(pvrSlave)

	adminCookieValue, err := c.Cookie(adminCookieKey)
	if err != nil {
		resp.Response401(c)
		return
	}
	value, err := svcSlave.GetCookieValue()
	if err != nil {
		resp.Response401(c)
		return
	}
	if value != adminCookieValue {
		resp.Response401(c)
		return
	} else {
		resp.Response200(c)
		return
	}
}
