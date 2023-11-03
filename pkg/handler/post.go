package handler

import (
	"fmt"
	"io"
	"os"

	"github.com/choigonyok/techlog/pkg/database"
	"github.com/choigonyok/techlog/pkg/model"
	resp "github.com/choigonyok/techlog/pkg/response"
	"github.com/choigonyok/techlog/pkg/service"
	"github.com/gin-gonic/gin"
)

// 작성된 게시글에 썸네일 추가
func WritePostImageHandler(c *gin.Context) {
}

// 게시글 작성 // DB에 저장할 때 tags Upper + 양사이드 whitespace 제거해야함 // ImagePath는 /assets/{filename} 형식으로 저장
func WritePostHandler(c *gin.Context) {

	// if !isCookieAdmin(c) {
	// 	http.Response500(c)
	// 	return
	// }
	// var data model.Post
	// if err := c.ShouldBindJSON(&data); err != nil {
	// 	http.Response500(c)
	// 	return
	// }
	// data.Text = strings.ReplaceAll(data.Text, `'`, `\'`)
	// err := model.AddPost(data.Tag, data.Title, data.Text, time.GetCurrentTimeByFormat("2006-01-02"))
	// if err != nil {
	// 	fmt.Println("ERROR #3 : ", err.Error())
	// 	http.Response500(c)
	// 	return
	// }
	// http.Response200(c)
}

// 게시글 삭제
func DeletePostHandler(c *gin.Context) {

}

// 게시글 내용 불러오기
func GetPostHandler(c *gin.Context) {
	pvr := database.NewMysqlProvider(database.GetConnector())
	svc := service.NewService(pvr)
	postID := c.Param("postid")
	posts, err := svc.GetPostByID(postID)
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

// 게시글 수정
func ModifyPostHandler(c *gin.Context) {

}

// 태그 클릭 시 게시글 출력
func GetEveryCardByTag(c *gin.Context) {
	pvr := database.NewMysqlProvider(database.GetConnector())
	svc := service.NewService(pvr)

	m := model.PostTags{}
	if err := c.ShouldBindJSON(&m); err != nil {
		resp.Response500(c)
		return
	}

	cards, err := svc.GetEveryCardByTag(m.Tags)
	if err != nil {
		resp.Response500(c)
		return
	}

	err = resp.ResponseDataWith200(c, cards)
	if err != nil {
		resp.Response500(c)
		return
	}
}

// 현재 존재하는 모든 카드 게시글 불러오기
func GetTags(c *gin.Context) {
	pvr := database.NewMysqlProvider(database.GetConnector())
	svc := service.NewService(pvr)
	tags, err := svc.GetEveryTags()
	if err != nil {
		resp.Response500(c)
		return
	}

	if resp.ResponseDataWith200(c, tags) != nil {
		resp.Response500(c)
	}
}

func GetImageByID(c *gin.Context) {
	imageName := c.Param("imageName")

	file, err := os.Open("assets/" + imageName)
	if err != nil {
		fmt.Println(err.Error())
		resp.Response500(c)
		return
	}
	defer file.Close()

	_, err = io.Copy(c.Writer, file)
	if err != nil {
		resp.Response500(c)
		return
	}
}
