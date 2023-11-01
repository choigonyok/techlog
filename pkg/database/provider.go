package database

import (
	"database/sql"
	"strconv"

	"github.com/choigonyok/techlog/pkg/model"
)

type Provider interface {
	GetVisitorDate() (string, error)
	ResetVisitorTodayAndDate(today string) error
	GetVisitor() (model.Visitor, error)
	UpdateVisitorToday(newToday int) error
}

type mysqlProvider struct {
	connector *sql.DB
}

func NewMysqlProvider(db *sql.DB) *mysqlProvider {
	return &mysqlProvider{
		connector: db,
	}
}

func (p *mysqlProvider) GetVisitorDate() (string, error) {
	var data string
	r, err := p.connector.Query(`SELECT date FROM visitor LIMIT 1`)
	r.Next()
	r.Scan(&data)
	defer r.Close()
	return data, err
}

func (p *mysqlProvider) ResetVisitorTodayAndDate(today string) error {
	_, err := p.connector.Exec(`UPDATE visitor SET today = 1, date = "` + today + `"`)
	return err
}

func (p *mysqlProvider) GetVisitor() (model.Visitor, error) {
	result := model.Visitor{}
	r, err := p.connector.Query(`SELECT today, total, date FROM visitor LIMIT 1`)
	r.Next()
	r.Scan(&result)
	defer r.Close()
	return result, err
}

func (p *mysqlProvider) UpdateVisitorToday(newToday int) error {
	_, err := p.connector.Exec(`UPDATE visitor SET today = ` + strconv.Itoa(newToday))
	return err
}
