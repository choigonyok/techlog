package model

import (
	"database/sql"
	"fmt"
	"strconv"

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
	Text      string    `json:"text"`
	WriteTime string 	`json:"writetime"`
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
	r.Next()
	r.Scan(&cookieValue)
	return cookieValue, nil
}

func UpdateCookieRecord() (uuid.UUID, error) {
	tx, err := db.Begin()
	if err != nil {
		return uuid.Nil, err
	}
	_, err =  db.Exec("DELETE FROM cookie")
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
func AddPost(tag, title, text, writetime string) error {
	_, err := db.Exec(`INSERT INTO post (tag, writetime, title, text) values ('` + tag + `', '`+writetime+`' ,'` + title + `','` + text + `')`)
	DBTest()
	return err
}

func UpdatePostImagePath(recentID int, imagename string) error {
	_, err := db.Query(`UPDATE post SET imgpath = '` + strconv.Itoa(recentID) + `-` + imagename + `' where id = ` + strconv.Itoa(recentID))
	return err
}

func UpdatePost(title, text, tag, postID string) error {
	_, err := db.Query(`UPDATE post SET title = '` + title + `', text = '` + text + `', tag ='` + tag + `'  where id = ` + postID)
	return err
}

func SelectPostByTag(tag string) ([]Post, error) {
	r, err := db.Query("SELECT id,tag,title,text,writetime,imgpath FROM post where tag LIKE '%" + tag + "%' order by writetime desc")
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
	r, err := db.Query("SELECT id FROM comment WHERE postid = " + postID)
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
	_, err = tx.Exec("DELETE FROM comment WHERE id = " + commentID)
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
	r, err := db.Query("SELECT id FROM comment order by id desc limit 1")
	if err != nil {
		return 0, err
	}
	var recentCommentID int
	for r.Next() {
		r.Scan(&recentCommentID)
	}
	return recentCommentID, nil
}

func InsertComment(postID, id, admin int, text, writerID, writerPW string) error {
	_, err := db.Query(`INSERT INTO comment(postid, text, writerid, writerpw, admin, id) values (` + strconv.Itoa(postID) + `,'` + text + `','` + writerID + `','` + writerPW + `',` + strconv.Itoa(admin) + `,` + strconv.Itoa(id) + `)`)
	return err
}

func GetCommentWriterPWByCommentID(commentID string) (string, error) {
	r, err := db.Query("SELECT writerpw FROM comment WHERE id =" + commentID)
	var writerPW string
	r.Next()
	err = r.Scan(&writerPW)
	return writerPW, err
}

func SelectNotAdminWriterComment(postID int) ([]Comment, error) {
	r, err := db.Query(`SELECT writerid, writerpw, text, admin, id FROM comment WHERE admin != 1`)
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
	r, err := db.Query(`SELECT writerid, writerpw, text, admin, id FROM comment WHERE postid = ` + strconv.Itoa(postID))
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
	r, err := db.Query("SELECT id, admin, replywriterid, writerpw, text FROM reply WHERE commentid = " + commentID + " order by id asc")
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
	r, err := db.Query("SELECT id FROM reply order by id desc limit 1")
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
	_, err := db.Query(`INSERT INTO reply (commentid, text, writerid, writerpw, admin, id) values (` + commentID + `,'` + replyText + `','` + writerID + `','` + writerPW + `',` + strconv.Itoa(isAdmin) + `,` + strconv.Itoa(recentReplyID) + `)`)
	return err
}

func GetReplyPWByReplyID(replyID string) (string, error) {
	r, err := db.Query("SELECT writerpw FROM reply WHERE id =" + replyID)
	if err != nil {
		return "", err
	}
	var replyPW string
	r.Next()
	r.Scan(&replyPW)
	return replyPW, nil
}

func DeleteReplyByReplyID(replyID string) error {
	_, err := db.Query("DELETE FROM reply WHERE id = " + replyID)
	return err
}

func GetEveryPost() ([]Post, error) {
	r, err := db.Query("SELECT id, tag,title,text,writetime,imgpath FROM post")
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
	r, err := db.Query("SELECT id, tag,title,text,writetime,imgpath FROM post where id = " + postID)
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

func DBTest() {
	r1, _ := db.Query("SELECT id, postid, admin, text, writerid, writerpw FROM comment")
	var comment Comment
	var comments []Comment
	for r1.Next() {
		r1.Scan(&comment.ID, &comment.PostID, &comment.Admin, &comment.Text, &comment.WriterID, &comment.WriterPW)
		comments = append(comments, comment)
	}
	fmt.Println("comment: ", comments)

	r2, _ := db.Query("SELECT id, tag, title, text, writetime, imgpath, imgnum FROM post")
	var post Post
	var posts []Post
	for r2.Next() {
		r2.Scan(&post.ID, &post.Tag, &post.Title, &post.Text, &post.WriteTime, &post.ImagePath, &post.ImageNum)
		posts = append(posts, post)
	}
	fmt.Println("post: ", posts)

	r3, _ := db.Query("SELECT id, admin, writerid, writerpw, text, commentid FROM reply")
	var reply Reply
	var replys []Reply
	for r3.Next() {
		r3.Scan(&reply.ID, &reply.Admin, &reply.WriterID, &reply.WriterPW, &reply.Text, &reply.CommentID)
		replys = append(replys, reply)
	}
	fmt.Println("reply: ", replys)

	r4, _ := db.Query("SELECT value FROM cookie")
	var cookie Cookie
	r4.Next()
	r4.Scan(&cookie.Value)
	fmt.Println("cookie: ", cookie.Value)

	r5, _ := db.Query("SELECT today, total FROM visitor")
	var visitor Visitor
	var visitors []Visitor
	for r5.Next() {
		r5.Scan(&visitor.Today, &visitor.Total)
		visitors = append(visitors, visitor)
	}
	fmt.Println("visitor: ", visitors)
}