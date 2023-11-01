package database

import (
	"database/sql"
	"fmt"
)

type Database struct {
	driverName   string
	password     string
	user         string
	port         string
	host         string
	databaseName string
}

var (
	DB *sql.DB
)

func (c *Database) Open() (*sql.DB, error) {
	var err error
	DB, err = sql.Open(c.driverName, c.user+":"+c.password+"@tcp("+c.host+")/"+c.databaseName)
	if err != nil {
		return nil, err
	}

	fmt.Println(c.user + ":" + c.password + "@tcp(" + c.host + ")/" + c.databaseName)
	return DB, nil
}

func (c *Database) Close(db *sql.DB) {
	db.Close()
}

func NewDatabase(driver, password, user, port, host, databasename string) *Database {
	newDatabase := &Database{
		driverName:   driver,
		password:     password,
		user:         user,
		port:         port,
		host:         host,
		databaseName: databasename,
	}
	return newDatabase
}

func GetConnector() *sql.DB {
	return DB
}
