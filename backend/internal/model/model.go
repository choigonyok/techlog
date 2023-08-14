package model

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

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

type Visitor struct {
	Today int `json:"today"`
	Total int `json:"total"`
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

func UpdateCookieRecord() (uuid.UUID, error) {
	tx, err := db.Begin()
	if err != nil {
		return uuid.Nil, err
	}
	_, err =  db.Exec("DELETE * FROM cookie")
	if err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}
	cookieValue := uuid.New()
	_, err = db.Exec(`INSERT INTO cookie (value) VALUES ("`+cookieValue.String()+`")`)
	if err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}
	err = tx.Commit()
	if err != nil {
		return uuid.Nil, err
	}
	return cookieValue, nil
}

func OpenDB(driverName, dataSourceName string) error {
	fmt.Println(driverName)
	fmt.Println(dataSourceName)
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

func SelectPostByTag(tag string) ([]Post, error) {
	r, err := db.Query("SELECT id,tag,title,body,datetime,imgpath FROM post where tag LIKE '%" + tag + "%' order by datetime desc")
	if err != nil {
		return nil, err
	}
	var data Post
	var datas []Post
	for r.Next() {
		r.Scan(&data.ID, &data.Tag, &data.Title, &data.Text, &data.WriteTime, &data.ImagePath)
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
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("DELETE FROM comments WHERE uniqueid = " + commentID)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = db.Exec("DELETE FROM reply WHERE commentid = " + commentID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
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
	r.Next()
	err = r.Scan(&writerPW)
	return writerPW, err
}

func SelectNotAdminWriterComment(postID int) ([]Comment, error) {
	r, err := db.Query(`SELECT writerid, writerpw, contents, isadmin, uniqueid FROM comments WHERE isadmin != 1`)
	if err != nil {
		return nil, err
	}
	datas := []Comment{}
	data := Comment{}
	for r.Next() {
		r.Scan(&data.WriterID, &data.WriterPW, &data.Text, &data.Admin, &data.ID)
		data.PostID = postID
		datas = append(datas, data)
	}
	return datas, nil
}

func SelectCommentByPostID(postID int) ([]Comment, error) {
	r, err := db.Query(`SELECT writerid, writerpw, contents, isadmin, uniqueid FROM comments WHERE id = ` + strconv.Itoa(postID))
	if err != nil {
		return nil, err
	}
	datas := []Comment{}
	data := Comment{}
	for r.Next() {
		r.Scan(&data.WriterID, &data.WriterPW, &data.Text, &data.Admin, &data.ID)
		data.PostID = postID
		datas = append(datas, data)
	}
	return datas, nil
}

func SelectReplyByCommentID(commentID string) ([]Reply, error) {
	r, err := db.Query("SELECT replyuniqueid, replyisadmin, replywriterid, replywriterpw, replycontents FROM reply WHERE commentid = " + commentID + " order by replyuniqueid asc")
	if err != nil {
		return nil, err
	}
	var datas []Reply
	var data Reply
	for r.Next() {
		r.Scan(&data.ID, &data.Admin ,&data.WriterID, &data.WriterPW, &data.Text)
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

func GetEveryPost() ([]Post, error) {
	r, err := db.Query("SELECT id, tag,title,body,datetime,imgpath FROM post")
	if err != nil {
		return nil, err
	}
	var datas []Post
	var data Post
	for r.Next() {
		r.Scan(&data.ID, &data.Tag, &data.Title, &data.Text, &data.WriteTime, &data.ImagePath)
		datas = append(datas, data)
	}
	return datas, nil
}
func GetPostByPostID(postID string) ([]Post, error) {
	r, err := db.Query("SELECT id, tag,title,body,datetime,imgpath FROM post where id = " + postID)
	if err != nil {
		return nil, err
	}
	var datas []Post
	var data Post
	for r.Next() {
		r.Scan(&data.ID, &data.Tag, &data.Title, &data.Text, &data.WriteTime, &data.ImagePath)
		datas = append(datas, data)
	}
	return datas, nil
}

func CountTodayVisit() error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	r, err := tx.Query("SELECT today, total FROM visitor")
	if err != nil {
		tx.Rollback()
		return err
	}
	var visitor Visitor
	r.Next()
	r.Scan(&visitor.Today, &visitor.Total)
	_, err = tx.Exec(`UPDATE visitor SET today = `+strconv.Itoa(visitor.Today+1)+`, total = `+strconv.Itoa(visitor.Total+1))
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}