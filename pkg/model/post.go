package model

import (
	"time"
)

type Post struct {
	ID        int
	Tags      []string
	Title     string
	Text      string
	WriteTime time.Time
	ImagePath string
	ImageID   int
}
