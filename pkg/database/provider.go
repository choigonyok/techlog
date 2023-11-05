package database

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/choigonyok/techlog/pkg/model"
)

type Provider interface {
	ResetVisitorTodayAndDate(today string) error
	GetVisitor() (model.Visitor, error)
	UpdateVisitorToday(newToday, newTotal int) error
	GetEveryTag() ([]model.PostTags, error)
	GetEveryCardByTag(string) ([]model.PostCard, error)
	GetEveryCard() ([]model.PostCard, error)
	GetPostByID(postID string) ([]model.Post, error)
	GetThumbnailNameByPostID(postID string) (string, error)
	SetNewCookieValueByUniqueID(uniqueID string) error
	GetCookieValue() (string, error)
	UpdatePost(post model.Post) error
	DeletePostByPostID(postID string) ([]string, error)
	CreatePost(post model.Post) (int, error)
	StoreImage(image model.PostImage) error
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
	r, err := p.connector.Query(`SELECT today, total, date FROM visitor`)
	r.Next()
	r.Scan(&result.Today, &result.Total, &result.Date)
	defer r.Close()
	return result, err
}

func (p *MysqlProvider) UpdateVisitorToday(newToday, newTotal int) error {
	_, err := p.connector.Exec(`UPDATE visitor SET today = ` + strconv.Itoa(newToday) + `, total = ` + strconv.Itoa(newTotal))
	return err
}

func (p *MysqlProvider) GetEveryTag() ([]model.PostTags, error) {
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

func (p *MysqlProvider) GetEveryCardByTag(tag string) ([]model.PostCard, error) {
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

func (p *MysqlProvider) GetEveryCard() ([]model.PostCard, error) {
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

func (p *MysqlProvider) GetPostByID(postID string) ([]model.Post, error) {
	post := model.Post{}
	posts := []model.Post{}

	r, err := p.connector.Query(`SELECT id, tags, title, text, writeTime FROM post WHERE id = ` + postID)
	for r.Next() {
		r.Scan(&post.ID, &post.Tags, &post.Title, &post.Text, &post.WriteTime)
		posts = append(posts, post)
	}
	defer r.Close()
	return posts, err
}

func (p *MysqlProvider) GetThumbnailNameByPostID(postID string) (string, error) {
	var thumbnailName string
	r, err := p.connector.Query(`SELECT imageName FROM image WHERE thumbnail = 1 AND postID = ` + postID)
	r.Next()
	r.Scan(&thumbnailName)
	defer r.Close()
	return thumbnailName, err
}

func (p *MysqlProvider) SetNewCookieValueByUniqueID(uniqueID string) error {
	value, err := p.GetCookieValue()
	if err != nil {
		return err
	}
	if value == "" {
		_, err = p.connector.Exec(`INSERT INTO cookie (value) VALUES ("` + uniqueID + `")`)
	} else {
		_, err = p.connector.Exec(`UPDATE cookie SET value = "` + uniqueID + `"`)
	}

	return err
}

func (p *MysqlProvider) GetCookieValue() (string, error) {
	var value string
	r, err := p.connector.Query(`SELECT value FROM cookie`)
	r.Next()
	r.Scan(&value)
	defer r.Close()
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

	p.connector.Exec(`INSERT INTO post (tags, title, text, writeTime) VALUES ('` + post.Tags + `', '` + post.Title + `', '` + post.Text + `', '` + post.WriteTime + `')`)

	r, err := p.connector.Query(`SELECT id FROM post WHERE tags = '` + post.Tags + `' AND title = '` + post.Title + `' AND text = '` + post.Text + `' AND writeTime = '` + post.WriteTime + `' ORDER BY id DESC LIMIT 1`)

	r.Next()
	r.Scan(&postID)
	r.Close()

	return postID, err
}

func (p *MysqlProvider) StoreImage(image model.PostImage) error {
	_, err := p.connector.Exec(`INSERT INTO image (postID, imageName, thumbnail) VALUES (` + strconv.Itoa(image.PostID) + `, '` + image.ImageName + `', '` + image.Thumbnail + `')`)

	return err
}
