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
	db     *sql.DB
	readDB *sql.DB
)

// Open opens database connector with specific use, password, host, db, driver
func (c *Database) Open() (*sql.DB, error) {
	tempDB, err := sql.Open(c.driverName, "root:"+c.password+"@tcp(mysql-ha-0)/"+c.databaseName)
	fmt.Println(c.user + ":" + c.password + "@tcp(" + c.host + ")/" + c.databaseName)
	db = tempDB
	return db, err
}

func (c *Database) OpenReadDB() (*sql.DB, error) {
	tempDB, err := sql.Open(c.driverName, "root:"+c.password+"@tcp(mysql-ha-haproxy:3306)/"+c.databaseName)
	fmt.Println(c.user + ":" + c.password + "@tcp(" + c.host + ")/" + c.databaseName)
	readDB = tempDB
	return readDB, err
}

// Close closed opened database connector
func (c *Database) Close(db *sql.DB) {
	db.Close()
}

func (c *Database) CloseReadDB(db *sql.DB) {
	readDB.Close()
}

// NewDatabase returns new database object
func New(driver, password, user, port, host, databasename string) *Database {
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

// GetConnector returns connector of opened database
func GetConnector() *sql.DB {
	return db
}

func GetReadConnector() *sql.DB {
	return readDB
}
