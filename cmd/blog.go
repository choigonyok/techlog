package main

import (
	"fmt"
	"log"
	"os"

	"github.com/choigonyok/techlog/pkg/database"
	"github.com/choigonyok/techlog/pkg/server"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	databaseDriver = "mysql"
)

func main() {
	godotenv.Load(".env")

	var databasePassword = os.Getenv("DB_PASSWORD")
	var databaseUser = os.Getenv("DB_USER")
	var databasePort = os.Getenv("DB_PORT")
	var databaseHost = os.Getenv("DB_HOST")
	var databaseName = os.Getenv("DB_NAME")
	database := database.New(databaseDriver, databasePassword, databaseUser, databasePort, databaseHost, databaseName)
	db, err := database.Open()
	if err != nil {
		fmt.Println("Database Open Error")
	} else {
		fmt.Println("Database Opened")
	}
	defer database.Close(db)

	fmt.Println("ping start")
	err = db.Ping()
	if err != nil {
		fmt.Println("ping fail")
	}
	fmt.Println("ping end")

	server, err := server.New()
	if err != nil {
		log.Fatal("server creating error...", err)
	}
	err = server.Start()
	if err != nil {
		log.Fatal("server starting error...", err.Error())
	}
}
