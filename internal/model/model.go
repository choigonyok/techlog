package model

import (
	"database/sql"
	"strconv"
	"strings"
	"time"
)

type TagData struct {
	Tags string `json:"Tag"`
}
type IdData struct {
	Id int
}
type TagButtonData struct {
	Tagname string
}
type SendData struct {
	Id        int
	Tag       string
	Title     string
	Body      string
	Datetime  string
	ImagePath string
	Comments  []string
	WriterID  []string
	WriterPW  []string
}
type RecieveData struct {
	Body     string `json:"body"`
	Datetime time.Time `json:"datetime"`
	Tag      string `json:"tag"`
	Title    string `json:"title"`
}
type LoginData struct {
	Id       string `json:"id"`
	Password string `json:"pw"`
}
type CommentData struct {
	Comments  string `json:"comments"`
	PostId    int `json:"postid"`
	CommentID string `json:"comid"`
	CommentPW string `json:"compw"`
	IsAdmin   int    `json:"isadmin"`
	ID        int    `json:"uniqueid"`
}
type ReplyData struct {
	Reply         string `json:"comments"`
	CommentID     string `json:"commentid"`
	ReplyID       string `json:"comid"`
	ReplyPW       string `json:"compw"`
	ReplyIsAdmin  int    `json:"replyisadmin"`
	ReplyUniqueID int    `json:"replyuniqueid"`
}

var db *sql.DB

func OpenDB(driverName, dataSourceName string) error {
	database, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return err
	}

	db = database

	// DB와 서버가 연결 되었는지 확인
	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}

func CloseDB() error {
	err := db.Close()
	return err
}

func GetRecentPostID() (int, error){
	var idnum int
r, err := db.Query("SELECT id FROM post order by id desc limit 1")
if err != nil {
	return 0, err
}
			for r.Next() {
				r.Scan(
					&idnum)
			}
return idnum, nil
		}
func AddPost(postID int, tag, title, body string, datetime time.Time ) error {
	_, err := db.Query(`INSERT INTO post (id, tag, datetime, title, body) values (` + strconv.Itoa(postID) + `, '` + tag + `','` + datetime.Format("2006-01-02") + `','` + title + `','` + body + `')`)
	return err
}

func UpdatePostImagePath(recentID int, imagename string) error {
	_, err := db.Query(`UPDATE post SET imgpath = "` + strconv.Itoa(recentID) + "-" + imagename + `" where id = ` + strconv.Itoa(recentID))
	return err
}

func UpdatePost(title, body, tag, postID string, datetime time.Time) error {
	_, err := db.Query(`UPDATE post SET title = '` + title + `',body = '` + body + `',tag='` + tag + `',datetime='` + datetime.Format("2006-01-02") + `' where id = ` + postID)
	return err
}

func SelectPostByTag(tag string) ([]SendData, error){
	r, err := db.Query("SELECT id,tag,title,body,datetime,imgpath FROM post where tag LIKE '%" + tag + "%' order by datetime desc")
	if err != nil {
		return nil, err
	}
	var data SendData
	var datas []SendData
	for r.Next() {
		r.Scan(&data.Id, &data.Tag, &data.Title, &data.Body, &data.Datetime, &data.ImagePath)
		data.Datetime = strings.TrimSuffix(data.Datetime, " 00:00:00")
		data.Datetime = strings.ReplaceAll(data.Datetime, "-", "/")
		data.Tag = strings.ToUpper(data.Tag)
		data.Title = strings.ToUpper(data.Title)		
		datas = append(datas, data)
	}
	return datas, nil
}

func GetEveryTagAsString() (string, error) {
	r, err := db.Query("SELECT tag FROM post group by tag")
	if err != nil {
		return "", err
	}
	tagdata := TagButtonData{}
	sum := ""
	for r.Next() {
		r.Scan(&tagdata.Tagname)
		sum += " " + tagdata.Tagname
	}
	return sum, nil
}

func DeleteRecentPost() error {
	_, err := db.Query("DELETE FROM post ORDER BY id DESC LIMIT 1")
	return err
}

func DeletePostByPostID(postID string) error {
	_, err := db.Query("DELETE FROM post WHERE id = " + postID)
	return err
}

func SelectEveryCommentIDByPostID(postID string) ([]string, error){
	r, err := db.Query("SELECT uniqueid FROM comments WHERE id = " + postID)
	if err != nil {
		return nil, err
	}
	var commentsSlice []string
	var tempString string
	for r.Next() {
		r.Scan(&tempString)
		commentsSlice = append(commentsSlice, tempString)
	}
	return commentsSlice, nil
}

func DeleteEveryCommentByCommentID(commentID string) error {
	_, err := db.Query("DELETE FROM comments WHERE uniqueid = " + commentID)
		return err
}

func DeleteEveryReplyByCommentID(commentID string) error {
	_, err := db.Query("DELETE FROM reply WHERE commentid = " + commentID)
	return err
}

func GetRecentCommentID() (int, error) {
	r, err := db.Query("SELECT uniqueid FROM comments order by uniqueid desc limit 1")
	if err != nil {
		return 0, err
	}
	var recentCommentID int
	for r.Next() {
		r.Scan(&recentCommentID)
	}
	return recentCommentID, nil
}

func InsertComment(postID, commentID,isAdmin int, commentText, writerID, writerPW string) error {
	_, err := db.Query(`INSERT INTO comments(id, contents, writerid, writerpw, isadmin, uniqueid) values (` + strconv.Itoa(postID)	+ `,'` + commentText + `','` + writerID + `','` + writerPW + `',` + strconv.Itoa(isAdmin) + `,` + strconv.Itoa(commentID) + `)`)
	return err
}

func GetCommentWriterPWByCommentID(commentID string) (string, error) {
	r, err := db.Query("SELECT writerpw FROM comments WHERE uniqueid =" + commentID)
	var writerPW string
	// r.Next()
	err = r.Scan(&writerPW)
	return writerPW, err
}

func SelectNotAdminWriterComment(postID int) ([]CommentData, error){
	r, err := db.Query(`SELECT writerid, writerpw, contents, isadmin, uniqueid FROM comments WHERE isadmin != 1`)
	if err != nil {
		return nil, err
	}
	datas := []CommentData{}
	data := CommentData{}
	for r.Next() {
		r.Scan(&data.CommentID, &data.CommentPW, &data.Comments, &data.IsAdmin, &data.ID)
		data.PostId = postID
		datas = append(datas, data)
	}
	return datas, nil
}

func SelectCommentByPostID(postID int) ([]CommentData, error) {
	r, err := db.Query(`SELECT writerid, writerpw, contents, isadmin, uniqueid FROM comments WHERE id = ` + strconv.Itoa(postID))
	if err != nil {
		return nil, err
	}
	datas := []CommentData{}
	data := CommentData{}
	for r.Next() {
		r.Scan(&data.CommentID, &data.CommentPW,&data.Comments, &data.IsAdmin, &data.ID)
		data.PostId = postID
		datas = append(datas, data)
	}
	return datas, nil
}

func SelectReplyByCommentID(commendID string) ([]ReplyData, error) {
	r, err := db.Query("SELECT replyuniqueid, replyisadmin, replywriterid, replywriterpw, replycontents FROM reply WHERE commentid = " + commentID + " order by replyuniqueid asc")
	if err != nil {
		return nil, err
	}
	datas := []ReplyData{}
	data := ReplyData{}
	for r.Next() {
		r.Scan(&data.ReplyUniqueID, &data.ReplyID, &data.ReplyPW, &data.Reply)
		datas = append(datas, data)
	}
	return datas, nil
}

func GetRecentReplyID() (int, error){
	r, err := db.Query("SELECT replyuniqueid FROM reply order by replyuniqueid desc limit 1")
	if err != nil {
		return 0, err
	}
	var recentReplyID int
	for r.Next() {
		r.Scan(&recentReplyID)
	}
	return recentReplyID, nil
}

func InsertReply(isAdmin, recentReplyID int, commentID, replyText, writerID, writerPW string) error {
	_, err := db.Query(`INSERT INTO reply (commentid, replycontents, replywriterid, replywriterpw, replyisadmin, replyuniqueid) values (` + strconv.Itoa(commentID) + `,'` + replyText + `','` + writerID + `','` + writerPW + `',` + strconv.Itoa(isAdmin) + `,` + strconv.Itoa(recentReplyID) + `)`)
	return err
}
