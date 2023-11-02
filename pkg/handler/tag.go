package handler

import (
	"github.com/choigonyok/techlog/pkg/database"
	"github.com/choigonyok/techlog/pkg/model"
	resp "github.com/choigonyok/techlog/pkg/response"
	"github.com/choigonyok/techlog/pkg/service"
	"github.com/gin-gonic/gin"
)

// 태그 클릭 시 게시글 출력
func GetEveryCardByTag(c *gin.Context) {
	pvr := database.NewMysqlProvider(database.GetConnector())
	svc := service.NewService(pvr)

	m := model.PostTags{}
	if err := c.ShouldBindJSON(&m); err != nil {
		resp.Response500(c)
		return
	}

	posts, err := svc.GetEveryCardByTag(m.Tags)
	if err != nil {
		resp.Response500(c)
		return
	}

	err = resp.ResponseDataWith200(c, posts)
	if err != nil {
		resp.Response500(c)
		return
	}
}

// 현재 존재하는 모든 태그 불러오기
func GetTags(c *gin.Context) {
	pvr := database.NewMysqlProvider(database.GetConnector())
	svc := service.NewService(pvr)
	tags, err := svc.GetEveryTag()
	if err != nil {
		resp.Response500(c)
		return
	}

	if resp.ResponseDataWith200(c, tags) != nil {
		resp.Response500(c)
	}
}
