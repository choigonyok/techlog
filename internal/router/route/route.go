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

	h := &Routes{
		// Post
		{
			Path:    prefix + "post",
			Method:  POST,
			Handler: m[handler.POST_HANDLER].(*handler.PostHandler).CreatePost,
			// handler.CreatePost
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
		// {
		// 	Path:    prefix + "posts/count",
		// 	Method:  GET,
		// 	Handler: handler.GetEveryPostCount,
		// },

		// Visitor
		{
			Path:    prefix + "visitor",
			Method:  GET,
			Handler: m[handler.VISITOR_HANDLER].(*handler.VisitorHandler).GetVisitorCounts,
		},

		// Deprecated: this is replaced with google SSO with Oauth2 Proxy
		// Login
		// {
		// 	Path:    prefix + "login",
		// 	Method:  POST,
		// 	Handler: handler.VerifyAdminIDAndPW,
		// },
		// {
		// 	Path:    prefix + "login",
		// 	Method:  GET,
		// 	Handler: handler.VerifyAdminUser,
		// },
		// 	// Comment
		// 	{
		// 		Path:    prefix + "posts/:postid/comment",
		// 		Method:  POST,
		// 		Handler: handler.CreateComment,
		// 	},
		// 	{
		// 		Path:    prefix + "comments",
		// 		Method:  GET,
		// 		Handler: handler.GetComments,
		// 	},
		// 	{
		// 		Path:    prefix + "posts/:postid/comments",
		// 		Method:  GET,
		// 		Handler: handler.GetCommentsByPostID,
		// 	},
		// 	{
		// 		Path:    prefix + "comments/:commentid",
		// 		Method:  DELETE,
		// 		Handler: handler.DeleteCommentByCommentID,
		// 	},

		// 	// Reply
		// 	{
		// 		Path:    prefix + "posts/:postid/replies",
		// 		Method:  GET,
		// 		Handler: handler.GetRepliesByPostID,
		// 	},
		// 	{
		// 		Path:    prefix + "posts/:postid/comments/:commentid/reply",
		// 		Method:  POST,
		// 		Handler: handler.CreateReply,
		// 	},
		// 	{
		// 		Path:    prefix + "posts/:postid/comments/:commentid/replies/:replyid",
		// 		Method:  DELETE,
		// 		Handler: handler.DeleteReplyByReplyID,
		// 	},
	}
	return h
}
