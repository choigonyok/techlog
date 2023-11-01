package model

type Reply struct {
	ID        int
	Admin     bool
	WriterID  string
	WriterPW  string
	Text      string
	CommentID int
}
