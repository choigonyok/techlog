package main

import (
	"fmt"
	"os"

	"github.com/choigonyok/techlog/pkg/database"
	"github.com/choigonyok/techlog/pkg/middleware"
	"github.com/choigonyok/techlog/pkg/router"
	"github.com/choigonyok/techlog/pkg/server"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	databaseDriver   = "mysql"
	handlerPrefix    = "/api/"
	listenAddress    = "0.0.0.0:8080" // should be 0.0.0.0, not localhost
	allowOrigin      = []string{"http://localhost", "http://localhost:3000", "http://frontend", "http://fronend:3000"}
	allowMethods     = []string{"GET", "POST", "DELETE", "PUT"}
	allowHeaders     = []string{"Content-type"}
	allowCredentials = true
)

func main() {
	godotenv.Load(".env")

	var databasePassword = os.Getenv("DB_PASSWORD")
	var databaseUser = os.Getenv("DB_USER")
	var databasePort = os.Getenv("DB_PORT")
	var databaseHost = os.Getenv("DB_HOST")
	var databaseName = os.Getenv("DB_NAME")

	database := database.NewDatabase(databaseDriver, databasePassword, databaseUser, databasePort, databaseHost, databaseName)
	db, _ := database.Open()
	defer database.Close(db)

	middleware := &middleware.Middleware{}
	middleware.AllowConfig(allowOrigin, allowMethods, allowHeaders, allowCredentials)

	router := router.NewRouter(middleware)
	routes := router.NewRoutes(handlerPrefix)
	router.SetRoutes(routes)
	httpHandler := router.GetHTTPHandler()

	server := server.New(httpHandler, listenAddress)
	err := server.Start()
	if err != nil {
		fmt.Println("server starting error...", err.Error())
	}
}
