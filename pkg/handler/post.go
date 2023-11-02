package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// 작성된 게시글에 썸네일 추가
func WritePostImageHandler(c *gin.Context) {
}

// 게시글 작성 // DB에 저장할 때 tags Upper + 양사이드 whitespace 제거해야함
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
	fmt.Println("CALLED")
	fmt.Println("CALLED")
	fmt.Println("CALLED")
	fmt.Println("CALLED")
	fmt.Println("CALLED")
}

// 게시글 수정
func ModifyPostHandler(c *gin.Context) {

}

// 게시글 썸네일 불러오기
func GetThumbnailHandler(c *gin.Context) {

}
