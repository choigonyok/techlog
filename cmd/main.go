package main

import (
	"os"

	"github.com/choigonyok/blog-project-backend/internal/controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	controller.ConnectDB(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER")+":"+os.Getenv("DB_PASSWORD")+"@/"+os.Getenv("DB_NAME"))
	defer controller.UnConnectDB()

	eg := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"https://www.choigonyok.com", "http://www.choigonyok.com"}
	config.AllowMethods = []string{"POST", "DELETE", "GET", "PUT"}
	config.AllowHeaders = []string{"cookie", "Content-type"}
	config.AllowCredentials = true
	eg.Use(cors.New(config))

	eg.POST("/api/post/:param", controller.WritePostHandler)                        // 게시글 작성
	eg.GET("/api/cookie", controller.GetTodayAndTotalVisitorNumHandler)             // today, total 방문자 수 확인
	eg.POST("/api/mod/:param", controller.ModifyPostHandler)                        // 게시글 수정
	eg.POST("/api/login", controller.CheckIDAndPW)                                  // 로그인
	eg.POST("/api/tag", controller.GetPostsByTagHandler)                            // 태그 클릭시 포스트 출력
	eg.GET("/api/tag", controller.GetEveryTagHandler)                               // 현재 존재하는 모든 태그 불러오기
	eg.DELETE("/api/post/delete:deleteid", controller.DeletePostHandler)            // 게시글 삭제
	eg.PUT("/api/comments", controller.AddCommentHandler)                           // 댓글 달기
	eg.GET("/api/post/comments/pw/:uniqueid", controller.GetCommentPWHandler)       // 댓글 작성시 생성한 비밀번호 불러오기
	eg.DELETE("/api/post/comments/:postid", controller.DeleteCommentByAdminHandler) // 관리자 권한 댓글 삭제
	eg.POST("/api/post/comments", controller.DeleteCommentHandler)                  // 댓글 작성자의 댓글삭제
	eg.GET("/api/post/comments/:postid", controller.GetCommentHandler)              // 게시글 별 댓글 불러오기
	eg.GET("/api/reply/:commentid", controller.GetReplyHandler)                     // 댓글 별 답글 불러오기
	eg.PUT("/api/reply/:commentid", controller.AddReplyHandler)                     // 답글 작성
	eg.POST("/api/reply", controller.DeleteReplyHandler)                            // 답글 삭제
	eg.GET("/api/post/:postid", controller.GetPostHandler)                          // 게시글 내용 불러오기
	eg.GET("/api/IMAGES/:imgname", controller.GetThumbnailHandler)                  // 게시글 썸네일 불러오기
	eg.Run(":8080")
}
