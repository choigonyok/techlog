package handler

import (
	"net/http"
	"os"
	"strings"

	"github.com/choigonyok/techlog/pkg/data"
	"github.com/choigonyok/techlog/pkg/database"
	resp "github.com/choigonyok/techlog/pkg/response"
	"github.com/choigonyok/techlog/pkg/service"
	"github.com/gin-gonic/gin"
)

var (
	visitTimeCookieKey     = data.EncodeBase64("visitTime")
	adminCookieKey         = data.EncodeBase64("admin")
	defaultCookieAliveTime = 60 * 60 * 12 // second
	defaultCookeiPath      = "/"
	defaultCookieDomain    = os.Getenv("HOST")
)

type Cookie interface {
	setCookie()
	verifyCookieValue()
}
type AdminCookie struct{}
type VisitTimeCookie struct{}

// AdminCookie.setCookie stores cookie to database, and sets cookie to client
func (ck *AdminCookie) setCookie(c *gin.Context, secure, httpOnly bool) {
	pvr := database.NewMysqlProvider(database.GetConnector())
	svc := service.NewService(pvr)

	uniqueID := data.CreateRandomString()
	_, err := svc.GetCookieValue()
	if err != nil {
		err := svc.SetCookieValueByUniqueID(uniqueID)
		if err != nil {
			resp.Response500(c, err)
			return
		}
	} else {
		err := svc.UpdateCookieValueByUniqueID(uniqueID)
		if err != nil {
			resp.Response500(c, err)
			return
		}
	}

	c.SetCookie(adminCookieKey, uniqueID, defaultCookieAliveTime, defaultCookeiPath, defaultCookieDomain, secure, httpOnly)
}

// verifyCookieValue verify cookie is exist & cookie has correct value
func (ck *AdminCookie) verifyCookieValue(c *gin.Context, value string) bool {
	var cookieValue string
	var err error
	cookieValue, err = c.Cookie(adminCookieKey)

	if err == http.ErrNoCookie {
		return false
	} else {
		if result := strings.Compare(cookieValue, value); result != 0 {
			return false
		}
	}
	return true
}

// setCookie sets cookie to verify day's first time visit of visitor
func (ck *VisitTimeCookie) setCookie(c *gin.Context, today string, secure, httpOnly bool) {
	c.SetCookie(visitTimeCookieKey, today, 0, defaultCookeiPath, defaultCookieDomain, secure, httpOnly)
}

// verifyCookieValue verify cookie is exist & cookie has correct value
func (ck *VisitTimeCookie) verifyCookieValue(c *gin.Context, value string) bool {
	var cookieValue string
	var err error
	cookieValue, err = c.Cookie(visitTimeCookieKey)

	if err == http.ErrNoCookie {
		return false
	} else {
		if result := strings.Compare(cookieValue, value); result != 0 {
			return false
		}
	}
	return true
}
