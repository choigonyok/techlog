package main

import (
	"os"

	"github.com/choigonyok/blog-project/backend/internal/controller"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func originConfig() cors.Config{
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{os.Getenv("ORIGIN")} 
	config.AllowMethods= []string{"GET", "POST", "DELETE", "PUT"}
	config.AllowHeaders = []string{"Content-type"}
	config.AllowCredentials = true
	return config
}

func main() {
	godotenv.Load(".env")
	
	controller.ConnectDB(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER")+":"+os.Getenv("DB_PASSWORD")+"@tcp("+os.Getenv("DB_HOST")+")/"+os.Getenv("DB_NAME"))
	defer controller.UnConnectDB()

	eg := gin.Default()
	config := originConfig()
	eg.Use(cors.New(config))	

	eg.POST("/api/post/image", controller.WritePostHandler)      // 작성된 게시글에 썸네일 추가
	eg.POST("/api/post", controller.WritePostImageHandler)       // 게시글 작성
	eg.DELETE("/api/post/:postid", controller.DeletePostHandler) // 게시글 삭제
	eg.GET("/api/post/:postid", controller.GetPostHandler)       // 게시글 내용 불러오기
	eg.PUT("/api/post/:postid", controller.ModifyPostHandler)    // 게시글 수정

	eg.DELETE("/api/comment", controller.DeleteCommentHandler)                // 댓글 작성자의 댓글삭제
	eg.GET("/api/comment/:postid", controller.GetCommentHandler)              // 게시글 별 댓글 불러오기
	eg.POST("/api/comment", controller.AddCommentHandler)                     // 댓글 달기
	eg.GET("/api/comment/pw/:commentid", controller.GetCommentPWHandler)      // 댓글 작성시 생성한 비밀번호 불러오기
	eg.DELETE("/api/comment/:postid", controller.DeleteCommentByAdminHandler) // 관리자 권한 댓글 삭제

	eg.GET("/api/reply/:commentid", controller.GetReplyHandler)  // 댓글 별 답글 불러오기
	eg.POST("/api/reply/:commentid", controller.AddReplyHandler) // 답글 작성
	eg.DELETE("/api/reply", controller.DeleteReplyHandler)       // 답글 삭제

	eg.GET("/api/cookie", controller.GetTodayAndTotalVisitorNumHandler) // today, total 방문자 수 확인

	eg.POST("/api/login", controller.CheckAdminIDAndPW) // 로그인

	eg.POST("/api/tag", controller.GetPostsByTagHandler) // 태그 클릭시 포스트 출력
	eg.GET("/api/tag", controller.GetEveryTagHandler)    // 현재 존재하는 모든 태그 불러오기

	eg.GET("/api/assets/:name", controller.GetThumbnailHandler) // 게시글 썸네일 불러오기

	eg.Run(":8080")
}
