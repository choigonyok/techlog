package service

import (
	"strings"

	"github.com/choigonyok/techlog/pkg/data"
	"github.com/choigonyok/techlog/pkg/model"
	"github.com/choigonyok/techlog/pkg/time"
)

// GetTags returns every tag in blog without redundancy
func (svc *Service) GetTags() ([]string, error) {
	result := []string{}
	tags, err := svc.provider.GetTags()
	m := make(map[string]bool)

	for _, v := range tags {
		separateTags := strings.Split(v.Tags, " ")
		for _, tag := range separateTags {
			pureTag := strings.ReplaceAll(tag, " ", "")
			if !m[pureTag] {
				result = append(result, pureTag)
				m[tag] = true
			}
		}
	}
	return result, err
}

// GetPostsByTag returns posts having same tag with input tag parameter
func (svc *Service) GetPostsByTag(tag string) ([]model.PostCard, error) {
	switch tag {
	case "ALL":
		return svc.provider.GetPosts()
	default:
		return svc.provider.GetPostsByTag(tag)
	}
}

// GetPosts returns every post's data except image data
func (svc *Service) GetPosts() ([]model.PostCard, error) {
	return svc.provider.GetPosts()
}

// GetPostByID returns a post data from database by post id
func (svc *Service) GetPostByID(postID string) (model.Post, error) {
	post, err := svc.provider.GetPostByID(postID)
	return data.UnMarshalPostDatabaseFmt(post), err
}

// GetThumbnailNameByPostID returns post's thumbnail image file name from database
func (svc *Service) GetThumbnailNameByPostID(postID string) (string, error) {
	return svc.provider.GetThumbnailNameByPostID(postID)
}

// UpdatePost updates post data in to database
func (svc *Service) UpdatePost(post model.Post) error {
	return svc.provider.UpdatePost(post)
}

// DeletePostByPostID removes stored post data from database by post id
func (svc *Service) DeletePostByPostID(postID string) ([]string, error) {
	return svc.provider.DeletePostByPostID(postID)
}

// CreatePost stores post's text datas in to database
func (svc *Service) CreatePost(post model.Post) (int, error) {
	post = data.MarshalPostToDatabaseFmt(post)
	post.WriteTime = time.GetCurrentTimeByFormat("2006-01-02")
	return svc.provider.CreatePost(post)
}

// StoreImage stores post's image data in to database
func (svc *Service) StoreImage(image model.PostImage) error {
	switch image.Thumbnail {
	case "true":
		image.Thumbnail = "1"
	case "false":
		image.Thumbnail = "0"
	}

	return svc.provider.StoreImage(image)
}

func (svc *Service) StoreInitialPosts(post model.Post) error {
	post = data.MarshalPostToDatabaseFmt(post)
	return svc.provider.StoreInitialPosts(post)
}
