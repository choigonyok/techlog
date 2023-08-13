package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/choigonyok/blog-project-backend/internal/model"
	"github.com/gin-gonic/gin"
)

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

func CheckIDAndPW(c *gin.Context) {
	data := model.LoginData{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		fmt.Println("ERROR #1 : ", err.Error())
	}
	if data.Id == os.Getenv("BLOG_ID") && data.Password == os.Getenv("BLOG_PW") {
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.SetCookie("admin", "authorized", 60*60*12, "/", "choigonyok.com", false, true)
		c.String(http.StatusOK, "COOKIE SENDED")
	}
}

func WritePostHandler(c *gin.Context) {
	_, err := c.Cookie("admin")
	if err == http.ErrNoCookie {
		c.String(http.StatusUnauthorized, "You are not administrator")
	} else {
		param := c.Param("param")
		switch param {
		//게시글 작성 요청
		case "post":
			var data model.RecieveData
			if err := c.ShouldBindJSON(&data); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			recentPostID, err := model.GetRecentPostID()
			if err != nil {
				fmt.Println("ERROR #2 : ", err.Error())
			}
			if strings.Contains(data.Body, `'`) {
				c.Writer.WriteHeader(http.StatusBadRequest)
			} else {
				err := model.AddPost(recentPostID+1, data.Tag, data.Title, data.Body, data.Datetime)
				if err != nil {
					fmt.Println("ERROR #3 : ", err.Error())
				}
				if err != nil {
					fmt.Println("DB POST ADD ERROR")
				}
				c.Writer.WriteHeader(http.StatusOK)
			}
		//위에서 작성된 게시글에 이미지 추가 요청
		case "img":
			imgfile, err := c.MultipartForm()
			if err != nil {
				c.String(http.StatusBadRequest, "IMG PARSING ERROR")
			}
			//방금 만들어진 post의 id 확인
			recentID, err := model.GetRecentPostID()
			if err != nil {
				fmt.Println("ERROR #4 : ", err.Error())
			}

			wholeimg := imgfile.File["file"]
			for _, v := range wholeimg {
				no_space_filename := strings.ReplaceAll(v.Filename, " ", "")
				err = c.SaveUploadedFile(v, "IMAGES/"+strconv.Itoa(recentID)+"-"+no_space_filename)
				if err != nil {
					c.String(http.StatusBadRequest, "IMG UPLOAD ERROR")
				}
			}
			//DB에는 대표이미지(썸네일)만 저장
			no_space_thimbnail := strings.ReplaceAll(wholeimg[0].Filename, " ", "")
			err = model.UpdatePostImagePath(recentID, no_space_thimbnail)
			if err != nil {
				c.String(http.StatusBadRequest, "POST IMAGE ATTACH ERROR")
			}
		}
	}
}

func GetTodayAndTotalVisitorNumHandler(c *gin.Context) {
}

func ModifyPostHandler(c *gin.Context) {
	_, err := c.Cookie("admin")
	if err == http.ErrNoCookie {
		c.String(http.StatusUnauthorized, "You are not administrator")
	} else {
		param := c.Param("param")
		var recieved_data model.RecieveData
		err = c.ShouldBindJSON(&recieved_data)
		if err != nil {
			fmt.Println("JSON BINDING ERROR")
		}
		if strings.Contains(recieved_data.Body, `'`) {
			c.Writer.WriteHeader(http.StatusBadRequest)
		} else {
			err := model.UpdatePost(recieved_data.Title, recieved_data.Body,recieved_data.Tag,param, recieved_data.Datetime)	
			if err != nil {
				fmt.Println("ERROR #5 : ", err.Error())
			}
			c.Writer.WriteHeader(http.StatusOK)
		}
	}
}

func GetPostsByTagHandler(c *gin.Context) {
	var data model.TagData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if data.Tags == "ALL" {
		data.Tags = ""
	}
	mainTag, _, _ := strings.Cut(data.Tags, " ")
	datas, err := model.SelectPostByTag(mainTag)
	if err != nil {
		fmt.Println("ERROR #6 : ", err.Error())
	}

	// JSON 응답 생성
	marshaledDatas, err := json.Marshal(datas)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// JSON 헤더 설정 및 응답 전송

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Write(marshaledDatas)
}


func GetEveryTagHandler(c *gin.Context) {
	tagString, err := model.GetEveryTagAsString()
	if err != nil {
		log.Fatalln("TAG BUTTON FINDING ERROR!!")
	}
	tagString = strings.ReplaceAll(tagString, " / ", " ")
	_, tagString, ok := strings.Cut(tagString, " ")
	tagCount := strings.Count(tagString, " ") //모든 포스트의 총 tag 합계-중복포함
	posttagdata := []model.TagButtonData{}
	temp := model.TagButtonData{}

	if !ok {
		fmt.Println("STRING ERROR 1")
	}
	for i := 0; i < tagCount; i++ {
		b, a, ok := strings.Cut(tagString, " ")
		if !ok {
			fmt.Println("TAG COUNT ERROR OCCURED")
		}
		temp.Tagname = strings.ToUpper(b)
		tagString = a
		posttagdata = append(posttagdata, temp)
	}
	temp.Tagname = strings.ToUpper(tagString)
	posttagdata = append(posttagdata, temp)
	// realdata[0] = posttagdata[0].Tagname
	ret := []model.TagButtonData{}
	m := make(map[model.TagButtonData]int)
	for i, v := range posttagdata {
		if _, ok := m[v]; !ok {
			m[v] = i
			ret = append(ret, v)
		}
	}

	// JSON 응답 생성
	response, err := json.Marshal(ret)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// JSON 헤더 설정 및 응답 전송
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Write(response)
}

func DeletePostHandler(c *gin.Context) {
	_, err := c.Cookie("admin") // 쿠키 : admin 권한 확인
	if err == http.ErrNoCookie {
		c.String(http.StatusUnauthorized, "You are not administrator")
	} else {
		deleteid := c.Param("deleteid")
		if deleteid == "0" {
			err := model.DeleteRecentPost()
			if err != nil {
				fmt.Println("ERROR #7 : ", err.Error())
			} // 이미지 없이 작성된 글 삭제
		} else {
			err := model.DeletePostByPostID(deleteid)
			if err != nil {
				fmt.Println("ERROR #8 : ", err.Error())
			} // DB에서 레코드 먼저 지우기
			list, err := os.ReadDir("assets")
			if err != nil {
				fmt.Println("ERROR #9 : ", err.Error())
			}
			// 그 다음 해당 게시글과 관계된 이미지 전체 삭제
			for _, v := range list {
				if strings.HasPrefix(v.Name(), deleteid+"-") {
					os.Remove("IMAGES/" + v.Name())
					if err != nil {
						fmt.Println("IMAGE DELETE ERROR")
					}
				}
			}
			commentsIDs, err := model.SelectEveryCommentIDByPostID(deleteid)
			if err != nil {
				fmt.Println("ERROR #10 : ", err.Error())
			}
			for _, v := range commentsIDs {
				err = model.DeleteEveryCommentByCommentID(v)
				if err != nil {
					fmt.Println("ERROR #12 : ", err.Error())
				}
				err = model.DeleteEveryReplyByCommentID(v)
				if err != nil {
					fmt.Println("ERROR #16 : ", err.Error())
				}
			}
		}
		c.Writer.WriteHeader(http.StatusOK)
	}
}

func AddCommentHandler(c *gin.Context) {
	data := model.CommentData{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		fmt.Println("COMMENTS JSON BINDING ERROR")
		c.Writer.WriteHeader(http.StatusInternalServerError)
	}
	// 빈 텍스트가 있는지 확인
	isComEmpty, err := regexp.MatchString("^$", data.Comments)
	if err != nil {
		fmt.Println("REGEXP ERROR4")
	}
	isIDEmpty, err := regexp.MatchString("^$", data.CommentID)
	if err != nil {
		fmt.Println("REGEXP ERROR5")
	}
	if isComEmpty || isIDEmpty {
		c.Writer.WriteHeader(http.StatusLengthRequired)
	} else {
		// password 숫자 맞는지 확인
		isPwNumber, err := regexp.MatchString("^[0-9]+$", data.CommentPW)
		if err != nil {
			fmt.Println("REGEXP ERROR6")
		}
		if !isPwNumber {
			c.Writer.WriteHeader(http.StatusNotAcceptable)
		} else {
			_, err := c.Cookie("admin") // 쿠키 : admin 권한 확인
			if err == http.ErrNoCookie {
				data.IsAdmin = 0
			} else {
				data.IsAdmin = 1
			}
			if strings.Contains(data.CommentID, `'`) || strings.Contains(data.Comments, `'`) {
				c.Writer.WriteHeader(http.StatusBadRequest)
			} else {
				recentCommentID, err := model.GetRecentCommentID()
				if err != nil {
					fmt.Println("ERROR #12 : ", err.Error())
				}
				err = model.InsertComment(data.PostId, recentCommentID+1, data.IsAdmin, data.Comments, data.CommentID, data.CommentPW)
				if err != nil {
					fmt.Println("ERROR #13 : ", err.Error())
					c.Writer.WriteHeader(http.StatusInternalServerError)
				} else {
					c.Writer.WriteHeader(http.StatusOK)
				}
			}
		}
	}
}

func GetCommentPWHandler(c *gin.Context) {
	commentID := c.Param("uniqueid")
	writerPW, err := model.GetCommentWriterPWByCommentID(commentID)
	if err != nil {
		fmt.Println("ERROR #14 : ", err.Error())
	}
	temp_struct := struct {
		ComPW string `json:"compw"`
	}{
		writerPW,
	}
	data, err := json.Marshal(temp_struct)
	if err != nil {
		fmt.Println("COMMENT PW MARSHALING ERROR")
	}
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(data)
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
