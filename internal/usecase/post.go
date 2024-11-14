package usecase

import (
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"

	"github.com/choigonyok/techlog/internal/model"
	repo "github.com/choigonyok/techlog/internal/repository"
	"github.com/choigonyok/techlog/pkg/image"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type PostUsecase struct {
	postRepository *repo.PostRepository
}

func NewPostUsecase() *PostUsecase {
	return &PostUsecase{
		postRepository: repo.NewPostRepository(),
	}
}

func (u *PostUsecase) updateMarkdownImage(post *model.Post) ([]*model.Image, map[string]bool, error) {
	re := regexp.MustCompile(`!\[.*?\]\((.*?)\)`)
	matches := re.FindAllStringSubmatch(post.Text, -1)
	images := []*model.Image{}
	notModifiedImageIds := make(map[string]bool)

	for _, m := range matches {
		imageName := m[1][strings.LastIndex(m[1], "/")+1:]
		if ok := strings.HasPrefix(m[1], "https://"+image.S3BucketName+".s3."+image.Region+".amazonaws.com/"); ok {
			notModifiedImageIds[imageName] = true
			continue
		}
		id := uuid.NewString()
		post.Text = strings.Replace(post.Text, m[0], "![IMAGE](https://"+image.S3BucketName+".s3."+image.Region+".amazonaws.com/"+id+")", -1)
		req, _ := http.NewRequest(http.MethodGet, m[1], nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, nil, err
		}
		if err := image.Upload(resp.Body, id); err != nil {
			return nil, nil, err
		}

		image := model.Image{
			ID:     id,
			PostID: post.ID,
			Name:   imageName,
		}
		images = append(images, &image)
	}

	return images, notModifiedImageIds, nil
}
func (u *PostUsecase) CreatePost(post *model.Post, images []*multipart.FileHeader) error {
	post.ID = uuid.NewString()

	imgs, _, err := u.updateMarkdownImage(post)
	if err != nil {
		return err
	}

	if err := u.postRepository.CreatePost(post); err != nil {
		return err
	}

	for _, i := range imgs {
		err = u.postRepository.CreateImage(i)
		if err != nil {
			return err
		}
	}

	if err := u.postRepository.CreateTags(post.Tags, post.ID); err != nil {
		return err
	}

	for _, img := range images {
		i := model.Image{}
		i.Name = img.Filename
		i.PostID = post.ID
		i.ID = uuid.NewString()
		file, _ := img.Open()
		file.Seek(0, 0)

		err := image.Upload(file, i.ID)
		if err != nil {
			return err
		}

		err = u.postRepository.CreateImage(&i)
		if err != nil {
			return err
		}

		if post.ThumbnailName == i.Name {
			if err := u.postRepository.CreateThumbnail(&i); err != nil {
				return err
			}
		}
	}

	return nil
}

func (u *PostUsecase) UpdateImages(images []model.Image, postId string) error {
	if err := u.postRepository.DeleteImages(postId); err != nil {
		return err
	}
	for _, image := range images {
		if err := u.postRepository.CreateImage(&image); err != nil {
			return err
		}
	}
	return nil
}

func (u *PostUsecase) UpdatePost(post *model.Post, images []*multipart.FileHeader) error {
	thumbnailId := u.postRepository.GetThumbnailId(post.ID)
	inputImages, notModified, err := u.updateMarkdownImage(post)
	if err != nil {
		return err
	}

	if post.ThumbnailName != "" {
		if err := u.postRepository.DeleteThumbnail(post.ID); err != nil {
			return err
		}
	} else {
		notModified[thumbnailId] = true
	}

	imgs, err := u.postRepository.GetImages(post.ID)
	if err != nil {
		return err
	}

	for _, i := range imgs {
		if notModified[i.ID] {
			continue
		}

		if err := image.Remove(i.ID); err != nil {
			return err
		}
	}

	notModifiedArray := []string{}
	for k, _ := range notModified {
		notModifiedArray = append(notModifiedArray, k)
	}

	if err := u.postRepository.DeleteImagesWithException(post.ID, notModifiedArray); err != nil {
		return err
	}

	for _, i := range inputImages {
		err = u.postRepository.CreateImage(i)
		if err != nil {
			return err
		}
	}

	if err := u.postRepository.UpdatePost(post); err != nil {
		return err
	}
	if err := u.postRepository.DeleteTags(post.ID); err != nil {
		return err
	}
	if err := u.postRepository.CreateTags(post.Tags, post.ID); err != nil {
		return err
	}

	for _, img := range images {
		i := model.Image{}
		i.Name = img.Filename
		i.PostID = post.ID
		i.ID = uuid.NewString()

		err = u.postRepository.CreateImage(&i)
		if err != nil {
			return err
		}
		if post.ThumbnailName == i.Name {
			file, _ := img.Open()
			file.Seek(0, 0)

			if err := image.Upload(file, i.ID); err != nil {
				return err
			}

			if err := image.Remove(thumbnailId); err != nil {
				return err
			}

			if err := u.postRepository.CreateThumbnail(&i); err != nil {
				return err
			}
		}
	}

	return nil
}

func (u *PostUsecase) GetPosts(tag string) ([]*model.Post, error) {
	if tag == "" || tag == "ALL" {
		return u.postRepository.GetPosts()
	}
	return u.postRepository.GetPostsWithTag(tag)
}

func (u *PostUsecase) GetPost(postId string) (*model.Post, error) {
	tags, err := u.postRepository.GetTagsByPostId(postId)
	if err != nil {
		return nil, err
	}
	p, err := u.postRepository.GetPost(postId)
	if err != nil {
		return nil, err
	}
	p.Tags = *tags

	return p, nil
}

func (u *PostUsecase) GetTags() ([]*model.Tag, error) {
	tags, err := u.postRepository.GetTags()
	if err != nil {
		return nil, err
	}

	m := make(map[string]int)

	for _, tag := range *tags {
		if _, ok := m[tag]; !ok {
			counts := u.postRepository.GetCounts(tag)
			m[tag] = counts
		}
	}
	m["ALL"] = u.postRepository.GetTotalCounts()

	uniqueTags := []*model.Tag{}
	for k, v := range m {
		t := model.Tag{
			Name:  k,
			Count: v,
		}
		uniqueTags = append(uniqueTags, &t)
	}

	return uniqueTags, nil
}

func (u *PostUsecase) GetThumbnail(postId string) *model.Image {
	id := u.postRepository.GetThumbnailId(postId)
	return u.postRepository.GetImage(id)
}

func (u *PostUsecase) DeletePost(postId string) error {
	images, err := u.postRepository.GetImages(postId)
	if err != nil {
		return err
	}
	for _, i := range images {
		if err := image.Remove(i.ID); err != nil {
			return err
		}
	}

	if err := u.postRepository.DeleteThumbnail(postId); err != nil {
		return err
	}
	if err := u.postRepository.DeleteTags(postId); err != nil {
		return err
	}
	if err := u.postRepository.DeleteImages(postId); err != nil {
		return err
	}
	if err := u.postRepository.DeletePost(postId); err != nil {
		return err
	}
	return nil
}

func (u *PostUsecase) GetImages(postId string) ([]*model.Image, error) {
	return u.postRepository.GetImages(postId)
}

func (u *PostUsecase) GetImage(imageId string) *model.Image {
	return u.postRepository.GetImage(imageId)
}
