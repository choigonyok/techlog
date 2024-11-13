package usecase

import (
	"mime/multipart"

	repo "github.com/choigonyok/techlog/internal/repository"
	"github.com/choigonyok/techlog/pkg/model"
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

func (u *PostUsecase) CreatePost(post *model.Post, images []*multipart.FileHeader) error {
	post.ID = uuid.NewString()
	if err := u.postRepository.CreatePost(post); err != nil {
		return err
	}
	if err := u.postRepository.CreateTags(&post.Tags, post.ID); err != nil {
		return err
	}

	for _, image := range images {
		i := model.Image{}
		i.Name = image.Filename
		i.PostID = post.ID
		i.ID = uuid.NewString()

		// err = image.Upload(v, image.ImageName)
		// if err != nil {
		// 	fmt.Println(err.Error())
		// }

		err := u.postRepository.CreateImage(i)
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
		if err := u.postRepository.CreateImage(image); err != nil {
			return err
		}
	}
	return nil
}

func (u *PostUsecase) UpdatePost(post *model.Post, images []*multipart.FileHeader) error {
	if err := u.postRepository.UpdatePost(post); err != nil {
		return err
	}
	if err := u.postRepository.DeleteTags(post.ID); err != nil {
		return err
	}
	if err := u.postRepository.CreateTags(&post.Tags, post.ID); err != nil {
		return err
	}
	if err := u.postRepository.DeleteThumbnail(post.ID); err != nil {
		return err
	}
	if err := u.postRepository.DeleteImages(post.ID); err != nil {
		return err
	}

	for _, image := range images {
		i := model.Image{}
		i.Name = image.Filename
		i.PostID = post.ID
		i.ID = uuid.NewString()

		// err = image.Upload(v, image.ImageName)
		// if err != nil {
		// 	fmt.Println(err.Error())
		// }
		err := u.postRepository.CreateImage(i)
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

func (u *PostUsecase) GetPosts(tag string) ([]*model.Post, error) {
	if tag == "ALL" {
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
