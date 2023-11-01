package handler

import "github.com/gin-gonic/gin"

type Handler struct {
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

func NewHandlers(prefix string) []Handler {
	h := []Handler{
		{
			Path:    prefix + "post/image",
			Method:  POST,
			Handler: WritePostImageHandler,
		},
		{
			Path:    prefix + "post",
			Method:  POST,
			Handler: WritePostHandler,
		},
		{
			Path:    prefix + "post/:postid",
			Method:  DELETE,
			Handler: DeletePostHandler,
		},
		{
			Path:    prefix + "post/:postid",
			Method:  GET,
			Handler: GetPostHandler,
		},
		{
			Path:    prefix + "post/:postid",
			Method:  PUT,
			Handler: ModifyPostHandler,
		},
		{
			Path:    prefix + "comment",
			Method:  DELETE,
			Handler: DeleteCommentHandler,
		},
		{
			Path:    prefix + "comment/:postid",
			Method:  GET,
			Handler: GetCommentHandler,
		},
		{
			Path:    prefix + "comment",
			Method:  POST,
			Handler: AddCommentHandler,
		},
		{
			Path:    prefix + "comment/pw/:commentid",
			Method:  GET,
			Handler: GetCommentPWHandler,
		},
		{
			Path:    prefix + "comment/:postid",
			Method:  DELETE,
			Handler: DeleteCommentByAdminHandler,
		},
		{
			Path:    prefix + "reply/:commentid",
			Method:  GET,
			Handler: GetReplyHandler,
		},
		{
			Path:    prefix + "reply/:commentid",
			Method:  POST,
			Handler: AddReplyHandler,
		},
		{
			Path:    prefix + "reply",
			Method:  DELETE,
			Handler: DeleteReplyHandler,
		},
		{
			Path:    prefix + "visitor",
			Method:  GET,
			Handler: GetTodayAndTotalVisitorNumHandler,
		},
		{
			Path:    prefix + "login",
			Method:  POST,
			Handler: CheckAdminIDAndPWHandler,
		},
		{
			Path:    prefix + "login",
			Method:  GET,
			Handler: CheckCookieHandelr,
		},
		{
			Path:    prefix + "tag",
			Method:  POST,
			Handler: GetPostsByTagHandler,
		},
		{
			Path:    prefix + "tag",
			Method:  GET,
			Handler: GetEveryTagHandler,
		},
		{
			Path:    prefix + "assets/:name",
			Method:  GET,
			Handler: GetThumbnailHandler,
		},
	}
	return h
}
