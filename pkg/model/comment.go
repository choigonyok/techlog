package model

type Comment struct {
	ID       int
	Text     string
	WriterID string
	WriterPW string
	Admin    bool
	PostID   int
}
