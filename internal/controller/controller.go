package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/choigonyok/blog-project-backend/internal/model"
	"github.com/gin-gonic/gin"
)

func isCookieAdmin(c *gin.Context) bool {
	inputValue, CookieErr := c.Cookie("admin")
	cookieValue, err := model.GetCookieValue(inputValue)
	if err != nil {
		return false
	}
	if CookieErr != nil || cookieValue != inputValue {
		return false
	}
	return true
}

func isCookieValid(c *gin.Context) bool {
	visitTime, err := c.Cookie("visitTime")
	if err == http.ErrNoCookie {
		return false
	}
	isValid := strings.Contains(visitTime, getTimeNow().Format("2006-01-02"))
	return isValid
}

func ConnectDB(driverName, dbData string) {
	err := model.OpenDB(driverName, dbData)
	if err != nil {
		fmt.Println("ERROR #73 : ", err.Error())
	}
}

func UnConnectDB() {
	err := model.CloseDB()
	if err != nil {
		fmt.Println("ERROR #74 : ", err.Error())
	}
}

func CheckAdminIDAndPW(c *gin.Context) {
	data := struct {
		ID       string `json:"id"`
		Password string `json:"pw"`
	}{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		fmt.Println("ERROR #1 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if data.ID != os.Getenv("BLOG_ID") || data.Password != os.Getenv("BLOG_PW") {
		fmt.Println("ERROR #31 : ", err.Error())
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	
	cookieValue, err := model.UpdateCookieRecord()
	if err != nil {
		fmt.Println("ERROR #30 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.SetCookie("admin", cookieValue.String(), 60*60*12, "/", os.Getenv("ORIGIN"), false, true)
	c.Writer.WriteHeader(http.StatusOK)
}

func WritePostHandler(c *gin.Context) {
	if !isCookieAdmin(c) {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	var data model.Post
	if err := c.ShouldBindJSON(&data); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	recentPostID, err := model.GetRecentPostID()
	if err != nil {
		fmt.Println("ERROR #2 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	data.Text = strings.ReplaceAll(data.Text, `'`, `\'`)
	err = model.AddPost(recentPostID+1, data.Tag, data.Title, data.Text, data.WriteTime.Format("2006-01-02"))
	if err != nil {
		fmt.Println("ERROR #3 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}

func WritePostImageHandler(c *gin.Context) {
	if !isCookieAdmin(c) {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	imgFile, err := c.MultipartForm()
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	recentID, err := model.GetRecentPostID()
	if err != nil {
		fmt.Println("ERROR #4 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	everyIMG := imgFile.File["file"]
	for _, v := range everyIMG {
		noSpaceImageName := strings.ReplaceAll(v.Filename, " ", "")
		err = c.SaveUploadedFile(v, "assets/"+strconv.Itoa(recentID)+"-"+noSpaceImageName)
		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	noSpaceThumbnailName := strings.ReplaceAll(everyIMG[0].Filename, " ", "")
	err = model.UpdatePostImagePath(recentID, noSpaceThumbnailName)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func DeletePostHandler(c *gin.Context) {
	if !isCookieAdmin(c) {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	postID := c.Param("postid")

	// 이미지 없이 작성된 글 삭제
	if postID == "0" {
		err := model.DeleteRecentPost()
		if err != nil {
			fmt.Println("ERROR #7 : ", err.Error())
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		err := model.DeletePostByPostID(postID)
		if err != nil {
			fmt.Println("ERROR #8 : ", err.Error())
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		imageList, err := os.ReadDir("assets")
		if err != nil {
			fmt.Println("ERROR #9 : ", err.Error())
		}
		for _, v := range imageList {
			if strings.HasPrefix(v.Name(), postID+"-") {
				os.Remove("assets/" + v.Name())
				if err != nil {
					fmt.Println("ERROR #26 : ", err.Error())
					c.Writer.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
		}
		commentsIDs, err := model.SelectEveryCommentIDByPostID(postID)
		if err != nil {
			fmt.Println("ERROR #10 : ", err.Error())
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		for _, v := range commentsIDs {
			err = model.DeleteEveryCommentByCommentID(v)
			if err != nil {
				fmt.Println("ERROR #12 : ", err.Error())
				c.Writer.WriteHeader(http.StatusInternalServerError)
				return
			}
			err = model.DeleteEveryReplyByCommentID(v)
			if err != nil {
				fmt.Println("ERROR #16 : ", err.Error())
				c.Writer.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
	c.Writer.WriteHeader(http.StatusOK)
}

func getTimeNow() time.Time {
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		fmt.Println("ERROR #32 : ", err.Error())
	}
	now := time.Now()
	t := now.In(loc)
	return t
}

func GetTodayAndTotalVisitorNumHandler(c *gin.Context) {
	if !isCookieValid(c) {
		err := model.CountTodayVisit()
		if err != nil {
			fmt.Println("ERROR #33 : ", err.Error())
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		c.SetCookie("visitTime", getTimeNow().String(), 0, "/", os.Getenv("ORIGIN"), false, true)
	}
	c.Writer.WriteHeader(http.StatusOK)
}

func ModifyPostHandler(c *gin.Context) {
	if !isCookieAdmin(c) {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	postID := c.Param("postid")
	var PostData model.Post
	err := c.ShouldBindJSON(&PostData)
	if err != nil {
		fmt.Println("ERROR #22 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	PostData.Text = strings.ReplaceAll(PostData.Text, `'`, `\'`)
	err = model.UpdatePost(PostData.Title, PostData.Text, PostData.Tag, postID, PostData.WriteTime)
	if err != nil {
		fmt.Println("ERROR #5 : ", err.Error())
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}

func GetPostsByTagHandler(c *gin.Context) {
	var data model.Post
	if err := c.ShouldBindJSON(&data); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if data.Tag == "ALL" {
		data.Tag = ""
	}
	mainTag, _, _ := strings.Cut(data.Tag, " ")
	datas, err := model.SelectPostByTag(mainTag)
	if err != nil {
		fmt.Println("ERROR #6 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	marshaledData, err := json.Marshal(datas)
	if err != nil {
		fmt.Println("ERROR #23 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Write(marshaledData)
}

// 코드리뷰
// 코드리뷰
// 코드리뷰
// 코드리뷰
// 코드리뷰
func GetEveryTagHandler(c *gin.Context) {
	tagString, err := model.GetEveryTagAsString()
	if err != nil {
		fmt.Println("ERROR #24 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	tagString = strings.ReplaceAll(tagString, " / ", " ")
	_, tagString, ok := strings.Cut(tagString, " ")
	tagCount := strings.Count(tagString, " ") //모든 포스트의 총 tag 합계-중복포함
	posts := []model.Post{}
	post := model.Post{}

	if !ok {
		fmt.Println("STRING ERROR 1")
	}
	for i := 0; i < tagCount; i++ {
		b, a, ok := strings.Cut(tagString, " ")
		if !ok {
			fmt.Println("TAG COUNT ERROR OCCURED")
		}
		post.Tag = strings.ToUpper(b)
		tagString = a
		posts = append(posts, post)
	}
	post.Tag = strings.ToUpper(tagString)
	posts = append(posts, post)
	ret := []model.Post{}
	tagMap := make(map[model.Post]int)
	for i, v := range posts {
		if _, ok := tagMap[v]; !ok {
			tagMap[v] = i
			ret = append(ret, v)
		}
	}
	marshaledData, err := json.Marshal(ret)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Write(marshaledData)
}

func AddCommentHandler(c *gin.Context) {
	var data model.Comment
	err := c.ShouldBindJSON(&data)
	if err != nil {
		fmt.Println("ERROR #27 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	// 빈 텍스트가 있는지 확인
	isTextEmpty, textErr := regexp.MatchString("^$", data.Text)
	isIDEmpty, idErr := regexp.MatchString("^$", data.WriterID)
	isPwValid, pwErr := regexp.MatchString("^[0-9]+$", data.WriterPW)
	if textErr != nil || idErr != nil || pwErr != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !isTextEmpty || !isIDEmpty || !isPwValid {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if !isCookieAdmin(c) {
		data.Admin = 1
	} else {
		data.Admin = 0
	}
	strings.ReplaceAll(data.WriterID, `'`, `\'`)
	strings.ReplaceAll(data.WriterPW, `'`, `\'`)
	recentCommentID, err := model.GetRecentCommentID()
	if err != nil {
		fmt.Println("ERROR #12 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = model.InsertComment(data.PostID, recentCommentID+1, data.Admin, data.Text, data.WriterID, data.WriterPW)
	if err != nil {
		fmt.Println("ERROR #13 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}

func GetCommentPWHandler(c *gin.Context) {
	commentID := c.Param("commentid")
	writerPW, err := model.GetCommentWriterPWByCommentID(commentID)
	if err != nil {
		fmt.Println("ERROR #14 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	TempWriterPW := struct {
		WriterPW string `json:"writerpw"`
	}{
		writerPW,
	}
	marshaledData, err := json.Marshal(TempWriterPW)
	if err != nil {
		fmt.Println("ERROR #28 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.Write(marshaledData)
}

func DeleteCommentByAdminHandler(c *gin.Context) {
	param := c.Param("postid")
	_, err := c.Cookie("admin") // 쿠키 : admin 권한 확인
	if err == http.ErrNoCookie {
		c.String(http.StatusUnauthorized, "THERE IS NO COOKIE")
	} else {
		err := model.DeleteEveryReplyByCommentID(param)
		if err != nil {
			fmt.Println("ERROR #15 : ", err.Error())
		}
		err = model.DeleteEveryCommentByCommentID(param)
		if err != nil {
			fmt.Println("ERROR #17 : ", err.Error())
		} else {
			c.Writer.WriteHeader(http.StatusOK)
		}
	}
}

func DeleteCommentHandler(c *gin.Context) {
	comid := c.Query("comid")
	inputpw := c.Query("inputpw")
	writerPW, err := model.GetCommentWriterPWByCommentID(comid)
	if err != nil {
		fmt.Println("ERROR #18 : ", err.Error())
	}
	if inputpw == writerPW {
		err := model.DeleteEveryReplyByCommentID(comid)
		if err != nil {
			fmt.Println("ERROR #19 : ", err.Error())
		}
		err2 := model.DeleteEveryCommentByCommentID(comid)
		if err2 != nil {
			fmt.Println("ERROR #20 : ", err.Error())
		}
		c.Writer.WriteHeader(http.StatusOK)
	} else {
		c.Writer.WriteHeader(http.StatusBadRequest)
	}
}

func GetCommentHandler(c *gin.Context) {
	param := c.Param("postid")
	postID, _ := strconv.Atoi(param)
	if param == "0" { // 전체 댓글 중 admin 댓글이 아닌 것
		commentSlice, err := model.SelectNotAdminWriterComment(postID)
		if err != nil {
			fmt.Println("ERROR #21 : ", err.Error())
		}
		marshaledData, err := json.Marshal(commentSlice)
		if err != nil {
			fmt.Println("ERROR #23 : ", err.Error())
		}
		c.Writer.Write(marshaledData)
	} else { // 해당 post 댓글
		commentSlice, err := model.SelectCommentByPostID(postID)
		if err != nil {
			fmt.Println("ERROR #22 : ", err.Error())
		}
		marshaledData, err := json.Marshal(commentSlice)
		if err != nil {
			fmt.Println("ERROR #23 : ", err.Error())
		}
		c.Writer.Write(marshaledData)
	}
}

// 여기서부터 시작
// 여기서부터 시작
// 여기서부터 시작
func GetReplyHandler(c *gin.Context) {
	commentID := c.Param("commentid")
	replySlice, err := model.SelectReplyByCommentID(commentID)
	if err != nil {
		fmt.Println("ERROR #24 : ", err.Error())
	}
	marshaledData, err := json.Marshal(replySlice)
	if err != nil {
		fmt.Println("ERROR #25 : ", err.Error())
	}
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(marshaledData)
}

func AddReplyHandler(c *gin.Context) {
	commentid := c.Param("commentid")
	data := model.ReplyData{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		fmt.Println("REPLY JSON BINDING ERROR")
		c.Writer.WriteHeader(http.StatusInternalServerError)
	}
	// 빈 텍스트가 있는지 확인
	isReplyEmpty, err := regexp.MatchString("^$", data.Reply)
	if err != nil {
		fmt.Println("REGEXP ERROR1")
	}
	isIDEmpty, err := regexp.MatchString("^$", data.ReplyID)
	if err != nil {
		fmt.Println("REGEXP ERROR2")
	}
	if isReplyEmpty || isIDEmpty {
		c.Writer.WriteHeader(http.StatusLengthRequired)
	} else {
		// password 숫자 맞는지 확인
		isPwNumber, err := regexp.MatchString("^[0-9]+$", data.ReplyPW)
		if err != nil {
			fmt.Println("REGEXP ERROR3")
		}
		if !isPwNumber {
			c.Writer.WriteHeader(http.StatusNotAcceptable)
		} else {
			_, err := c.Cookie("admin") // 쿠키 : admin 권한 확인
			if err == http.ErrNoCookie {
				data.ReplyIsAdmin = 0
			} else {
				data.ReplyIsAdmin = 1
			}
			if strings.Contains(data.Reply, `'`) || strings.Contains(data.ReplyPW, `'`) || strings.Contains(data.ReplyID, `'`) {
				c.Writer.WriteHeader(http.StatusBadRequest)
			} else {
				recentReplyID, err := model.GetRecentReplyID()
				if err != nil {
					fmt.Println("ERROR #26 : ", err.Error())
				}
				err = model.InsertReply(data.ReplyIsAdmin, recentReplyID+1, commentid, data.Reply, data.ReplyID, data.ReplyPW)
				if err != nil {
					fmt.Println("ERROR #27 : ", err.Error())
					c.Writer.WriteHeader(http.StatusInternalServerError)
				} else {
					c.Writer.WriteHeader(http.StatusOK)
				}
			}
		}
	}
}

func DeleteReplyHandler(c *gin.Context) {
	inputpw := c.Query("inputpw")
	replyid := c.Query("replyid")
	replyPW, err := model.GetReplyPWByReplyID(replyid)
	if err != nil {
		fmt.Println("ERROR #28 : ", err.Error())
	}
	if replyPW == inputpw {
		err := model.DeleteReplyByReplyID(replyid)
		if err != nil {
			fmt.Println("ERROR #29 : ", err.Error())
		}
		c.Writer.WriteHeader(http.StatusOK)
	} else {
		c.Writer.WriteHeader(http.StatusBadRequest)
	}
}

func GetPostHandler(c *gin.Context) {
	postid := c.Param("postid")
	if postid == "all" {
		postData, err := model.GetEveryPost()
		if err != nil {
			fmt.Println("ERROR #30 : ", err.Error())
		}
		marshaledData, err := json.Marshal(postData)
		if err != nil {
			fmt.Println("ERROR #32 : ", err.Error())
			return
		}
		c.Writer.Write(marshaledData)
	} else {
		postData, err := model.GetPostByPostID(postid)
		if err != nil {
			fmt.Println("ERROR #31 : ", err.Error())
		}
		marshaledData, err := json.Marshal(postData)
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}
		c.Writer.Write(marshaledData)
	}
}

func GetThumbnailHandler(c *gin.Context) {
	imgname := c.Param("imgname")
	imgname = strings.ReplaceAll(imgname, " ", "")
	file, err := os.Open("IMAGES/" + imgname)
	if err != nil {
		// 파일 열기에 실패한 경우 에러 처리
		http.Error(c.Writer, "파일 오픈 실패", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		// 파일 전송에 실패한 경우 에러 처리
		http.Error(c.Writer, "Failed to send file", http.StatusInternalServerError)
		return
	}
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

}
