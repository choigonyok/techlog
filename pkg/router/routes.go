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
		// Post
		{
			Path:    prefix + "post",
			Method:  POST,
			Handler: handler.CreatePost,
		},
		{
			Path:    prefix + "posts",
			Method:  GET,
			Handler: handler.GetPosts,
		},
		{
			Path:    prefix + "/posts/:postid/thumbnail",
			Method:  GET,
			Handler: handler.GetThumbnailByPostID,
		},
		{
			Path:    prefix + "posts/:postid",
			Method:  GET,
			Handler: handler.GetPost,
		},
		{
			Path:    prefix + "posts/:postid",
			Method:  PUT,
			Handler: handler.UpdatePostByPostID,
		},
		{
			Path:    prefix + "posts/:postid",
			Method:  DELETE,
			Handler: handler.DeletePostByPostID,
		},

		// Visitor
		{
			Path:    prefix + "visitor",
			Method:  GET,
			Handler: handler.GetVisitorCounts,
		},

		// Login
		{
			Path:    prefix + "login",
			Method:  POST,
			Handler: handler.VerifyAdminIDAndPW,
		},
		{
			Path:    prefix + "login",
			Method:  GET,
			Handler: handler.VerifyAdminUser,
		},

		// Tag
		{
			Path:    prefix + "tags",
			Method:  GET,
			Handler: handler.GetTags,
		},

		// Comment
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
	}
	return h
}
