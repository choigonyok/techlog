package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/choigonyok/techlog/internal/model"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository() *PostRepository {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", TMP_DB_HOST, TMP_DB_PORT, TMP_DB_USERNAME, TMP_DB_PASSWORD, TMP_POST_DB_DATABASE)

	cli, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("TEST CREATE POSTGRESQL CLIENT ERROR: ", err)
	}

	return &PostRepository{
		db: cli,
	}
}

func (repo *PostRepository) Get(column, conditionKey, conditionValue string) any {
	r := repo.db.QueryRow(`SELECT ` + column + ` FROM visitor WHERE ` + conditionKey + ` = '` + conditionValue + `'`)

	count := 0
	r.Scan(&count)

	return count
}

func (repo *PostRepository) CreatePost(post *model.Post) error {
	query := `INSERT INTO posts (id, title, text, write_time, subtitle) VALUES ($1, $2, $3, $4, $5)`
	_, err := repo.db.Exec(query, post.ID, post.Title, post.Text, post.WriteTime, post.Subtitle)
	return err
}

func (repo *PostRepository) UpdatePost(post *model.Post) error {
	query := `UPDATE posts SET title = $1, text = $2, subtitle = $3 WHERE id = $4`
	_, err := repo.db.Exec(query, post.Title, post.Text, post.Subtitle, post.ID)

	return err
}

func (repo *PostRepository) CreateTags(tags []string, postId string) error {
	for _, tag := range tags {
		if _, err := repo.db.Exec(`INSERT INTO tags (name, post_id) VALUES ('` + tag + `', '` + postId + `')`); err != nil {
			return err
		}

	}
	return nil
}

func (repo *PostRepository) CreateImage(image *model.Image) error {
	_, err := repo.db.Exec(`INSERT INTO images (id, post_id, name) VALUES ('` + image.ID + `', '` + image.PostID + `', '` + image.Name + `')`)
	return err
}

func (repo *PostRepository) GetPostsWithTag(tag string) ([]*model.Post, error) {
	tags := strings.Join(strings.Split(tag, ","), "','")

	r, err := repo.db.Query(`SELECT post_id FROM tags WHERE name IN ('` + tags + `')`)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	ids := []string{}
	for r.Next() {
		id := ""
		r.Scan(&id)
		ids = append(ids, id)
	}

	str := strings.Join(ids, "', '")

	r2, err := repo.db.Query(`SELECT id, title, write_time, subtitle FROM posts WHERE id IN ('` + str + `') ORDER BY id desc`)
	if err != nil {
		return nil, err
	}
	defer r2.Close()

	posts := []*model.Post{}
	for r2.Next() {
		p := model.Post{}
		r2.Scan(&p.ID, &p.Title, &p.WriteTime, &p.Subtitle)

		r, err := repo.db.Query(`SELECT name FROM tags WHERE post_id = '` + p.ID + `'`)
		if err != nil {
			return nil, err
		}

		tags := []string{}
		for r.Next() {
			tag := ""
			r.Scan(&tag)
			tags = append(tags, tag)
		}

		r.Close()
		p.Tags = tags
		posts = append(posts, &p)
	}

	return posts, nil
}

func (repo *PostRepository) GetPosts() ([]*model.Post, error) {
	r, err := repo.db.Query(`SELECT id, title, write_time, subtitle FROM posts ORDER BY id desc`)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	posts := []*model.Post{}
	for r.Next() {
		p := model.Post{}
		r.Scan(&p.ID, &p.Title, &p.WriteTime, &p.Subtitle)

		tags, err := repo.GetTagsByPostId(p.ID)
		if err != nil {
			return nil, err
		}
		p.Tags = *tags
		posts = append(posts, &p)
	}

	return posts, nil
}

func (repo *PostRepository) GetTags() (*[]string, error) {
	r, err := repo.db.Query(`SELECT name FROM tags`)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	tags := []string{}
	for r.Next() {
		tag := ""
		r.Scan(&tag)
		tags = append(tags, tag)
	}

	return &tags, nil
}

func (repo *PostRepository) GetTagsByPostId(postId string) (*[]string, error) {
	r, err := repo.db.Query(`SELECT name FROM tags WHERE post_id = '` + postId + `'`)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	tags := []string{}
	for r.Next() {
		tag := ""
		r.Scan(&tag)
		tags = append(tags, tag)
	}

	return &tags, nil
}

func (repo *PostRepository) GetPost(postId string) (*model.Post, error) {
	r := repo.db.QueryRow(`SELECT id, title, write_time, subtitle, text FROM posts WHERE id = '` + postId + `'`)

	p := model.Post{}
	r.Scan(&p.ID, &p.Title, &p.WriteTime, &p.Subtitle, &p.Text)

	return &p, nil
}

func (repo *PostRepository) GetThumbnailId(postId string) string {
	r := repo.db.QueryRow(`SELECT image_id FROM thumbnails WHERE post_id = '` + postId + `'`)
	id := ""
	r.Scan(&id)

	return id
}

func (repo *PostRepository) GetImage(imageId string) *model.Image {
	r := repo.db.QueryRow(`SELECT name, post_id FROM images WHERE id = '` + imageId + `'`)
	img := model.Image{}
	r.Scan(&img.Name, &img.PostID)
	img.ID = imageId

	return &img
}

func (repo *PostRepository) GetImages(postId string) ([]*model.Image, error) {
	r, err := repo.db.Query(`SELECT id, name FROM images WHERE post_id = '` + postId + `'`)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	images := []*model.Image{}
	for r.Next() {
		image := model.Image{}
		r.Scan(&image.ID, &image.Name)
		image.PostID = postId
		images = append(images, &image)
	}

	return images, nil
}

func (repo *PostRepository) DeleteThumbnail(postId string) error {
	_, err := repo.db.Exec(`DELETE FROM thumbnails WHERE post_id = '` + postId + `'`)
	return err
}

func (repo *PostRepository) DeleteTags(postId string) error {
	_, err := repo.db.Exec(`DELETE FROM tags WHERE post_id = '` + postId + `'`)
	return err
}

func (repo *PostRepository) DeleteImages(postId string) error {
	_, err := repo.db.Exec(`DELETE FROM images WHERE post_id = '` + postId + `'`)
	return err
}

func (repo *PostRepository) DeleteImagesWithException(postId string, id []string) error {
	_, err := repo.db.Exec(`DELETE FROM images WHERE post_id = '` + postId + `' and id NOT IN ('` + strings.Join(id, "','") + `')`)
	return err
}

func (repo *PostRepository) DeletePost(postId string) error {
	_, err := repo.db.Exec(`DELETE FROM posts WHERE id = '` + postId + `'`)
	return err
}

func (repo *PostRepository) CreateThumbnail(image *model.Image) error {
	_, err := repo.db.Exec(`INSERT INTO thumbnails (image_id, post_id) VALUES ('` + image.ID + `', '` + image.PostID + `')`)
	return err
}

func (repo *PostRepository) GetCounts(tag string) int {
	r := repo.db.QueryRow(`SELECT COUNT(*) FROM tags WHERE name = '` + tag + `'`)
	counts := 0
	r.Scan(&counts)
	return counts
}

func (repo *PostRepository) GetTotalCounts() int {
	r := repo.db.QueryRow(`SELECT COUNT(*) FROM posts`)
	counts := 0
	r.Scan(&counts)
	return counts
}
