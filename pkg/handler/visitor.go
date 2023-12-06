package handler

import (
	"github.com/choigonyok/techlog/pkg/database"
	resp "github.com/choigonyok/techlog/pkg/response"
	"github.com/choigonyok/techlog/pkg/service"
	"github.com/choigonyok/techlog/pkg/time"
	"github.com/gin-gonic/gin"
)

// GetVisitorCounts returns today/total visitor counts
func GetVisitorCounts(c *gin.Context) {
	pvrMaster := database.NewMysqlProvider(database.GetConnector())
	svcMaster := service.NewService(pvrMaster)
	pvrSlave := database.NewMysqlProvider(database.GetReadConnector())
	svcSlave := service.NewService(pvrSlave)
	today := time.GetCurrentTimeByFormat("2006-01-02")

	cookie := &VisitTimeCookie{}

	if !cookie.verifyCookieValue(c, today) {
		err := svcMaster.AddTodayAndTotal()
		if err != nil {
			resp.Response500(c, err)
			return
		}
		cookie.setCookie(c, today, false, true)
	}

	date, err := svcSlave.GetDate()
	if err != nil {
		resp.Response500(c, err)
		return
	}
	if today != date {
		err := svcMaster.ResetToday(today)
		if err != nil {
			resp.Response500(c, err)
			return
		}
	}

	todayCount, totalCount, err := svcSlave.GetCounts()
	if err != nil {
		resp.Response500(c, err)
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
