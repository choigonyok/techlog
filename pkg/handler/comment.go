package handler

import (
	"strconv"

	"github.com/choigonyok/techlog/pkg/database"
	"github.com/choigonyok/techlog/pkg/model"
	resp "github.com/choigonyok/techlog/pkg/response"
	"github.com/choigonyok/techlog/pkg/service"
	"github.com/gin-gonic/gin"
)

// DeleteCommentByCommentID deletes comment by admin user or verified password
func DeleteCommentByCommentID(c *gin.Context) {
	pvr := database.NewMysqlProvider(database.GetConnector())
	svc := service.NewService(pvr)
	commentID := c.Param("commentid")
	// admin user delete
	userType := c.Query("type")
	if userType == "admin" {
		VerifyAdminUser(c)
		if err := svc.DeleteCommentByCommentID(commentID); err != nil {
			resp.Response500(c, err)
		} else {
			resp.Response200(c)
			return
		}
	}
	// common user delete
	inputPassword := c.Query("password")
	password, err := svc.GetCommentPasswordByCommentID(commentID)
	if err != nil {
		resp.Response500(c, err)
		return
	}
	if password == inputPassword {
		err := svc.DeleteCommentByCommentID(commentID)
		if err != nil {
			resp.Response500(c, err)
			return
		}
	} else {
		resp.Response400(c)
		return
	}
}

// GetCommentsByPostID returns comments in specific post
func GetCommentsByPostID(c *gin.Context) {
	pvr := database.NewMysqlProvider(database.GetConnector())
	svc := service.NewService(pvr)
	postID := c.Param("postid")

	comments, err := svc.GetCommentsByPostID(postID)
	if err != nil {
		resp.Response500(c, err)
		return
	}
	err = resp.ResponseDataWith200(c, comments)
	if err != nil {
		resp.Response500(c, err)
		return
	}
}

// CreateComment creates new comment
func CreateComment(c *gin.Context) {
	pvr := database.NewMysqlProvider(database.GetConnector())
	svc := service.NewService(pvr)
	postID := c.Param("postid")
	adminCookie := AdminCookie{}
	comment := model.Comment{}

	err := c.ShouldBindJSON(&comment)
	if err != nil {
		resp.Response500(c, err)
		return
	}
	comment.PostID, _ = strconv.Atoi(postID)

	clientCookieValue, _ := c.Cookie(adminCookieKey)
	isAdmin := adminCookie.verifyCookieValue(c, clientCookieValue)
	if isAdmin {
		comment.Admin = true
	} else {
		comment.Admin = false
	}

	err = svc.CreateComment(comment)
	if err != nil {
		resp.Response500(c, err)
		return
	}
}

// GetComments returns every comments
func GetComments(c *gin.Context) {
	pvr := database.NewMysqlProvider(database.GetConnector())
	svc := service.NewService(pvr)

	comments, err := svc.GetComments()
	if err != nil {
		resp.Response500(c, err)
		return
	}
	resp.ResponseDataWith200(c, comments)
}
