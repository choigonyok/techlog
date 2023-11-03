package router

import (
	"github.com/choigonyok/techlog/pkg/handler"
	"github.com/gin-gonic/gin"
)

type Route struct {
	Path    string
	Method  string
	Handler gin.HandlerFunc
}

const (
	POST   = "post"
	GET    = "get"
	PUT    = "put"
	DELETE = "delete"
)

func (r *Router) NewRoutes(prefix string) []Route {
	h := []Route{
		{
			Path:    prefix + "post/image",
			Method:  POST,
			Handler: handler.WritePostImageHandler,
		},
		{
			Path:    prefix + "/post/:postid/thumbnail",
			Method:  GET,
			Handler: handler.GetThumbnailByPostID,
		},
		{
			Path:    prefix + "post",
			Method:  POST,
			Handler: handler.WritePostHandler,
		},
		{
			Path:    prefix + "post/:postid",
			Method:  DELETE,
			Handler: handler.DeletePostHandler,
		},
		{
			Path:    prefix + "post/:postid",
			Method:  GET,
			Handler: handler.GetPost,
		},
		{
			Path:    prefix + "post/:postid",
			Method:  PUT,
			Handler: handler.ModifyPostHandler,
		},
		{
			Path:    prefix + "comment",
			Method:  DELETE,
			Handler: handler.DeleteCommentHandler,
		},
		{
			Path:    prefix + "comment/:postid",
			Method:  GET,
			Handler: handler.GetCommentHandler,
		},
		{
			Path:    prefix + "comment",
			Method:  POST,
			Handler: handler.AddCommentHandler,
		},
		{
			Path:    prefix + "comment/pw/:commentid",
			Method:  GET,
			Handler: handler.GetCommentPWHandler,
		},
		{
			Path:    prefix + "comment/:postid",
			Method:  DELETE,
			Handler: handler.DeleteCommentByAdminHandler,
		},
		{
			Path:    prefix + "reply/:commentid",
			Method:  GET,
			Handler: handler.GetReplyHandler,
		},
		{
			Path:    prefix + "reply/:commentid",
			Method:  POST,
			Handler: handler.AddReplyHandler,
		},
		{
			Path:    prefix + "reply",
			Method:  DELETE,
			Handler: handler.DeleteReplyHandler,
		},
		{
			Path:    prefix + "visitor",
			Method:  GET,
			Handler: handler.GetVisitorCounts,
		},
		{
			Path:    prefix + "login",
			Method:  POST,
			Handler: handler.CheckAdminIDAndPWHandler,
		},
		{
			Path:    prefix + "login",
			Method:  GET,
			Handler: handler.CheckCookieHandelr,
		},
		{
			Path:    prefix + "tag",
			Method:  POST,
			Handler: handler.GetEveryCardByTag,
		},
		{
			Path:    prefix + "tag",
			Method:  GET,
			Handler: handler.GetTags,
		},
	}
	return h
}
