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
	pvrSlave := database.NewMysqlProvider(database.GetReadConnector())
	svcSlave := service.NewService(pvrSlave)
	postID := c.Param("postid")

	replies, err := svcSlave.GetRepliesByPostID(postID)
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
	pvrMaster := database.NewMysqlProvider(database.GetConnector())
	svcMaster := service.NewService(pvrMaster)
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

	err = svcMaster.CreateReply(reply)
	if err != nil {
		resp.Response500(c, err)
		return
	}
}

// DeleteReplyByReplyID deletes a reply by password
func DeleteReplyByReplyID(c *gin.Context) {
	pvrMaster := database.NewMysqlProvider(database.GetConnector())
	svcMaster := service.NewService(pvrMaster)
	pvrSlave := database.NewMysqlProvider(database.GetReadConnector())
	svcSlave := service.NewService(pvrSlave)
	replyID := c.Param("replyid")
	inputPassword := c.Query("password")

	password, err := svcSlave.GetReplyPasswordByReplyID(replyID)
	if err != nil {
		resp.Response500(c, err)
		return
	}

	if password == inputPassword {
		err := svcMaster.DeleteReplyByReplyID(replyID)
		if err != nil {
			resp.Response500(c, err)
			return
		}
	} else {
		resp.Response400(c)
	}
}
