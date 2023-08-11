package model

import (
	"database/sql"
	"strconv"
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
	PostId    string `json:"postid"`
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

func UpdatePostImagePath (recentID int, imagename string) error {
	_, err := db.Query(`UPDATE post SET imgpath = "` + strconv.Itoa(recentID) + "-" + imagename + `" where id = ` + strconv.Itoa(recentID))
	return err
}
