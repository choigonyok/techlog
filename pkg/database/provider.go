package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/choigonyok/techlog/pkg/model"
)

type Provider interface {
	ResetVisitorTodayAndDate(today string) error
	GetVisitor() (model.Visitor, error)
	UpdateVisitorToday(newToday, newTotal int) error
	GetTags() ([]model.PostTags, error)
	GetPostsByTag(string) ([]model.PostCard, error)
	GetPosts() ([]model.PostCard, error)
	GetPostByID(postID string) (model.Post, error)
	GetImageNamesByPostID(postID string) ([]string, error)
	GetThumbnailNameByPostID(postID string) (string, error)
	GetPostImageNameByImageID(imageID string) (string, error)
	SetCookieValueByUniqueID(uniqueID string) error
	UpdateCookieValueByUniqueID(uniqueID string) error
	GetCookieValue() (string, error)
	UpdatePost(post model.Post) error
	DeletePostByPostID(postID string) ([]string, error)
	CreatePost(post model.Post) (int, error)
	StoreImage(image model.PostImage) error
	GetComments() ([]model.Comment, error)
	GetCommentsByPostID(postID string) ([]model.Comment, error)
	GetCommentPasswordByCommentID(commentID string) (string, error)
	DeleteCommentByCommentID(commentID string) error
	CreateComment(comment model.Comment, admin string) error
	GetRepliesByPostID(postID string) ([]model.Reply, error)
	GetReplyPasswordByReplyID(replyID string) (string, error)
	DeleteReplyByReplyID(replyID string) error
	CreateReply(reply model.Reply) error
	StoreInitialPost(post model.Post) error
	DeleteImagesByImageName(name string) error
	IsDatabaseEmpty() bool
	GetImagesByPostID(postID string) ([]model.PostImage, error)
	DeleteImagesByPostID(postID string) error
	CreatePostImageByPostID(postID string, image model.PostImage) error
	DeletePostImagesByPostID(postID string) error
}

type MysqlProvider struct {
	connector *sql.DB
}

func NewMysqlProvider(db *sql.DB) *MysqlProvider {
	return &MysqlProvider{
		connector: db,
	}
}

func (p *MysqlProvider) ResetVisitorTodayAndDate(today string) error {
	_, err := p.connector.Exec(`UPDATE visitor SET today = 1, date = "` + today + `"`)
	return err
}

func (p *MysqlProvider) GetVisitor() (model.Visitor, error) {
	result := model.Visitor{}
	r := p.connector.QueryRow(`SELECT today, total, date FROM visitor`)
	err := r.Scan(&result.Today, &result.Total, &result.Date)
	return result, err
}

func (p *MysqlProvider) UpdateVisitorToday(newToday, newTotal int) error {
	_, err := p.connector.Exec(`UPDATE visitor SET today = ` + strconv.Itoa(newToday) + `, total = ` + strconv.Itoa(newTotal))
	return err
}

func (p *MysqlProvider) GetTags() ([]model.PostTags, error) {
	tag := model.PostTags{}
	tags := []model.PostTags{}
	r, err := p.connector.Query(`SELECT tags FROM post`)
	for r.Next() {
		r.Scan(&tag.Tags)
		tags = append(tags, tag)
	}
	defer r.Close()
	return tags, err
}

func (p *MysqlProvider) GetPostsByTag(tag string) ([]model.PostCard, error) {
	card := model.PostCard{}
	cards := []model.PostCard{}
	r, err := p.connector.Query(`SELECT id, tags, title, writeTime FROM post WHERE tags LIKE "%` + tag + `%" ORDER BY id desc`)
	for r.Next() {
		r.Scan(&card.ID, &card.Tags, &card.Title, &card.WriteTime)
		cards = append(cards, card)
	}
	defer r.Close()
	return cards, err
}

func (p *MysqlProvider) GetPosts() ([]model.PostCard, error) {
	card := model.PostCard{}
	cards := []model.PostCard{}
	r, err := p.connector.Query(`SELECT id, tags, title, writeTime FROM post ORDER BY id desc`)
	for r.Next() {
		r.Scan(&card.ID, &card.Tags, &card.Title, &card.WriteTime)
		cards = append(cards, card)
	}
	defer r.Close()
	return cards, err
}

func (p *MysqlProvider) GetPostByID(postID string) (model.Post, error) {
	post := model.Post{}

	r := p.connector.QueryRow(`SELECT id, tags, title, text, writeTime FROM post WHERE id = ` + postID)
	err := r.Scan(&post.ID, &post.Tags, &post.Title, &post.Text, &post.WriteTime)

	return post, err
}

func (p *MysqlProvider) GetThumbnailNameByPostID(postID string) (string, error) {
	var thumbnailName string
	r := p.connector.QueryRow(`SELECT imageName FROM image WHERE thumbnail = 1 AND postID = ` + postID)
	err := r.Scan(&thumbnailName)
	return thumbnailName, err
}

func (p *MysqlProvider) SetCookieValueByUniqueID(uniqueID string) error {
	_, err := p.connector.Exec(`INSERT INTO cookie (value) VALUES ("` + uniqueID + `")`)
	return err
}

func (p *MysqlProvider) UpdateCookieValueByUniqueID(uniqueID string) error {
	_, err := p.connector.Exec(`UPDATE cookie SET value = "` + uniqueID + `"`)
	return err
}

func (p *MysqlProvider) GetCookieValue() (string, error) {
	var value string
	r := p.connector.QueryRow(`SELECT value FROM cookie`)
	err := r.Scan(&value)
	return value, err
}

func (p *MysqlProvider) UpdatePost(post model.Post) error {
	_, err := p.connector.Exec(`UPDATE post SET title = '` + post.Title + `', text = '` + post.Text + `', tags = '` + post.Tags + `' WHERE id = ` + strconv.Itoa(post.ID))

	return err
}

func (p *MysqlProvider) DeletePostByPostID(postID string) ([]string, error) {
	imageNames := []string{}
	var imageName string

	tx, err := p.connector.Begin()
	if err != nil {
		return nil, errors.New("FAILED TO CREATE TX")
	}
	r, err := tx.Query(`SELECT imageName FROM image WHERE postID = ` + postID)
	if err != nil {
		tx.Rollback()
	}
	for r.Next() {
		r.Scan(&imageName)
		imageNames = append(imageNames, imageName)
	}
	defer r.Close()

	_, err = tx.Exec(`DELETE FROM post WHERE id = ` + postID)
	if err != nil {
		tx.Rollback()
	}
	err = tx.Commit()
	if err != nil {
		return nil, errors.New("FAILED TO COMMIT TX")
	}

	return imageNames, err
}

func (p *MysqlProvider) CreatePost(post model.Post) (int, error) {
	var postID int

	tx, err := p.connector.Begin()
	if err != nil {
		fmt.Println("TX CREATE ERROR", err.Error())
	}
	_, err = tx.Exec(`INSERT INTO post (tags, title, text, writeTime) VALUES ('` + post.Tags + `', '` + post.Title + `', '` + post.Text + `', '` + post.WriteTime + `')`)
	if err != nil {
		tx.Rollback()
	}
	r := tx.QueryRow(`SELECT id FROM post ORDER BY id DESC LIMIT 1`)
	if err != nil {
		tx.Rollback()
	}
	err = r.Scan(&postID)
	if err != nil {
		tx.Rollback()
	}
	err = tx.Commit()
	return postID, err
}

func (p *MysqlProvider) StoreImage(image model.PostImage) error {
	_, err := p.connector.Exec(`INSERT INTO image (postID, imageName, thumbnail) VALUES (` + strconv.Itoa(image.PostID) + `, '` + image.ImageName + `', '` + image.Thumbnail + `')`)

	return err
}

func (p *MysqlProvider) GetComments() ([]model.Comment, error) {
	comment := model.Comment{}
	comments := []model.Comment{}

	r, err := p.connector.Query(`SELECT id, text, writerID, admin, postID FROM comment`)

	for r.Next() {
		r.Scan(&comment.ID, &comment.Text, &comment.WriterID, &comment.Admin, &comment.PostID)
		comments = append(comments, comment)
	}
	r.Close()

	return comments, err
}

func (p *MysqlProvider) GetCommentsByPostID(postID string) ([]model.Comment, error) {
	comment := model.Comment{}
	comments := []model.Comment{}

	r, err := p.connector.Query(`SELECT id, text, writerID, admin, postID FROM comment WHERE postID = ` + postID)

	for r.Next() {
		r.Scan(&comment.ID, &comment.Text, &comment.WriterID, &comment.Admin, &comment.PostID)
		comments = append(comments, comment)
	}
	r.Close()

	return comments, err
}

func (p *MysqlProvider) GetCommentPasswordByCommentID(commentID string) (string, error) {
	commentPW := ""
	r := p.connector.QueryRow(`SELECT writerPW FROM comment WHERE id = ` + commentID)
	err := r.Scan(&commentPW)
	return commentPW, err
}

func (p *MysqlProvider) DeleteCommentByCommentID(commentID string) error {
	_, err := p.connector.Exec(`DELETE FROM comment WHERE id =` + commentID)
	return err
}

func (p *MysqlProvider) CreateComment(comment model.Comment, admin string) error {
	fmt.Println(comment)
	_, err := p.connector.Exec(`INSERT INTO comment (text, writerID, writerPW, admin, postID) VALUES ('` + comment.Text + `', '` + comment.WriterID + `', '` + comment.WriterPW + `', ` + admin + `, ` + strconv.Itoa(comment.PostID) + `)`)
	return err
}

func (p *MysqlProvider) GetRepliesByPostID(postID string) ([]model.Reply, error) {
	reply := model.Reply{}
	replies := []model.Reply{}

	r, err := p.connector.Query(`SELECT id, admin, writerID, text, commentID FROM reply WHERE postID = ` + postID)
	for r.Next() {
		r.Scan(&reply.ID, &reply.Admin, &reply.WriterID, &reply.Text, &reply.CommentID)
		replies = append(replies, reply)
	}

	return replies, err
}

func (p *MysqlProvider) DeleteReplyByReplyID(replyID string) error {
	_, err := p.connector.Exec(`DELETE FROM reply WHERE id =` + replyID)
	return err
}

func (p *MysqlProvider) GetReplyPasswordByReplyID(replyID string) (string, error) {
	password := ""
	r := p.connector.QueryRow(`SELECT writerPW FROM reply WHERE id = ` + replyID)
	err := r.Scan(&password)
	return password, err
}

func (p *MysqlProvider) CreateReply(reply model.Reply) error {
	_, err := p.connector.Exec(`INSERT INTO reply (admin, writerID, writerPW, text, commentID, postID) VALUES (` + reply.Admin + `, '` + reply.WriterID + `', '` + reply.WriterPW + `', '` + reply.Text + `', ` + strconv.Itoa(reply.CommentID) + `, ` + strconv.Itoa(reply.PostID) + `)`)

	return err
}

func (p *MysqlProvider) StoreInitialPost(post model.Post) error {
	_, err := p.connector.Exec(`INSERT INTO post (id, tags, title, text, writeTime) VALUES (` + strconv.Itoa(post.ID) + `, '` + post.Tags + `', '` + post.Title + `', '` + post.Text + `', '` + post.WriteTime + `')`)

	return err
}

func (p *MysqlProvider) GetImagesByPostID(postID string) ([]model.PostImage, error) {
	image := model.PostImage{}
	images := []model.PostImage{}
	r, err := p.connector.Query(`SELECT id, imageName, thumbnail FROM image WHERE postID = ` + postID + ` ORDER BY thumbnail DESC`)
	for r.Next() {
		r.Scan(&image.ID, &image.ImageName, &image.Thumbnail)
		images = append(images, image)
	}
	defer r.Close()
	return images, err
}

func (p *MysqlProvider) DeleteImagesByPostID(postID string) error {
	_, err := p.connector.Exec(`DELETE FROM image WHERE postID = ` + postID)
	return err
}

func (p *MysqlProvider) CreatePostImageByPostID(postID string, image model.PostImage) error {
	_, err := p.connector.Exec(`INSERT INTO image (postID, imageName, thumbnail) VALUES (` + postID + `, "` + image.ImageName + `", ` + image.Thumbnail + `)`)
	return err
}

func (p *MysqlProvider) DeletePostImagesByPostID(postID string) error {
	_, err := p.connector.Exec(`DELETE FROM image WHERE postID = '` + postID + `'`)
	return err
}

func (p *MysqlProvider) GetPostImageNameByImageID(imageID string) (string, error) {
	imageName := ""
	r, err := p.connector.Query(`SELECT imageName FROM image WHERE id = '` + imageID + `'`)
	r.Next()
	r.Scan(&imageName)
	return imageName, err
}

func (p *MysqlProvider) IsDatabaseEmpty() bool {
	i := 0
	r := p.connector.QueryRow(`SELECT id FROM post LIMIT 1`)
	if err := r.Scan(i); err == sql.ErrNoRows {
		return true
	}
	return false
}

func (p *MysqlProvider) GetImageNamesByPostID(postID string) ([]string, error) {
	imageNames := []string{}
	imageName := ""
	r, err := p.connector.Query(`SELECT imageName FROM image WHERE postID = ` + postID)
	for r.Next() {
		r.Scan(&imageName)
		imageNames = append(imageNames, imageName)
	}
	return imageNames, err
}

func (p *MysqlProvider) DeleteImagesByImageName(name string) error {
	_, err := p.connector.Exec(`DELETE FROM image WHERE imageName = "` + name + `"`)
	return err
}
