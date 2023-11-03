package model

import (
	"time"
)

type Post struct {
	ID            int       `json:"id"`
	Tags          string    `json:"tags"`
	Title         string    `json:"title"`
	WriteTime     time.Time `json:"writeTime"`
	ThumbnailPath string    `json:"thumbnailPath"`
	Text          string    `json:"text"`
}

type PostTags struct {
	Tags string `json:"tags"`
}

type PostCard struct {
	ID            int       `json:"id"`
	Tags          string    `json:"tags"`
	Title         string    `json:"title"`
	WriteTime     time.Time `json:"writeTime"`
	ThumbnailPath string    `json:"thumbnailPath"`
}
