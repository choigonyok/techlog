package handler

import (
	"fmt"
	"strconv"

	"github.com/choigonyok/techlog/pkg/database"
	"github.com/choigonyok/techlog/pkg/model"
	resp "github.com/choigonyok/techlog/pkg/response"
	"github.com/choigonyok/techlog/pkg/service"
	"github.com/gin-gonic/gin"
)

// GetRepliesByPostID returns every reply in a post
func GetRepliesByPostID(c *gin.Context) {
	pvr := database.NewMysqlProvider(database.GetConnector())
	svc := service.NewService(pvr)
	postID := c.Param("postid")

	replies, err := svc.GetRepliesByPostID(postID)
	if err != nil {
		resp.Response500(c, err)
		return
	}

	err = resp.ResponseDataWith200(c, replies)
	if err != nil {
		resp.Response500(c, err)
		return
	}
}

// CreateReply creates new reply
func CreateReply(c *gin.Context) {
	pvr := database.NewMysqlProvider(database.GetConnector())
	svc := service.NewService(pvr)
	commentID := c.Param("commentid")
	postID := c.Param("postid")
	adminCookie := AdminCookie{}
	reply := model.Reply{}

	err := c.ShouldBindJSON(&reply)
	if err != nil {
		resp.Response500(c, err)
		return
	}

	fmt.Println(reply)
	fmt.Println(reply)
	fmt.Println(reply)
	fmt.Println(reply)
	reply.CommentID, _ = strconv.Atoi(commentID)
	reply.PostID, _ = strconv.Atoi(postID)

	clientCookieValue, _ := c.Cookie(adminCookieKey)
	isAdmin := adminCookie.verifyCookieValue(c, clientCookieValue)
	if isAdmin {
		reply.Admin = "1"
	} else {
		reply.Admin = "0"
	}

	err = svc.CreateReply(reply)
	if err != nil {
		resp.Response500(c, err)
		return
	}
}

// DeleteReplyByReplyID deletes a reply by password
func DeleteReplyByReplyID(c *gin.Context) {
	pvr := database.NewMysqlProvider(database.GetConnector())
	svc := service.NewService(pvr)
	replyID := c.Param("replyid")
	inputPassword := c.Query("password")

	password, err := svc.GetReplyPasswordByReplyID(replyID)
	if err != nil {
		resp.Response500(c, err)
		return
	}

	if password == inputPassword {
		err := svc.DeleteReplyByReplyID(replyID)
		if err != nil {
			resp.Response500(c, err)
			return
		}
	} else {
		resp.Response400(c)
	}
}
