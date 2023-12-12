package route

import (
	"github.com/choigonyok/techlog/pkg/handler"
	"github.com/gin-gonic/gin"
)

type Routes []struct {
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

func New(prefix string) *Routes {
	h := &Routes{
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
			Path:    prefix + "posts/:postid/thumbnail",
			Method:  GET,
			Handler: handler.GetThumbnailByPostID,
		},
		{
			Path:    prefix + "posts/:postid/images/:imageid",
			Method:  GET,
			Handler: handler.GetPostImageByImageID,
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
			Path:    prefix + "posts/:postid/images",
			Method:  PUT,
			Handler: handler.UpdatePostImagesByPostID,
		},
		{
			Path:    prefix + "posts/:postid",
			Method:  DELETE,
			Handler: handler.DeletePostByPostID,
		},
		{
			Path:    prefix + "posts/:postid/images",
			Method:  GET,
			Handler: handler.GetImagesByPostID,
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

		// Tagee
		{
			Path:    prefix + "tags",
			Method:  GET,
			Handler: handler.GetTags,
		},

		// Comment
		{
			Path:    prefix + "posts/:postid/comment",
			Method:  POST,
			Handler: handler.CreateComment,
		},
		{
			Path:    prefix + "comments",
			Method:  GET,
			Handler: handler.GetComments,
		},
		{
			Path:    prefix + "posts/:postid/comments",
			Method:  GET,
			Handler: handler.GetCommentsByPostID,
		},
		{
			Path:    prefix + "comments/:commentid",
			Method:  DELETE,
			Handler: handler.DeleteCommentByCommentID,
		},

		// Reply
		{
			Path:    prefix + "posts/:postid/replies",
			Method:  GET,
			Handler: handler.GetRepliesByPostID,
		},
		{
			Path:    prefix + "posts/:postid/comments/:commentid/reply",
			Method:  POST,
			Handler: handler.CreateReply,
		},
		{
			Path:    prefix + "posts/:postid/comments/:commentid/replies/:replyid",
			Method:  DELETE,
			Handler: handler.DeleteReplyByReplyID,
		},
	}
	return h
}
