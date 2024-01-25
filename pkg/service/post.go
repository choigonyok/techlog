package service

import (
	"strconv"
	"strings"

	"github.com/choigonyok/techlog/pkg/data"
	"github.com/choigonyok/techlog/pkg/model"
)

// GetTags returns every tag in blog without redundancy and number of posts each tag contains
func (svc *Service) GetTagsAndPostNum() ([]model.PostTag, error) {
	results := []model.PostTag{}
	m := make(map[string]bool)

	tags, err := svc.provider.GetTags()

	for _, v := range tags {
		separateTags := strings.Split(v.Tags, " ")
		for _, tag := range separateTags {
			pureTag := strings.ReplaceAll(tag, " ", "")
			if !m[pureTag] {
				num := svc.provider.GetNumberOfPostsByTag(pureTag)
				temp := model.PostTag{
					Tag: pureTag,
					Num: strconv.Itoa(num),
				}
				results = append(results, temp)
				m[tag] = true
			}
		}
	}
	return results, err
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

func (svc *Service) GetImageNamesByPostID(postID string) ([]string, error) {
	return svc.provider.GetImageNamesByPostID(postID)
}

func (svc *Service) DeletePostImagesByPostID(postID string) error {
	return svc.provider.DeletePostImagesByPostID(postID)
}

// GetThumbnailNameByPostID returns post's thumbnail image file name from database
func (svc *Service) GetThumbnailNameByPostID(postID string) (string, error) {
	return svc.provider.GetThumbnailNameByPostID(postID)
}

func (svc *Service) GetPostImageNameByImageID(imageID string) (string, error) {
	return svc.provider.GetPostImageNameByImageID(imageID)
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

func (svc *Service) DeleteImagesByImageName(name string) error {
	return svc.provider.DeleteImagesByImageName(name)
}

func (svc *Service) StoreInitialPost(post model.Post) error {
	post = data.MarshalPostToDatabaseFmt(post)
	return svc.provider.StoreInitialPost(post)
}

func (svc *Service) IsDatabaseEmpty() bool {
	return svc.provider.IsDatabaseEmpty()
}

func (svc *Service) StoreInitialPostImages(post model.Post) error {
	if post.ThumbnailPath == "" {
		return nil
	}

	imageNames := strings.Split(post.ThumbnailPath, " ")
	image := model.PostImage{}

	for i, imageName := range imageNames {
		image.PostID = post.ID
		image.ImageName = imageName
		if i == 0 {
			image.Thumbnail = "1"
		} else {
			image.Thumbnail = "0"
		}

		if err := svc.provider.StoreImage(image); err != nil {
			return err
		}
	}

	return nil
}

func (svc *Service) GetImagesByPostID(postID string) ([]model.PostImage, error) {
	return svc.provider.GetImagesByPostID(postID)
}

func (svc *Service) DeleteImagesByPostID(postID string) error {
	return svc.provider.DeleteImagesByPostID(postID)
}

func (svc *Service) CreatePostImagesByPostID(postID string, images []model.PostImage) error {
	for _, image := range images {
		if err := svc.provider.CreatePostImageByPostID(postID, image); err != nil {
			return err
		}
	}
	return nil
}
