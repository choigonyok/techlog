package database

import (
	"database/sql"
	"fmt"
)

type Connector struct {
	DriverName   string
	Password     string
	User         string
	Port         string
	Host         string
	DatabaseName string
}

func (c *Connector) Open() (*sql.DB, error) {
	db, err := sql.Open(c.DriverName, c.User+":"+c.Password+"@tcp("+c.Host+")/"+c.DatabaseName)

	fmt.Println(c.User + ":" + c.Password + "@tcp(" + c.Host + ")/" + c.DatabaseName)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (c *Connector) Close(db *sql.DB) {
	db.Close()
}
