package model

type Reply struct {
	ID        int    `json:"id"`
	Admin     string `json:"admin"`
	WriterID  string `json:"writerID"`
	WriterPW  string `json:"writerPW"`
	Text      string `json:"text"`
	CommentID int    `json:"commentID"`
	PostID    int    `json:"postID"`
}
