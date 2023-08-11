package controller

import (
	"database/sql"
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
	_, err := c.Cookie("visitor")
	if err == http.ErrNoCookie {
		visitnum += 1
		totalnum += 1
		c.SetCookie("visitor", "OK", 60, "/", "", false, true)
	}
	temp := struct {
		VisitNumber int
		TotalNumber int
	}{
		VisitNumber: visitnum,
		TotalNumber: totalnum,
	}
	data, err := json.Marshal(temp)
	if err != nil {
		fmt.Println("VISITOR NUM MARSHALING ERROR")
	}
	c.Writer.Header().Set("Content-type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(data)
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
			_, err = db.Query(`UPDATE post SET title = '` + recieved_data.Title + `',body = '` + recieved_data.Body + `',tag='` + recieved_data.Tag + `',datetime='` + recieved_data.Datetime + `' where id = ` + param)
			if err != nil {
				fmt.Println("MODIFY POST : DB READING ERROR")
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
	main_tag, _, _ := strings.Cut(data.Tags, " ")
	r, err := db.Query("SELECT id,tag,title,body,datetime,imgpath FROM post where tag LIKE '%" + main_tag + "%' order by datetime desc")
	if err != nil {
		log.Fatalln("TAG FINDING ERROR!!")
	}
	postdata := []model.SendData{}
	var temp model.SendData
	for r.Next() {
		r.Scan(&temp.Id, &temp.Tag, &temp.Title, &temp.Body, &temp.Datetime, &temp.ImagePath)
		temp.Datetime = strings.TrimSuffix(temp.Datetime, " 00:00:00")
		temp.Datetime = strings.ReplaceAll(temp.Datetime, "-", "/")
		temp.Tag = strings.ToUpper(temp.Tag)
		temp.Title = strings.ToUpper(temp.Title)
		postdata = append(postdata, temp)
	}

	// JSON 응답 생성
	response, err := json.Marshal(postdata)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// JSON 헤더 설정 및 응답 전송

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Write(response)

}

func GetEveryTagHandler(c *gin.Context) {
	r, err := db.Query("SELECT tag FROM post group by tag")
	if err != nil {
		log.Fatalln("TAG BUTTON FINDING ERROR!!")
	}
	tagdata := model.TagButtonData{}
	sum := ""
	for r.Next() {
		r.Scan(&tagdata.Tagname)
		sum += " " + tagdata.Tagname
	}
	sum = strings.ReplaceAll(sum, " / ", " ")
	_, sum, ok := strings.Cut(sum, " ")
	tagnum := strings.Count(sum, " ") //모든 포스트의 총 tag 합계-중복포함
	posttagdata := []model.TagButtonData{}
	temp := model.TagButtonData{}

	if !ok {
		fmt.Println("STRING ERROR 1")
	}
	for i := 0; i < tagnum; i++ {
		b, a, ok := strings.Cut(sum, " ")
		if !ok {
			fmt.Println("TAG COUNT ERROR OCCURED")
		}
		temp.Tagname = strings.ToUpper(b)
		sum = a
		posttagdata = append(posttagdata, temp)
	}
	temp.Tagname = strings.ToUpper(sum)
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
			_, err = db.Query("DELETE FROM post ORDER BY id DESC LIMIT 1")
			if err != nil {
				fmt.Println("DELETE NOT HAVE IMG POST ERROR")
			} // 이미지 없이 작성된 글 삭제
		} else {
			_, err = db.Query("DELETE FROM post WHERE id = " + deleteid)
			if err != nil {
				fmt.Println("DELETE POST ERROR")
			} // DB에서 레코드 먼저 지우기
			list, err := os.ReadDir("IMAGES")
			if err != nil {
				fmt.Println("OPENING IMG LIST ERROR")
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
			r, err := db.Query("SELECT uniqueid FROM comments WHERE id = " + deleteid)
			if err != nil {
				fmt.Println("READING COMMENT'S ID TO DELETE REPLY ERROR")
			}
			comment_id := ""
			for r.Next() {
				r.Scan(&comment_id)
			}
			_, err = db.Query("DELETE FROM comments WHERE id = " + deleteid)
			if err != nil {
				fmt.Println("DELETE COMMENTS INCLUDED POST ERROR")
			}
			_, err = db.Query("DELETE FROM reply WHERE commentid = " + comment_id)
			if err != nil {
				fmt.Println("DELETE COMMENTS INCLUDED POST ERROR")
			}
		}
		c.Writer.WriteHeader(http.StatusOK)
	}
}

func AddCommentHandler(c *gin.Context) {
	data := model.CommentData{}
	err = c.ShouldBindJSON(&data)
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
				var com_id int
				r, err := db.Query("SELECT uniqueid FROM comments order by uniqueid desc limit 1")
				if err != nil {
					fmt.Println("FINDING COMMENT'S UNIQUE ID IN DB")
				}
				for r.Next() {
					r.Scan(&com_id)
				}
				com_id += 1
				_, err = db.Query(`INSERT INTO comments(id, contents, writerid, writerpw, isadmin, uniqueid) values (` + data.PostId + `,'` + data.Comments + `','` + data.CommentID + `','` + data.CommentPW + `',` + strconv.Itoa(data.IsAdmin) + `,` + strconv.Itoa(com_id) + `)`)
				if err != nil {
					fmt.Println("DB COMMENTS INPUT ERROR")
					c.Writer.WriteHeader(http.StatusInternalServerError)
				} else {
					c.Writer.WriteHeader(http.StatusOK)
				}
			}
		}
	}
}

func GetCommentPWHandler(c *gin.Context) {
	param := c.Param("uniqueid")
	r, err := db.Query("SELECT writerpw FROM comments WHERE uniqueid =" + param)
	if err != nil {
		fmt.Println("FAILED TO LOAD COMMENT'S PASSWORD")
	}
	temp := ""
	for r.Next() {
		r.Scan(&temp)
	}
	temp_struct := struct {
		ComPW string `json:"compw"`
	}{
		temp,
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
	_, err = c.Cookie("admin") // 쿠키 : admin 권한 확인
	if err == http.ErrNoCookie {
		c.String(http.StatusUnauthorized, "THERE IS NO COOKIE")
	} else {
		_, err = db.Query("DELETE FROM reply WHERE commentid = " + param)
		if err != nil {
			fmt.Println("INCLUDING REPLY DELETE ERROR")
		}
		_, err = db.Query("DELETE FROM comments WHERE uniqueid = " + param)
		if err != nil {
			fmt.Println("COMMENT DELETE ERROR")
		} else {
			c.Writer.WriteHeader(http.StatusOK)
		}
	}
}

func DeleteCommentHandler(c *gin.Context) {
	comid := c.Query("comid")
	inputpw := c.Query("inputpw")
	writerpw := ""
	r, err := db.Query("SELECT writerpw FROM comments WHERE uniqueid = " + comid)
	if err != nil {
		fmt.Println("WRITERPW LOADING FROM DB ERROR")
	}
	for r.Next() {
		r.Scan(&writerpw)
	}
	if inputpw == writerpw {
		_, err = db.Query("DELETE FROM reply WHERE commentid = " + comid)
		if err != nil {
			fmt.Println("INCLUDING REPLY DELETE ERROR")
		}
		_, err = db.Query("DELETE FROM comments WHERE uniqueid = " + comid)
		if err != nil {
			fmt.Println("COMMENT DELETE ERROR")
		}
		c.Writer.WriteHeader(http.StatusOK)
	} else {
		c.Writer.WriteHeader(http.StatusBadRequest)
	}
}

func GetCommentHandler(c *gin.Context) {
	param := c.Param("postid")
	data := []model.CommentData{}
	temp := model.CommentData{}
	if param == "0" { // 전체 댓글 중 admin 댓글이 아닌 것
		r, err := db.Query(`SELECT writerid, writerpw, contents, isadmin, uniqueid FROM comments WHERE isadmin != 1`)
		if err != nil {
			fmt.Println("DB COMMENTS INPUT ERROR")
		}
		for r.Next() {
			r.Scan(&temp.CommentID, &temp.CommentPW, &temp.Comments, &temp.IsAdmin, &temp.ID)
			temp.PostId = param
			data = append(data, temp)
		}
	} else { // 해당 post 댓글
		r, err := db.Query(`SELECT writerid, writerpw, contents, isadmin, uniqueid FROM comments WHERE id = ` + param)
		if err != nil {
			fmt.Println("DB COMMENTS INPUT ERROR")
		}
		for r.Next() {
			r.Scan(&temp.CommentID, &temp.CommentPW, &temp.Comments, &temp.IsAdmin, &temp.ID)
			temp.PostId = param
			data = append(data, temp)
		}
	}
	real_data, err := json.Marshal(data)
	if err != nil {
		fmt.Println("COMMENTS DATA JSON BINDING ERROR")
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Write(real_data)
}

func GetReplyHandler(c *gin.Context) {
	commentid := c.Param("commentid")
	r, err := db.Query("SELECT replyuniqueid, replyisadmin, replywriterid, replywriterpw, replycontents FROM reply WHERE commentid = " + commentid + " order by replyuniqueid asc")
	if err != nil {
		fmt.Println("CAN NOT READ REPLY DATA FROM DB")
	}
	data := []model.ReplyData{}
	temp := model.ReplyData{}
	for r.Next() {
		r.Scan(&temp.ReplyUniqueID, &temp.ReplyIsAdmin, &temp.ReplyID, &temp.ReplyPW, &temp.Reply)
		data = append(data, temp)
	}
	senddata, err := json.Marshal(data)
	if err != nil {
		fmt.Println("REPLY DATA MARSHALING ERROR")
	}
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(senddata)
}

func AddReplyHandler(c *gin.Context) {
	commentid := c.Param("commentid")
	data := model.ReplyData{}
	err = c.ShouldBindJSON(&data)
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
				var reply_id int
				r, err := db.Query("SELECT replyuniqueid FROM reply order by replyuniqueid desc limit 1")
				if err != nil {
					fmt.Println("FINDING REPLY'S UNIQUE ID IN DB")
				}
				for r.Next() {
					r.Scan(&reply_id)
				}
				reply_id += 1
				_, err = db.Query(`INSERT INTO reply (commentid, replycontents, replywriterid, replywriterpw, replyisadmin, replyuniqueid) values (` + commentid + `,'` + data.Reply + `','` + data.ReplyID + `','` + data.ReplyPW + `',` + strconv.Itoa(data.ReplyIsAdmin) + `,` + strconv.Itoa(reply_id) + `)`)
				if err != nil {
					fmt.Println("DB REPLY INPUT ERROR")
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
	r, err := db.Query("SELECT replywriterpw FROM reply WHERE replyuniqueid =" + replyid)
	if err != nil {
		fmt.Println("READING REPLY UNIQUE ID ERROR")
	}
	compare_pw := ""
	for r.Next() {
		r.Scan(&compare_pw)
	}
	if compare_pw == inputpw {
		_, err = db.Query("DELETE FROM reply WHERE replyuniqueid = " + replyid)
		if err != nil {
			fmt.Println("DELETE REPLY IN DB ERROR")
		}
		c.Writer.WriteHeader(http.StatusOK)
	} else {
		c.Writer.WriteHeader(http.StatusBadRequest)
	}
}

func GetPostHandler(c *gin.Context) {
	postid := c.Param("postid")
	var r2 *sql.Rows
	if postid == "all" {
		r2, err = db.Query("SELECT id, tag,title,body,datetime,imgpath FROM post")
		if err != nil {
			log.Fatalln("ID FINDING ERROR!!")
		}
	} else {
		r2, err = db.Query("SELECT id, tag,title,body,datetime,imgpath FROM post where id = " + postid)
		if err != nil {
			log.Fatalln("ID FINDING ERROR!!")
		}
	}
	var data []model.SendData
	var temp model.SendData
	for r2.Next() {
		r2.Scan(&temp.Id, &temp.Tag, &temp.Title, &temp.Body, &temp.Datetime, &temp.ImagePath)
		temp.Datetime = strings.TrimSuffix(temp.Datetime, " 00:00:00")
		temp.Datetime = strings.ReplaceAll(temp.Datetime, "-", "/")
		temp.Tag = strings.ToUpper(temp.Tag)
		temp.Title = strings.ToUpper(temp.Title)
		data = append(data, temp)
	}
	var response []byte
	if postid == "all" {
		// JSON 응답 생성
		response, err = json.Marshal(data)
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}

	} else {
		// JSON 응답 생성
		response, err = json.Marshal(temp)
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// JSON 헤더 설정 및 응답 전송
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Write(response)
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
