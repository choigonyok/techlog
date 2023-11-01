package database

import (
	"database/sql"
	"fmt"
)

type Connector struct {
	driverName   string
	password     string
	user         string
	port         string
	host         string
	databaseName string
}

func (c *Connector) Open() (*sql.DB, error) {
	db, err := sql.Open(c.driverName, c.user+":"+c.password+"@tcp("+c.host+")/"+c.databaseName)

	fmt.Println(c.user + ":" + c.password + "@tcp(" + c.host + ")/" + c.databaseName)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (c *Connector) Close(db *sql.DB) {
	db.Close()
}

func New(driver, password, user, port, host, databasename string) *Connector {
	newConnector := &Connector{
		driverName:   driver,
		password:     password,
		user:         user,
		port:         port,
		host:         host,
		databaseName: databasename,
	}
	return newConnector
}
