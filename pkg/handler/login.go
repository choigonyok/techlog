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
	cookie.setCookie(c, true, false)
	resp.Response200(c)
}

// VerifyAdminUser checks the client has already logged in
func VerifyAdminUser(c *gin.Context) {
	pvr := database.NewMysqlProvider(database.GetConnector())
	svc := service.NewService(pvr)

	adminCookieValue, err := c.Cookie(adminCookieKey)
	if err != nil {
		resp.Response401(c, err)
		return
	}

	isAdmin, err := svc.VerifyAdminByCookieValue(adminCookieValue)
	if err != nil {
		resp.Response500(c, err)
		return
	} else if !isAdmin {
		resp.Response500(c, err)
		return
	}
}
