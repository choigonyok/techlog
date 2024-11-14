package model

type Post struct {
	ID            string   `json:"id"`
	Tags          []string `json:"tags"`
	Title         string   `json:"title"`
	Subtitle      string   `json:"subtitle"`
	WriteTime     string   `json:"writeTime"`
	ThumbnailName string   `json:"thumbnail_name"`
	Text          string   `json:"text"`
	Images        []*Image `json:"images"`
}

type PostWrite struct {
	Tags     string `json:"tags"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Text     string `json:"text"`
}

type PostTags struct {
	Tags string `json:"tags"`
}

type PostTag struct {
	Tag string `json:"tag"`
	Num string `json:"num"`
}

type PostCard struct {
	ID            int    `json:"id"`
	Tags          string `json:"tags"`
	Title         string `json:"title"`
	Subtitle      string `json:"subtitle"`
	WriteTime     string `json:"writeTime"`
	ThumbnailPath string `json:"thumbnailPath"`
}

type Image struct {
	ID     string `json:"id"`
	PostID string `json:"post_id"`
	Name   string `json:"name"`
}

type Tag struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}
