package model

import (
	"database/sql"
	"strconv"
	"strings"
	"time"
)

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

type CommentData struct {
	Comments  string `json:"comments"`
	PostId    int    `json:"postid"`
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

type Comment struct {
	ID       int    `json:"id"`
	PostID   int    `json:"postid"`
	Admin    int    `json:"admin"`
	Text     string `json:"text"`
	WriterID string `json:"writerid"`
	WriterPW string `json:"writerpw"`
}

type Post struct {
	ID        int       `json:"id"`
	Tag       string    `json:"tag"`
	Title     string    `json:"title"`
	Text      string    `json:"string"`
	WriteTime time.Time `json:"writetime"`
	ImagePath string    `json:"imagepath"`
	ImageNum  int       `json:"imagenum"`
}

type Reply struct {
	ID        int    `json:"id"`
	Admin     int    `json:"admin"`
	WriterID  string `json:"writerid"`
	WriterPW  string `json:"writerpw"`
	Text      string `json:"text"`
	CommentID int    `json:"commentid"`
}

type Cookie struct {
	Value string `json:"value"`
}

var db *sql.DB

func GetCookieValue(inputValue string) (string, error) {
	r, err := db.Query("SELECT value FROM cookie")
	if err != nil {
		return "", err
	}
	var cookieValue string
	r.Scan(&cookieValue)
	return cookieValue, nil
}

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

func GetRecentPostID() (int, error) {
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
func AddPost(postID int, tag, title, body string, datetime string) error {
	_, err := db.Query(`INSERT INTO post (id, tag, datetime, title, body) values (` + strconv.Itoa(postID) + `, '` + tag + `','` + datetime + `','` + title + `','` + body + `')`)
	return err
}

func UpdatePostImagePath(recentID int, imagename string) error {
	_, err := db.Query(`UPDATE post SET imgpath = "` + strconv.Itoa(recentID) + "-\"" + imagename + `" where id = ` + strconv.Itoa(recentID))
	return err
}

func UpdatePost(title, body, tag, postID string, datetime time.Time) error {
	_, err := db.Query(`UPDATE post SET title = '` + title + `',body = '` + body + `',tag='` + tag + `',datetime='` + datetime.Format("2006-01-02") + `' where id = ` + postID)
	return err
}

func SelectPostByTag(tag string) ([]SendData, error) {
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
	tagdata := Post{}
	sum := ""
	for r.Next() {
		r.Scan(&tagdata.Tag)
		sum += " " + tagdata.Tag
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

func SelectEveryCommentIDByPostID(postID string) ([]string, error) {
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

func InsertComment(postID, commentID, isAdmin int, commentText, writerID, writerPW string) error {
	_, err := db.Query(`INSERT INTO comments(id, contents, writerid, writerpw, isadmin, uniqueid) values (` + strconv.Itoa(postID) + `,'` + commentText + `','` + writerID + `','` + writerPW + `',` + strconv.Itoa(isAdmin) + `,` + strconv.Itoa(commentID) + `)`)
	return err
}

func GetCommentWriterPWByCommentID(commentID string) (string, error) {
	r, err := db.Query("SELECT writerpw FROM comments WHERE uniqueid =" + commentID)
	var writerPW string
	// r.Next()
	err = r.Scan(&writerPW)
	return writerPW, err
}

func SelectNotAdminWriterComment(postID int) ([]CommentData, error) {
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
		r.Scan(&data.CommentID, &data.CommentPW, &data.Comments, &data.IsAdmin, &data.ID)
		data.PostId = postID
		datas = append(datas, data)
	}
	return datas, nil
}

func SelectReplyByCommentID(commentID string) ([]ReplyData, error) {
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

func GetRecentReplyID() (int, error) {
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
	_, err := db.Query(`INSERT INTO reply (commentid, replycontents, replywriterid, replywriterpw, replyisadmin, replyuniqueid) values (` + commentID + `,'` + replyText + `','` + writerID + `','` + writerPW + `',` + strconv.Itoa(isAdmin) + `,` + strconv.Itoa(recentReplyID) + `)`)
	return err
}

func GetReplyPWByReplyID(replyID string) (string, error) {
	r, err := db.Query("SELECT replywriterpw FROM reply WHERE replyuniqueid =" + replyID)
	if err != nil {
		return "", err
	}
	var replyPW string
	r.Next()
	r.Scan(&replyPW)
	return replyPW, nil
}

func DeleteReplyByReplyID(replyID string) error {
	_, err := db.Query("DELETE FROM reply WHERE replyuniqueid = " + replyID)
	return err
}

func GetEveryPost() ([]SendData, error) {
	r, err := db.Query("SELECT id, tag,title,body,datetime,imgpath FROM post")
	if err != nil {
		return nil, err
	}
	var datas []SendData
	var data SendData
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
func GetPostByPostID(postID string) ([]SendData, error) {
	r, err := db.Query("SELECT id, tag,title,body,datetime,imgpath FROM post where id = " + postID)
	if err != nil {
		return nil, err
	}
	var datas []SendData
	var data SendData
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
