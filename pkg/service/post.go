package service

import (
	"strings"

	"github.com/choigonyok/techlog/pkg/model"
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

func (svc *Service) GetPostByID(postID string) ([]model.Post, error) {
	m, _ := svc.provider.GetPostByID(postID)
	return m, nil
}

// func (svc *Service) GetImageNameByPostID(postID string) ([]model.PostImageResponse, error) {
// 	images := []model.PostImageResponse{}
// 	providedImages, err := svc.provider.GetImageNameByPostID(postID)
// 	for _, v := range providedImages {
// 		if v.Thumbnail == 1 {
// 			images = append(images, model.PostImageResponse{
// 				ID:        v.ID,
// 				PostID:    v.PostID,
// 				ImageName: v.ImageName,
// 				Thumbnail: true,
// 			})
// 		} else {
// 			images = append(images, model.PostImageResponse{
// 				ID:        v.ID,
// 				PostID:    v.PostID,
// 				ImageName: v.ImageName,
// 				Thumbnail: false,
// 			})
// 		}
// 	}
// 	return images, err
// }

func (svc *Service) GetThumbnailNameByPostID(postID string) (string, error) {
	return svc.provider.GetThumbnailNameByPostID(postID)
}
