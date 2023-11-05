package model

type Comment struct {
	ID       int    `json:"id"`
	Text     string `json:"text"`
	WriterID string `json:"writerID"`
	WriterPW string `json:"writerPW"`
	Admin    bool   `json:"admin"`
	PostID   int    `json:"postID"`
}

type CommentPassword struct {
	Password string `json:"password"`
}
