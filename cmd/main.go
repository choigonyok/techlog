package main

import (
	"fmt"

	"github.com/choigonyok/techlog/internal/server"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	databaseDriver = "mysql"
)

func main() {
	godotenv.Load(".env")

	server, err := server.New()
	if err != nil {
		fmt.Println("server creating error...", err)
	}
	err = server.Start()
	if err != nil {
		fmt.Println("server starting error...", err.Error())
	}
}
