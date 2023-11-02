package model

import (
	"time"
)

type Post struct {
	ID        int
	Tags      string
	Title     string
	Text      string
	WriteTime time.Time
}

type PostTags struct {
	Tags string `json:"tags"`
}

type PostCard struct {
	ID        int       `json:"id"`
	Tags      string    `json:"tags"`
	Title     string    `json:"title"`
	WriteTime time.Time `json:"writeTime"`
}
