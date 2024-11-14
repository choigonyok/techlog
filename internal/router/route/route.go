package route

import (
	"github.com/choigonyok/techlog/internal/handler"
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
	m := make(map[int]handler.Handler)
	m[handler.VISITOR_HANDLER] = handler.NewVisitorHandler()
	m[handler.POST_HANDLER] = handler.NewPostHandler()
	m[handler.AUTH_HANDLER] = handler.NewAuthHandler()

	h := &Routes{
		// Post
		{
			Path:    prefix + "post",
			Method:  POST,
			Handler: m[handler.POST_HANDLER].(*handler.PostHandler).CreatePost,
		},
		{
			Path:    prefix + "posts",
			Method:  GET,
			Handler: m[handler.POST_HANDLER].(*handler.PostHandler).GetPosts,
		},
		{
			Path:    prefix + "posts/:postId",
			Method:  GET,
			Handler: m[handler.POST_HANDLER].(*handler.PostHandler).GetPost,
		},
		{
			Path:    prefix + "posts/:postId",
			Method:  DELETE,
			Handler: m[handler.POST_HANDLER].(*handler.PostHandler).DeletePost,
		},
		{
			Path:    prefix + "posts/:postId",
			Method:  PUT,
			Handler: m[handler.POST_HANDLER].(*handler.PostHandler).UpdatePost,
		},

		// Tag
		{
			Path:    prefix + "tags",
			Method:  GET,
			Handler: m[handler.POST_HANDLER].(*handler.PostHandler).GetTags,
		},

		// Image
		{
			Path:    prefix + "posts/:postId/thumbnail",
			Method:  GET,
			Handler: m[handler.POST_HANDLER].(*handler.PostHandler).GetThumbnail,
		},
		{
			Path:    prefix + "posts/:postId/images/:imageId",
			Method:  GET,
			Handler: m[handler.POST_HANDLER].(*handler.PostHandler).GetImage,
		},
		{
			Path:    prefix + "posts/:postId/images",
			Method:  PUT,
			Handler: m[handler.POST_HANDLER].(*handler.PostHandler).UpdateImages,
		},
		{
			Path:    prefix + "posts/:postId/images",
			Method:  GET,
			Handler: m[handler.POST_HANDLER].(*handler.PostHandler).GetImages,
		},

		// Visitor
		{
			Path:    prefix + "visitor",
			Method:  GET,
			Handler: m[handler.VISITOR_HANDLER].(*handler.VisitorHandler).GetVisitorCounts,
		},

		// Login
		{
			Path:    prefix + "login",
			Method:  GET,
			Handler: m[handler.AUTH_HANDLER].(*handler.AuthHandler).Login,
		},
		{
			Path:    prefix + "github/callback",
			Method:  GET,
			Handler: m[handler.AUTH_HANDLER].(*handler.AuthHandler).Callback,
		},
		{
			Path:    prefix + "token",
			Method:  POST,
			Handler: m[handler.AUTH_HANDLER].(*handler.AuthHandler).Validate,
		},
	}
	return h
}
