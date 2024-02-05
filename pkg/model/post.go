package model

type Post struct {
	ID            int    `json:"id"`
	Tags          string `json:"tags"`
	Title         string `json:"title"`
	Subtitle      string `json:"subtitle"`
	WriteTime     string `json:"writeTime"`
	ThumbnailPath string `json:"thumbnailPath"`
	Text          string `json:"text"`
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

type PostImage struct {
	ID        int    `json:"id"`
	PostID    int    `json:"postid"`
	ImageName string `json:"imageName"`
	Thumbnail string `json:"thumbnail"`
}
