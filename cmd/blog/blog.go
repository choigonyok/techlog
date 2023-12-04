package main

import (
	"fmt"
	"log"
	"os"

	"github.com/choigonyok/techlog/pkg/database"
	"github.com/choigonyok/techlog/pkg/github"
	"github.com/choigonyok/techlog/pkg/server"
	"github.com/choigonyok/techlog/pkg/service"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	databaseDriver = "mysql"
)

func main() {
	godotenv.Load(".env")
	var databasePassword = os.Getenv("DB_PASSWORD")
	var databaseSlaveUser = os.Getenv("DB_MASTER_USER")
	var databaseMasterUser = os.Getenv("DB_SLAVE_USER")
	var databasePort = os.Getenv("DB_PORT")
	var databaseMasterHost = os.Getenv("DB_HOST_READ")
	var databaseSlaveHost = os.Getenv("DB_HOST_WRITE")
	var databaseName = os.Getenv("DB_NAME")
	var githubToken = os.Getenv("GITHUB_TOKEN")

	master := database.New(databaseDriver, databasePassword, databaseMasterUser, databasePort, databaseMasterHost, databaseName)
	db0, err := master.Open()
	if err != nil {
		fmt.Println("Master Open Error")
	}
	defer master.Close(db0)

	slave := database.New(databaseDriver, databasePassword, databaseSlaveUser, databasePort, databaseSlaveHost, databaseName)
	haproxy, err := slave.OpenReadDB()
	if err != nil {
		fmt.Println("Haproxy Open Error")
	}
	defer slave.CloseReadDB(haproxy)

	pvr := database.NewMysqlProvider(database.GetConnector())
	svc := service.NewService(pvr)

	fmt.Println("master ping start")
	err = db0.Ping()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("master ping end")

	fmt.Println("slave ping start")
	err = haproxy.Ping()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("slave ping end")

	fmt.Println("test1")
	github.SyncGithubToken(githubToken)
	posts := github.GetPostsFromGithubRepo()
	for _, post := range posts {
		svc.StoreInitialPost(post)
		svc.StoreInitialPostImages(post)
	}
	fmt.Println("test2")

	server, err := server.New()
	if err != nil {
		log.Fatal("server creating error...", err)
	}
	err = server.Start()
	if err != nil {
		log.Fatal("server starting error...", err.Error())
	}
}
