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

	"github.com/choigonyok/blog-project/backend/internal/model"
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
	_, err := c.Cookie("visitTime")
	if err == http.ErrNoCookie {
		return false
	}
	// isValid := strings.Contains(visitTime, getTimeNow().Format("2006-01-02"))
	// return isVali
	return true
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
		fmt.Println("ERROR #37: ", err.Error())
		return
	}
	data.Text = strings.ReplaceAll(data.Text, `'`, `\'`)
	err := model.AddPost(data.Tag, data.Title, data.Text, getTimeNow().Format("2006/01/02"))
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
		fmt.Println("ERROR #35 : ", err.Error())
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
			fmt.Println("ERROR #34 : ", err.Error())
			return
		}
	}
	noSpaceThumbnailName := strings.ReplaceAll(everyIMG[0].Filename, " ", "")
	err = model.UpdatePostImagePath(recentID, noSpaceThumbnailName)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		fmt.Println("ERROR #36 : ", err.Error())
		return
	}
}

func DeletePostHandler(c *gin.Context) {
	if !isCookieAdmin(c) {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	postID := c.Param("postid")
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
			err = model.DeleteCommentByCommentID(v)
			if err != nil {
				fmt.Println("ERROR #12 : ", err.Error())
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
	var cookieNeed bool
	cookieNeed = false
	today := getTimeNow().Format("2006-01-02")
	cookieValue, err := c.Cookie("visitTime")
	if err == http.ErrNoCookie {
		cookieNeed = true
	} else {
		isValid := strings.Contains(cookieValue, today)
		if !isValid {
			cookieNeed = true
		}
	}
	if cookieNeed {
		visitor,err := model.GetVisitorCount()
		if err != nil {
			fmt.Println("ERROR #33 : ", err.Error())
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = model.AddVisitorCount(visitor)
		c.SetCookie("visitTime", getTimeNow().String(), 0, "/", os.Getenv("ORIGIN"), false, true)
	}
	todayRecord, err := model.GetTodayRecord()
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if today != todayRecord {
		err := model.ResetTodayVisitorNum(today)
		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			fmt.Println("ERROR #40 : ", err.Error())
			return
		}	
	}
	
	visitor, err := model.GetVisitorCount()
	if err != nil {
		fmt.Println("ERROR #41 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	marshaledData, err := json.Marshal(visitor)
	if err != nil {
		fmt.Println("ERROR #40 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.Write(marshaledData)
}

func ModifyPostHandler(c *gin.Context) {
	if !isCookieAdmin(c) {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	postID := c.Param("postid")
	var data model.Post
	err := c.ShouldBindJSON(&data)
	if err != nil {
		fmt.Println("ERROR #22 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	data.Text = strings.ReplaceAll(data.Text, `'`, `\'`)
	err = model.UpdatePost(data.Title, data.Text, data.Tag, postID)
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
	for i := range datas {
		datas[i].Text = strings.ReplaceAll(datas[i].Text, `\'`, `'`)
		datas[i].Tag = strings.ToUpper(datas[i].Tag)
	}

	marshaledData, err := json.Marshal(datas)
	if err != nil {
		fmt.Println("ERROR #23 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.Write(marshaledData)
}

// 코드리뷰
func GetEveryTagHandler(c *gin.Context) {
	tagString, err := model.GetEveryTagAsString()
	if err != nil {
		fmt.Println("ERROR #24 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	tagString = strings.ReplaceAll(tagString, " / ", " ")
	_, tagString, _ = strings.Cut(tagString, " ")
	tagCount := strings.Count(tagString, " ")
	posts := []model.Post{}
	post := model.Post{}
	for i := 0; i < tagCount; i++ {
		b, a, ok := strings.Cut(tagString, " ")
		if !ok {
			fmt.Println("ERROR #42 : ", err.Error())
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
	isTextEmpty, textErr := regexp.MatchString("^$", data.Text)
	isIDEmpty, idErr := regexp.MatchString("^$", data.WriterID)
	isPwValid, pwErr := regexp.MatchString("^[0-9]+$", data.WriterPW)
	if textErr != nil || idErr != nil || pwErr != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if isTextEmpty || isIDEmpty || !isPwValid {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if isCookieAdmin(c) {
		data.Admin = 1
	} else {
		data.Admin = 0
	}
	data.WriterID = strings.ReplaceAll(data.WriterID, `'`, `\'`)
	data.Text = strings.ReplaceAll(data.Text, `'`, `\'`)
	err = model.InsertComment(data.PostID, data.Admin, data.Text, data.WriterID, data.WriterPW)
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
	postID := c.Param("postid")
	if !isCookieAdmin(c) {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	err := model.DeleteCommentByCommentID(postID)
	if err != nil {
		fmt.Println("ERROR #17 : ", err.Error())
	} else {
		c.Writer.WriteHeader(http.StatusOK)
	}
}

func DeleteCommentHandler(c *gin.Context) {
	comnmentID := c.Query("commentid")
	inputPW := c.Query("inputpw")
	writerPW, err := model.GetCommentWriterPWByCommentID(comnmentID)
	if err != nil {
		fmt.Println("ERROR #18 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if inputPW == writerPW {
		err := model.DeleteCommentByCommentID(comnmentID)
		if err != nil {
			fmt.Println("ERROR #20 : ", err.Error())
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
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
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		for _, v := range commentSlice {
			v.Text = strings.ReplaceAll(v.Text, `\'`, `'`)
			v.WriterID = strings.ReplaceAll(v.Text, `\'`, `'`)
		}
		marshaledData, err := json.Marshal(commentSlice)
		if err != nil {
			fmt.Println("ERROR #23 : ", err.Error())
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		c.Writer.Write(marshaledData)
	}
}

func GetReplyHandler(c *gin.Context) {
	commentID := c.Param("commentid")
	replySlice, err := model.SelectReplyByCommentID(commentID)
	if err != nil {
		fmt.Println("ERROR #24 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, v := range replySlice {
		v.Text = strings.ReplaceAll(v.Text,`\'`, `'`)
		v.WriterID = strings.ReplaceAll(v.Text,`\'`, `'`)
	}
	marshaledData, err := json.Marshal(replySlice)
	if err != nil {
		fmt.Println("ERROR #25 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(marshaledData)
}

func AddReplyHandler(c *gin.Context) {
	commentID := c.Param("commentid")
	data := model.Reply{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		fmt.Println("ERROR #43 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	isReplyEmpty, err := regexp.MatchString("^$", data.Text)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	isIDEmpty, err := regexp.MatchString("^$", data.WriterID)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if isReplyEmpty || isIDEmpty {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	isPWValid, err := regexp.MatchString("^[0-9]+$", data.WriterPW)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !isPWValid {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
		
	}
	if !isCookieAdmin(c) {
		data.Admin = 0
	} else {
		data.Admin = 1
	}
	data.Text = strings.ReplaceAll(data.Text, `'`, `\'`)
	data.WriterID = strings.ReplaceAll(data.WriterID, `'`, `\'`)
	err = model.InsertReply(data.Admin, commentID, data.Text, data.WriterID, data.WriterPW)
	if err != nil {
		fmt.Println("ERROR #27 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
	} else {
		c.Writer.WriteHeader(http.StatusOK)
	}
}

func DeleteReplyHandler(c *gin.Context) {
	inputPW := c.Query("inputpw")
	ID := c.Query("replyid")
	replyPW, err := model.GetReplyPWByReplyID(ID)
	if err != nil {
		fmt.Println("ERROR #28 : ", err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if replyPW == inputPW {
		err := model.DeleteReplyByReplyID(ID)
		if err != nil {
			fmt.Println("ERROR #29 : ", err.Error())
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		c.Writer.WriteHeader(http.StatusOK)
	} else {
		c.Writer.WriteHeader(http.StatusBadRequest)
	}
}

func GetPostHandler(c *gin.Context) {
	postID := c.Param("postid")
	if postID == "all" {
		datas, err := model.GetEveryPost()
		if err != nil {
			fmt.Println("ERROR #30 : ", err.Error())
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		for i := range datas {
			datas[i].Tag = strings.ToUpper(datas[i].Tag)
		}
		marshaledData, err := json.Marshal(datas)
		if err != nil {
			fmt.Println("ERROR #32 : ", err.Error())
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		c.Writer.Write(marshaledData)
	} else {
		datas, err := model.GetPostByPostID(postID)
		datas[0].Tag = strings.ToUpper(datas[0].Tag)
		datas[0].Text = strings.ReplaceAll(datas[0].Text, `\'`, `'`)
		if err != nil {
			fmt.Println("ERROR #31 : ", err.Error())
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		marshaledData, err := json.Marshal(datas)
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}
		c.Writer.Write(marshaledData)
	}
}

func GetThumbnailHandler(c *gin.Context) {
	imgName := c.Param("name")
	imgName = strings.ReplaceAll(imgName, " ", "")
	file, err := os.Open("assets/" + imgName)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
