package main

import (
	"fmt"
	"os"

	"github.com/choigonyok/techlog/pkg/database"
	"github.com/choigonyok/techlog/pkg/handler"
	"github.com/choigonyok/techlog/pkg/middleware"
	"github.com/choigonyok/techlog/pkg/router"
	"github.com/choigonyok/techlog/pkg/server"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	databaseDriver   = "mysql"
	databasePassword = os.Getenv("DB_PASSWORD")
	databaseUser     = os.Getenv("DB_USER")
	databasePort     = os.Getenv("DB_PORT")
	databaseHost     = os.Getenv("DB_HOST")
	databaseName     = os.Getenv("DB_NAME")
	handlerPrefix    = "/api/"
	listenAddress    = "0.0.0.0:8080" // should be 0.0.0.0, not localhost
	allowOrigin      = []string{"*"}
	allowMethods     = []string{"GET", "POST", "DELETE", "PUT"}
	allowHeaders     = []string{"Content-type"}
	allowCredentials = true
)

func main() {
	godotenv.Load(".env")

	connector := database.New(databaseDriver, databasePassword, databaseUser, databasePort, databaseHost, databaseName)

	db, _ := connector.Open()
	defer connector.Close(db)

	middleware := &middleware.Middleware{}
	middleware.AllowConfig(allowOrigin, allowMethods, allowHeaders, allowCredentials)

	// router.New(middleware)
	router := router.New(middleware)

	handlers := handler.NewHandlers(handlerPrefix)
	router.SetHandlers(handlers)
	httpHandlers := router.GetHTTPHandlers()

	server := server.New(httpHandlers, listenAddress)
	err := server.Start()
	if err != nil {
		fmt.Println("server starting error...", err.Error())
	}
}
