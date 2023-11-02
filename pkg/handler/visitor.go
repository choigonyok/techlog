package handler

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/choigonyok/techlog/pkg/data"
	"github.com/choigonyok/techlog/pkg/database"
	resp "github.com/choigonyok/techlog/pkg/response"
	"github.com/choigonyok/techlog/pkg/service"
	"github.com/choigonyok/techlog/pkg/time"
	"github.com/gin-gonic/gin"
)

var (
	cookieKey      = data.EncodeBase64("visitTime")
	cookieDomain   = os.Getenv("HOST")
	cookieSecure   = false
	cookieHttpOnly = true
)

// GetVisitorCounts returns today/total visitor counts
func GetVisitorCounts(c *gin.Context) {
	pvr := database.NewMysqlProvider(database.GetConnector())
	svc := service.NewService(pvr)
	today := time.GetCurrentTimeByFormat("2006-01-02")

	if !verifyCookieValue(c, today) {
		err := svc.AddTodayAndTotal()
		if err != nil {
			resp.Response500(c)
			fmt.Println(err.Error())
			return
		}
		setCookie(c, today, cookieSecure, cookieHttpOnly)
	}

	date, err := svc.GetDate()
	if err != nil {
		resp.Response500(c)
		fmt.Println(err.Error())
		return
	}
	if today != date {
		err := svc.ResetToday(today)
		if err != nil {
			resp.Response500(c)
			fmt.Println(err.Error())
			return
		}
	}

	todayCount, totalCount, err := svc.GetCounts()
	if err != nil {
		resp.Response500(c)
		fmt.Println(err.Error())
		return
	}

	err = resp.ResponseDataWith200(c, struct {
		Today int `json:"today"`
		Total int `json:"total"`
	}{
		Today: todayCount,
		Total: totalCount,
	})
	if err != nil {
		resp.Response500(c)
	}
}

// setCookie sets cookie to verify day's first time visit of visitor
func setCookie(c *gin.Context, today string, secure, httpOnly bool) {
	c.SetCookie(cookieKey, today, 0, "/", cookieDomain, secure, httpOnly)
}

// verifyCookieValue verify cookie is exist & cookie has correct value
func verifyCookieValue(c *gin.Context, value string) bool {
	cookieValue, err := c.Cookie(cookieKey)
	if err == http.ErrNoCookie {
		return false
	} else {
		if result := strings.Compare(cookieValue, value); result != 0 {
			return false
		}
	}
	return true
}
