package handler

import (
	"fmt"

	"github.com/choigonyok/techlog/pkg/database"
	resp "github.com/choigonyok/techlog/pkg/response"
	"github.com/choigonyok/techlog/pkg/service"
	"github.com/choigonyok/techlog/pkg/time"
	"github.com/gin-gonic/gin"
)

// GetVisitorCounts returns today/total visitor counts
func GetVisitorCounts(c *gin.Context) {
	pvr := database.NewMysqlProvider(database.GetConnector())
	svc := service.NewService(pvr)
	today := time.GetCurrentTimeByFormat("2006-01-02")

	cookie := &VisitTimeCookie{}

	if !cookie.verifyCookieValue(c, today) {
		err := svc.AddTodayAndTotal()
		if err != nil {
			resp.Response500(c, err)
			fmt.Println(err.Error())
			return
		}
		cookie.setCookie(c, today, false, true)
	}

	date, err := svc.GetDate()
	if err != nil {
		resp.Response500(c, err)
		fmt.Println(err.Error())
		return
	}
	if today != date {
		err := svc.ResetToday(today)
		if err != nil {
			resp.Response500(c, err)
			fmt.Println(err.Error())
			return
		}
	}

	todayCount, totalCount, err := svc.GetCounts()
	if err != nil {
		resp.Response500(c, err)
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
		resp.Response500(c, err)
	}
}
