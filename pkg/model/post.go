package model

type Post struct {
	ID            int    `json:"id"`
	Tags          string `json:"tags"`
	Title         string `json:"title"`
	WriteTime     string `json:"writeTime"`
	ThumbnailPath string `json:"thumbnailPath"`
	Text          string `json:"text"`
}

type PostWrite struct {
	Tags  string `json:"tags"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

type PostTags struct {
	Tags string `json:"tags"`
}

type PostCard struct {
	ID            int    `json:"id"`
	Tags          string `json:"tags"`
	Title         string `json:"title"`
	WriteTime     string `json:"writeTime"`
	ThumbnailPath string `json:"thumbnailPath"`
}

type PostImage struct {
	ID        int
	PostID    int
	ImageName string
	Thumbnail string
}
