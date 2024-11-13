package handler

import (
	"strconv"

	"github.com/choigonyok/techlog/internal/usecase"
	"github.com/gin-gonic/gin"
)

type VisitorHandler struct {
	gin.HandlerFunc
	usecase *usecase.VisitorUsecase
}

func NewVisitorHandler() *VisitorHandler {
	return &VisitorHandler{
		usecase: usecase.NewVisitorUsecase(),
	}
}

// GetVisitorCounts returns today/total visitor counts
func (h *VisitorHandler) GetVisitorCounts(c *gin.Context) {
	// today := time.GetCurrentTimeByFormat("2006-01-02")
	count := h.usecase.GetVisitorCount()

	c.Writer.Write([]byte(strconv.Itoa(count)))

	// cookie := &VisitTimeCookie{}

	// if !cookie.verifyCookieValue(c, today) {
	// 	err := svcMaster.AddTodayAndTotal()
	// 	if err != nil {
	// 		resp.Response500(c, err)
	// 		return
	// 	}
	// 	cookie.setCookie(c, today, false, true)
	// }

	// date, err := svcSlave.GetDate()
	// if err != nil {
	// 	resp.Response500(c, err)
	// 	return
	// }
	// if today != date {
	// 	err := svcMaster.ResetToday(today)
	// 	if err != nil {
	// 		resp.Response500(c, err)
	// 		return
	// 	}
	// }

	// todayCount, totalCount, err := svcSlave.GetCounts()
	// if err != nil {
	// 	resp.Response500(c, err)
	// 	return
	// }

	// err = resp.ResponseDataWith200(c, struct {
	// 	Today int `json:"today"`
	// 	Total int `json:"total"`
	// }{
	// 	Today: todayCount,
	// 	Total: totalCount,
	// })
	// if err != nil {
	// 	resp.Response500(c, err)
	// }
}
