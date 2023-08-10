package model

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
	Id int
	Tag string
	Title string
	Body string
	Datetime string
	ImagePath string
	Comments []string
	WriterID []string
	WriterPW []string
}
type RecieveData struct {
	Body string `json:"body"`
	Datetime string `json:"datetime"`
	Tag string `json:"tag"`
	Title string `json:"title"`
}
type LoginData struct {
	Id string `json:"id"`
	Password string `json:"pw"`
}
type CommentData struct {
	Comments string `json:"comments"`
	PostId string `json:"postid"`
	CommentID string `json:"comid"`
	CommentPW string `json:"compw"`
	IsAdmin int `json:"isadmin"`
	ID int `json:"uniqueid"`
}
type ReplyData struct {
	Reply string `json:"comments"`
	CommentID string `json:"commentid"`
	ReplyID string `json:"comid"`
	ReplyPW string `json:"compw"`
	ReplyIsAdmin int `json:"replyisadmin"`
	ReplyUniqueID int `json:"replyuniqueid"`
}