package database

import (
	"database/sql"
	"fmt"
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
	return tags, err
}

func (p *MysqlProvider) GetEveryCardByTag(tag string) ([]model.PostCard, error) {
	card := model.PostCard{}
	cards := []model.PostCard{}
	r, err := p.connector.Query(`SELECT id, tags, title, writeTime, thumbnailPath FROM post WHERE tags LIKE "%` + tag + `%" ORDER BY writeTime desc`)
	for r.Next() {
		r.Scan(&card.ID, &card.Tags, &card.Title, &card.WriteTime, &card.ThumbnailPath)
		cards = append(cards, card)
	}
	return cards, err
}

func (p *MysqlProvider) GetEveryCard() ([]model.PostCard, error) {
	card := model.PostCard{}
	cards := []model.PostCard{}
	r, err := p.connector.Query(`SELECT id, tags, title, writeTime, thumbnailPath FROM post WHERE tags ORDER BY writeTime desc`)
	for r.Next() {
		r.Scan(&card.ID, &card.Tags, &card.Title, &card.WriteTime, &card.ThumbnailPath)
		cards = append(cards, card)
	}
	return cards, err
}

func (p *MysqlProvider) GetPostByID(postID string) ([]model.Post, error) {
	post := model.Post{}
	posts := []model.Post{}

	r, err := p.connector.Query("SELECT id, tags, title, text, thumbnailPath, writeTime FROM post WHERE id = " + postID)
	for r.Next() {
		r.Scan(&post.ID, &post.Tags, &post.Title, &post.Text, &post.ThumbnailPath, &post.WriteTime)
		fmt.Println(post)
		posts = append(posts, post)
	}

	return posts, err
}
