package handler

import (
	"os"

	"github.com/choigonyok/techlog/pkg/database"
	"github.com/choigonyok/techlog/pkg/model"
	resp "github.com/choigonyok/techlog/pkg/response"
	"github.com/choigonyok/techlog/pkg/service"
	"github.com/gin-gonic/gin"
)

// login
func VerifyAdminIDAndPW(c *gin.Context) {
	user := model.User{}
	cookie := &AdminCookie{}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		resp.Response500(c)
		return
	}

	if user.ID != os.Getenv("BLOG_ID") || user.Password != os.Getenv("BLOG_PW") {
		resp.Response400(c)
		return
	}
	cookie.setCookie(c, true, false)
	resp.Response200(c)
}

// login
func VerifyAdminUser(c *gin.Context) {
	pvr := database.NewMysqlProvider(database.GetConnector())
	svc := service.NewService(pvr)

	adminCookieValue, err := c.Cookie(adminCookieKey)
	if err != nil {
		resp.Response401(c)
		return
	}

	isAdmin, err := svc.VerifyAdminByCookieValue(adminCookieValue)
	if err != nil {
		resp.Response500(c)
		return
	} else if isAdmin {
		resp.Response200(c)
		return
	} else {
		resp.Response401(c)
		return
	}
}
