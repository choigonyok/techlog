package service

import (
	"strings"

	"github.com/choigonyok/techlog/pkg/data"
	"github.com/choigonyok/techlog/pkg/model"
	"github.com/choigonyok/techlog/pkg/time"
)

func (svc *Service) GetEveryTags() ([]string, error) {
	result := []string{}
	tags, err := svc.provider.GetEveryTag()
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

func (svc *Service) GetEveryCardByTag(tag string) ([]model.PostCard, error) {
	switch tag {
	case "ALL":
		return svc.provider.GetEveryCard()
	default:
		return svc.provider.GetEveryCardByTag(tag)
	}
}

func (svc *Service) GetEveryCard() ([]model.PostCard, error) {
	return svc.provider.GetEveryCard()
}

func (svc *Service) GetPostByID(postID string) ([]model.Post, error) {
	m, _ := svc.provider.GetPostByID(postID)
	return m, nil
}

func (svc *Service) GetThumbnailNameByPostID(postID string) (string, error) {
	return svc.provider.GetThumbnailNameByPostID(postID)
}

func (svc *Service) UpdatePost(post model.Post) error {
	return svc.provider.UpdatePost(post)
}

func (svc *Service) DeletePostByPostID(postID string) ([]string, error) {
	return svc.provider.DeletePostByPostID(postID)
}

func (svc *Service) CreatePost(post model.Post) (int, error) {
	post = data.MarshalPostToDatabaseFmt(post)
	post.WriteTime = time.GetCurrentTimeByFormat("2006-01-02")
	return svc.provider.CreatePost(post)
}

func (svc *Service) StoreImage(image model.PostImage) error {
	switch image.Thumbnail {
	case "true":
		image.Thumbnail = "1"
	case "false":
		image.Thumbnail = "0"
	}

	return svc.provider.StoreImage(image)
}
